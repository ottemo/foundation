package giftcard_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/actors/product/review"
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

//var apiGetProductRating = review.APIGetProductRating
// PUT
var apiUpdateProductReview = review.APIUpdateReview

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

	giftcardMap := createGiftCard(t, context, visitor1["_id"], 10, "1234567890", time.Now(), "Test 1", "Test Name 1", "test1@test.com", "test1")

	checkGiftCardsCount(t, context, "guest without content", "", "0", false)
	checkGiftCardsCount(t, context, "visitor own without content", visitor1["_id"], "1", true)
	checkGiftCardsCount(t, context, "visitor other without content", visitor2["_id"], "0", false)
	checkGiftCardsCount(t, context, "admin without content", "admin", "1", true)

	checkGetGiftCard(t, context, "", giftcardMap["_id"], "guest", false)
	checkGetGiftCard(t, context, visitor1["_id"], giftcardMap["_id"], "owner", true)
	checkGetGiftCard(t, context, visitor2["_id"], giftcardMap["_id"], "other", false)
	checkGetGiftCard(t, context, "admin", giftcardMap["_id"], "admin", true)

	giftcardMap = createGiftCard(t, context, visitor2["_id"], 20, "09876654321", time.Now(), "Test 2", "Test Name 2", "test2@test.com", "test2")

	checkGiftCardsCount(t, context, "guest without content", "", "0", false)
	checkGiftCardsCount(t, context, "visitor own without content", visitor1["_id"], "1", true)
	checkGiftCardsCount(t, context, "visitor other without content", visitor2["_id"], "1", false)
	checkGiftCardsCount(t, context, "admin without content", "admin", "3", true)

	checkGetGiftCard(t, context, "", giftcardMap["_id"], "not aproved, guest", false)
	checkGetGiftCard(t, context, visitor1["_id"], giftcardMap["_id"], "not aproved, owner", true)
	checkGetGiftCard(t, context, visitor2["_id"], giftcardMap["_id"], "not aproved, other", false)
	checkGetGiftCard(t, context, "admin", giftcardMap["_id"], "not aproved, admin", true)

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

func updateByVisitorOwnReview(t *testing.T, context *testContext, visitorID interface{}, reviewID interface{}) map[string]interface{} {
	var reviewValue = "review text"

	context.GetSession().Set(api.ConstSessionKeyAdminRights, false)
	context.GetSession().Set(visitorInterface.ConstSessionKeyVisitorID, visitorID)

	context.ContextValues = map[string]interface{}{}
	context.RequestContent = map[string]interface{}{
		"review":        reviewValue,
		"unknown_field": "value",
	}

	context.RequestArguments = map[string]string{
		"reviewID": utils.InterfaceToString(reviewID),
	}

	updateResult, err := apiUpdateProductReview(context)
	if err != nil {
		t.Error(err)
	}
	var updateResultMap = utils.InterfaceToMap(updateResult)

	if utils.InterfaceToString(updateResultMap["approved"]) != "false" {
		t.Error("updated by visitor review should not be approved")
	}
	if utils.InterfaceToString(updateResultMap["review"]) != reviewValue {
		t.Error("updated by visitor review is incorrect")
	}

	return utils.InterfaceToMap(updateResult)
}

func updateByVisitorOtherReview(t *testing.T, context *testContext, visitorID interface{}, reviewID interface{}) {
	context.GetSession().Set(api.ConstSessionKeyAdminRights, false)
	context.GetSession().Set(visitorInterface.ConstSessionKeyVisitorID, visitorID)

	context.ContextValues = map[string]interface{}{}
	context.RequestContent = map[string]interface{}{}

	context.RequestArguments = map[string]string{
		"reviewID": utils.InterfaceToString(reviewID),
	}

	_, err := apiUpdateProductReview(context)
	if err == nil {
		t.Error("visitor can not update other visitor review")
	}
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
