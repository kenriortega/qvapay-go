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
	api := qvapay.NewPaymentAppClient(
		os.Getenv("APP_ID"),     // app_id
		os.Getenv("APP_SECRET"), // secret_id
		qvapay.BaseURL,          // constants url base https://qvapay.com/api
		false,                   // skip verificationSSL
		nil,                     // custom http.PaymentAppClient
		os.Stdout,               // debug io.Writter (os.Stdout)
	)

	info, err := api.GetInfo(context.Background())
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(info)
}
