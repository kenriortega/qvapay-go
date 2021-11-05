package qvapay

import (
	"fmt"
	"net/url"
)

const (
	ApiVersion = "v1"
	BaseURL    = "https://qvapay.com/api"
)

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
