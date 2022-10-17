package qvapay

import (
	"net/http"
	"os"
	// "github.com/hashicorp/go-retryablehttp"
)

// QueryParams ...
type QueryParams struct {
	Page int
}

// QvaClient is an interface that implements some method for Qvapay
type QvaClient interface {
	IQvaPay
}

// NewQvaPay constructor
func NewQvaPay(
	opts Options,
) QvaClient {

	c := &client{
		url:        opts.BaseURL,
		httpClient: opts.HttpClient,
		debug:      opts.Debug,
	}

	if opts.BaseURL == "" {
		c.url = os.Getenv("QVAPAY_API")
	}
	if opts.HttpClient != nil {
		c.httpClient = opts.HttpClient
	} else {
		c.httpClient = http.DefaultClient
	}

	return c
}
