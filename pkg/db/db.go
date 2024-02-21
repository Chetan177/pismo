package db

import "github.com/Chetan177/pismo/pkg/db/model"

// DB interface for DB implementation
type DB interface {
	CreateAccount(account *model.Account) (string, error)
	GetAccount(id string) (*model.Account, error)
	IsAccountExsists(id string) bool
	CreateTransaction(transaction *model.Transaction) (string, error)
	ShutDown()
}
