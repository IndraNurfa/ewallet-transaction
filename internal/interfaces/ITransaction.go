package interfaces

import (
	"context"
	"ewallet-transaction/internal/models"

	"github.com/gin-gonic/gin"
)

type ITransactionAPI interface {
	CreateTransaction(c *gin.Context)
}

type ITransactionService interface {
	CreateTransaction(ctx context.Context, req *models.Transaction) (models.CreateTransactionResponse, error)
}

type ITransactionRepo interface {
	CreateTransaction(ctx context.Context, trx *models.Transaction) error
}
