package emma

import (
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
		if _, err := subscribeToDefaultGroups(email); err != nil {
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

// Subscribe a user to a Emma by default Emma group IDs
func subscribeToDefaultGroups(email string) (interface{}, error) {

	// get default configured group IDs
	var defaultGroupIds = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaDefaultGroupIds))

	return subscribe(email, defaultGroupIds)
}

// Subscribe a user to a Emma by specifying Emma group IDs
func subscribe(email string, commaSeparatedGroupIDs string) (interface{}, error) {

	// Compose credentials
	// If emma is not enabled, ignore this request and do nothing
	var enabled = utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEmmaEnabled))
	if !enabled {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b3548446-1453-4862-a649-393fc0aafda1", "emma does not active")
	}

	var emmaCredentialsPtr, err = composeEmmaCredentials()
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "53ef2f7b-5fd2-46c1-a012-ce666296c6e2", "internal error: invalid Emma credentials: "+err.Error())
	}

	// Compose subscription params
	var srcGroupIDsList = strings.Split(commaSeparatedGroupIDs, ",")
	var groupIDsList []string
	for _, groupID := range srcGroupIDsList {
		groupID = strings.TrimSpace(groupID)
		if len(groupID) > 0 {
			groupIDsList = append(groupIDsList, groupID)
		}
	}

	var emmaSubscribeInfo = emmaSubscribeInfoType{
		GroupIDsList: groupIDsList,
		Email:        email,
	}

	return emmaService.subscribe(*emmaCredentialsPtr, emmaSubscribeInfo)
}

// composeEmmaCredentials returns emmaCredentialsType filled out with configured values
// This function declared as variable to support future testing substitution
var composeEmmaCredentials = func() (*emmaCredentialsType, error) {
	var accountID = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaAccountID))
	if accountID == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "88111f54-e8a1-4c43-bc38-0e660c4caa16", "account id was not specified")
	}

	var publicAPIKey = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaPublicAPIKey))
	if publicAPIKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1b5c42f5-d856-48c5-98a2-fd8b5929703c", "public api key was not specified")
	}

	var privateAPIKey = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaPrivateAPIKey))
	if privateAPIKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e0282f80-43b4-418e-a99b-60805e74c75d", "private api key was not specified")
	}

	var emmaCredentialsPtr = &emmaCredentialsType{
		AccountID:     accountID,
		PublicAPIKey:  publicAPIKey,
		PrivateAPIKey: privateAPIKey,
	}

	return emmaCredentialsPtr, nil
}
