package taxcloud

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
	apiLoginID     string // This is your website's API ID
	apiKey         string // This is your website's API KEY

	errorProcessor ErrorProcessorInterface
}

// responseType is internal transport type
type responseType struct {
	statusCode int		// HTTP Response Status Code
	body       []byte	// stores Response Body
}

// NewGateway creates gateway for TaxCloud API methods
func NewGateway(apiLoginID, apiKey string, errorProcessor ErrorProcessorInterface) *Gateway {
	return &Gateway{
		apiLoginID:     apiLoginID,
		apiKey:         apiKey,
		errorProcessor: errorProcessor,
	}
}

// makeHttpPostJsonBody converts request msg to http json body
func (g *Gateway) makeHttpPostJsonBody(jsonData interface{}) (*bytes.Buffer, error) {
	// add API params to body
	var dataMap map[string]interface{}

	marshaled, err := json.Marshal(jsonData)
	if err != nil {
		return nil, errors.New("411b6fbf-3eb4-4752-96ac-0d29806abee0: " + err.Error())
	}
	if err := json.Unmarshal(marshaled, &dataMap); err != nil {
		return nil, errors.New("e0b13c72-5fb9-4577-98c9-aed38ae2026e: " + err.Error())
	}

	dataMap["apiLoginID"] = g.apiLoginID
	dataMap["apiKey"] = g.apiKey

	// create Buffer from data map
	marshaled, err = json.Marshal(dataMap)
	if err != nil {
		return nil, errors.New("41aa4512-6fe0-4bf7-bada-f8afd40aac06:" + err.Error())
	}
	buf := bytes.NewBuffer(marshaled)

	return buf, nil
}

// httpPost executes POST to API
func (g *Gateway) httpPost(operation string, data interface{}) (*responseType, error) {
	body, err := g.makeHttpPostJsonBody(data)
	if err != nil {
		return nil, errors.New("255090c9-86ff-4315-a70f-6b2ca12e0738:" + err.Error())
	}

	request, err := http.NewRequest("POST", apiURL+"/"+operation, body)
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
			g.errorProcessor.Process("8dcb8fe5-563b-4b1b-b3f6-f1d57bbcd2eb", err)
		}
	}(httpResponse.Body)

	response := responseType{
		statusCode: httpResponse.StatusCode,
	}

	response.body, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, errors.New("be05564a-1e66-49b8-a049-2ba2c075bf13:" + err.Error())
	}

	return &response, nil
}
