package dto

import (
	"net/http"
	"testing"
)

func Test_should_return_error_when_transaction_type_is_not_valid(t *testing.T) {
	// Arrange
	request := TransactionRequest{
		TransactionType: "not valid",
	}
	// Act
	appError := request.Validate()

	//Assert
	if appError.Message != "Not valid transaction type. Only withdrawal or deposit are valid" {
		t.Error("Invalid message while testing transation type")
	}
	if appError.Code != http.StatusUnprocessableEntity {
		t.Error("Invalid error code while testing transation type")
	}
}

func Test_should_return_error_when_amount_is_less_than_zero(t *testing.T) {
	// arrange
	request := TransactionRequest{
		TransactionType: WITHDRAWAL,
		Amount:          -10,
	}

	// act
	appError := request.Validate()

	// assert
	if appError.Message != "Amount cannot be negative" {
		t.Error("Invalid message while testing negative amount")
	}
	if appError.Code != http.StatusUnprocessableEntity {
		t.Error("Invalid error code while testing negative amount")
	}
}
