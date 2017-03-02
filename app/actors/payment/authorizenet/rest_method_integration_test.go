package authorizenet_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ottemo/foundation/app/actors/order"
	"github.com/ottemo/foundation/app/actors/payment/authorizenet"
	"github.com/ottemo/foundation/app/actors/visitor/token"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
)

// TestPaymentMethodRESTGuestTransaction tests Authorize method for guest visitor transaction.
func TestPaymentMethodRESTGuestTransaction(t *testing.T) {
	// start app
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
	}

	db.RegisterOnDatabaseStart(func () error {
		testPaymentMethodRESTGuestTransaction(t)
		return nil
	})
}

// testPaymentMethodRESTGuestTransaction tests Authorize method for guest visitor transaction.
func testPaymentMethodRESTGuestTransaction(t *testing.T) {
	initConfigWithSandboxData(t)

	var paymentMethod = &authorizenet.RestMethod{}

	rand.Seed(time.Now().UnixNano())
	var orderInstance = &order.DefaultOrder{
		GrandTotal: float64(rand.Intn(100)), //100,
	}
	if err := orderInstance.SetID("id" + utils.InterfaceToString(orderInstance.GetGrandTotal())); err != nil {
		t.Error("orderInstance.SetID", err)
		return
	}

	var paymentInfo = map[string]interface{}{
		"cc": map[string]interface{}{
			"number":       "4111111111111111",
			"cvc":          "111",
			"expire_year":  "2025",
			"expire_month": "12",
		},
	}

	result, err := paymentMethod.Authorize(orderInstance, paymentInfo)
	if err != nil {
		t.Error(err)
	}

	var resultMap = utils.InterfaceToMap(result)

	if resultMap["transactionID"] == "" {
		t.Error("Incorrect transactionID [" + utils.InterfaceToString(resultMap["transactionID"]) + "]")
	}
	if resultMap["customerID"] != "" {
		t.Error("Incorrect customerID [" + utils.InterfaceToString(resultMap["customerID"]) + "]")
	}
	if resultMap["creditCardExp"] != "1225" {
		t.Error("Incorrect creditCardExp [" + utils.InterfaceToString(resultMap["creditCardExp"]) + "]")
	}
	if resultMap["creditCardType"] != "Visa" {
		t.Error("Incorrect creditCardType [" + utils.InterfaceToString(resultMap["creditCardType"]) + "]")
	}
	if resultMap["creditCardLastFour"] != "1111" {
		t.Error("Incorrect creditCardLastFour [" + utils.InterfaceToString(resultMap["creditCardLastFour"]) + "]")
	}
}

// testPaymentMethodRestGuestTokenizedTransaction tests Authorize method for guest visitor with token creation and
// creating transaction based on that token
func TestPaymentMethodRestGuestTokenizedTransaction(t *testing.T) {
	// start app
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
	}

	db.RegisterOnDatabaseStart(func () error {
		testPaymentMethodRestGuestTokenizedTransaction(t)
		return nil
	})
}

