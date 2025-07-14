package repository

import (
	"context"
	"ewallet-transaction/internal/models"

	"gorm.io/gorm"
)

type TransactionRepo struct {
	DB *gorm.DB
}

func (r *TransactionRepo) CreateTransaction(ctx context.Context, trx *models.Transaction) error {
	return r.DB.Create(trx).Error
}
