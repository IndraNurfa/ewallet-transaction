package services

import (
	"context"
	"encoding/json"
	"ewallet-transaction/constants"
	"ewallet-transaction/external"
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/interfaces"
	"ewallet-transaction/internal/models"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type TransactionService struct {
	TransactionRepo interfaces.ITransactionRepo
	External        interfaces.IExternal
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *models.Transaction) (models.CreateTransactionResponse, error) {
	var (
		resp models.CreateTransactionResponse
	)

	req.TransactionStatus = constants.TransactionStatusPending
	req.Reference = helpers.GenerateReference()

	jsonAdditionalInfo := map[string]interface{}{}
	if req.AddtionalInfo != "" {
		err := json.Unmarshal([]byte(req.AddtionalInfo), &jsonAdditionalInfo)
		if err != nil {
			return resp, errors.Wrap(err, "failed to unmarshal current additional info")
		}
	}

	err := s.TransactionRepo.CreateTransaction(ctx, req)
	if err != nil {
		return resp, errors.Wrap(err, "failed to insert create transaction")
	}

	resp.Reference = req.Reference
	resp.TransactionStatus = req.TransactionStatus

	return resp, nil
}

func (s *TransactionService) UpdateStatusTransaction(ctx context.Context, tokenData models.TokenData, req *models.UpdateStatusTransaction) error {

	// get transaction by reference
	trx, err := s.TransactionRepo.GetTransactionByReference(ctx, req.Reference, false)
	if err != nil {
		return errors.Wrap(err, "failed to get transaction")
	}

	// check transaction flow
	isValid := false
	mapStatusFlow := constants.MapTransactionStatusFlow[trx.TransactionStatus]
	for i := range mapStatusFlow {
		if mapStatusFlow[i] == req.TransactionStatus {
			isValid = true
		}
	}

	if !isValid {
		return fmt.Errorf("transaction status flow invalid. request satus = %s", req.TransactionStatus)
	}

	// request update balance
	reqUpdateBalance := external.UpdateBalance{
		Reference: req.Reference,
		Amount:    trx.Amount,
	}

	if req.TransactionStatus == constants.TransactionStatusReversed {
		reqUpdateBalance.Reference = "REVERSED-" + req.Reference
		now := time.Now()

		expiredReversalTime := trx.CreatedAt.Add(constants.MaximumReversalDuration)
		if now.After(expiredReversalTime) {
			return errors.New("reversal duration is already expired")
		}
	}

	var (
		errUpdateBalance error
		jwt              = tokenData.Token
	)

	switch trx.TransactionType {
	case constants.TransactionTypeTopup:
		switch req.TransactionStatus {
		case constants.TransactionStatusSuccess:
			_, errUpdateBalance = s.External.CreditBalance(ctx, jwt, reqUpdateBalance)
		case constants.TransactionStatusReversed:
			_, errUpdateBalance = s.External.DebitBalance(ctx, jwt, reqUpdateBalance)
		}
	case constants.TransactionTypePurchase:
		switch req.TransactionStatus {
		case constants.TransactionStatusSuccess:
			_, errUpdateBalance = s.External.DebitBalance(ctx, jwt, reqUpdateBalance)
		case constants.TransactionStatusReversed:
			_, errUpdateBalance = s.External.CreditBalance(ctx, jwt, reqUpdateBalance)
		}
	}

	if errUpdateBalance != nil {
		return errors.Wrap(errUpdateBalance, "failed to update balance")
	}

	// update additional info
	var (
		newAdditionalInfo     = map[string]interface{}{}
		currentAdditionalInfo = map[string]interface{}{}
	)

	if trx.AddtionalInfo != "" {
		err = json.Unmarshal([]byte(trx.AddtionalInfo), &currentAdditionalInfo)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal current additional info")
		}
	}

	if req.AddtionalInfo != "" {
		err = json.Unmarshal([]byte(req.AddtionalInfo), &newAdditionalInfo)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal new additional info")
		}
	}

	for key, val := range newAdditionalInfo {
		currentAdditionalInfo[key] = val
	}

	byteAdditionalInfo, err := json.Marshal(currentAdditionalInfo)
	if err != nil {
		return errors.Wrap(err, "failed to marshal current additional info")
	}

	// update status transaction
	err = s.TransactionRepo.UpdateStatusTransaction(ctx, req.Reference, req.TransactionStatus, string(byteAdditionalInfo))
	if err != nil {
		return errors.Wrap(err, "failed to update status transaction")
	}

	trx.TransactionStatus = req.TransactionStatus

	s.sendNotification(ctx, tokenData, trx)

	return nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, userID int) ([]models.Transaction, error) {
	return s.TransactionRepo.GetTransaction(ctx, userID)
}

func (s *TransactionService) GetTransactionDetail(ctx context.Context, reference string) (models.Transaction, error) {
	return s.TransactionRepo.GetTransactionByReference(ctx, reference, true)
}

