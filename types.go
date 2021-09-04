package qvapaygo

import (
	"context"
	"io"
	"net/http"
)

// Client is an interface that implements https://qvapay.com/api
type Client interface {
	// GetInfo returns the corresponding object info on fetch call, or an error.
	// URL GET https://qvapay.com/api/v1/info?app_id={app_id}&app_secret={app_secret}
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
	GetInfo(
		ctx context.Context,
	) (AppInfoResponse, error)

	// CreateInvoice ...
	// GET https://qvapay.com/api/v1/create_invoice?app_id={app_id}&app_secret={app_secret}&amount={amount}&description={description}&remote_id={remote_id}&signed={remote_id}
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
	CreateInvoice(
		ctx context.Context,
		amount float32,
		description string,
		remoteID string,
	) (InvoiceResponse, error)

	// GetTrasactions ...
	// GET https://qvapay.com/api/v1/transactions?app_id={app_id}&app_secret={app_secret}
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
	GetTransactions(
		ctx context.Context,
	) error

	// GetTransaction ...
	// GET https://qvapay.com/api/v1/transaction/6507ee0d-db6c-4aa9-b59a-75dc7f6eab52?app_id={app_id}&app_secret={app_secret}
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
	GetTransaction(
		ctx context.Context,
		id string,
	) error

	// GetBalance ...
	// GET https://qvapay.com/api/v1/balance?app_id={app_id}&app_secret={app_secret}
	// {
	// 	"66.00"
	// }
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
