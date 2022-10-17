package qvapay

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"
)

// PaymentAppClient is an interface that implements https://qvapay.com/api
type PaymentAppClient interface {
	IQvaPay
}

func NewPaymentAppClient(
	opts Options,
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
		appID:      opts.AppID,
		appSecret:  opts.SecretID,
		url:        opts.BaseURL,
		httpClient: opts.HttpClient,
		debug:      opts.Debug,
	}

	if opts.AppID == "" {
		c.appID = os.Getenv("APP_ID")
	}
	if opts.SecretID == "" {
		c.appSecret = os.Getenv("APP_SECRET")
	}
	if opts.BaseURL == "" {
		c.url = BaseURL
	}
	if opts.HttpClient != nil {
		c.httpClient = opts.HttpClient
	} else {
		c.httpClient = http.DefaultClient
	}
	if opts.SkipVerify {
		// #nosec
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		c.httpClient.Transport = tr
	}
	c.httpClient.Transport = tr

	return c
}
