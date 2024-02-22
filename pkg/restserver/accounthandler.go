package restserver

import (
	"net/http"

	"github.com/Chetan177/pismo/pkg/db/model"

	"github.com/labstack/echo/v4"
)

// createAccount handler for post /accounts
func (r *RestServer) createAccount(c echo.Context) error {
	req := &CreateAccountRequest{}
	err := c.Bind(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	if req.DocumentNumber == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "document number is missing"})
	}
	accountData := &model.Account{
		DocumentNumber: req.DocumentNumber,
	}

	id, err := r.db.CreateAccount(accountData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, &CreateAccountResponse{AccountID: id})
}

// getAccount handler for get /accounts/:accId
func (r *RestServer) getAccount(c echo.Context) error {
	accID := c.Param("accId")
	if accID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "accId is not present in query params"})
	}
	accountData, err := r.db.GetAccount(accID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, &GetAccountResponse{
		AccountID:      accountData.Id.Hex(),
		DocumentNumber: accountData.DocumentNumber,
	})
}
