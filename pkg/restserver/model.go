package restserver

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type CreateAccountResponse struct {
	AccountID string `json:"account_id"`
}

type GetAccountResponse struct {
	AccountID      string `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type CreateTransactionRequest struct {
	AccountID       string  `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float32 `json:"amount"`
}

type CreateTransactionResponse struct {
	TransactionId string `json:"transaction_id"`
}
