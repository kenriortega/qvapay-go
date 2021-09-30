package qvapay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	RouteInfo    = "info"
	RouteInvoice = "create_invoice"
	RouteTxs     = "transactions"
	RouteTx      = "transaction"
	RouteBalance = "balance"
)

// AppInfoResponse it`s an object that show general datail about your app
type AppInfoResponse struct {
	UserID   int    `json:"user_id,omitempty"`
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Callback string `json:"callback,omitempty"`
	Logo     string `json:"logo,omitempty"`
	Uuid     string `json:"uuid,omitempty"`
	Active   int    `json:"active,omitempty"`
	Enabled  int    `json:"enabled,omitempty"`
	Secret   string `json:"secret,omitempty"`
}

// InvoiceResponse object
type InvoiceResponse struct {
	AppID           string `json:"app_id,omitempty"`
	Amount          string `json:"amount,omitempty"`
	Desciption      string `json:"desciption,omitempty"`
	RemoteID        string `json:"remote_id,omitempty"`
	Signed          string `json:"signed,omitempty"`
	TransactionUUID string `json:"transation_uuid,omitempty"` // report typo miss c (transaction_uuid)
	URL             string `json:"url,omitempty"`
	SignedUrl       string `json:"signedUrl,omitempty"`
}

// TransactionsResponse reposnses
type TransactionsResponse struct {
	CurrentPage  int           `json:"current_page,omitempty"`
	Data         []Transaction `json:"data,omitempty"`
	FristPageURL string        `json:"frist_page_url,omitempty"`
	From         int           `json:"from,omitempty"`
	LastPage     int           `json:"last_page,omitempty"`
	LastPageURL  string        `json:"last_page_url,omitempty"`
	NextPageURL  string        `json:"next_page_url,omitempty"`
	Path         string        `json:"path,omitempty"`
	PerPage      int           `json:"per_page,omitempty"`
	PrevPageURL  string        `json:"prev_page_url,omitempty"`
	To           int           `json:"to,omitempty"`
	Total        int           `json:"total,omitempty"`
}

// TransactionReponse object
type TransactionReponse struct {
	ID                string `json:"uuid,omitempty"`
	UserID            int    `json:"user_id,omitempty"`
	AppID             int    `json:"app_id,omitempty"`
	Amount            string `json:"amount,omitempty"`
	Description       string `json:"description,omitempty"`
	RemoteID          string `json:"remote_id,omitempty"`
	Status            string `json:"status,omitempty"`
	PaidByUserID      int    `json:"paid_by_user_id,omitempty"`
	Signed            int    `json:"signed,omitempty"`
	CreatedAt         string `json:"created_at,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty"`
	TransactionPaidBy `json:"paid_by,omitempty"`
	App               `json:"app,omitempty"`
	Owner             `json:"owner,omitempty"`
}

// Models

// Trasaction object
type Transaction struct {
	ID           string `json:"uuid,omitempty"`
	UserID       int    `json:"user_id,omitempty"`
	AppID        int    `json:"app_id,omitempty"`
	Amount       string `json:"amount,omitempty"`
	Description  string `json:"description,omitempty"`
	RemoteID     string `json:"remote_id,omitempty"`
	Status       string `json:"status,omitempty"`
	PaidByUserID int    `json:"paid_by_user_id,omitempty"`
	Signed       int    `json:"signed,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

// TransactionPaidBy object
type TransactionPaidBy struct {
	Name string `json:"name,omitempty"`
	Logo string `json:"logo,omitempty"`
}

// App object
type App struct {
	UserID   int    `json:"user_id,omitempty"`
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Callback string `json:"callback,omitempty"`
	Logo     string `json:"logo,omitempty"`
	Uuid     string `json:"uuid,omitempty"`
	Active   int    `json:"active,omitempty"`
	Enabled  int    `json:"enabled,omitempty"`
}

// Owner object
type Owner struct {
	ID       string `json:"uuid,omitempty"`
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
	Lastname string `json:"lastname,omitempty"`
	Logo     string `json:"logo,omitempty"`
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
	v.Add("app_secret", c.appSecret)
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
func (c *client) CreateInvoice(ctx context.Context, amount float64,
	description string,
	remoteID string,
) (*InvoiceResponse, error) {
	requestUrl, err := url.Parse(fmt.Sprintf("%s/%s/%s", c.url, ApiVersion, RouteInvoice))
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("app_id", c.appID)
	v.Add("app_secret", c.appSecret)
	v.Add("amount", fmt.Sprintf("%f", amount))
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
	v.Add("app_secret", c.appSecret)
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
	v.Add("app_secret", c.appSecret)
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
	v.Add("app_secret", c.appSecret)
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

	result, err := strconv.ParseFloat(respParsered, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}
