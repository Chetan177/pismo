package api

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" validate:"required"`
}

type CreateAccountResponse struct {
	AccountID string `json:"account_id"`
}

type GetAccountResponse struct {
	AccountID      string `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type CreateTransactionRequest struct {
	AccountID       string  `json:"account_id" validate:"required"`
	OperationTypeID int     `json:"operation_type_id" validate:"required,gte=1,lte=4"`
	Amount          float32 `json:"amount" validate:"required,gt=0"`
}

type CreateTransactionResponse struct {
	TransactionId string `json:"transaction_id"`
}

// DB models

type Account struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
	DocumentNumber string             `bson:"document_number"`
	TimeStamp      string             `bson:"time_stamp"`
}

type Transaction struct {
	Id              primitive.ObjectID `bson:"_id,omitempty"`
	AccountID       string             `bson:"account_id"`
	OperationTypeId int                `bson:"operation_type_id"`
	Amount          float32            `bson:"amount"`
	TimeStamp       string             `bson:"time_stamp"`
	Balance         float32            `bson:"balance"`
}
