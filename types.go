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
	) (AppInfoResponse, error)

	// CreateInvoice ...
	CreateInvoice(
		ctx context.Context,
		amount float32,
		description string,
		remoteID string,
	) (InvoiceResponse, error)

	// GetTrasactions ...
	GetTransactions(
		ctx context.Context,
	) error

	// GetTransaction ...
	GetTransaction(
		ctx context.Context,
		id string,
	) error

	// GetBalance ...
	GetBalance(
		ctx context.Context,
	) (BalanceQvaPay, error)
}

// Client
// client represents a Mtss client. If the Debug field is set to an io.Writer
// (for example os.Stdout), then the client will dump API requests and responses
// to it.  To use a non-default HTTP client (for example, for testing, or to set
// a timeout), assign to the HTTPClient field. To set a non-default URL (for
// example, for testing), assign to the URL field.
type client struct {
	url        string
	httpClient *http.Client
	debug      io.Writer
}

// AppInfoResponse it`s an object that show general datail about your app
type AppInfoResponse struct {
	UserID   string `json:"user_id,omitempty"`
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Callback string `json:"callback,omitempty"`
	Logo     string `json:"logo,omitempty"`
	Uuid     string `json:"uuid,omitempty"`
	Secret   string `json:"secret,omitempty"`
	Active   int    `json:"active,omitempty"`
	Enabled  int    `json:"enabled,omitempty"`
}

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

type BalanceQvaPay string
