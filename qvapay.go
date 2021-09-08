package qvapaygo

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func NewClient(
	// APP_ID
	appID string,
	// SECRET_ID
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
		secretID:   secretID,
		url:        baseURL,
		httpClient: httpClient,
		debug:      debug,
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

// GetInfo returns the corresponding object info on fetch call, or an error.
//
// URL GET https://qvapay.com/api/v1/info?app_id={app_id}&app_secret={app_secret}
//
// E.g:
//
// {
// 	"user_id":1,
// 	"name":"my_website",
// 	"url":"https:\/\/www.website.com",
// 	"desc":"WebSite",
// 	"callback":"https:\/\/www.website.com\/webhook",
// 	"logo":"",
// 	"uuid":"123456789",
// 	"secret":"123456987",
// 	"active":1,
// 	"enabled":1
// }
func (c *client) GetInfo(ctx context.Context) (*AppInfoResponse, error) {

	requestUrl, err := url.Parse(fmt.Sprintf("%s/%s/%s", c.url, ApiVersion, RouteInfo))
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("app_id", c.appID)
	v.Add("app_secret", c.secretID)
	requestUrl.RawQuery = v.Encode()

	status, res, err := c.apiCall(
		ctx,
		http.MethodGet,
		requestUrl.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status %d: %q", status, res)
	}
	result := AppInfoResponse{}
	err = json.NewDecoder(strings.NewReader(res)).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decoding error for data %s: %v", res, err)
	}
	return &result, nil
}

// CreateInvoice ...
//
// GET https://qvapay.com/api/v1/create_invoice?app_id={app_id}&app_secret={app_secret}&amount={amount}&description={description}&remote_id={remote_id}&signed={remote_id}
//
// E.g:
//
// {
// 	"app_id": "c2ffb4b5-0c73-44f8-b947-53eeddb0afc6",
// 	"amount": "25.60",
// 	"description": "Enanitos verdes",
// 	"remote_id": "BRID56568989",
// 	"signed": "1",
// 	"transation_uuid": "543105f4-b50a-4141-8ede-0ecbbaf5bc87",
// 	"url": "http://qvapay.com/pay/b9330412-2e3d-4fe8-a531-b2be5f68ff4c",
// 	"signedUrl": "http://qvapay.com/pay/b9330412-2e3d-4fe8-a531-b2be5f68ff4c?expires=1610255133&signature=c35db0f1f9e810fd51748aaf69f0981b8d5f83949b7082eeb28c56857b91072b"
// }
func (c *client) CreateInvoice(ctx context.Context, amount float32,
	description string,
	remoteID string,
) (*InvoiceResponse, error) {
	requestUrl, err := url.Parse(fmt.Sprintf("%s/%s/%s", c.url, ApiVersion, RouteInvoice))
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("app_id", c.appID)
	v.Add("app_secret", c.secretID)
	v.Add("amount", fmt.Sprintf("%d", amount))
	v.Add("description", description)
	v.Add("remote_id", remoteID)
	v.Add("signed", remoteID)
	requestUrl.RawQuery = v.Encode()

	status, res, err := c.apiCall(
		ctx,
		http.MethodGet,
		requestUrl.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status %d: %q", status, res)
	}
	result := InvoiceResponse{}
	err = json.NewDecoder(strings.NewReader(res)).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decoding error for data %s: %v", res, err)
	}
	return &result, nil
}

// GetTrasactions ...
//
// GET https://qvapay.com/api/v1/transactions?app_id={app_id}&app_secret={app_secret}
//
// E.g:
//
// {
// 	"current_page": 1,
// 	"data": [
// 		{
// 			"uuid": "b9330412-2e3d-4fe8-a531-b2be5f68ff4c",
// 			"user_id": 1,
// 			"app_id": 1,
// 			"amount": "25.60",
// 			"description": "Enanitos verdes",
// 			"remote_id": "BRID56568989",
// 			"status": "pending",
// 			"paid_by_user_id": 0,
// 			"created_at": "2021-01-10T04:35:33.000000Z",
// 			"updated_at": "2021-01-10T04:35:33.000000Z",
// 			"signed": 0
// 		},
// 	],
// 	"first_page_url": "http://qvapay.com/api/v1/transactions?page=1",
// 	"from": 1,
// 	"last_page": 1,
// 	"last_page_url": "http://qvapay.com/api/v1/transactions?page=1",
// 	"next_page_url": null,
// 	"path": "http://qvapay.com/api/v1/transactions",
// 	"per_page": 15,
// 	"prev_page_url": null,
// 	"to": 9,
// 	"total": 9
// 	}
func (c *client) GetTransactions(ctx context.Context) (*TransactionsResponse, error) {
	requestUrl, err := url.Parse(fmt.Sprintf("%s/%s/%s", c.url, ApiVersion, RouteTxs))
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("app_id", c.appID)
	v.Add("app_secret", c.secretID)
	requestUrl.RawQuery = v.Encode()

	status, res, err := c.apiCall(
		ctx,
		http.MethodGet,
		requestUrl.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status %d: %q", status, res)
	}
	result := TransactionsResponse{}
	err = json.NewDecoder(strings.NewReader(res)).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decoding error for data %s: %v", res, err)
	}
	return &result, nil
}

