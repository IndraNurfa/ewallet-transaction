package interfaces

import (
	"context"
	"ewallet-transaction/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type ITransactionAPI interface {
	CreateTransaction(c *gin.Context)
	UpdateStatusTransaction(c *gin.Context)
}

type ITransactionService interface {
	CreateTransaction(ctx context.Context, req *models.Transaction) (models.CreateTransactionResponse, error)
	UpdateStatusTransaction(ctx context.Context, tokenData *models.TokenData, req *models.UpdateStatusTransaction) error
}

type ITransactionRepo interface {
	CreateTransaction(ctx context.Context, trx *models.Transaction) error
	GetTransactionByReference(ctx context.Context, reference string, includeRefund bool) (models.Transaction, error)
	UpdateStatusTransaction(ctx context.Context, reference, status, additional_info string, now time.Time) error
}
