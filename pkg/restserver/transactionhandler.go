package restserver

import (
	"net/http"

	"github.com/Chetan177/pismo/pkg/db/model"

	"github.com/labstack/echo/v4"
)

// enum for operational type
const (
	normalPurchase = iota + 1
	purchaseWithInstallment
	withdrawal
	creditVoucher
)

// createTransaction handler for post /transactions
func (r *RestServer) createTransaction(c echo.Context) error {
	req := new(CreateTransactionRequest)
	err := c.Bind(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	accExists := r.db.IsAccountExsists(req.AccountID)
	if !accExists {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "account don't exists"})
	}

	switch req.OperationTypeID {
	case purchaseWithInstallment, normalPurchase, withdrawal:
		req.Amount = -1 * req.Amount
	}

	transData := &model.Transaction{
		AccountID:       req.AccountID,
		OperationTypeId: req.OperationTypeID,
		Amount:          req.Amount,
	}
	id, err := r.db.CreateTransaction(transData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, &CreateTransactionResponse{TransactionId: id})
}
