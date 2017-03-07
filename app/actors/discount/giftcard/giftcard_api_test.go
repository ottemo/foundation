package giftcard_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/actors/visitor"
	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/env/logger"

	visitorInterface "github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/app/actors/discount/giftcard"
	"os"
)

//--------------------------------------------------------------------------------------------------------------
// api.InterfaceSession test implementation
//--------------------------------------------------------------------------------------------------------------

type testSession struct {
	_test_data_ map[string]interface{}
}

func (it *testSession) Close() error {
	return nil
}
func (it *testSession) Get(key string) interface{} {
	return it._test_data_[key]
}
func (it *testSession) GetID() string {
	return "ApplicationSession GetID"
}
func (it *testSession) IsEmpty() bool {
	return true
}
func (it *testSession) Set(key string, value interface{}) {
	it._test_data_[key] = value
}
func (it *testSession) Touch() error {
	return nil
}

//--------------------------------------------------------------------------------------------------------------
// api.InterfaceApplicationContext test implementation
//--------------------------------------------------------------------------------------------------------------

type testContext struct {
	//ResponseWriter    http.ResponseWriter
	//Request           *http.Request
	Request string
	//RequestParameters map[string]string
	RequestArguments map[string]string
	RequestContent   interface{}
	//RequestFiles      map[string]io.Reader

	Session       api.InterfaceSession
	ContextValues map[string]interface{}
	//Result        interface{}
}

func (it *testContext) GetRequestArguments() map[string]string {
	return it.RequestArguments
}
func (it *testContext) GetContextValues() map[string]interface{} {
	return it.ContextValues
}
func (it *testContext) GetContextValue(key string) interface{} {
	return it.ContextValues[key]
}
func (it *testContext) GetRequest() interface{} {
	return it.Request
}
func (it *testContext) GetRequestArgument(name string) string {
	return it.RequestArguments[name]
}
func (it *testContext) GetRequestContent() interface{} {
	return it.RequestContent
}
func (it *testContext) GetRequestContentType() string {
	return "request content type"
}
func (it *testContext) GetRequestFile(name string) io.Reader {
	return nil
}
func (it *testContext) GetRequestFiles() map[string]io.Reader {
	return nil
}
func (it *testContext) GetRequestSettings() map[string]interface{} {
	return map[string]interface{}{}
}
func (it *testContext) GetRequestSetting(name string) interface{} {
	return "request setting"
}
func (it *testContext) GetResponse() interface{} {
	return "response"
}
func (it *testContext) GetResponseContentType() string {
	return "response content type"
}
func (it *testContext) GetResponseResult() interface{} {
	return "response result"
}
func (it *testContext) GetResponseSetting(name string) interface{} {
	return "response setting"
}
func (it *testContext) GetResponseWriter() io.Writer {
	return nil
}
func (it *testContext) GetSession() api.InterfaceSession {
	return it.Session
}
func (it *testContext) SetContextValue(key string, value interface{}) {
	//return it.Session
}
func (it *testContext) SetResponseContentType(mimeType string) error {
	return nil
}
func (it *testContext) SetResponseResult(value interface{}) error {
	return nil
}
func (it *testContext) SetResponseSetting(name string, value interface{}) error {
	return nil
}
func (it *testContext) SetResponseStatus(code int) {
	//return nil
}
func (it *testContext) SetResponseStatusBadRequest()          {}
func (it *testContext) SetResponseStatusForbidden()           {}
func (it *testContext) SetResponseStatusNotFound()            {}
func (it *testContext) SetResponseStatusInternalServerError() {}
func (it *testContext) SetSession(session api.InterfaceSession) error {
	it.Session = session
	return nil
}

//--------------------------------------------------------------------------------------------------------------
// test functions
//--------------------------------------------------------------------------------------------------------------

// redeclare api in case of admin check
// POST
var apiCreateGiftCard = giftcard.Create

// GET
var apiListGiftCards = giftcard.GetList
var apiGetGiftCardByID = giftcard.GetSingleID
var apiGetGiftCardByCode = giftcard.GetSingleCode
var apiGiftCardCodeUnique = giftcard.IfGiftCardCodeUnique

//var apiGetProductRating = review.APIGetProductRating
// PUT
var apiUpdateGiftCard = giftcard.Edit

func TestMain(m *testing.M) {
	err := test.StartAppInTestingMode()
	if err != nil {
		fmt.Println("Unable to start app in testing mode:", err)
	}

	os.Exit(m.Run())
}


