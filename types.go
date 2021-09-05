package qvapaygo

import (
	"context"
	"io"
	"net/http"
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
		amount float32,
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
// client represents a Mtss client. If the Debug field is set to an io.Writer
// (for example os.Stdout), then the client will dump API requests and responses
// to it.  To use a non-default HTTP client (for example, for testing, or to set
// a timeout), assign to the HTTPClient field. To set a non-default URL (for
// example, for testing), assign to the URL field.
type client struct {
	appID      string
	secretID   string
	url        string
	httpClient *http.Client
	debug      io.Writer
}

// Objects Responses

// AppInfoResponse it`s an object that show general datail about your app
type AppInfoResponse struct {
	UserID   string `json:"user_id,omitempty"`
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
	AppID           string  `json:"app_id,omitempty"`
	Amount          float32 `json:"amount,omitempty"`
	Desciption      string  `json:"desciption,omitempty"`
	RemoteID        string  `json:"remote_id,omitempty"`
	Signed          string  `json:"signed,omitempty"`
	TransactionUUID string  `json:"transation_uuid,omitempty"` // report typo miss c (transaction_uuid)
	URL             string  `json:"url,omitempty"`
	SignedUrl       string  `json:"signedUrl,omitempty"`
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
	UserID            string `json:"user_id,omitempty"`
	AppID             string `json:"app_id,omitempty"`
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
	UserID       string `json:"user_id,omitempty"`
	AppID        string `json:"app_id,omitempty"`
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
	UserID   string `json:"user_id,omitempty"`
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

const (
	ApiVersion   = "v1"
	BaseURL      = "https://qvapay.com/api"
	RouteInfo    = "info"
	RouteInvoice = "create_invoice"
	RouteTxs     = "transactions"
	RouteTx      = "transaction"
	RouteBalance = "balance"
)
