package taxcloud

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
	response, err := g.httpPost("Authorized", authorizedParams)
	if err != nil {
		return nil, err
	}

	var authorizedResponse ResponseBase
	err = json.Unmarshal(response.body, &authorizedResponse)
	if err != nil {
		return nil, errors.New("84883788-fb9a-4ee8-8dc6-491a7a0927fe:" + err.Error())
	}

	if authorizedResponse.ResponseType != MessageTypeOK {
		messages, err := json.Marshal(authorizedResponse.Messages)
		if err != nil {
			return nil, errors.New("60ebfa08-95b4-4862-a90d-88d4ed628d7f: " + string(response.body))
		}
		return nil, errors.New("9c0f7821-0fc1-4379-856f-9e05781649ba: " + string(messages))
	}

	return &authorizedResponse, nil
}
