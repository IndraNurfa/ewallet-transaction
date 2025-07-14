package constants

const (
	SuccessMessage      = "success"
	ErrFailedBadRequest = "data tidak sesuai"
	ErrServerError      = "terjadi kesalahan pada server"
)

const (
	TransactionStatusPending  = "PENDING"
	TransactionStatusSuccess  = "SUCCESS"
	TransactionStatusFailed   = "FAILED"
	TransactionStatusReversed = "REVERSERD"
)

const (
	TransactionTypeTopup    = "TOPUP"
	TransactionTypePurchase = "PURCHASE"
	TransactionTypeRefund   = "REFUND"
)

var MapTransactionType = map[string]bool{
	TransactionTypeTopup:    true,
	TransactionTypePurchase: true,
	TransactionTypeRefund:   true,
}