// GetTransaction ...
//
// GET https://qvapay.com/api/v1/transaction/6507ee0d-db6c-4aa9-b59a-75dc7f6eab52?app_id={app_id}&app_secret={app_secret}
//
// E.g response:
//
// {
// 	"uuid": "6507ee0d-db6c-4aa9-b59a-75dc7f6eab52",
// 	"user_id": 1,
// 	"app_id": 1,
// 	"amount": "30.00",
// 	"description": "QVAPAY-APP",
// 	"remote_id": "15803",
// 	"status": "pending",
// 	"paid_by_user_id": 0,
// 	"signed": 0,
// 	"created_at": "2021-02-06T18:10:09.000000Z",
// 	"updated_at": "2021-02-06T18:10:09.000000Z",
// 	"paid_by": {
// 		"name": "QvaPay",
// 		"logo": "apps/qvapay.jpg"
// 	},
// 	"app": {
// 		"user_id": 1,
// 		"name": "QvaPay-app",
// 		"url": "https://qvapay.com",
// 		"desc": "Pasarela de pagos con criptomoendas",
// 		"callback": "https://qvapay.com/api/callback",
// 		"success_url": "",
// 		"cancel_url": "",
// 		"logo": "apps/L0YTTe3YdYz9XUh2B78OPdMPNVpt4aVci8FV5y3B.png",
// 		"uuid": "9955dd29-082f-470b-992d-f4f0f25ea164",
// 		"active": 1,
// 		"enabled": 1,
// 		"created_at": "2021-01-12T01:34:21.000000Z",
// 		"updated_at": "2021-01-12T01:34:21.000000Z"
// 	},
// 	"owner": {
// 		"uuid": "796a9e01-3d67-4a42-9dc2-02a5d069fa23",
// 		"username": "qvapay-owner",
// 		"name": "QvaPay",
// 		"lastname": "Pasarela Pagos",
// 		"logo": "profiles/zV93I93mbarZo0fKgwGcpWFWDn41UYfAgj7wNCbf.jpg"
// 	}
// }
func (c *client) GetTransaction(ctx context.Context, id string) (*TransactionReponse, error) {
	requestUrl, err := url.Parse(fmt.Sprintf("%s/%s/%s/%s", c.url, ApiVersion, RouteTx, id))
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("app_id", c.appID)
	v.Add("app_secret", c.secretID)
	requestUrl.RawQuery = v.Encode()

	status, res, err := c.apiCall(
		ctx,
		http.MethodGet,
		requestUrl.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status %d: %q", status, res)
	}
	result := TransactionReponse{}
	err = json.NewDecoder(strings.NewReader(res)).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decoding error for data %s: %v", res, err)
	}
	return &result, nil
}

// GetBalance ...
//
// GET https://qvapay.com/api/v1/balance?app_id={app_id}&app_secret={app_secret}
//
// {
// 	"66.00"
// }
func (c *client) GetBalance(ctx context.Context) (float64, error) {
	requestUrl, err := url.Parse(fmt.Sprintf("%s/%s/%s", c.url, ApiVersion, RouteBalance))
	if err != nil {
		return 0, err
	}
	v := url.Values{}
	v.Set("app_id", c.appID)
	v.Add("app_secret", c.secretID)
	requestUrl.RawQuery = v.Encode()

	status, res, err := c.apiCall(
		ctx,
		http.MethodGet,
		requestUrl.String(),
		nil,
	)
	if err != nil {
		return 0, err
	}
	if status != http.StatusOK {
		return 0, fmt.Errorf("unexpected response status %d: %q", status, res)
	}
	firstParser := strings.ReplaceAll(res, `{"`, "")
	respParsered := strings.ReplaceAll(firstParser, `"}`, "")

	result, err := strconv.ParseFloat(respParsered, 32)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Helpers functions

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
