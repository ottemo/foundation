package taxcloud

import (
	"encoding/json"
	"errors"
)

type CapturedRequestType struct {
	OrderID string `json:"orderID"`
}

func (g *Gateway) Captured(capturedRequest CapturedRequestType) (*ResponseBase, error) {
	response, err := g.httpPost("Captured", capturedRequest)
	if err != nil {
		return nil, err
	}

	var capturedResponse ResponseBase
	err = json.Unmarshal(response.body, &capturedResponse)
	if err != nil {
		return nil, errors.New("ab8f5f62-4f8f-4340-a41b-8edd5e709f39:" + err.Error())
	}

	if capturedResponse.ResponseType != MessageTypeOK {
		messages, err := json.Marshal(capturedResponse.Messages)
		if err != nil {
			return nil, errors.New("ceb5307d-e672-4978-9882-553a41a8ad44: " + string(response.body))
		}
		return nil, errors.New("d8308775-26ad-447d-8157-4f20eaafb7a8: " + string(messages))
	}

	return &capturedResponse, nil
}
