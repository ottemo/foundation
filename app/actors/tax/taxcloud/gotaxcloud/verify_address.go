package gotaxcloud

import (
	"encoding/json"
	"errors"
)

type VerifiedAddress struct {
	ErrNumber      string
	ErrDescription string

	Address
}

func (g *Gateway) VerifyAddress(address Address) (*VerifiedAddress, error) {
	responsePtr, err := g.httpPost("VerifyAddress", address)
	if err != nil {
		return nil, errors.New("fcbbeef3-1388-4e27-b817-cb94beecadb6: " + err.Error())
	}

	var verifiedAddress VerifiedAddress
	err = json.Unmarshal(*responsePtr, &verifiedAddress)
	if err != nil {
		return nil, errors.New("27fab46b-8dd5-46f7-aeba-90154e684d80: " + err.Error())
	}

	if verifiedAddress.ErrNumber != "0" {
		return nil, errors.New("0ba62835-c08a-40ed-b74a-6213c32514b5: " + verifiedAddress.ErrDescription)
	}

	return &verifiedAddress, nil
}
