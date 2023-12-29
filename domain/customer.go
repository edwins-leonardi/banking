package domain

import (
	"github.com/edwins-leonardi/banking-lib/errs"
	"github.com/edwins-leonardi/banking/dto"
)

type Customer struct {
	Id          string `db:"customer_id"`
	Name        string
	City        string
	ZipCode     string
	DateOfBirth string `db:"date_of_birth"`
	Status      string
}

func (c Customer) ToDto() dto.CustomerResponse {
	return dto.CustomerResponse{
		Id:          c.Id,
		Name:        c.Name,
		City:        c.City,
		ZipCode:     c.ZipCode,
		DateOfBirth: c.DateOfBirth,
		Status:      c.statusAsText(),
	}
}

func (c Customer) statusAsText() string {
	statusValue := "active"
	if c.Status == "0" {
		statusValue = "inactive"
	}
	return statusValue
}

type StatusFilter struct {
	Active bool
}

type CustomerRepository interface {
	FindAll(filter *StatusFilter) ([]Customer, *errs.AppError)
	ById(id string) (*Customer, *errs.AppError)
}
