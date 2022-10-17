// https://documenter.getpostman.com/view/8765260/TzzHnDGw#55629b48-821c-49cf-af14-ce76cfe0d65f
package qvapay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Offers
// curl --location --request GET 'https://qvapay.com/api/p2p/index'
func (c *client) Offers(ctx context.Context, query QueryParams) (map[string]any, error) {
	requestUrl, err := url.Parse(c.url + "/p2p/index")
	if err != nil {
		return nil, err
	}
	ParseUrlQueryParams(query, requestUrl)

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
		return nil, HandleAPIErrorResponse(res)
	}
	result := map[string]any{}

	err = json.NewDecoder(strings.NewReader(res)).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decoding error for data %s: %v", res, err)
	}
	return result, nil
}
