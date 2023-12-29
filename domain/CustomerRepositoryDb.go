package domain

import (
	"database/sql"
	"log"

	"github.com/edwins-leonardi/banking-lib/errs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

func (db CustomerRepositoryDb) FindAll(filter *StatusFilter) ([]Customer, *errs.AppError) {
	filterSql := ""
	if filter != nil {
		if filter.Active {
			filterSql = " where status = 1"
		} else {
			filterSql = " where status = 0"
		}
	}
	findAllSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers" + filterSql

	var customers []Customer = make([]Customer, 0)
	err := db.client.Select(&customers, findAllSql)
	if err != nil {
		log.Println("Error while querying customer table", err.Error())
		return nil, errs.NewUnexpectedError("unexpected error retrieving data")
	}

	return customers, nil
}

func (db CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	sqlQuery := "select customer_id, name, city, zipcode, date_of_birth, status" +
		" from customers where customer_id = ?"

	var c Customer
	err := db.client.Get(&c, sqlQuery, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("customer not found")
		}
		log.Println("Error while scanning customer row result", err.Error())
		return nil, errs.NewUnexpectedError("unexpected error")
	}
	return &c, nil
}

func NewCustomerRepositoryDb(client *sqlx.DB) CustomerRepositoryDb {
	return CustomerRepositoryDb{
		client: client,
	}
}
