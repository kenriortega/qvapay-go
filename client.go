package qvapay

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const (
	ApiVersion = "v1"
	BaseURL    = "https://qvapay.com/api"
)

// Client is an interface that implements https://qvapay.com/api
type Client interface {
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

// Client
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

func NewAPIClient(
	// APP_ID
	appID string,
	// APP_SECRET
	secretID string,
	// mtss API's base url
	baseURL string,
	// skipVerify
	skipVerify bool,
	//optional, defaults to http.DefaultClient
	httpClient *http.Client,
	debug io.Writer,
) Client {

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
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.httpClient.Transport = tr
	}
	return c
}

// dumpResponse writes the raw response data to the debug output, if set, or
// standard error otherwise.
func (c *client) dumpResponse(resp *http.Response) {
	// ignore errors dumping response - no recovery from this
	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalf("dumpResponse: " + err.Error())
	}
	fmt.Fprintln(c.debug, string(responseDump))
	fmt.Fprintln(c.debug)
}

// apiCall define how you can make a call to Mtss API
func (c *client) apiCall(
	ctx context.Context,
	method string,
	URL string,
	data []byte,
) (statusCode int, response string, err error) {

	req, err := http.NewRequest(method, URL, bytes.NewBuffer(data))
	if err != nil {
		return 0, "", fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Set("User-Agent", "qvapaygo-client/0.0")
	if c.debug != nil {
		requestDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return 0, "", fmt.Errorf("error dumping HTTP request: %v", err)
		}
		fmt.Fprintln(c.debug, string(requestDump))
		fmt.Fprintln(c.debug)
	}
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("HTTP request failed with: %v", err)
	}
	defer resp.Body.Close()
	if c.debug != nil {
		c.dumpResponse(resp)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", fmt.Errorf("HTTP request failed: %v", err)
	}
	return resp.StatusCode, string(res), nil
}

type APIPagingParams struct {
	Count int
	Page  int
	Order string
	From  string
	To    string
}

type APIQueryParams APIPagingParams

func formatParams(v url.Values, query APIPagingParams) url.Values {

	if query.Page > 0 {
		v.Add("page", fmt.Sprintf("%d", query.Page))
	}

	v.Encode()
	return v
}
