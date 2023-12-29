package domain

import (
	"github.com/edwins-leonardi/banking/dto"
)

const WITHDRAWAL = "withdrawal"

type Transaction struct {
	TransactionId   string  `json:"transaction_id"`
	AccountId       string  `json:"account_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	TransactionDate string  `json:"transaction_date"`
}

func (t Transaction) IsWithdrawal() bool {
	return t.TransactionType == WITHDRAWAL
}

func (a Transaction) ToResponseDto() dto.TransactionResponse {
	return dto.TransactionResponse{
		TransactionId: a.TransactionId,
		Balance:       a.Amount,
	}
}
