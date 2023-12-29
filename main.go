package main

import (
	"github.com/edwins-leonardi/banking-lib/logger"
	"github.com/edwins-leonardi/banking/app"
)

func main() {
	logger.Info("Starting our application...")
	app.Start()
}
