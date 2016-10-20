package review_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/actors/product"
	"github.com/ottemo/foundation/app/actors/product/review"
	"github.com/ottemo/foundation/app/actors/visitor"
	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/utils"

	visitorInterface "github.com/ottemo/foundation/app/models/visitor"
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
var apiCreateProductReview = review.APICreateProductReview

// GET
var apiListReviews = review.APIListReviews
var apiGetReview = review.APIGetReview

//var apiGetProductRating = review.APIGetProductRating
// PUT
var apiUpdateProductReview = review.APIUpdateReview

// DELETE
var apiDeleteProductReview = review.APIDeleteProductReview

func TestReviewAPI(t *testing.T) {

	_ = fmt.Sprint("")

	// start app
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
	}

	// init session
	session := new(testSession)
	session._test_data_ = map[string]interface{}{}

	// init context
	context := new(testContext)
	context.SetSession(session)

	// scenario
	var visitor1 = createVisitor(t, context, "1")
	var visitor2 = createVisitor(t, context, "2")
	_ = visitor1
	_ = visitor2

	var product3 = createProduct(t, context, "3")
	var product4 = createProduct(t, context, "4")
	_ = product3
	_ = product4

	//--------------------------------------------------------------------------------------------------------------
	// Count
	//--------------------------------------------------------------------------------------------------------------

	// admin could retrieve all reviews
	// admin could delete reviews
	deleteExistingReviewsByAdmin(t, context)

	// visitor could create review
	var reviewMap = createReview(t, context, visitor1["_id"], product3["_id"], "")

	// NOT approved
	checkReviewsCount(t, context, "guest without content", "", "", "0")
	checkReviewsCount(t, context, "visitor own without content", visitor1["_id"], "", "1")
	checkReviewsCount(t, context, "visitor other without content", visitor2["_id"], "", "0")
	checkReviewsCount(t, context, "admin without content", "admin", "", "1")

	// admin could update review
	// admin could approve review
	reviewMap = approveReview(t, context, reviewMap["_id"])

	// APPROVED WITHOUT content
	checkReviewsCount(t, context, "guest without content approved", "", "", "0")
	checkReviewsCount(t, context, "visitor own without content approved", visitor1["_id"], "", "1")
	checkReviewsCount(t, context, "visitor other without content approved", visitor2["_id"], "", "0")
	checkReviewsCount(t, context, "admin without content approved", "admin", "", "1")

	// logged in visitor could update his/her review
	reviewMap = updateByVisitorOwnReview(t, context, visitor1["_id"], reviewMap["_id"])
	reviewMap = approveReview(t, context, reviewMap["_id"])

	// APPROVED WITH content
	checkReviewsCount(t, context, "guest content approved", "", "", "1")
	checkReviewsCount(t, context, "visitor own content approved", visitor1["_id"], "", "1")
	checkReviewsCount(t, context, "visitor other content approved", visitor2["_id"], "", "0")
	checkReviewsCount(t, context, "admin content approved", "admin", "", "1")

	// logged in visitor could not update other visitor review
	updateByVisitorOtherReview(t, context, visitor2["_id"], reviewMap["_id"])

	//--------------------------------------------------------------------------------------------------------------
	// Single record
	//--------------------------------------------------------------------------------------------------------------

	deleteExistingReviewsByAdmin(t, context)
	reviewMap = createReview(t, context, visitor1["_id"], product3["_id"], "")

	checkGetReview(t, context, "", reviewMap["_id"], "not aproved, guest", false)
	checkGetReview(t, context, visitor1["_id"], reviewMap["_id"], "not aproved, owner", true)
	checkGetReview(t, context, visitor2["_id"], reviewMap["_id"], "not aproved, other", false)
	checkGetReview(t, context, "admin", reviewMap["_id"], "not aproved, admin", true)

	reviewMap = updateByVisitorOwnReview(t, context, visitor1["_id"], reviewMap["_id"])

	checkGetReview(t, context, "", reviewMap["_id"], "not aproved, guest", false)
	checkGetReview(t, context, visitor1["_id"], reviewMap["_id"], "not aproved, owner", true)
	checkGetReview(t, context, visitor2["_id"], reviewMap["_id"], "not aproved, other", false)
	checkGetReview(t, context, "admin", reviewMap["_id"], "not aproved, admin", true)

	reviewMap = approveReview(t, context, reviewMap["_id"])

	checkGetReview(t, context, "", reviewMap["_id"], "aproved, guest", true)
	checkGetReview(t, context, visitor1["_id"], reviewMap["_id"], "aproved, owner", true)
	checkGetReview(t, context, visitor2["_id"], reviewMap["_id"], "aproved, other", true)
	checkGetReview(t, context, "admin", reviewMap["_id"], "aproved, admin", true)

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

