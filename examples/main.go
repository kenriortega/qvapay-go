package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kenriortega/qvapay-go"
)

func main() {
	api := qvapay.NewAPIClient(
		os.Getenv("APP_ID"),     // app_id
		os.Getenv("APP_SECRET"), // secret_id
		qvapay.BaseURL,          // constants url base https://qvapay.com/api
		false,                   // skip verificationSSL
		nil,                     // custom http.Client
		nil,                     // debug io.Writter (os.Stdout)
	)

	info, err := api.GetInfo(context.Background())
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(info)
}
