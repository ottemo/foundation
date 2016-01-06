package mailchimp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ottemo/foundation/app"
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

	if enabled := utils.InterfaceToBool(env.ConfigGetValue(ConstMailchimpEnabled)); !enabled {
		//If mailchimp is not enabled, ignore
		return nil
	}

	if payload, err := json.Marshal(registration); err == nil {
		if baseURL := utils.InterfaceToString(env.ConfigGetValue(ConstMailchimpBaseURL)); baseURL == "" {
			return env.ErrorDispatch(errors.New("Base URL for MailChimp must be defined in the Dashboard or ottemo.ini file"))
		} else if _, err := sendRequest(fmt.Sprintf(baseURL+"lists/%s/members/", listID), payload); err != nil {
			sendEmail(payload)
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
	if apiKey = utils.InterfaceToString(env.ConfigGetValue(ConstMailchimpAPIKey)); apiKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "415B50E2-E469-44F4-A179-67C72F3D9631", "MailChimp API key must be defined in the Dashboard or your ottemo.ini file")
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

// sendEmail will send an email notification to the specificed user in the dashboard
// -- if an addition to a mailchimp list fails for any reason
func sendEmail(payload []byte) error {

	// populate the email template
	emailTemplate := utils.InterfaceToString(env.ConfigGetValue(ConstMailchimpEmailTemplate))
	if emailTemplate == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6ce922b1-2fe7-4451-9621-1ecd3dc0e45c", "Email Template must be defined in the Dashboard or your ottemo.ini file")
	}

	// configure the email address to send errors when adding visitor email addresses to mailchimp
	supportAddress := utils.InterfaceToString(env.ConfigGetValue(ConstMailchimpSupportAddress))
	if supportAddress == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2869ffed-fee9-4e03-9e0f-1be31ffef093", "The Mailchimp Support Email address must be defined in the Dashboard or your ottemo.ini file")
	}

	// configure the subject of the email for mailchimp errrors
	subjectLine := utils.InterfaceToString(env.ConfigGetValue(ConstMailchimpSubjectLine))
	if subjectLine == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9a1cb487-484f-4b0c-b4c4-815d5313ff68", "The Mailchimp Support Email subject must be defined in the Dashboard or your ottemo.ini file")
	}

	var registration map[string]interface{}
	if err := json.Unmarshal(payload, &registration); err != nil {
		return env.ErrorDispatch(err)
	}
	if email, err := utils.TextTemplate(emailTemplate, registration); err == nil {
		if err := app.SendMail(supportAddress, subjectLine, email); err != nil {
			return env.ErrorDispatch(err)
		}
	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}
