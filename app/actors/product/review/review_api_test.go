package review_test

import (
	"fmt"
	"io"
	"strings"
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

	// var
	var numberOfUsers = 2
	var numberOfProducts = 2
	var newVisitors []interface{}
	var newProducts []interface{}
	var newReviews []interface{}

	// scenario
	createUsers(t, context, numberOfUsers, &newVisitors)
	createProducts(t, context, numberOfProducts, &newProducts)

	// admin could retrieve all reviews
	// admin could delete reviews
	deleteExistingReviewsByAdmin(t, context)

	// visitor could create review
	createReviewsAllUsersAllProducts(t, context, newProducts, newVisitors, &newReviews)

	// guest couldn't retrieve unapproved reviews
	checkReviewsCount(t, context, "guest unapproved", false, "", "0", "")

	// admin could update review
	// admin could approve review
	approveSomeReviewsByAdmin(t, context, newReviews)

	// guest could retrieve approved reviews for product
	checkReviewsCount(t, context, "guest product approved", false, "", "1", utils.InterfaceToString(utils.InterfaceToMap(newProducts[0])["_id"]))

	// admin could retrieve all reviews
	checkReviewsCount(t, context, "admin all", true, "", "4", "")

	// logged in visitor could get list of his/her reviews
	checkReviewsCount(t, context, "logged personal", false, utils.InterfaceToString(utils.InterfaceToMap(newVisitors[0])["_id"]), "2", "")

	// logged in visitor could get list of approved non empty reviews for product
	checkReviewsCount(
		t,
		context,
		"logged for product",
		false,
		utils.InterfaceToString(utils.InterfaceToMap(newVisitors[1])["_id"]),
		"1",
		utils.InterfaceToString(utils.InterfaceToMap(newProducts[0])["_id"]))

	// logged in visitor could update his/her review
	updateReviewByUser(
		t,
		context,
		utils.InterfaceToString(utils.InterfaceToMap(newVisitors[0])["_id"]),
		utils.InterfaceToString(utils.InterfaceToMap(newReviews[0])["_id"]))
}

func createUsers(t *testing.T, context *testContext, numberOfUsers int, newVisitors *[]interface{}) {
	context.GetSession().Set(api.ConstSessionKeyAdminRights, true)
	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{}

	for i := 0; i < numberOfUsers; i++ {
		context.RequestContent = map[string]interface{}{
			"email": "user" + utils.InterfaceToString(time.Now().Unix()) + utils.InterfaceToString(i) + "@test.com",
		}
		newVisitor, err := visitor.APICreateVisitor(context)
		if err != nil {
			t.Error(err)
		}
		*newVisitors = append(*newVisitors, newVisitor)
	}

}

func createProducts(t *testing.T, context *testContext, numberOfProducts int, newProducts *[]interface{}) {
	context.GetSession().Set(api.ConstSessionKeyAdminRights, true)
	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{}

	for i := 0; i < numberOfProducts; i++ {
		context.RequestContent = map[string]interface{}{
			"sku":  "sku" + utils.InterfaceToString(time.Now().Unix()) + utils.InterfaceToString(i),
			"name": "product name" + utils.InterfaceToString(time.Now().Unix()) + utils.InterfaceToString(i),
		}

		newProduct, err := product.APICreateProduct(context)
		if err != nil {
			t.Error(err)
		}
		*newProducts = append(*newProducts, newProduct)
	}
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

func createReviewsAllUsersAllProducts(
	t *testing.T,
	context *testContext,
	newProducts []interface{},
	newVisitors []interface{},
	newReviews *[]interface{}) {

	context.GetSession().Set(api.ConstSessionKeyAdminRights, false)
	context.ContextValues = map[string]interface{}{}

	var counter = 0
	for _, productItem := range newProducts {
		var productMap = utils.InterfaceToMap(productItem)

		for _, visitorItem := range newVisitors {
			var visitorMap = utils.InterfaceToMap(visitorItem)

			context.GetSession().Set(visitorInterface.ConstSessionKeyVisitorID, visitorMap["_id"])
			context.RequestArguments = map[string]string{
				"productID": utils.InterfaceToString(productMap["_id"]),
			}
			context.RequestContent = map[string]interface{}{
				"review": strings.Repeat("r", counter),
			}
			counter++

			review, err := apiCreateProductReview(context)
			if err != nil {
				t.Error(err)
			}
			reviewMap := utils.InterfaceToMap(review)

			if utils.InterfaceToString(reviewMap["approved"]) != "false" {
				t.Error("New review should not be approved.")
			}

			*newReviews = append(*newReviews, review)
		}
	}
}

func checkReviewsCount(
	t *testing.T,
	context *testContext,
	msg string,
	isAdmin bool,
	visitorID string,
	requiredCount string,
	productID string) {

	context.GetSession().Set(api.ConstSessionKeyAdminRights, isAdmin)
	context.GetSession().Set(visitorInterface.ConstSessionKeyVisitorID, visitorID)

	context.ContextValues = map[string]interface{}{}
	context.RequestContent = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"action": "count",
	}
	if productID != "" {
		context.RequestArguments["product_id"] = productID
	}
	if visitorID != "" {
		context.RequestArguments["visitor_id"] = visitorID
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

func approveSomeReviewsByAdmin(t *testing.T, context *testContext, newReviews []interface{}) {
	context.GetSession().Set(api.ConstSessionKeyAdminRights, true)
	context.ContextValues = map[string]interface{}{}

	var approvedProductID = ""
	var approvedVisitorID = ""
	for _, reviewItem := range newReviews {
		var reviewMap = utils.InterfaceToMap(reviewItem)

		if approvedVisitorID != reviewMap["visitor_id"] && approvedProductID != reviewMap["product_id"] {
			context.RequestArguments = map[string]string{
				"reviewID": utils.InterfaceToString(reviewMap["_id"]),
			}
			context.RequestContent = map[string]interface{}{
				"approved": true,
			}

			updateResult, err := apiUpdateProductReview(context)
			if err != nil {
				t.Error(err)
			}
			_ = updateResult
		}

		if approvedVisitorID == "" && approvedProductID == "" {
			approvedVisitorID = utils.InterfaceToString(reviewMap["visitor_id"])
			approvedProductID = utils.InterfaceToString(reviewMap["product_id"])
		}
	}
}

func updateReviewByUser(t *testing.T, context *testContext, visitorID string, reviewID string) {
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
}
