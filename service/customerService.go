package service

import (
	"github.com/edwins-leonardi/banking-lib/errs"
	"github.com/edwins-leonardi/banking/domain"
	"github.com/edwins-leonardi/banking/dto"
)

//go:generate mockgen -destination=../mocks/service/mockCustomerService.go github.com/edwins-leonardi/banking/service CustomerService
type CustomerService interface {
	GetAllCustomers(filter *domain.StatusFilter) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func (s DefaultCustomerService) GetAllCustomers(filter *domain.StatusFilter) ([]dto.CustomerResponse, *errs.AppError) {
	customers, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	customersDto := make([]dto.CustomerResponse, 0)
	for _, c := range customers {
		customersDto = append(customersDto, c.ToDto())
	}
	return customersDto, nil
}

func (s DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}
	response := c.ToDto()
	return &response, nil
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{
		repo: repository,
	}
}
