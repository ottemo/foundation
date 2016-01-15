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

// checkoutSuccessHandler handles the checkout success event to begin the subscription process if an order meets the
// requirements
func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {

	//If MailChimp is not enabled, ignore this handler and do nothing
	if enabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathMailchimpEnabled)); !enabled {
		return true
	}

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

	var triggerSKU, listID string

	if triggerSKU = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpSKU)); triggerSKU == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b8c7217c-b509-11e5-aa09-28cfe917b6c7", "Mailchimp Trigger SKU list may not be empty.")
	}

	if listID = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpList)); listID == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9b5aefcc-b50b-11e5-9689-28cfe917b6c7", "Mailchimp List ID may not be empty.")
	}

	// inspect for sku
	if orderHasSKU := containsItem(order, triggerSKU); orderHasSKU {

		var registration Registration
		registration.EmailAddress = utils.InterfaceToString(order.Get("customer_email"))
		registration.Status = ConstMailchimpSubscribeStatus

		// split Order.CustomerName into sub-parts
		customerName := utils.InterfaceToString(order.Get("customer_name"))
		firstName, lastName := splitName(customerName)
		registration.MergeFields = map[string]string{
			"FNAME": firstName,
			"LNAME": lastName,
		}

		// subscribe to specified list
		if err := Subscribe(listID, registration); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// splitName will take a fullname as a string and split it into first name and last names
func splitName(name string) (string, string) {

	var firstName, lastName string

	fullName := strings.SplitN(name, " ", 2)

	if len(fullName) == 2 {
		firstName = fullName[0]
		lastName = fullName[1]
	} else if len(fullName) == 1 {
		firstName = fullName[0]
		lastName = ""
	} else {
		firstName = ""
		lastName = ""
	}

	return firstName, lastName
}

//Subscribe a user to a MailChimp list when passed:
//    -- listID string - a MailChimp list id
//    -- registration Registration - a struct to holded needed data to subscribe to a list
func Subscribe(listID string, registration Registration) error {

	var payload []byte
	var baseURL string
	var err error

	// load the base url
	if baseURL = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpBaseURL)); baseURL == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3d314122-b50f-11e5-8846-28cfe917b6c7", "MailChimp Base URL may not be empty.")
	}

	// marshal the json payload
	if payload, err = json.Marshal(registration); err != nil {
		return env.ErrorDispatch(err)
	}

	// subscribe to mailchimp
	if _, err = sendRequest(fmt.Sprintf(baseURL+"lists/%s/members", listID), payload); err != nil {
		sendEmail(payload)
		return env.ErrorDispatch(err)
	}

	return nil
}

// containsItem will inspect an order for a sku in the trigger list
func containsItem(order order.InterfaceOrder, triggerList string) bool {

	skuList := strings.Split(triggerList, ",")

	// trim possible whitespace from user entry
	for index, val := range skuList {
		skuList[index] = strings.TrimSpace(val)
	}

	for _, item := range order.GetItems() {
		if inList := utils.IsInListStr(item.GetSku(), skuList); inList {
			return true
		}
	}
	return false
}

// sendRequest handles the logic for making a json request to the MailChimp server
func sendRequest(url string, payload []byte) (map[string]interface{}, error) {

	var apiKey string
	if apiKey = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpAPIKey)); apiKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "415B50E2-E469-44F4-A179-67C72F3D9631", "MailChimp API key may not be empty.")
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
	if response.StatusCode != http.StatusOK {

		var status string
		if response == nil {
			status = "nil"
		} else {
			status = response.Status
		}

		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "dc1dc0ce-0918-4eff-a6ce-575985a1bc58", "Unable to subscribe visitor to MailChimp list, response code returned was "+status)
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
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6ce922b1-2fe7-4451-9621-1ecd3dc0e45c", "Mailchimp Email template may not be empty.")
	}

	// configure the email address to send errors when adding visitor email addresses to MailChimp
	supportAddress := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpSupportAddress))
	if supportAddress == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2869ffed-fee9-4e03-9e0f-1be31ffef093", "Mailchimp Support Email address may not be empty.")
	}

	// configure the subject of the email for MailChimp errrors
	subjectLine := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailchimpSubjectLine))
	if subjectLine == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9a1cb487-484f-4b0c-b4c4-815d5313ff68", "MailChimp Support Email Subject may not be empty.")
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
