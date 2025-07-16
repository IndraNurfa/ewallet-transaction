package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Transaction struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	Amount            float64   `json:"amount" gorm:"column:amount;type:decimal(15,2)" validate:"required"`
	TransactionType   string    `json:"transaction_type" gorm:"column:transaction_type;type:enum('TOPUP','PURCHASE','REFUND')" validate:"required"`
	TransactionStatus string    `json:"transaction_status" gorm:"column:transaction_status;type:enum('PENDING','SUCCESS','FAILED','REVERSED')"`
	Reference         string    `json:"reference" gorm:"column:reference;type:varchar(255)"`
	Description       string    `json:"description" gorm:"column:description;type:varchar(255)" validate:"required"`
	AddtionalInfo     string    `json:"additional_info" gorm:"column:additional_info;type:text"`
	CreatedAt         time.Time `json:"date"`
	CreatedBy         string    `json:"-" gorm:"column:created_by;type:varchar(255)"`
	UpdatedAt         time.Time `json:"-"`
	UpdatedBy         string    `json:"-" gorm:"column:updated_by;type:varchar(255)"`
}

func (*Transaction) TableName() string {
	return "transactions"
}

func (l Transaction) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type CreateTransactionResponse struct {
	Reference         string `json:"reference"`
	TransactionStatus string `json:"transaction_status"`
}

type UpdateStatusTransaction struct {
	Reference         string `json:"reference"`
	TransactionStatus string `json:"transaction_status" validate:"required"`
	AddtionalInfo     string `json:"additional_info"`
}

func (l UpdateStatusTransaction) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type RefundTransaction struct {
	Reference     string `json:"reference" validate:"required"`
	Description   string `json:"description" validate:"required"`
	AddtionalInfo string `json:"additional_info"`
}

func (l RefundTransaction) Validate() error {
	v := validator.New()
	return v.Struct(l)
}
