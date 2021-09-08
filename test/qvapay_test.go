package test

import (
	"context"
	"net/http"
	"net/http/httptest"
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