func TestGiftCardAPI(t *testing.T) {

	_ = fmt.Sprint("")

	initConfig(t)

	// init session
	session := new(testSession)
	session._test_data_ = map[string]interface{}{}

	// init context
	context := new(testContext)
	if err := context.SetSession(session); err != nil {
		t.Error(err)
	}

	// scenario
	var visitor1 = createVisitor(t, context, "1")
	var visitor2 = createVisitor(t, context, "2")
	_ = visitor1
	_ = visitor2


	//--------------------------------------------------------------------------------------------------------------
	// Create
	//--------------------------------------------------------------------------------------------------------------

	giftcardMap1 := createGiftCard(t, context, visitor1["_id"], 10, "1234567890", time.Now(), "Test 1", "Test Name 1", "test1@test.com", "test1")

	checkGiftCardsCount(t, context, "guest without content", "", "0", false)
	checkGiftCardsCount(t, context, "visitor own without content", visitor1["_id"], "1", true)
	checkGiftCardsCount(t, context, "visitor other without content1", visitor2["_id"], "0", false)
	checkGiftCardsCount(t, context, "admin without content", "admin", "1", true)

	checkGetGiftCard(t, context, "", giftcardMap1["_id"], "guest", false)
	checkGetGiftCard(t, context, visitor1["_id"], giftcardMap1["_id"], "owner1", true)
	checkGetGiftCard(t, context, visitor2["_id"], giftcardMap1["_id"], "other", false)
	checkGetGiftCard(t, context, "admin", giftcardMap1["_id"], "admin", true)

	checkGetGiftCardByCode(t, context, "", giftcardMap1["code"], "guest", true)
	checkGetGiftCardByCode(t, context, visitor1["_id"], giftcardMap1["code"], "owner1", true)
	checkGetGiftCardByCode(t, context, visitor2["_id"], giftcardMap1["code"], "other", true)
	checkGetGiftCardByCode(t, context, "admin", giftcardMap1["code"], "admin", true)

	checkGiftCardCode(t, context, "", giftcardMap1["code"], "guest", false)
	checkGiftCardCode(t, context, visitor1["_id"], giftcardMap1["code"], "owner1", false)
	checkGiftCardCode(t, context, visitor2["_id"], giftcardMap1["code"], "other", false)
	checkGiftCardCode(t, context, "admin", giftcardMap1["code"], "admin", false)

	checkGiftCardCode(t, context, "", "1234567891", "guest", true)
	checkGiftCardCode(t, context, visitor1["_id"], "1234567891", "owner1", true)
	checkGiftCardCode(t, context, visitor2["_id"], "1234567891", "other", true)
	checkGiftCardCode(t, context, "admin", "1234567891", "admin", true)

	//updateByVisitorGiftCard(t, context, visitor1["_id"], giftcardMap1["_id"], false)
	//updateByVisitorGiftCard(t, context, visitor2["_id"], giftcardMap1["_id"], false)
	//updateByVisitorGiftCard(t, context, "admin", giftcardMap1["_id"], true)

	giftcardMap2 := createGiftCard(t, context, visitor2["_id"], 20, "09876654321", time.Now(), "Test 2", "Test Name 2", "test2@test.com", "test2")

	checkGiftCardsCount(t, context, "guest without content", "", "0", false)
	checkGiftCardsCount(t, context, "visitor own without content", visitor1["_id"], "1", true)
	checkGiftCardsCount(t, context, "visitor other without content2", visitor2["_id"], "1", true)
	checkGiftCardsCount(t, context, "admin without content", "admin", "2", true)

	checkGetGiftCard(t, context, "", giftcardMap2["_id"], "guest", false)
	checkGetGiftCard(t, context, visitor1["_id"], giftcardMap2["_id"], "owner2", false)
	checkGetGiftCard(t, context, visitor2["_id"], giftcardMap2["_id"], "other", true)
	checkGetGiftCard(t, context, "admin", giftcardMap2["_id"], "admin", true)

	checkGetGiftCardByCode(t, context, "", giftcardMap2["code"], "guest", true)
	checkGetGiftCardByCode(t, context, visitor1["_id"], giftcardMap2["code"], "owner", true)
	checkGetGiftCardByCode(t, context, visitor2["_id"], giftcardMap2["code"], "other", true)
	checkGetGiftCardByCode(t, context, "admin", giftcardMap2["code"], "admin", true)

	checkGiftCardCode(t, context, "", giftcardMap2["code"], "guest", false)
	checkGiftCardCode(t, context, visitor1["_id"], giftcardMap2["code"], "owner1", false)
	checkGiftCardCode(t, context, visitor2["_id"], giftcardMap2["code"], "other", false)
	checkGiftCardCode(t, context, "admin", giftcardMap2["code"], "admin", false)

	checkGiftCardCode(t, context, "", "19876654321", "guest", true)
	checkGiftCardCode(t, context, visitor1["_id"], "19876654321", "owner1", true)
	checkGiftCardCode(t, context, visitor2["_id"], "19876654321", "other", true)
	checkGiftCardCode(t, context, "admin", "19876654321", "admin", true)

	//updateByVisitorGiftCard(t, context, visitor1["_id"], giftcardMap2["_id"], false)
	//updateByVisitorGiftCard(t, context, visitor2["_id"], giftcardMap2["_id"], false)
	//updateByVisitorGiftCard(t, context, "admin", giftcardMap2["_id"], true)

}

