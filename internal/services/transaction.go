package services

import (
	"context"
	"ewallet-transaction/constants"
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/interfaces"
	"ewallet-transaction/internal/models"

	"github.com/pkg/errors"
)

type TransactionService struct {
	TransactionRepo interfaces.ITransactionRepo
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *models.Transaction) (models.CreateTransactionResponse, error) {
	var (
		resp models.CreateTransactionResponse
	)

	req.TransactionStatus = constants.TransactionStatusPending
	req.Reference = helpers.GenerateReference()

	err := s.TransactionRepo.CreateTransaction(ctx, req)
	if err != nil {
		return resp, errors.Wrap(err, "failed to insert create transaction")
	}

	resp.Reference = req.Reference
	resp.TransactionStatus = req.TransactionStatus

	return resp, nil
}
