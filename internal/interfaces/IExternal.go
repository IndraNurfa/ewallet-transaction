package interfaces

import (
	"context"
	"ewallet-transaction/external"
	"ewallet-transaction/internal/models"
)

type IExternal interface {
	ValidateToken(ctx context.Context, token string) (models.TokenData, error)
	CreditBalance(ctx context.Context, token string, req external.UpdateBalance) (*external.UpdateBalanceResponse, error)
	DebitBalance(ctx context.Context, token string, req external.UpdateBalance) (*external.UpdateBalanceResponse, error)
}
