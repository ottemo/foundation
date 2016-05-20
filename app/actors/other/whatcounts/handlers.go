package whatcounts

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

	var triggerSKU string

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
func Subscribe(reg Registration) error {

	var baseURL, realm, apiKey, listID, noConfirmation, forceSubscribe string
	var err error

	// load the base url
	baseURL = "https://secure.whatcounts.com"
	// if baseURL = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsBaseURL)); baseURL == "" {
	// 	return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "0221bf42-b4b1-4d44-824f-17a9b4c22d76", "WhatCounts Base URL may not be empty.")
	// }

	// load the realm
	realm = "karigran"
	// if realm = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsRealm)); realm == "" {
	// 	return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cd9e37e9-6322-4534-9cd5-63f806a4f00e", "WhatCounts Base URL may not be empty.")
	// }

	// load the API key
	apiKey = "ompegrab4293"

	// load the whatcounts list id
	listID = "5"
	// if listID = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsList)); listID == "" {
	// 	return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b08e7d67-1b75-4d27-9c46-4234eb47ed90", "Whatcounts List ID may not be empty.")
	// }

	// flag to send a confirmation email
	if utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsNoConfirm)) == "1" {
		noConfirmation = "1"
	} else {
		noConfirmation = "0"
	}
	fmt.Printf("value of noConfirmation: %v\n", noConfirmation)

	// flag to force subscribe to whatcounts list
	if utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsForceSub)) == "1" {
		forceSubscribe = "1"
	} else {
		forceSubscribe = "0"
	}
	fmt.Printf("value of forceSubscribe: %v\n", forceSubscribe)

	// subscribe to whatcounts
	if _, err = sendRequest(
		fmt.Sprintf(baseURL + "/bin/api_web?r=" + realm + "&p=" + apiKey + "&c=sub&list_id=" + listID +
			"&format=99&force_sub=" + forceSubscribe + "&override_confirmation=" + noConfirmation +
			"&data=email,first,last" + reg.EmailAddress + "," + reg.FirstName + "," + reg.LastName)); err != nil {
		// sendEmail(payload)

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
func sendRequest(url string) (string, error) {

	var apiKey string

	// load the api key
	if apiKey = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathWhatcountsAPIKey)); apiKey == "" {
		return "", env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a4bef3b9-e1c2-43ce-b2e2-11f1b08a5682", "MailChimp API key may not be empty.")
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

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

		return "", env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1a600d0d-fad4-4927-ad4d-20acb464d7b1", "Unable to subscribe visitor to WhatCounts list, response code returned was "+status)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return string(responseBody), nil
}
