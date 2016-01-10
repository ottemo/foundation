package mailchimp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

//Subscribe a user to a MailChimp list when passed:
//    -- listID string - a MailChimp list id
//    -- registration Registration - a struct to holded needed data to subscribe to a list
func Subscribe(listID string, registration Registration) error {

	if enabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathMailchimpEnabled)); !enabled {
		//If MailChimp is not enabled, ignore
		return nil
	}

	if payload, err := json.Marshal(registration); err == nil {
		if baseURL := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpBaseURL)); baseURL == "" {
			return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3d314122-b50f-11e5-8846-28cfe917b6c7", "Base URL for MailChimp must be defined in the dashboard")
		} else if _, err := sendRequest(fmt.Sprintf(baseURL+"lists/%s/members/", listID), payload); err != nil {
			sendEmail(payload)
			return env.ErrorDispatch(err)
		}

	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}

// checkoutSuccessHandler handles the checkout success event to begin the subscription process if an order meets the
// requirements
func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {

	// grab the order off event map
	var checkoutOrder order.InterfaceOrder
	if eventItem, present := eventData["order"]; present {
		if typedItem, ok := eventItem.(order.InterfaceOrder); ok {
			checkoutOrder = typedItem
		}
	}

	// inspect the order only if not nil
	if checkoutOrder != nil {
		go processOrder(checkoutOrder)
	}

	return true

}

// processOrder is called from the checkout handler to process the order and call Subscribe if the trigger sku is in the
// order
func processOrder(order order.InterfaceOrder) error {

	const givenNameSplit = 1
	const givenName = 0
	const surnameStart = 1

	// TODO: should we support a comma delimited list of trigger skus or maybe json definition which allows
	// multiple lists and multple skus?
	var triggerSKU string
	if triggerSKU = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpSKU)); triggerSKU == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b8c7217c-b509-11e5-aa09-28cfe917b6c7", "Trigger SKU for MailChimp must be defined in the dashboard")
	}

	// inpsect for sku
	if orderContainsSku := containsItem(order, triggerSKU); orderContainsSku {

		var listID string
		if listID = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpList)); listID == "" {
			return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9b5aefcc-b50b-11e5-9689-28cfe917b6c7", "The Mailchimp List ID is not defined in dashboard.")
		}

		var registration Registration
		registration.EmailAddress = utils.InterfaceToString(order.Get("customer_email"))
		registration.Status = ConstMailchimpSubscribeStatus

		// split Order.CustomerName into sub-parts
		splitName := strings.SplitN(utils.InterfaceToString(order.Get("CustomerName")), " ", givenNameSplit)
		registration.MergeFields = map[string]string{
			"FNAME": splitName[givenName],
			"LNAME": splitName[surnameStart],
		}
		// subscribe to specified list
		if err := Subscribe(listID, registration); err != nil {
			return env.ErrorDispatch(err)
		}
	}
	return nil
}

// containsItem will inspect an order for a provided sku
func containsItem(order order.InterfaceOrder, sku string) bool {
	for _, item := range order.GetItems() {
		if item.GetSku() == sku {
			return true
		}
	}
	return false
}

// sendRequest handles the logic for making a json request to the MailChimp server
func sendRequest(url string, payload []byte) (map[string]interface{}, error) {

	var apiKey string
	if apiKey = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpAPIKey)); apiKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "415B50E2-E469-44F4-A179-67C72F3D9631", "MailChimp API key must be defined in the dashboard or your ottemo.ini file")
	}

	buf := bytes.NewBuffer([]byte(payload))
	request, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth("key", apiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	// require http response code of 200 or error out
	if err != nil || response.StatusCode != 200 {

		var status string
		if response == nil {
			status = "nil"
		} else {
			status = response.Status
		}

		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "dc1dc0ce-0918-4eff-a6ce-575985a1bc58", "Unable to subscribe to MailChimp, response code returned was "+status)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return jsonResponse, nil
}

// sendEmail will send an email notification to the specificed user in the dashboard,
// if a subscribe to a MailChimp list fails for any reason
func sendEmail(payload []byte) error {

	// populate the email template
	emailTemplate := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpEmailTemplate))
	if emailTemplate == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6ce922b1-2fe7-4451-9621-1ecd3dc0e45c", "Email template must be defined in the dashboard")
	}

	// configure the email address to send errors when adding visitor email addresses to MailChimp
	supportAddress := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpSupportAddress))
	if supportAddress == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2869ffed-fee9-4e03-9e0f-1be31ffef093", "The MailChimp support email address must be defined in the dashboard")
	}

	// configure the subject of the email for MailChimp errrors
	subjectLine := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpSubjectLine))
	if subjectLine == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9a1cb487-484f-4b0c-b4c4-815d5313ff68", "The MailChimp support email subject must be defined in the dashboard")
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
