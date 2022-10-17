package qvapay

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

// Transaction  object
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

type APIQueryParams struct {
	Page int
}

func formatParams(v url.Values, query APIQueryParams) url.Values {

	if query.Page > 0 {
		v.Add("page", fmt.Sprintf("%d", query.Page))
	} else {
		v.Add("page", fmt.Sprintf("%d", 1))
	}

	v.Encode()
	return v
}

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

// ToJSON ...
func (ai *AppInfoResponse) ToJSON() string {
	bytes, err := json.Marshal(ai)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
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

func (i *InvoiceResponse) ToJSON() string {
	bytes, err := json.Marshal(i)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
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

func (txs *TransactionsResponse) ToJSON() string {
	bytes, err := json.Marshal(txs)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
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

func (tx *TransactionReponse) ToJSON() string {
	bytes, err := json.Marshal(tx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
}

// GetInfo returns the corresponding object info on fetch call, or an error.
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

// GetTransactions ...
func (c *client) GetTransactions(ctx context.Context, query APIQueryParams) (*TransactionsResponse, error) {
	requestUrl, err := url.Parse(fmt.Sprintf("%s/%s/%s", c.url, ApiVersion, RouteTxs))
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("app_id", c.appID)
	v.Add("app_secret", c.appSecret)
	v = formatParams(v, query)
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
