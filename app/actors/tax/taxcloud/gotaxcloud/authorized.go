package gotaxcloud

import (
	"encoding/json"
	"errors"
	"time"
)

type AuthorizedParams struct {
	CustomerID     string    `json:"customerID"`
	CartID         string    `json:"cartID"`
	OrderID        string    `json:"orderID"`
	DateAuthorized time.Time `json:"dateAuthorized"`
}

func (g *Gateway) Authorized(authorizedParams AuthorizedParams) (*ResponseBase, error) {
	responsePtr, err := g.httpPost("Authorized", authorizedParams)
	if err != nil {
		return nil, err
	}

	var authorizedResponse ResponseBase
	err = json.Unmarshal(*responsePtr, &authorizedResponse)
	if err != nil {
		return nil, errors.New("84883788-fb9a-4ee8-8dc6-491a7a0927fe:" + err.Error())
	}

	if err = authorizedResponse.check(responsePtr); err != nil {
		return nil, err
	}

	return &authorizedResponse, nil
}
