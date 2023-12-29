package service

import (
	"net/http"
	"testing"
	"time"

	"github.com/edwins-leonardi/banking-lib/errs"
	"github.com/edwins-leonardi/banking/domain"
	"github.com/edwins-leonardi/banking/dto"
	mock_domain "github.com/edwins-leonardi/banking/mocks/domain"
	"go.uber.org/mock/gomock"
)

var mockRepo *mock_domain.MockAccountRepository
var service AccountService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)

	mockRepo = mock_domain.NewMockAccountRepository(ctrl)
	service = NewAccountService(mockRepo)

	return func() {
		service = nil
		mockRepo = nil
		ctrl.Finish()
	}
}

func Test_should_return_validation_error_response_when_request_is_not_valid(t *testing.T) {
	// arrang
	dto := dto.NewAccountRequest{
		CustomerId:  "123",
		AccountType: "savings",
		Amount:      10.97,
	}
	service = NewAccountService(nil)

	// act
	_, err := service.NewAccount(dto)

	// assert
	if err == nil {
		t.Error("failed while testing new account validation")
	}
}

func Test_should_return_an_error_from_repository_account_cannot_be_created(t *testing.T) {
	// arrange
	teardown := setup(t)
	defer teardown()
	dto := dto.NewAccountRequest{
		CustomerId:  "123",
		AccountType: "savings",
		Amount:      3000,
	}
	mockRepo.EXPECT().Save(gomock.Any()).Return(nil, errs.NewUnexpectedError("error trying to save new account"))

	// act
	_, err := service.NewAccount(dto)
	if err.Code != http.StatusInternalServerError {
		t.Error("failed while testing new account failure")
	}
}

func Test_should_return_new_account_when_account_created_successfully(t *testing.T) {
	// arrange
	teardown := setup(t)
	defer teardown()
	dto := dto.NewAccountRequest{
		CustomerId:  "123",
		AccountType: "savings",
		Amount:      3000,
	}
	account := domain.Account{
		AccountId:   "42",
		CustomerId:  dto.CustomerId,
		OpeningDate: time.Now().Format("2006-01-02 15:04:05"),
		AccountType: dto.AccountType,
		Amount:      dto.Amount,
		Status:      "1",
	}
	mockRepo.EXPECT().Save(gomock.Any()).Return(&account, nil)

	// act
	accountResponse, err := service.NewAccount(dto)
	if err != nil {
		t.Error("failed while testing new account created successfully")
	}
	if account.AccountId != accountResponse.AccountId {
		t.Error("failed while testing new account created successfully")
	}
}
