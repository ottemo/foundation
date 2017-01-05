package braintree_test

import (
	"testing"

	"github.com/ottemo/foundation/app/actors/order"
	"github.com/ottemo/foundation/app/actors/payment/braintree"
	_ "github.com/ottemo/foundation/app/actors/visitor" // required to initialize Visitor Address Model
	"github.com/ottemo/foundation/app/actors/visitor/token"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"
)

// TestPaymentMethodCcGuestTransaction tests Authorize method for guest visitor transaction.
func TestPaymentMethodCcGuestTransaction(t *testing.T) {
	// start app
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
	}

	initConfig(t)

	var paymentMethod = &braintree.CreditCardMethod{}
	var orderInstance = &order.DefaultOrder{
		GrandTotal: 100,
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

	if resultMap["transactionID"] != "" {
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

// TestPaymentMethodCcGuestTokenizedTransaction tests Authorize method for guest visitor with token creation and
// creating transaction based on that token
func TestPaymentMethodCcGuestTokenizedTransaction(t *testing.T) {
	// start app
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
	}

	initConfig(t)

	var paymentMethod = &braintree.CreditCardMethod{}

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
	var orderInstance = &order.DefaultOrder{
		GrandTotal: 100,
	}

	var visitorCardInstance = &token.DefaultVisitorCard{}
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

// TestPaymentMethodCcVisitorTokenizedTransaction tests Authorize method for registered visitor with token creation and
// creating transaction based on that token
func TestPaymentMethodCcVisitorTokenizedTransaction(t *testing.T) {
	// start app
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
	}

	initConfig(t)

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
	defer visitorModel.Delete()

	var paymentMethod = &braintree.CreditCardMethod{}

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
	var orderInstance = &order.DefaultOrder{
		GrandTotal: 100,
	}

	var visitorCardInstance = &token.DefaultVisitorCard{}
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

func initConfig(t *testing.T) {
	var config = env.GetConfig()
	if err := config.SetValue(braintree.ConstGeneralConfigPathEnvironment, braintree.ConstEnvironmentSandbox); err != nil {
		t.Error(err)
	}
	if err := config.SetValue(braintree.ConstGeneralConfigPathMerchantID, "vgysg32p79zh9vwr"); err != nil {
		t.Error(err)
	}
	if err := config.SetValue(braintree.ConstGeneralConfigPathPublicKey, "pgzz3pvzy8gwhc7s"); err != nil {
		t.Error(err)
	}
	if err := config.SetValue(braintree.ConstGeneralConfigPathPrivateKey, "2a7363cc16ae440b67e2d5621c70baea"); err != nil {
		t.Error(err)
	}
}
