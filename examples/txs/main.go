package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kenriortega/qvapay-go"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(err.Error())
	}
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu) // Try to use all available CPUs.
}

func main() {
	api := qvapay.NewPaymentAppClient(
		os.Getenv("APP_ID"),     // app_id
		os.Getenv("APP_SECRET"), // secret_id
		qvapay.BaseURL,          // constants url base https://qvapay.com/api
		false,                   // skip verificationSSL
		nil,                     // custom http.PaymentAppClient
		nil,                     // debug io.Writter (os.Stdout)
	)

	tx, err := api.GetTransactions(
		context.Background(),
		qvapay.APIQueryParams{Page: 1},
	)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(tx)
}
