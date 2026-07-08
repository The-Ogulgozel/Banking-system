package router

import (
	"github.com/The-Ogulgozel/Banking-system/internal/config"
	"github.com/The-Ogulgozel/Banking-system/internal/handlers"
	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	AccountsHandler    *handlers.AccountsHandler
	TransactionHandler *handlers.TransactionHandler
}

func NewRouter(cfg *config.Config, deps *RouterDeps) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		account := api.Group("/accounts")
		{
			account.POST("/", deps.AccountsHandler.Create)
			account.GET("/:id", deps.AccountsHandler.GetByID)
			account.GET("", deps.AccountsHandler.ListAccounts)
			account.DELETE("/:id", deps.AccountsHandler.Delete)

			account.POST("/:id/deposit", deps.TransactionHandler.Deposit)
			account.POST("/:id/withdraw", deps.TransactionHandler.Withdraw)
			account.POST("/:id/transfer", deps.TransactionHandler.Transfer)
			account.GET("/:id/transactions", deps.TransactionHandler.ListTransactionByID)
		}

	}
	return r
}
