package qvapay

import (
	"io"
	"net/http"
)

type (
	Options struct {
		//  API's base url
		BaseURL string
		//optional, defaults to http.DefaultClient
		HttpClient *http.Client
		//optional for debuging
		Debug io.Writer
		// App endpoints
		AppID      string
		SecretID   string
		SkipVerify bool
	}
)
