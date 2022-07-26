package qvapay

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"time"
)

// PaymentAppClient is an interface that implements https://qvapay.com/api
type PaymentAppClient interface {
	// GetInfo returns the corresponding object info on fetch call, or an error.
	GetInfo(
		ctx context.Context,
	) (*AppInfoResponse, error)

	// CreateInvoice ...
	CreateInvoice(
		ctx context.Context,
		amount float64,
		description string,
		remoteID string,
	) (*InvoiceResponse, error)

	// GetTrasactions ...
	GetTransactions(
		ctx context.Context,
		query APIQueryParams,
	) (*TransactionsResponse, error)

	// GetTransaction ...
	GetTransaction(
		ctx context.Context,
		id string,
	) (*TransactionReponse, error)

	// GetBalance ...
	GetBalance(
		ctx context.Context,
	) (float64, error)
}

// PaymentAppClient
// client represents a qvapay client. If the Debug field is set to an io.Writer
// (for example os.Stdout), then the client will dump API requests and responses
// to it.  To use a non-default HTTP client (for example, for testing, or to set
// a timeout), assign to the HTTPClient field. To set a non-default URL (for
// example, for testing), assign to the URL field.
type client struct {
	appID      string
	appSecret  string
	url        string
	httpClient *http.Client
	debug      io.Writer
}

func NewPaymentAppClient(
	appID string,
	secretID string,
	baseURL string,
	skipVerify bool,
	httpClient *http.Client,
	debug io.Writer,
) PaymentAppClient {
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          256,
		MaxIdleConnsPerHost:   256,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    true,
	}
	c := &client{
		appID:      appID,
		appSecret:  secretID,
		url:        baseURL,
		httpClient: httpClient,
		debug:      debug,
	}

	if appID == "" {
		c.appID = os.Getenv("APP_ID")
	}
	if secretID == "" {
		c.appSecret = os.Getenv("APP_SECRET")
	}
	if baseURL == "" {
		c.url = BaseURL
	}
	if httpClient != nil {
		c.httpClient = httpClient
	} else {
		c.httpClient = http.DefaultClient
	}
	if skipVerify {
		// #nosec
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		c.httpClient.Transport = tr
	}
	c.httpClient.Transport = tr

	return c
}
