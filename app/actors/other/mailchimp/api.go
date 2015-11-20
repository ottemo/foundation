package mailchimp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Registration is a struct to hold a single registation for a Mailchimp mailing list.
type Registration struct {
	EmailAddress string            `json:"email_address"`
	Status       string            `json:"status"`
	MergeFields  map[string]string `json:"merge_fields"`
}

//Subscribe a user to a list
func Subscribe(listID string, registration Registration) error {

	if enabled := utils.InterfaceToBool(env.ConfigGetValue(MailchimpEnabledConfig)); !enabled {
		//If mailchimp is not enabled, ignore
		return nil
	}

	if payload, err := json.Marshal(registration); err == nil {
		if baseURL := utils.InterfaceToString(env.ConfigGetValue(MailchimpBaseURLConfig)); baseURL == "" {
			return env.ErrorDispatch(errors.New("Base URL for MailChimp must be defined in your ottemo.ini file"))
		} else if _, err := sendRequest(fmt.Sprintf(baseURL+"lists/%s/members/", listID), payload); err != nil {
			return env.ErrorDispatch(err)
		}

	} else {
		return err
	}

	return nil
}

// handles the logic for making a json request to the mailchimp server
func sendRequest(url string, payload []byte) (map[string]interface{}, error) {

	var apiKey string
	if apiKey = utils.InterfaceToString(env.ConfigGetValue(MailchimpAPIKeyConfig)); apiKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "415B50E2-E469-44F4-A179-67C72F3D9631", "MailChimp API key must be defined in your ottemo.ini file")
	}

	buf := bytes.NewBuffer([]byte(payload))
	request, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "806c965f-92d4-4f78-b47c-ecafa91b23c1", "Unable to create POST request for MailChimp")
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth("key", apiKey)

	var responseBody interface{}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {

		var status string
		if response == nil {
			status = "nil"
		} else {
			status = response.Status
		}

		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "dc1dc0ce-0918-4eff-a6ce-575985a1bc58", "Unable to subscribe to MailChimp, response code returned was "+status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	responseBody = body

	jsonBody, err := utils.DecodeJSONToStringKeyMap(responseBody)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return jsonBody, nil
}