// testPaymentMethodRestGuestTokenizedTransaction tests Authorize method for guest visitor with token creation and
// creating transaction based on that token
func testPaymentMethodRestGuestTokenizedTransaction(t *testing.T) {
	initConfigWithSandboxData(t)

	var paymentMethod = &authorizenet.RestMethod{}

	// create token
	// take into account - in this test we use another card to avoid duplicate transaction gateway rejection
	var paymentInfo = map[string]interface{}{
		"cc": map[string]interface{}{
			"number":       "6011111111111117",
			"cvc":          "123",
			"expire_year":  "2024",
			"expire_month": "11",
		},
		"extra": map[string]interface{}{
			"billing_name": "First Last",
			"email":        "test@example.com",
		},
		checkout.ConstPaymentActionTypeKey: checkout.ConstPaymentActionTypeCreateToken,
	}

	result, err := paymentMethod.Authorize(nil, paymentInfo)
	if err != nil {
		t.Error(err)
	}

	var resultMap = utils.InterfaceToMap(result)

	if resultMap["transactionID"] == "" {
		t.Error("Incorrect transactionID [" + utils.InterfaceToString(resultMap["transactionID"]) + "]")
	}
	if resultMap["customerID"] == "" {
		t.Error("Incorrect customerID [" + utils.InterfaceToString(resultMap["customerID"]) + "]")
	}
	if resultMap["creditCardExp"] != "1124" {
		t.Error("Incorrect creditCardExp [" + utils.InterfaceToString(resultMap["creditCardExp"]) + "]")
	}
	if resultMap["creditCardType"] != "Discover" {
		t.Error("Incorrect creditCardType [" + utils.InterfaceToString(resultMap["creditCardType"]) + "]")
	}
	if resultMap["creditCardLastFour"] != "1117" {
		t.Error("Incorrect creditCardLastFour [" + utils.InterfaceToString(resultMap["creditCardLastFour"]) + "]")
	}

	// authorize tokenized transaction
	rand.Seed(time.Now().UnixNano())
	var orderInstance = &order.DefaultOrder{
		GrandTotal: float64(rand.Intn(100)), //100,
	}
	if err := orderInstance.SetID("id" + utils.InterfaceToString(orderInstance.GetGrandTotal())); err != nil {
		t.Error("orderInstance.SetID", err)
		return
	}

	var visitorCardInstance = &token.DefaultVisitorCard{
		ExpirationMonth: 11,
		ExpirationYear:  2024,
	}
	if err := visitorCardInstance.Set("token_id", resultMap["transactionID"]); err != nil {
		t.Error(err)
	}
	if err := visitorCardInstance.Set("customer_id", resultMap["customerID"]); err != nil {
		t.Error(err)
	}

	paymentInfo = map[string]interface{}{
		"cc": visitorCardInstance,
	}

	result, err = paymentMethod.Authorize(orderInstance, paymentInfo)
	if err != nil {
		t.Error(err)
	}

	resultMap = utils.InterfaceToMap(result)

	if resultMap["transactionID"] != visitorCardInstance.GetToken() {
		t.Error("Incorrect transactionID [" + utils.InterfaceToString(resultMap["transactionID"]) + "]")
	}
	if resultMap["customerID"] != visitorCardInstance.GetCustomerID() {
		t.Error("Incorrect customerID [" + utils.InterfaceToString(resultMap["customerID"]) + "]")
	}
	if resultMap["creditCardExp"] != "1124" {
		t.Error("Incorrect creditCardExp [" + utils.InterfaceToString(resultMap["creditCardExp"]) + "]")
	}
	if resultMap["creditCardType"] != "Discover" {
		t.Error("Incorrect creditCardType [" + utils.InterfaceToString(resultMap["creditCardType"]) + "]")
	}
	if resultMap["creditCardLastFour"] != "1117" {
		t.Error("Incorrect creditCardLastFour [" + utils.InterfaceToString(resultMap["creditCardLastFour"]) + "]")
	}
}

// TestPaymentMethodRestVisitorTokenizedTransaction tests Authorize method for registered visitor with token creation and
// creating transaction based on that token
func TestPaymentMethodRestVisitorTokenizedTransaction(t *testing.T) {
	// start app
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
	}

	db.RegisterOnDatabaseStart(func () error {
		testPaymentMethodRestVisitorTokenizedTransaction(t)
		return nil
	})
}

// testPaymentMethodRestVisitorTokenizedTransaction tests Authorize method for registered visitor with token creation and
// creating transaction based on that token
func testPaymentMethodRestVisitorTokenizedTransaction(t *testing.T) {
	initConfigWithSandboxData(t)

	var visitorData = map[string]interface{}{
		"id": "fake_visitor_id",
	}

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		t.Error(err)
		return
	}
	err = visitorModel.FromHashMap(visitorData)
	if err != nil {
		t.Error(err)
		return
	}
	err = visitorModel.Save()
	if err != nil {
		t.Error(err)
		return
	}
	defer func (v visitor.InterfaceVisitor){
		if err := v.Delete(); err != nil {
			t.Error(err)
		}
	}(visitorModel)

	var paymentMethod = &authorizenet.RestMethod{}

	// create token
	// take into account - in this test we use another card to avoid duplicate transaction gateway rejection
	var paymentInfo = map[string]interface{}{
		"cc": map[string]interface{}{
			"number":       "3530111333300000",
			"cvc":          "135",
			"expire_year":  "2023",
			"expire_month": "10",
		},
		"extra": map[string]interface{}{
			"visitor_id": visitorModel.GetID(),
		},
		checkout.ConstPaymentActionTypeKey: checkout.ConstPaymentActionTypeCreateToken,
	}

	result, err := paymentMethod.Authorize(nil, paymentInfo)
	if err != nil {
		t.Error(err)
	}

	var resultMap = utils.InterfaceToMap(result)

	if resultMap["transactionID"] == "" {
		t.Error("Incorrect transactionID [" + utils.InterfaceToString(resultMap["transactionID"]) + "]")
	}
	if resultMap["customerID"] == "" {
		t.Error("Incorrect customerID [" + utils.InterfaceToString(resultMap["customerID"]) + "]")
	}
	if resultMap["creditCardExp"] != "1023" {
		t.Error("Incorrect creditCardExp [" + utils.InterfaceToString(resultMap["creditCardExp"]) + "]")
	}
	if resultMap["creditCardType"] != "JCB" {
		t.Error("Incorrect creditCardType [" + utils.InterfaceToString(resultMap["creditCardType"]) + "]")
	}
	if resultMap["creditCardLastFour"] != "0000" {
		t.Error("Incorrect creditCardLastFour [" + utils.InterfaceToString(resultMap["creditCardLastFour"]) + "]")
	}

	// authorize tokenized transaction
	rand.Seed(time.Now().UnixNano())
	var orderInstance = &order.DefaultOrder{
		GrandTotal: float64(rand.Intn(100)), //100,
	}
	if err := orderInstance.SetID("id" + utils.InterfaceToString(orderInstance.GetGrandTotal())); err != nil {
		t.Error("orderInstance.SetID", err)
		return
	}

	var visitorCardInstance = &token.DefaultVisitorCard{
		ExpirationMonth: 10,
		ExpirationYear:  2023,
	}
	if err := visitorCardInstance.Set("token_id", resultMap["transactionID"]); err != nil {
		t.Error(err)
	}
	if err := visitorCardInstance.Set("customer_id", resultMap["customerID"]); err != nil {
		t.Error(err)
	}

	paymentInfo = map[string]interface{}{
		"cc": visitorCardInstance,
	}

	result, err = paymentMethod.Authorize(orderInstance, paymentInfo)
	if err != nil {
		t.Error(err)
	}

	resultMap = utils.InterfaceToMap(result)

	if resultMap["transactionID"] != visitorCardInstance.GetToken() {
		t.Error("Incorrect transactionID [" + utils.InterfaceToString(resultMap["transactionID"]) + "]")
	}
	if resultMap["customerID"] != visitorCardInstance.GetCustomerID() {
		t.Error("Incorrect customerID [" + utils.InterfaceToString(resultMap["customerID"]) + "]")
	}
	if resultMap["creditCardExp"] != "1023" {
		t.Error("Incorrect creditCardExp [" + utils.InterfaceToString(resultMap["creditCardExp"]) + "]")
	}
	if resultMap["creditCardType"] != "JCB" {
		t.Error("Incorrect creditCardType [" + utils.InterfaceToString(resultMap["creditCardType"]) + "]")
	}
	if resultMap["creditCardLastFour"] != "0000" {
		t.Error("Incorrect creditCardLastFour [" + utils.InterfaceToString(resultMap["creditCardLastFour"]) + "]")
	}
}

