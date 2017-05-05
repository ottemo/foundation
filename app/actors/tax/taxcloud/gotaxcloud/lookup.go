package gotaxcloud

import (
	"encoding/json"
	"errors"
)

type LookupParams struct {
	CustomerID        string `json:"customerID"`
	DeliveredBySeller bool   `json:"deliveredBySeller"`

	Origin      Address `json:"origin"`
	Destination Address `json:"destination"`

	CartID    string     `json:"cartID"`
	CartItems []CartItem `json:"cartItems"`
}

type CartItemResponse struct {
	CartItemIndex int
	TaxAmount     float64
}

type LookupResponse struct {
	ResponseBase

	CartID            string `json:"cartID"`
	CartItemsResponse []CartItemResponse
}

func (g *Gateway) Lookup(lookupRequest LookupParams) (*LookupResponse, error) {
	responsePtr, err := g.httpPost("Lookup", lookupRequest)
	if err != nil {
		return nil, err
	}

	var lookupResponse LookupResponse
	err = json.Unmarshal(*responsePtr, &lookupResponse)
	if err != nil {
		return nil, errors.New("2bbc507c-cc90-4fe3-a4cc-a21895a9bf0a: " + err.Error())
	}

	if err = lookupResponse.check(responsePtr); err != nil {
		return nil, err
	}

	return &lookupResponse, nil
}
