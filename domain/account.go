package domain

import (
	"time"

	"github.com/edwins-leonardi/banking-lib/errs"
	"github.com/edwins-leonardi/banking/dto"
)

type Account struct {
	AccountId   string  `db:"account_id"`
	CustomerId  string  `db:"customer_id"`
	OpeningDate string  `db:"opening_date"`
	AccountType string  `db:"account_type"`
	Amount      float64 `db:"amount"`
	Status      string  `db:"status"`
}

//go:generate mockgen -destination=../mocks/domain/mockAccountRepository.go github.com/edwins-leonardi/banking/domain AccountRepository
type AccountRepository interface {
	Save(account Account) (*Account, *errs.AppError)
	FindBy(accountId string) (*Account, *errs.AppError)
	SaveTransaction(t Transaction) (*Transaction, *errs.AppError)
}

func (a Account) ToNewAccountResponseDto() *dto.NewAccountResponse {
	return &dto.NewAccountResponse{
		AccountId: a.AccountId,
	}
}

func (a Account) CanWithdraw(amount float64) bool {
	return a.Amount >= amount
}

func NewAccount(customerId string, accountType string, amount float64) Account {
	return Account{
		CustomerId:  customerId,
		OpeningDate: time.Now().Format("2006-01-02 15:04:05"),
		AccountType: accountType,
		Amount:      amount,
		Status:      "1",
	}
}
