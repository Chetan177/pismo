package api

import "github.com/labstack/echo/v4"

func (s *Server) loadAccountsRoutes(g *echo.Group) {
	accGroup := g.Group("/accounts")
	accGroup.GET("/:accId", s.getAccount)
	accGroup.POST("", s.createAccount)
}

func (s *Server) loadTransactionRoutes(g *echo.Group) {
	transactionGroup := g.Group("/transactions")
	transactionGroup.POST("", s.createTransaction)
}