func createProduct(t *testing.T, context *testContext, counter string) map[string]interface{} {
	context.GetSession().Set(api.ConstSessionKeyAdminRights, true)
	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{}

	context.RequestContent = map[string]interface{}{
		"sku":  "sku" + utils.InterfaceToString(time.Now().Unix()) + counter,
		"name": "product name" + utils.InterfaceToString(time.Now().Unix()) + counter,
	}

	newProduct, err := product.APICreateProduct(context)
	if err != nil {
		t.Error(err)
	}

	return utils.InterfaceToMap(newProduct)
}

func deleteExistingReviewsByAdmin(t *testing.T, context *testContext) {
	context.GetSession().Set(api.ConstSessionKeyAdminRights, true)
	context.ContextValues = map[string]interface{}{}
	context.RequestContent = map[string]interface{}{}

	reviews, err := apiListReviews(context)
	if err != nil {
		t.Error(err)
	}

	reviewsMap := utils.InterfaceToArray(reviews)
	for _, review := range reviewsMap {
		reviewMap := utils.InterfaceToMap(review)
		context.RequestArguments = map[string]string{
			"reviewID": utils.InterfaceToString(reviewMap["_id"]),
		}
		if _, err := apiDeleteProductReview(context); err != nil {
			t.Error(err)
		}
	}
}

func createReview(t *testing.T, context *testContext, visitorID interface{}, productID interface{}, reviewValue string) map[string]interface{} {

	context.GetSession().Set(api.ConstSessionKeyAdminRights, false)
	context.ContextValues = map[string]interface{}{}
	context.GetSession().Set(visitorInterface.ConstSessionKeyVisitorID, visitorID)
	context.RequestArguments = map[string]string{
		"productID": utils.InterfaceToString(productID),
	}
	context.RequestContent = map[string]interface{}{
		"review": reviewValue,
	}

	review, err := apiCreateProductReview(context)
	if err != nil {
		t.Error(err)
	}
	reviewMap := utils.InterfaceToMap(review)

	if utils.InterfaceToString(reviewMap["approved"]) != "false" {
		t.Error("New review should not be approved.")
	}

	return (utils.InterfaceToMap(reviewMap))
}

func checkReviewsCount(
	t *testing.T,
	context *testContext,
	msg string,
	visitorID interface{},
	productID interface{},
	requiredCount string) {

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
	if productID != "" {
		context.RequestArguments["product_id"] = utils.InterfaceToString(productID)
	}
	if visitorID != "" {
		context.RequestArguments["visitor_id"] = utils.InterfaceToString(visitorID)
	}

	countResult, err := apiListReviews(context)
	if err != nil {
		t.Error(err)
	}
	count := utils.InterfaceToString(countResult)

	if count != requiredCount {
		t.Error("Incorrect reviews count [" + count + "]. Shoud be " + requiredCount + ". [" + msg + "]")
	}
}

func approveReview(t *testing.T, context *testContext, reviewID interface{}) map[string]interface{} {
	context.GetSession().Set(api.ConstSessionKeyAdminRights, true)
	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"reviewID": utils.InterfaceToString(reviewID),
	}
	context.RequestContent = map[string]interface{}{
		"approved": true,
	}

	updateResult, err := apiUpdateProductReview(context)
	if err != nil {
		t.Error(err)
	}

	return (utils.InterfaceToMap(updateResult))
}

func updateByVisitorOwnReview(t *testing.T, context *testContext, visitorID interface{}, reviewID interface{}) map[string]interface{} {
	var reviewValue = "review text"

	context.GetSession().Set(api.ConstSessionKeyAdminRights, false)
	context.GetSession().Set(visitorInterface.ConstSessionKeyVisitorID, visitorID)

	context.ContextValues = map[string]interface{}{}
	context.RequestContent = map[string]interface{}{
		"review": reviewValue,
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

func checkGetReview(t *testing.T, context *testContext, visitorID interface{}, reviewID interface{}, msg string, canGet bool) {

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
		"reviewID": utils.InterfaceToString(reviewID),
	}

	_, err := apiGetReview(context)
	if err != nil {
		if canGet {
			t.Error(err)
		}
	} else if !canGet {
		t.Error(msg, ", should not be able to get review")
	}
}
