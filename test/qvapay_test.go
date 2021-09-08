package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	qvapaygo "github.com/kenriortega/qvapay-go"
	"github.com/stretchr/testify/assert"
)

const (
	appID    = "myAppID"
	secretID = "mySecretID"
)

func Test_Invalid_URL(t *testing.T) {

	client := qvapaygo.NewClient(appID, secretID, "ht&@-tp://:aa", true, nil, nil)

	actual, err := client.GetInfo(context.Background())
	assert.Error(t, err)
	assert.Empty(t, actual)

}

func Test_Get_Info(t *testing.T) {
	uuidExpected := "123456789"
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.Write([]byte(
				` {
					 	"user_id":1,
					 	"name":"my_website",
					 	"url":"https:\/\/www.website.com",
					 	"desc":"WebSite",
					 	"callback":"https:\/\/www.website.com\/webhook",
					 	"logo":"",
					 	"uuid":"123456789",
					 	"secret":"123456987",
					 	"active":1,
					 	"enabled":1
					 }`,
			))
		}),
	)
	defer s.Close()

	client := qvapaygo.NewClient(appID, secretID, s.URL, true, nil, nil)
	info, err := client.GetInfo(context.Background())
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(t, uuidExpected, info.Uuid)
}

func Test_Create_Invoice(t *testing.T) {
	amountInput := 25.60
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.Write([]byte(
				` {
					 	"app_id": "c2ffb4b5-0c73-44f8-b947-53eeddb0afc6",
					 	"amount": "25.60",
					 	"description": "Enanitos verdes",
					 	"remote_id": "BRID56568989",
					 	"signed": "1",
					 	"transation_uuid": "543105f4-b50a-4141-8ede-0ecbbaf5bc87",
					 	"url": "http://qvapay.com/pay/b9330412-2e3d-4fe8-a531-b2be5f68ff4c",
					 	"signedUrl": "http://qvapay.com/pay/b9330412-2e3d-4fe8-a531-b2be5f68ff4c?expires=1610255133&signature=c35db0f1f9e810fd51748aaf69f0981b8d5f83949b7082eeb28c56857b91072b"
					 }`,
			))
		}),
	)
	defer s.Close()
	client := qvapaygo.NewClient(appID, secretID, s.URL, true, nil, nil)
	invoice, err := client.CreateInvoice(context.Background(), amountInput, "Enanitos verdes", "BRID56568989")
	if err != nil {
		t.Fatalf(err.Error())
	}
	result, err := strconv.ParseFloat(invoice.Amount, 64)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, amountInput, result)

}

func Test_Get_Txs(t *testing.T) {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.Write([]byte(
				`   {
					 	"current_page": 1,
					 	"data": [
					 		{
					 			"uuid": "b9330412-2e3d-4fe8-a531-b2be5f68ff4c",
					 			"user_id": 1,
					 			"app_id": 1,
					 			"amount": "25.60",
					 			"description": "Enanitos verdes",
					 			"remote_id": "BRID56568989",
					 			"status": "pending",
					 			"paid_by_user_id": 0,
					 			"created_at": "2021-01-10T04:35:33.000000Z",
					 			"updated_at": "2021-01-10T04:35:33.000000Z",
					 			"signed": 0
					 		}
					 	],
					 	"first_page_url": "http://qvapay.com/api/v1/transactions?page=1",
					 	"from": 1,
					 	"last_page": 1,
					 	"last_page_url": "http://qvapay.com/api/v1/transactions?page=1",
					 	"next_page_url": null,
					 	"path": "http://qvapay.com/api/v1/transactions",
					 	"per_page": 15,
					 	"prev_page_url": null,
					 	"to": 9,
					 	"total": 1
					}`,
			))
		}),
	)
	defer s.Close()
	client := qvapaygo.NewClient(appID, secretID, s.URL, true, nil, nil)
	txs, err := client.GetTransactions(context.Background())
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, 1, len(txs.Data))

}

func Test_Get_Balance(t *testing.T) {

	expected := 66.0
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.Write([]byte(
				`{"66.0"}`,
			))
		}),
	)
	defer s.Close()

	client := qvapaygo.NewClient(appID, secretID, s.URL, true, nil, nil)
	balance, err := client.GetBalance(context.Background())
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(t, expected, balance)

}
