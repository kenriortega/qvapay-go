package qvapay

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	// "github.com/hashicorp/go-retryablehttp"
)

const (
	ApiVersion = "v1"
	BaseURL    = "https://qvapay.com/api"
)

type IQvaPay interface {
	// GetInfo returns the corresponding object info on fetch call, or an error.
	GetInfo(ctx context.Context) (*AppInfoResponse, error)
	// CreateInvoice ...
	CreateInvoice(
		ctx context.Context,
		amount float64,
		description string,
		remoteID string,
	) (*InvoiceResponse, error)
	// GetTransactions ...
	GetTransactions(ctx context.Context, query APIQueryParams) (*TransactionsResponse, error)
	// GetTransaction ...
	GetTransaction(ctx context.Context, id string) (*TransactionReponse, error)
	// GetBalance ...
	GetBalance(ctx context.Context) (float64, error)

	// qvapay v2

	// Offers ...
	Offers(ctx context.Context, query QueryParams) (map[string]any, error)
}

// QvaPayFactory it`s a constructor factory method
func QvaPayFactory(apiType string, opts Options) (IQvaPay, error) {
	switch apiType {
	case "qvapay":
		return NewQvaPay(opts), nil
	case "app":
		return NewPaymentAppClient(opts), nil
	}
	return nil, errors.New("error: Bad Payment method type passed")
}

type client struct {
	url        string
	httpClient *http.Client
	debug      io.Writer
	appID      string
	appSecret  string
}

type TransPortAuthBasic struct {
	Transport http.RoundTripper
	Token     string
}

func (t TransPortAuthBasic) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "Bearer "+t.Token)
	return t.Transport.RoundTrip(r)
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

// apiCall define how you can make a call to API
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
	req.Header.Add("User-Agent", "qvapay-go")
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
	defer DrainBody(resp.Body)
	if c.debug != nil {
		c.dumpResponse(resp)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", fmt.Errorf("HTTP request failed: %v", err)
	}
	return resp.StatusCode, string(res), nil
}

func DrainBody(respBody io.ReadCloser) {
	// Callers should close resp.Body when done reading from it.
	// If resp.Body is not closed, the Client's underlying RoundTripper
	// (typically Transport) may not be able to re-use a persistent TCP
	// connection to the server for a subsequent "keep-alive" request.
	if respBody != nil {
		// Drain any remaining Body and then close the connection.
		// Without this closing connection would disallow re-using
		// the same connection for future uses.
		//  - http://stackoverflow.com/a/17961593/4465767
		defer func(respBody io.ReadCloser) {
			err := respBody.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(respBody)
		_, _ = io.Copy(ioutil.Discard, respBody)
	}
}

// ParseUrlQueryParams ...
func ParseUrlQueryParams(query QueryParams, requestUrl *url.URL) {
	uv := url.Values{}
	if query.Page != 0 {
		uv.Add("page", strconv.Itoa(query.Page))
	}

	requestUrl.RawQuery = uv.Encode()
}

// APIError is used to describe errors from the API.
// See https://docs.blockfrost.io/#section/Errors
type APIError struct {
	ErrorMessage interface{} `json:"error"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error, %+v", e.ErrorMessage)
}

func HandleAPIErrorResponse(response string) error {
	var err error
	errorApi := &APIError{}
	if err = json.Unmarshal([]byte(response), errorApi); err != nil {
		return err
	}
	return errorApi
}
