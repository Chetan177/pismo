package restserver

import "github.com/labstack/echo/v4"

func (r *RestServer) loadAccountsRoutes(g *echo.Group) {
	accGroup := g.Group("/accounts")
	accGroup.GET("/:accId", r.getAccount)
	accGroup.POST("", r.createAccount)
}

func (r *RestServer) loadTransactionRoutes(g *echo.Group) {
	transactionGroup := g.Group("/transactions")
	transactionGroup.POST("", r.createTransaction)
}
