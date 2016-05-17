package whatcounts

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

	//If whatcounts is not enabled, ignore this handler and do nothing
	if enabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathWhatcountsEnabled)); !enabled {
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

	// load the trigger SKUs
	if triggerSKU = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsSKU)); triggerSKU == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "99bcf407-7a24-45b7-b2e0-610decffd7ce", "Whatcounts Trigger SKU list may not be empty.")
	}

	// inspect for sku
	if orderHasSKU := containsItem(order, triggerSKU); orderHasSKU {

		var registration Registration
		registration.EmailAddress = utils.InterfaceToString(order.Get("customer_email"))

		// split Order.CustomerName into sub-parts
		customerName := utils.InterfaceToString(order.Get("customer_name"))
		firstName, lastName := splitName(customerName)
		registration.FirstName = firstName
		registration.LastName = lastName

		// subscribe to specified list
		if err := Subscribe(registration); err != nil {
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
func Subscribe(registration Registration) error {

	var payload []byte
	var baseURL string
	var err error

	// load the base url
	if baseURL = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsBaseURL)); baseURL == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "0221bf42-b4b1-4d44-824f-17a9b4c22d76", "WhatCounts Base URL may not be empty.")
	}

	// load the realm
	if baseURL = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsRealm)); baseURL == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cd9e37e9-6322-4534-9cd5-63f806a4f00e", "WhatCounts Base URL may not be empty.")
	}

	// load the whatcounts list id
	if listID = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsList)); listID == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b08e7d67-1b75-4d27-9c46-4234eb47ed90", "Whatcounts List ID may not be empty.")
	}

	// flag to send a confirmation email
	noConfirmation = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsNoConfirm))

	// flag to force subscribe to whatcounts list
	forceSubscribe = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsForceSub))

	// marshal the json payload
	if payload, err = json.Marshal(registration); err != nil {
		return env.ErrorDispatch(err)
	}

	// subscribe to whatcounts
	if _, err = sendRequest(fmt.Sprintf(baseURL+"/bin/api_web?", listID), payload); err != nil {
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

	// load the api key
	if apiKey = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsAPIKey)); apiKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a4bef3b9-e1c2-43ce-b2e2-11f1b08a5682", "MailChimp API key may not be empty.")
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

		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1a600d0d-fad4-4927-ad4d-20acb464d7b1", "Unable to subscribe visitor to WhatCounts list, response code returned was "+status)
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

	var emailTemplate, supportAddress, subjectLine string

	// populate the email template
	if emailTemplate = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsEmailTemplate)); emailTemplate == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "901a60b4-a18f-46a1-992b-094724aa0407", "Whatcounts Email template may not be empty.")
	}

	// configure the email address to send errors when adding visitor email addresses to MailChimp
	if supportAddress = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsSupportAddress)); supportAddress == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f47cecad-3584-4bdf-a02a-ed73647e9c39", "Whatcounts Support Email address may not be empty.")
	}

	// configure the subject of the email for MailChimp errrors
	if subjectLine := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsSubjectLine)); subjectLine == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "ee67da86-97b1-43ef-aee3-aa5290a0f1e3", "WhatCounts Support Email Subject may not be empty.")
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
