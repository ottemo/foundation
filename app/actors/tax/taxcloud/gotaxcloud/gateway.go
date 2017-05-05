package gotaxcloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

const apiURL = "https://api.taxcloud.com/1.0/TaxCloud"

// ErrorProcessorInterface is interface for processing non-returnable errors
type ErrorProcessorInterface interface {
	Process(uniqueCode string, err error)
}

// gateway is internal type. Stores API credentials and contains API methods
type Gateway struct {
	apiLoginID string // This is your website's API ID
	apiKey     string // This is your website's API KEY

	ErrorProcessor *ErrorProcessorInterface
}

// NewGateway creates gateway for TaxCloud API methods
func NewGateway(apiLoginID, apiKey string) *Gateway {
	return &Gateway{
		apiLoginID:     apiLoginID,
		apiKey:         apiKey,
	}
}

// structToRequestParams1 converts params to map
func (g *Gateway) structToRequestParams(data interface{}) (map[string]interface{}, error) {
	var dataMap map[string]interface{}

	marshaled, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("411b6fbf-3eb4-4752-96ac-0d29806abee0: " + err.Error())
	}
	if err := json.Unmarshal(marshaled, &dataMap); err != nil {
		return nil, errors.New("e0b13c72-5fb9-4577-98c9-aed38ae2026e: " + err.Error())
	}

	if dataMap == nil {
		dataMap = map[string]interface{}{}
	}

	dataMap["apiLoginID"] = g.apiLoginID
	dataMap["apiKey"] = g.apiKey

	return dataMap, nil
}

// httpPost executes POST to API
func (g *Gateway) httpPost(operation string, params interface{}) (*[]byte, error) {
	requestParams, err := g.structToRequestParams(params)
	if err != nil {
		return nil, err
	}

	marshaled, err := json.Marshal(requestParams)
	if err != nil {
		return nil, errors.New("41aa4512-6fe0-4bf7-bada-f8afd40aac06:" + err.Error())
	}

	request, err := http.NewRequest("POST", apiURL+"/"+operation, bytes.NewBuffer(marshaled))
	if err != nil {
		return nil, errors.New("483b7974-e9cd-406f-ae69-9481cea8a078:" + err.Error())
	}

	request.Header.Set("Content-Type", "application/json")

	httpResponse, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.New("34d00613-402e-47e7-8a1a-fa07624ae58c:" + err.Error())
	}

	defer func(c io.ReadCloser) {
		if err := c.Close(); err != nil {
			if g.ErrorProcessor != nil {
				(*g.ErrorProcessor).Process("8dcb8fe5-563b-4b1b-b3f6-f1d57bbcd2eb", err)
			}
		}
	}(httpResponse.Body)

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, errors.New("be05564a-1e66-49b8-a049-2ba2c075bf13:" + err.Error())
	}

	return &responseBody, nil
}
