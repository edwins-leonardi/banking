package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/edwins-leonardi/banking-lib/logger"
	"github.com/edwins-leonardi/banking/domain"
	"github.com/edwins-leonardi/banking/service"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func sanityCheck() {
	hasError := false
	if os.Getenv("DB_USER") == "" {
		logger.Error("Required env variable DB_USER not defined")
		hasError = true
	}
	if os.Getenv("DB_PASSWORD") == "" {
		logger.Error("Required env variable DB_PASSWORD not defined")
		hasError = true
	}
	if os.Getenv("DB_HOSTNAME") == "" {
		logger.Error("Required env variable DB_HOSTNAME not defined")
		hasError = true
	}
	if os.Getenv("DB_PORT") == "" {
		logger.Error("Required env variable DB_PORT not defined")
		hasError = true
	}
	if os.Getenv("DB_NAME") == "" {
		logger.Error("Required env variable DB_NAME not defined")
		hasError = true
	}
	if hasError {
		log.Fatal("one or more required env vars not defined")
	}
}

func Start() {
	sanityCheck()
	router := mux.NewRouter()

	dbClient := getDbClient()
	ch := CustomerHandlers{
		service: service.NewCustomerService(domain.NewCustomerRepositoryDb(dbClient)),
	}
	ah := AccountHandler{service.NewAccountService(domain.NewAccountRepositoryDb(dbClient))}
	router.
		HandleFunc("/customers", ch.getAllCustomers).
		Methods(http.MethodGet).
		Name("GetAllCustomers")
	router.
		HandleFunc("/customers/{customer_id:[0-9]+}", ch.getCustomer).
		Methods(http.MethodGet).
		Name("GetCustomer")
	router.
		HandleFunc("/customers/{customer_id:[0-9]+}/account", ah.NewAccount).
		Methods(http.MethodPost).
		Name("NewAccount")
	router.
		HandleFunc("/customers/{customer_id:[0-9]+}/account/{account_id:[0-9]+}", ah.MakeTransaction).
		Methods(http.MethodPost).
		Name("NewTransaction")

	am := AuthMiddleware{domain.NewAuthRepository()}
	router.Use(am.authorizationHandler())

	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = "localhost"
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("Server Started at port: " + port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))
}

func getDbClient() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHostname := os.Getenv("DB_HOSTNAME")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHostname, dbPort, dbName)
	client, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client
}