func createVisitor(t *testing.T, context *testContext, counter string) map[string]interface{} {
	context.GetSession().Set(api.ConstSessionKeyAdminRights, true)
	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{}

	context.RequestContent = map[string]interface{}{
		"email": "user" + utils.InterfaceToString(time.Now().Unix()) + counter + "@test.com",
	}
	newVisitor, err := visitor.APICreateVisitor(context)
	if err != nil {
		t.Error(err)
	}

	return utils.InterfaceToMap(newVisitor)
}

func createGiftCard(t *testing.T, context *testContext, visitorID interface{}, amount int, code string, delivery_date time.Time, message string, name string, recipient_mailbox string, sku string) map[string]interface{} {

	context.GetSession().Set(api.ConstSessionKeyAdminRights, false)
	context.ContextValues = map[string]interface{}{}
	context.GetSession().Set(visitorInterface.ConstSessionKeyVisitorID, visitorID)
	context.RequestArguments = map[string]string{}
	context.RequestContent = map[string]interface{}{
		"code": utils.InterfaceToString(code),
		"delivery_date": utils.InterfaceToTime(delivery_date),
		"message": utils.InterfaceToString(message),
		"amount": utils.InterfaceToInt(amount),
		"name": utils.InterfaceToString(name),
		"recipient_mailbox": utils.InterfaceToString(recipient_mailbox),
		"sku": utils.InterfaceToString(sku),
	}

	giftCard, err := apiCreateGiftCard(context)
	if err != nil {
		t.Error(err)
	}
	giftCardMap := utils.InterfaceToMap(giftCard)

	if utils.InterfaceToString(giftCardMap["_id"]) == "" {
		t.Error("New gift card should have id.")
	}

	return (utils.InterfaceToMap(giftCardMap))
}

func updateByVisitorGiftCard(t *testing.T, context *testContext, visitorID interface{}, giftCardID interface{}, canUpdate bool) map[string]interface{} {
	giftCardAmount := 30
	customMessage := "test updated"
	recipientEmail := "newtest@test.com"
	giftCardUniqueCode := utils.InterfaceToString(time.Now().UnixNano())

	var isAdmin = false
	if utils.InterfaceToString(visitorID) == "admin" {
		isAdmin = true
		visitorID = ""
	}

	context.GetSession().Set(api.ConstSessionKeyAdminRights, isAdmin)
	context.GetSession().Set(visitorInterface.ConstSessionKeyVisitorID, visitorID)

	context.ContextValues = map[string]interface{}{}
	context.RequestContent = map[string]interface{}{
		"amount": giftCardAmount,
		"message": customMessage,
		"recipient_mailbox": recipientEmail,
		"code": giftCardUniqueCode,
	}

	context.RequestArguments = map[string]string{
		"id": utils.InterfaceToString(giftCardID),
	}

	updateResult, err := apiUpdateGiftCard(context)
	var updateGiftCardMap = utils.InterfaceToMap(updateResult)
	if err != nil {
		if canUpdate {
			t.Error(err)
			return updateGiftCardMap
		}
	} else if !canUpdate {
		t.Error(visitorID, ", should not be able to get giftcard")
	}


	if utils.InterfaceToInt(updateGiftCardMap["amount"]) != giftCardAmount {
		t.Error("updated gift card amount (" + utils.InterfaceToString(updateGiftCardMap["amount"]) + ") should be " + utils.InterfaceToString(giftCardAmount))
	}
	if utils.InterfaceToString(updateGiftCardMap["recipient_mailbox"]) != recipientEmail {
		t.Error("updated gift card recipient_mailbox (" + utils.InterfaceToString(updateGiftCardMap["recipient_mailbox"]) + ") should be " + recipientEmail)
	}
	if utils.InterfaceToString(updateGiftCardMap["code"]) != giftCardUniqueCode {
		t.Error("updated gift card code (" + utils.InterfaceToString(updateGiftCardMap["code"]) + ") should be " + giftCardUniqueCode)
	}
	if utils.InterfaceToString(updateGiftCardMap["message"]) != customMessage {
		t.Error("updated gift card message (" + utils.InterfaceToString(updateGiftCardMap["message"]) + ") should be " + customMessage)
	}

	return utils.InterfaceToMap(updateResult)
}

