# qvapay-go
A simple non-official paymentClient for qvapay service with go, for our comunity

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)


## Setup

You can install this package by using the go get tool and installing:

```bash
go get github.com/kenriortega/qvapay-go
```
​
## Sign up on **QvaPay**
Create your account to process payments through **QvaPay** at [qvapay.com/register](https://qvapay.com/register).


## Using the paymentClient

First create your **QvaPay** paymentClient using your app credentials.

```go
paymentClient := qvapaygo.NewClient(
    os.Getenv("APP_ID"), // app_id
    os.Getenv("SECRET_ID"), // secret_id
    qvapaygo.BaseURL, // constants url base https://qvapay.com/api
    false, // skip verificationSSL
    nil, // custom http.PaymentAppClient
    nil, // debug io.Writter (os.Stdout)
)

```
### Get your app info
```go
...
info, err := paymentClient.GetInfo(context.Background())
if err != nil {
    log.Fatalf(err.Error())
}
fmt.Println(info)

```
### Get your account balance
```go
...
balance, err := paymentClient.GetBalance(context.Background())
if err != nil {
    log.Fatalf(err.Error())
}
fmt.Println(balance)

```
### Create an invoice

```go
...
invoice, err := paymentClient.CreateInvoice(
    context.Background(),
    25.60,
    "Enanitos verdes",
    "BRID56568989",
)
if err != nil {
    log.Fatalf(err.Error())
}
fmt.Println(invoice)
```
### Get transaction

```go
...
inputId := "6507ee0d-db6c-4aa9-b59a-75dc7f6eab52"
tx, err := paymentClient.GetTransaction(context.Background(), inputId)
if err != nil {
    log.Fatalf(err.Error())
}
fmt.Println(tx)
```
### Get transactions
```go
...
txs, err := paymentClient.GetTransactions(context.Background())
if err != nil {
    log.Fatalf(err.Error())
}
fmt.Println(txs)
```


You can also read the **QvaPay API** documentation: [qvapay.com/docs](https://qvapay.com/docs).
​