package app

import (
	"encoding/json"
	"net/http"

	"github.com/edwins-leonardi/banking-lib/errs"
	"github.com/edwins-leonardi/banking/domain"
	"github.com/edwins-leonardi/banking/service"
	"github.com/gorilla/mux"
)

type CustomerHandlers struct {
	service service.CustomerService
}

func (ch *CustomerHandlers) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	var filter *domain.StatusFilter
	status := r.URL.Query().Get("status")
	if len(status) > 0 {
		if status == "active" {
			filter = &domain.StatusFilter{Active: true}
		} else if status == "inactive" {
			filter = &domain.StatusFilter{Active: false}
		}
	}
	customers, err := ch.service.GetAllCustomers(filter)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.AsMessage())
	} else {
		writeResponse(w, http.StatusOK, customers)
	}
}

func (ch *CustomerHandlers) getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId, ok := vars["customer_id"]
	if ok {
		customer, err := ch.service.GetCustomer(customerId)
		if err != nil {
			writeResponse(w, err.Code, err.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, customer)
		}
	} else {
		writeResponse(w, http.StatusBadRequest, errs.NewUnexpectedError("Missing value for customerId param"))
	}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
