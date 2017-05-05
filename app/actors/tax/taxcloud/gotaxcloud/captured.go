package gotaxcloud

import (
	"encoding/json"
	"errors"
)

type CapturedRequestType struct {
	OrderID string `json:"orderID"`
}

func (g *Gateway) Captured(capturedRequest CapturedRequestType) (*ResponseBase, error) {
	responsePtr, err := g.httpPost("Captured", capturedRequest)
	if err != nil {
		return nil, err
	}

	var capturedResponse ResponseBase
	err = json.Unmarshal(*responsePtr, &capturedResponse)
	if err != nil {
		return nil, errors.New("ab8f5f62-4f8f-4340-a41b-8edd5e709f39:" + err.Error())
	}

	if err = capturedResponse.check(responsePtr); err != nil {
		return nil, err
	}

	return &capturedResponse, nil
}
