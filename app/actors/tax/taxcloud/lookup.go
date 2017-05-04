package taxcloud

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
	response, err := g.httpPost("Lookup", lookupRequest)
	if err != nil {
		return nil, err
	}

	var lookupResponse LookupResponse
	err = json.Unmarshal(response.body, &lookupResponse)
	if err != nil {
		return nil, errors.New("2bbc507c-cc90-4fe3-a4cc-a21895a9bf0a: " + err.Error())
	}

	if lookupResponse.ResponseType != MessageTypeOK {
		messages, err := json.Marshal(lookupResponse.Messages)
		if err != nil {
			return nil, errors.New("fb1778f3-391b-4c81-8732-0d5ca0e5f1d1: " + string(response.body))
		}
		return nil, errors.New("aa4cf66c-d605-4a25-9ffe-7f4beca0d325: " + string(messages))
	}

	return &lookupResponse, nil
}
