package model

import "go.mongodb.org/mongo-driver/bson/primitive"

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
}