// TestPaymentMethodRestConfigurationReload checks if package reloads configuration
func TestPaymentMethodRestConfigurationReload(t *testing.T) {
	// start app
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
	}

	db.RegisterOnDatabaseStart(func () error {
		testPaymentMethodRestConfigurationReload(t)
		return nil
	})
}

// testPaymentMethodRestConfigurationReload checks if package reloads configuration
func testPaymentMethodRestConfigurationReload(t *testing.T) {
	initConfig(t, "", "")

	var err error

	var paymentMethod = &authorizenet.RestMethod{}
	var orderInstance = &order.DefaultOrder{
		GrandTotal: 100,
	}

	var paymentInfo = map[string]interface{}{
		"cc": map[string]interface{}{
			"number":       "378282246310005",
			"cvc":          "0005",
			"expire_year":  "2025",
			"expire_month": "05",
		},
	}

	_, err = paymentMethod.Authorize(orderInstance, paymentInfo)
	if err == nil {
		t.Error("libriry should not work with empty environment.")
	}

	initConfig(t, "invalid", "")
	_, err = paymentMethod.Authorize(orderInstance, paymentInfo)
	if err == nil {
		t.Error("libriry should not work with invalid loginID.")
	}

	initConfig(t, "7YAG5ym6P4r4", "")
	_, err = paymentMethod.Authorize(orderInstance, paymentInfo)
	if err == nil {
		t.Error("libriry should not work with invalid transaction key.")
	}

	initConfigWithSandboxData(t)

	_, err = paymentMethod.Authorize(orderInstance, paymentInfo)
	if err != nil {
		t.Error(err)
	}
}

// initConfig initializes configuration with test credentials
func initConfigWithSandboxData(t *testing.T) {
	initConfig(t, "7YAG5ym6P4r4", "72H7Bd638fgswBu8")
}

// initConfig initializes configuration by parameters
func initConfig(t *testing.T, loginID, transactionKey string) {
	if config := env.GetConfig(); config != nil {
		if err := env.GetConfig().SetValue(authorizenet.ConstConfigPathAuthorizeNetRestAPITest, true); err != nil {
			t.Error(err)
			return
		}
		if err := env.GetConfig().SetValue(authorizenet.ConstConfigPathAuthorizeNetRestAPIAPILoginID, loginID); err != nil {
			t.Error(err)
			return
		}
		if err := env.GetConfig().SetValue(authorizenet.ConstConfigPathAuthorizeNetRestAPITransactionKey, transactionKey); err != nil {
			t.Error(err)
			return
		}
	} else {
		t.Error("Test error: unable to get enviromnemt configuration.")
	}

	return
}