func checkGetGiftCard(t *testing.T, context *testContext, visitorID interface{}, giftCardID interface{}, msg string, canGet bool) {

	var isAdmin = false
	if utils.InterfaceToString(visitorID) == "admin" {
		isAdmin = true
		visitorID = ""
	}

	context.GetSession().Set(api.ConstSessionKeyAdminRights, isAdmin)
	context.GetSession().Set(visitorInterface.ConstSessionKeyVisitorID, visitorID)

	context.ContextValues = map[string]interface{}{}
	context.RequestContent = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"id": utils.InterfaceToString(giftCardID),
	}

	result, err := apiGetGiftCardByID(context)
	row := utils.InterfaceToMap(result)

	if err != nil {
		if canGet {
			t.Error(err)
		}
	} else if len(row) != 0 && !canGet {
		t.Error(msg, ", should not be able to get giftcard")
	}
}

func checkGiftCardCode(t *testing.T, context *testContext, visitorID interface{}, giftCardCode interface{}, msg string, isUnique bool) {

	var isAdmin = false
	if utils.InterfaceToString(visitorID) == "admin" {
		isAdmin = true
		visitorID = ""
	}

	context.GetSession().Set(api.ConstSessionKeyAdminRights, isAdmin)
	context.GetSession().Set(visitorInterface.ConstSessionKeyVisitorID, visitorID)

	context.ContextValues = map[string]interface{}{}
	context.RequestContent = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"giftcode": utils.InterfaceToString(giftCardCode),
	}

	result, err := apiGiftCardCodeUnique(context)
	unique := utils.InterfaceToBool(result)

	if err != nil {
		if isUnique {
			t.Error(err)
		}
	} else if unique && !isUnique {
		t.Error(msg, ", giftcard code not unique")
	}
}

func checkGetGiftCardByCode(t *testing.T, context *testContext, visitorID interface{}, giftCardCode interface{}, msg string, canGet bool) {

	var isAdmin = false
	if utils.InterfaceToString(visitorID) == "admin" {
		isAdmin = true
		visitorID = ""
	}

	context.GetSession().Set(api.ConstSessionKeyAdminRights, isAdmin)
	context.GetSession().Set(visitorInterface.ConstSessionKeyVisitorID, visitorID)

	context.ContextValues = map[string]interface{}{}
	context.RequestContent = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"giftcode": utils.InterfaceToString(giftCardCode),
	}

	result, err := apiGetGiftCardByCode(context)
	row := utils.InterfaceToMap(result)

	if err != nil {
		if canGet {
			t.Error(err)
		}
	} else if len(row) != 0 && !canGet {
		t.Error(msg, ", should not be able to get giftcard")
	}
}

func checkGiftCardsCount(
t *testing.T,
context *testContext,
msg string,
visitorID interface{},
requiredCount string,
canGet bool) {

	var isAdmin = false
	if utils.InterfaceToString(visitorID) == "admin" {
		isAdmin = true
		visitorID = ""
	}

	context.GetSession().Set(api.ConstSessionKeyAdminRights, isAdmin)
	context.GetSession().Set(visitorInterface.ConstSessionKeyVisitorID, visitorID)

	context.ContextValues = map[string]interface{}{}
	context.RequestContent = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"action": "count",
	}

	countResult, err := apiListGiftCards(context)
	if countResult == nil && !canGet {
		return
	}
	count := utils.InterfaceToString(countResult)
	if err != nil {
		if canGet {
			t.Error(err)
			return
		}
	} else if count != "0" && !canGet {
		fmt.Println(count)
		t.Error(msg, ", should not be able to get giftcard count")
	}

	if count != requiredCount {
		t.Error("Incorrect giftcards count [" + count + "]. Shoud be " + requiredCount + ". [" + msg + "]")
	}
}

func initConfig(t *testing.T) {
	var config = env.GetConfig()
	if err := config.SetValue(logger.ConstConfigPathErrorLogLevel, 10); err != nil {
		t.Error(err)
	}
}
