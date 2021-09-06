package test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	qvapaygo "github.com/kenriortega/qvapay-go"
)

const (
	appID    = "myAppID"
	secretID = "mySecretID"
	urlBase  = "https://qvapay.com/api"
)

var (
	balanceResult = 66.0
)

func validateAnything(*testing.T, []byte) {
}
func validateEmptyBody(t *testing.T, body []byte) {
	if len(body) > 0 {
		t.Errorf("expected empty body, but got %q", body)
	}
}

func responseServer(
	t *testing.T,
	wantMethod string,
	wantURL string,
	validate func(*testing.T, []byte),
	status int,
	filename string,
) *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wantMethod != r.Method {
			t.Errorf("want %q request, got %q", wantMethod, r.Method)
		}
		if r.URL.String() != wantURL {
			t.Errorf("want %q, got %q", wantURL, r.URL.String())
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		validate(t, body)
		w.WriteHeader(status)
		data, err := os.Open(fmt.Sprintf("./%s", filename))
		if err != nil {
			t.Fatal(err)
		}
		defer data.Close()
		io.Copy(w, data)
	}))
}

func Test_Get_Balance(t *testing.T) {
	t.Parallel()
	requestUrl, err := url.Parse(
		fmt.Sprintf("%s/%s/%s", urlBase, qvapaygo.ApiVersion, qvapaygo.RouteBalance),
	)
	if err != nil {
		t.Error(err)
	}
	v := url.Values{}
	v.Set("app_id", appID)
	v.Add("app_secret", secretID)
	requestUrl.RawQuery = v.Encode()
	response := responseServer(
		t,
		http.MethodGet,
		requestUrl.String(),
		validateEmptyBody,
		http.StatusOK,
		"GetBalance.json",
	)
	defer response.Close()
	client := qvapaygo.NewClient(appID, secretID, urlBase, true, response.Client(), nil)
	balance, err := client.GetBalance(context.Background())
	if err != nil {
		t.Fatalf(err.Error())
	}
	if balance != balanceResult {
		t.Fatalf(err.Error())
	}

}