func (s *TransactionService) RefundTransaction(ctx context.Context, tokenData *models.TokenData, req *models.RefundTransaction) (models.CreateTransactionResponse, error) {
	var (
		resp models.CreateTransactionResponse
		now  = time.Now()
	)

	trx, err := s.TransactionRepo.GetTransactionByReference(ctx, req.Reference, false)
	if err != nil {
		return resp, errors.Wrap(err, "failed to get transaction")
	}

	if trx.TransactionStatus != constants.TransactionStatusSuccess && trx.TransactionType != constants.TransactionTypePurchase {
		return resp, errors.New("current transaction is not success or transaction type is not purchase")
	}

	refundReference := "REFUND-" + req.Reference
	reqCreditBalance := external.UpdateBalance{
		Reference: refundReference,
		Amount:    trx.Amount,
	}

	_, err = s.External.CreditBalance(ctx, tokenData.Token, reqCreditBalance)
	if err != nil {
		return resp, errors.Wrap(err, "failed to credit balance")
	}

	transaction := models.Transaction{
		UserID:            int(tokenData.UserID),
		Amount:            trx.Amount,
		TransactionType:   constants.TransactionTypeRefund,
		TransactionStatus: constants.TransactionStatusSuccess,
		Reference:         refundReference,
		Description:       req.Description,
		AddtionalInfo:     req.AddtionalInfo,
		CreatedAt:         now,
		CreatedBy:         tokenData.FullName,
		UpdatedAt:         now,
		UpdatedBy:         tokenData.FullName,
	}

	err = s.TransactionRepo.CreateTransaction(ctx, &transaction)
	if err != nil {
		return resp, errors.Wrap(err, "failed to insert new transaction refund")
	}

	resp.Reference = refundReference
	resp.TransactionStatus = transaction.TransactionStatus

	s.sendNotification(ctx, *tokenData, transaction)

	return resp, nil

}

func (s *TransactionService) sendNotification(ctx context.Context, tokenData models.TokenData, trx models.Transaction) {
	if trx.TransactionType == constants.TransactionTypePurchase && trx.TransactionStatus == constants.TransactionStatusSuccess {
		err := s.External.SendNotification(ctx, tokenData.Email, "purchase_success", map[string]string{
			"full_name":   tokenData.FullName,
			"amount":      fmt.Sprintf("%.0f", trx.Amount),
			"reference":   trx.Reference,
			"description": trx.Description,
			"date":        trx.CreatedAt.Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			helpers.Logger.Warn("Failed to send notification: ", err)
		}
	} else if trx.TransactionType == constants.TransactionTypePurchase && trx.TransactionStatus == constants.TransactionStatusFailed {
		err := s.External.SendNotification(ctx, tokenData.Email, "purchase_failed", map[string]string{
			"full_name": tokenData.FullName,
			"amount":    fmt.Sprintf("%.0f", trx.Amount),
			"status":    "Purchase Failed",
			"reason":    trx.AddtionalInfo,
			"date":      trx.CreatedAt.Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			helpers.Logger.Warn("Failed to send notification: ", err)
		}
	} else if trx.TransactionType == constants.TransactionTypeTopup && trx.TransactionStatus == constants.TransactionStatusSuccess {
		err := s.External.SendNotification(ctx, tokenData.Email, "topup_success", map[string]string{
			"full_name": tokenData.FullName,
			"amount":    fmt.Sprintf("%.0f", trx.Amount),
			"reference": trx.Reference,
			"date":      trx.CreatedAt.Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			helpers.Logger.Warn("Failed to send notification: ", err)
		}
	} else if trx.TransactionType == constants.TransactionTypeTopup && trx.TransactionStatus == constants.TransactionStatusFailed {
		err := s.External.SendNotification(ctx, tokenData.Email, "topup_failed", map[string]string{
			"full_name": tokenData.FullName,
			"amount":    fmt.Sprintf("%.0f", trx.Amount),
			"status":    "TopUp Failed",
			"reason":    trx.AddtionalInfo,
			"date":      trx.CreatedAt.Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			helpers.Logger.Warn("Failed to send notification: ", err)
		}
	} else if trx.TransactionType == constants.TransactionTypeRefund && trx.TransactionStatus == constants.TransactionStatusSuccess {
		err := s.External.SendNotification(ctx, tokenData.Email, "refund", map[string]string{
			"full_name":   tokenData.FullName,
			"amount":      fmt.Sprintf("%.0f", trx.Amount),
			"reference":   trx.Reference,
			"description": trx.Description,
			"date":        trx.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			helpers.Logger.Warn("Failed to send notification: ", err)
		}
	} else if trx.TransactionType == constants.TransactionTypePurchase && trx.TransactionStatus == constants.TransactionStatusReversed {
		err := s.External.SendNotification(ctx, tokenData.Email, "purchase_reversed", map[string]string{
			"full_name": tokenData.FullName,
			"amount":    fmt.Sprintf("%.0f", trx.Amount),
			"reference": trx.Reference,
			"reason":    trx.AddtionalInfo,
			"date":      trx.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			helpers.Logger.Warn("Failed to send notification: ", err)
		}
	}
}
