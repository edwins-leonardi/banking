package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edwins-leonardi/banking-lib/errs"
	"github.com/edwins-leonardi/banking/dto"
	mock_service "github.com/edwins-leonardi/banking/mocks/service"
	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
)

var router *mux.Router
var ch CustomerHandlers
var mockService *mock_service.MockCustomerService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)

	mockService = mock_service.NewMockCustomerService(ctrl)

	ch = CustomerHandlers{
		service: mockService,
	}
	router = mux.NewRouter()
	router.HandleFunc("/customers", ch.getAllCustomers)

	return func() {
		router = nil
		mockService = nil
		ctrl.Finish()
	}
}

func Test_should_return_customers_with_status_code_200(t *testing.T) {
	// Arrange
	teardown := setup(t)
	defer teardown()

	dummyCustomers := []dto.CustomerResponse{
		{Id: "1001", Name: "John", City: "Dublin", ZipCode: "4103", DateOfBirth: "17/08/1983", Status: "Active"},
		{Id: "1002", Name: "Janice", City: "Dublin", ZipCode: "1983", DateOfBirth: "29/09/1993", Status: "Active"},
	}
	mockService.EXPECT().GetAllCustomers(nil).Return(dummyCustomers, nil)
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the 200 status code")
	}
}

func Test_should_return_status_code_500_with_error_message(t *testing.T) {
	// arrange
	teardown := setup(t)
	defer teardown()

	mockService.EXPECT().GetAllCustomers(nil).Return(nil, errs.NewUnexpectedError("unexpected error retrieving data"))
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// assert
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the 500 status code")
	}
}
