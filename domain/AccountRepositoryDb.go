package domain

import (
	"database/sql"
	"strconv"

	"github.com/edwins-leonardi/banking-lib/errs"
	"github.com/edwins-leonardi/banking-lib/logger"
	"github.com/jmoiron/sqlx"
)

type AccountRepositoryDb struct {
	client *sqlx.DB
}

func (db AccountRepositoryDb) Save(a Account) (*Account, *errs.AppError) {
	insertSQL := "INSERT INTO accounts (customer_id, opening_date, account_type, amount, status) values (?,?,?,?,?)"

	result, err := db.client.Exec(insertSQL, a.CustomerId, a.OpeningDate, a.AccountType, a.Amount, a.Status)
	if err != nil {
		logger.Error("Error trying to save new account " + err.Error())
		return nil, errs.NewUnexpectedError("error trying to save new account")
	}
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error retrieving id for new account " + err.Error())
		return nil, errs.NewUnexpectedError("error trying to save new account")
	}
	a.AccountId = strconv.FormatInt(id, 10)
	return &a, nil
}

func (db AccountRepositoryDb) FindBy(accountId string) (*Account, *errs.AppError) {
	selectSQL := "SELECT account_id, customer_id, opening_date, account_type, amount, status FROM accounts where account_id = ?"
	var a Account
	err := db.client.Get(&a, selectSQL, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("account not found: " + accountId + " Error: " + err.Error())
			return nil, errs.NewNotFoundError("account not found")
		}
		logger.Error("Error while scanning account row result account_id = " + accountId + " Error: " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected error retrieving account")
	}
	return &a, nil
}

func NewAccountRepositoryDb(client *sqlx.DB) AccountRepositoryDb {
	return AccountRepositoryDb{client: client}
}

func (db AccountRepositoryDb) SaveTransaction(t Transaction) (*Transaction, *errs.AppError) {
	logger.Info("Starting transaction")
	tx, err := db.client.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	insertSQL := "INSERT INTO transactions (account_id, amount, transaction_type, transaction_date) values (?,?,?,?)"

	result, _ := tx.Exec(insertSQL, t.AccountId, t.Amount, t.TransactionType, t.TransactionDate)
	if t.IsWithdrawal() {
		logger.Info("Withdrawing from account. Amount of " + strconv.FormatFloat(t.Amount, 'E', -1, 64))
		_, err = tx.Exec("UPDATE accounts SET amount = amount - ? where account_id = ?", t.Amount, t.AccountId)
	} else {
		logger.Info("Deposit to account. Amount of " + strconv.FormatFloat(t.Amount, 'E', -1, 64))
		_, err = tx.Exec("UPDATE accounts SET amount = amount + ? where account_id = ?", t.Amount, t.AccountId)
	}
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction " + err.Error())
		return nil, errs.NewUnexpectedError("Error processing transaction")
	}
	err = tx.Commit()
	logger.Info("Comitting transaction")
	if err != nil {
		logger.Error("Error while commiting transaction " + err.Error())
		return nil, errs.NewUnexpectedError("Error processing transaction")
	}
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error retrieving id for transaction " + err.Error())
		return nil, errs.NewUnexpectedError("Error processing transaction")
	}
	account, appErr := db.FindBy(t.AccountId)
	if appErr != nil {
		return nil, appErr
	}
	t.TransactionId = strconv.FormatInt(id, 10)
	t.Amount = account.Amount
	return &t, nil
}
