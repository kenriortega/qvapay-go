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
	api := qvapay.NewPaymentAppClient(qvapay.Options{
		BaseURL:    qvapay.BaseURL,          // constants url base https://qvapay.com/api
		HttpClient: nil,                     // custom http.PaymentAppClient
		Debug:      os.Stdout,               // debug io.Writter (os.Stdout)
		AppID:      os.Getenv("APP_ID"),     // app_id
		SecretID:   os.Getenv("APP_SECRET"), // secret_id
		SkipVerify: false,                   // skip verificationSSL
	},
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
