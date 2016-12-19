package emma

import (
	"bytes"
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

	//If emma is not enabled, ignore this handler and do nothing
	if enabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEmmaEnabled)); !enabled {
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
func processOrder(checkoutOrder order.InterfaceOrder) error {

	var triggerSKU string

	// load the trigger SKUs
	if triggerSKU = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaSKU)); triggerSKU == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ea659e2a-d52d-4d7d-8b94-17283f3c2d3d", "Emma Trigger SKU list may not be empty.")
	}

	// inspect for sku
	if orderHasSKU := containsItem(checkoutOrder, triggerSKU); orderHasSKU {

		email := utils.InterfaceToString(checkoutOrder.Get("customer_email"))

		// subscribe to specified list
		if _, err := subscribe(email); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// containsItem will inspect an order for a sku in the trigger list
func containsItem(checkoutOrder order.InterfaceOrder, triggerList string) bool {

	skuList := strings.Split(triggerList, ",")

	// trim possible whitespace from user entry
	for index, val := range skuList {
		skuList[index] = strings.TrimSpace(val)
	}

	for _, item := range checkoutOrder.GetItems() {
		if inList := utils.IsInListStr(item.GetSku(), skuList); inList {
			return true
		}
	}
	return false
}

// Subscribe a user to a Emma
func subscribe(email string) (interface{}, error) {

	//If emma is not enabled, ignore this request and do nothing
	if enabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEmmaEnabled)); !enabled {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b3548446-1453-4862-a649-393fc0aafda1", "emma does not active")
	}

	var accountId = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaAccountID))
	if accountId == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "88111f54-e8a1-4c43-bc38-0e660c4caa16", "account id was not specified")
	}

	var publicApiKey = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaPublicAPIKey))
	if publicApiKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1b5c42f5-d856-48c5-98a2-fd8b5929703c", "public api key was not specified")
	}

	var privateApiKey = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaPrivateAPIKey))
	if privateApiKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e0282f80-43b4-418e-a99b-60805e74c75d", "private api key was not specified")
	}
	var url = ConstEmmaApiUrl + accountId + "/members/add"

	var defaultGroupIds = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaDefaultGroupIds))
	defaultGroupIdsList := strings.Split(defaultGroupIds, ",")

	postData := map[string]interface{}{"email": email}
	if defaultGroupIds != "" {
		postData["group_ids"] = defaultGroupIdsList
	}

	postDataJson := utils.EncodeToJSONString(postData)

	buf := bytes.NewBuffer([]byte(postDataJson))
	request, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(publicApiKey, privateApiKey)

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

		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cad8ad77-dd4c-440c-ada2-1e315b706175", "Unable to subscribe visitor to Emma list, response code returned was "+status)
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

	var result = "Error occurred"
	if isAdded, isset := jsonResponse["added"]; isset {
		result = "E-mail was added successfully"
		if isAdded == false {
			result = "E-mail already added"
		}
	}

	return result, nil
}
