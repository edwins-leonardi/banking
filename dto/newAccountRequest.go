package dto

import (
	"strings"

	"github.com/edwins-leonardi/banking-lib/errs"
)

type NewAccountRequest struct {
	CustomerId  string  `json:"customer_Id"`
	AccountType string  `json:"account_type"`
	Amount      float64 `json:"amount"`
}

func (req NewAccountRequest) Validate() *errs.AppError {
	if req.Amount < 2000 {
		return errs.NewValidationError("You must deposit at least 2000.00 to open an account")
	}
	if strings.ToLower(req.AccountType) != "savings" && strings.ToLower(req.AccountType) != "checking" {
		return errs.NewValidationError("Not valid account type. Only 'checking' or 'savings' accounts are valid")
	}
	return nil
}
