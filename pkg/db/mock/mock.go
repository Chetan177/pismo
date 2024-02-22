package mock

import (
	"errors"
	"github.com/Chetan177/pismo/pkg/db/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type Mock struct {
	dataMap map[string]interface{}
}

func NewMockDB(dataMap map[string]interface{}) *Mock {
	return &Mock{
		dataMap: dataMap,
	}
}

func (m *Mock) CreateAccount(account *model.Account) (string, error) {
	if strings.Contains(account.DocumentNumber, "fail") {
		return "", errors.New("failed case")
	}
	account.Id = primitive.NewObjectID()
	m.dataMap[account.Id.Hex()] = account
	return account.Id.Hex(), nil

}

func (m *Mock) GetAccount(id string) (*model.Account, error) {
	if strings.Contains(id, "fail") {
		return nil, errors.New("failed")
	}
	data, ok := m.dataMap[id].(*model.Account)
	if !ok {
		return nil, errors.New("failed")
	}
	return data, nil

}

func (m *Mock) IsAccountExsists(id string) bool {
	if strings.Contains(id, "fail") {
		return false
	}
	_, ok := m.dataMap[id]
	return ok
}

func (m *Mock) CreateTransaction(transaction *model.Transaction) (string, error) {
	if strings.Contains(transaction.AccountID, "fail") {
		return "", errors.New("failed case")
	}
	transaction.Id = primitive.NewObjectID()
	m.dataMap[transaction.Id.Hex()] = transaction
	return transaction.Id.Hex(), nil
}

func (m *Mock) ShutDown() {
}
