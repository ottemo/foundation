package taxcloud

import (
	"errors"
)

type MessageType int
const (
	MessageTypeError MessageType = iota
	MessageTypeWarning
	MessageTypeInformational
	MessageTypeOK
)

type Address struct {
	Address1 string `json:"address1"` // This is the numbered street address
	Address2 string `json:"address2"` // This is a second address line (Not Required)
	City string `json:"city"` // This is the city name (Not Required)
	State string `json:"state"` // This is the two character state abbreviation
	Zip5 string `json:"zip5"` // This is the US zip code
	Zip4 string `json:"zip4"` // This is the Plus4 zip code (Required, but can be empty)
}

type CartItem struct {
	Index int
	ItemID string
	TIC int
	Price float64
	Qty int
}

type ResponseBase struct {
	ResponseType MessageType
	Messages []ResponseMessage
}

type ResponseMessage struct {
	ResponseType MessageType
	Message string
}

func (r *ResponseBase) check(body interface{}) error {
	if r.ResponseType == MessageTypeError {
		errorMsg := ""
		for _, msg := range r.Messages {
			errorMsg += msg.Message + "; "
		}
		return errors.New("3ccabeb8-a5b2-4b94-a2d5-f9ff95fee663: ERRORS response: " + errorMsg)
	}

	return nil
}

