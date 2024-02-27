package api

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	normalPurchase = iota + 1
	purchaseWithInstallment
	withdrawal
	creditVoucher
)

func (s *Server) bindAndValidate(c echo.Context, req interface{}) error {
	err := c.Bind(req)
	if err != nil {
		return err
	}
	err = s.validator.Struct(req)
	if err != nil {
		return err
	}
	return nil
}

// createAccount handler for post /accounts
func (s *Server) createAccount(c echo.Context) error {
	req := &CreateAccountRequest{}
	err := s.bindAndValidate(c, req)
	if err != nil {
		return s.logAndReturnResponse(c, err, http.StatusBadRequest, err.Error())
	}
	accountData := &Account{}

	filter := bson.M{"document_number": req.DocumentNumber}
	err = s.accountCollection.FindOne(context.Background(), filter).Decode(accountData)

	if err == mongo.ErrNoDocuments {
		accountData.DocumentNumber = req.DocumentNumber
		accountData.TimeStamp = time.Now().Format(time.RFC3339)
		result, err := s.accountCollection.InsertOne(context.Background(), accountData)
		if err != nil {
			return s.logAndReturnResponse(c, err, http.StatusInternalServerError, "not able to create account")
		}
		accountData.Id = result.InsertedID.(primitive.ObjectID)

	} else if err != nil {
		return s.logAndReturnResponse(c, err, http.StatusInternalServerError, "not able to create account")
	}

	return c.JSON(http.StatusOK, &CreateAccountResponse{AccountID: accountData.Id.Hex()})
}

// getAccount handler for get /accounts/:accId
func (s *Server) getAccount(c echo.Context) error {
	accID := c.Param("accId")
	if accID == "" {
		return s.logAndReturnResponse(c, nil, http.StatusBadRequest, "accId is not present in query params")
	}
	objectID, err := primitive.ObjectIDFromHex(accID)
	if err != nil {
		return s.logAndReturnResponse(c, err, http.StatusBadRequest, "accId is invalid")
	}

	filter := bson.M{"_id": objectID}
	accountData := Account{}
	err = s.accountCollection.FindOne(context.Background(), filter).Decode(&accountData)
	if err == mongo.ErrNoDocuments {
		return s.logAndReturnResponse(c, err, http.StatusNotFound, "accId not found")
	} else if err != nil {
		return s.logAndReturnResponse(c, err, http.StatusInternalServerError, "not able to create account")
	}

	return c.JSON(http.StatusOK, &GetAccountResponse{
		AccountID:      accountData.Id.Hex(),
		DocumentNumber: accountData.DocumentNumber,
	})
}

// createTransaction handler for post /transactions
func (s *Server) createTransaction(c echo.Context) error {
	req := new(CreateTransactionRequest)
	err := s.bindAndValidate(c, req)
	if err != nil {
		return s.logAndReturnResponse(c, err, http.StatusBadRequest, err.Error())
	}

	objectID, err := primitive.ObjectIDFromHex(req.AccountID)
	if err != nil {
		return s.logAndReturnResponse(c, err, http.StatusBadRequest, "accId is invalid")
	}

	filter := bson.M{"_id": objectID}
	accountData := &Account{}
	err = s.accountCollection.FindOne(context.Background(), filter).Decode(accountData)
	if err == mongo.ErrNoDocuments {
		return s.logAndReturnResponse(c, err, http.StatusNotFound, "accId not found")
	} else if err != nil {
		return s.logAndReturnResponse(c, err, http.StatusInternalServerError, "not able to create transaction")
	}

	switch req.OperationTypeID {
	case purchaseWithInstallment, normalPurchase, withdrawal:
		req.Amount = -1 * req.Amount
	case creditVoucher:
		result, err := s.processDischarges(req)
		if err != nil {
			return s.logAndReturnResponse(c, err, http.StatusInternalServerError, "not able to create transaction")
		}
		output := result.(*mongo.InsertOneResult)
		return c.JSON(http.StatusOK, &CreateTransactionResponse{TransactionId: output.InsertedID.(primitive.ObjectID).Hex()})
	}

	transData := Transaction{
		AccountID:       req.AccountID,
		OperationTypeId: req.OperationTypeID,
		Amount:          req.Amount,
		TimeStamp:       time.Now().Format(time.RFC3339),
		Balance:         req.Amount,
	}

	result, err := s.transactionCollection.InsertOne(context.Background(), transData)
	if err != nil {
		return s.logAndReturnResponse(c, err, http.StatusInternalServerError, "not able to create transaction")
	}
	transData.Id = result.InsertedID.(primitive.ObjectID)

	return c.JSON(http.StatusOK, &CreateTransactionResponse{TransactionId: transData.Id.Hex()})
}

func (s *Server) processDischarges(req *CreateTransactionRequest) (interface{}, error) {
	filter := bson.M{
		"account_id": req.AccountID,
		"balance":    bson.E{Key: "$lt", Value: 0},
	}

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.db.StartSession()
	if err != nil {
		return nil, err
	}

	result, err := session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		return s.processDischargeTranscation(ctx, filter, req)
	}, txnOptions)

	return result, nil
}

func (s *Server) processDischargeTranscation(ctx mongo.SessionContext, filter bson.M, newTarn *CreateTransactionRequest) (interface{}, error) {

	var transcations []Transaction
	cursor, err := s.transactionCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &transcations); err != nil {
		return nil, err
	}

	currBalance, err := s.doDischarge(transcations, newTarn)
	if err != nil {
		return nil, err
	}
	transData := Transaction{
		AccountID:       newTarn.AccountID,
		OperationTypeId: newTarn.OperationTypeID,
		Amount:          currBalance,
		TimeStamp:       time.Now().Format(time.RFC3339),
	}
	result, err := s.transactionCollection.InsertOne(context.Background(), transData)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Server) doDischarge(transcations []Transaction, newTarn *CreateTransactionRequest) (float32, error) {
	currBalance := newTarn.Amount
	flag := false
	for _, t := range transcations {
		t.Balance = currBalance - t.Balance
		if t.Balance > 0 {
			currBalance = t.Balance
			t.Balance = 0
		} else if t.Balance <= 0 {
			currBalance = 0
			flag = true
		}

		_, err := s.transactionCollection.UpdateByID(context.Background(), t.Id, t)
		if err != nil {
			return currBalance, err
		}

		if flag {
			break
		}
	}
	return currBalance, nil
}

func (s *Server) logAndReturnResponse(c echo.Context, err error, code int, message string) error {
	log.Printf("error api: %s, message: %s, status: %d, error: %+v ", c.Request().URL.String(), message, code, err)

	switch err.(type) {
	case nil:

	case *echo.HTTPError:
		errData := err.(*echo.HTTPError)
		message = errData.Message.(string)
	case validator.ValidationErrors:
		errData := err.(validator.ValidationErrors)
		message = ""
		for _, v := range errData {
			if message != "" {
				message += " | "
			}
			message += fmt.Sprintf("Validation failed for %s on tags: %s", v.Field(), v.Tag())
		}
	}
	return c.JSON(code, map[string]string{"message": message})
}
