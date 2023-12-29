package domain

type CustomerRepositoryStub struct {
	customers []Customer
}

func (s CustomerRepositoryStub) FindAll() ([]Customer, error) {
	return s.customers, nil
}

func NewCustomerRepositoryStub() CustomerRepositoryStub {
	customers := []Customer{
		{Id: "1001", Name: "John", City: "Dublin", ZipCode: "4103", DateOfBirth: "17/08/1983", Status: "Active"},
		{Id: "1002", Name: "Janice", City: "Dublin", ZipCode: "1983", DateOfBirth: "29/09/1993", Status: "Active"},
	}
	return CustomerRepositoryStub{customers: customers}
}
