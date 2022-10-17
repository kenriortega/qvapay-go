package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/kenriortega/qvapay-go"
)

func init() {
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

	info, err := api.GetInfo(context.Background())
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(info)
}
