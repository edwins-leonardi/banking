package dto

import "github.com/edwins-leonardi/banking-lib/errs"

type TransactionRequest struct {
	AccountId       string  `json:"account_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	CustomerId      string  `json:"customer_id"`
}

const WITHDRAWAL = "withdrawal"

func (req TransactionRequest) Validate() *errs.AppError {
	if req.TransactionType != "withdrawal" && req.TransactionType != "deposit" {
		return errs.NewValidationError("Not valid transaction type. Only withdrawal or deposit are valid")
	}
	if req.Amount < 0.0 {
		return errs.NewValidationError("Amount cannot be negative")
	}
	return nil
}

func (req TransactionRequest) IsTransactionTypeWithdrawal() bool {
	return req.TransactionType == WITHDRAWAL
}
