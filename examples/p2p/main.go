package main

import (
	"context"
	"fmt"
	"github.com/kenriortega/qvapay-go"
	"log"
	"os"
)

func main() {
	fmt.Println("Get deposit accounts")

	qvaApi, _ := qvapay.QvaPayFactory("qvapay", qvapay.Options{
		BaseURL:    qvapay.BaseURL, // constants url base https://qvapay.com/api
		HttpClient: nil,            // custom http.PaymentAppClient
		Debug:      os.Stdout,      // debug io.Writter (os.Stdout)
	})

	offers, err := qvaApi.Offers(context.Background(), qvapay.QueryParams{Page: 2})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(offers)
}
