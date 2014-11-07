package visitor

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// REST API registration function
func setupAPI() error {

	// Dashboard API
	err := api.GetRestService().RegisterAPI("visitor", "POST", "create", restCreateVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "PUT", "update", restUpdateVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "PUT", "update/:id", restUpdateVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "DELETE", "delete/:id", restDeleteVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "load/:id", restGetVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "list", restListVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "POST", "list", restListVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "count", restCountVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "attribute/list", restListVisitorAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "DELETE", "attribute/remove/:attribute", restRemoveVisitorAttribute)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "POST", "attribute/add", restAddVisitorAttribute)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Storefront API
	err = api.GetRestService().RegisterAPI("visitor", "POST", "register", restRegister)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "validate/:key", restValidate)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "forgot-password/:email", restForgotPassword)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "info", restInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "logout", restLogout)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "POST", "login", restLogin)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "POST", "login-facebook", restLoginFacebook)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "POST", "login-google", restLoginGoogle)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "order/list", restListVisitorOrders)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "POST", "order/list", restListVisitorOrders)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "order/details/:id", restVisitorOrderDetails)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor", "POST", "sendmail", restVisitorSendMail)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API used to create new visitor
//   - visitor attributes must be included in POST form
//   - email attribute required
func restCreateVisitor(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(reqData, "email") {
		return nil, env.ErrorNew("'email' was not specified")
	}

	// create operation
	//-----------------
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		err := visitorModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}
	visitorModel.Set("created_at", time.Now())

	err = visitorModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorModel.ToHashMap(), nil
}

// WEB REST API used to update existing visitor
//   - visitor id must be specified in request URI
//   - visitor attributes must be included in POST form
func restUpdateVisitor(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	visitorID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {

		sessionValue := params.Session.Get(visitor.SESSION_KEY_VISITOR_ID)
		sessionVisitorID, ok := sessionValue.(string)
		if !ok {
			return nil, env.ErrorNew("you are not logined in")
		}
		visitorID = sessionVisitorID
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if err := api.ValidateAdminRights(params); err != nil {
		if visitor.GetCurrentVisitorId(params) != visitorID {
			return nil, env.ErrorDispatch(err)
		}

		if _, present := reqData["is_admin"]; present {
			return nil, env.ErrorDispatch(err)
		}
	}

	// update operation
	//-----------------
	visitorModel, err := visitor.LoadVisitorById(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		err := visitorModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = visitorModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorModel.ToHashMap(), nil
}

// WEB REST API used to delete visitor
//   - visitor id must be specified in request URI
func restDeleteVisitor(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew("visitor id was not specified")
	}

	// delete operation
	//-----------------
	visitorModel, err := visitor.GetVisitorModelAndSetId(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.Delete()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// WEB REST API function used to obtain visitor information
//   - visitor id must be specified in request URI
func restGetVisitor(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew("visitor id was not specified")
	}

	// get operation
	//--------------
	visitorModel, err := visitor.LoadVisitorById(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorModel.ToHashMap(), nil
}

// WEB REST API function used to obtain visitors count in model collection
func restCountVisitors(params *api.T_APIHandlerParams) (interface{}, error) {

	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorCollectionModel, err := visitor.GetVisitorCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := visitorCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(params, dbCollection)

	return dbCollection.Count()
}

// WEB REST API function used to get visitors list
func restListVisitors(params *api.T_APIHandlerParams) (interface{}, error) {
	// check request params
	//---------------------
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		if params.Request.Method == "POST" {
			return nil, env.ErrorNew("unexpected request content")
		}

		reqData = make(map[string]interface{})
	}

	// operation start
	//----------------
	visitorCollectionModel, err := visitor.GetVisitorCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// limit parameter handle
	visitorCollectionModel.ListLimit(api.GetListLimit(params))

	// filters handle
	api.ApplyFilters(params, visitorCollectionModel.GetDBCollection())

	// extra parameter handle
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := visitorCollectionModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	return visitorCollectionModel.List()
}

// WEB REST API function used to obtain visitor attributes information
func restListVisitorAttributes(params *api.T_APIHandlerParams) (interface{}, error) {
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attrInfo := visitorModel.GetAttributesInfo()
	return attrInfo, nil
}

// WEB REST API function used to add new custom attribute to visitor model
func restAddVisitorAttribute(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attributeName, isSpecified := reqData["Attribute"]
	if !isSpecified {
		return nil, env.ErrorNew("attribute name was not specified")
	}

	attributeLabel, isSpecified := reqData["Label"]
	if !isSpecified {
		return nil, env.ErrorNew("attribute label was not specified")
	}

	// make product attribute operation
	//---------------------------------
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attribute := models.T_AttributeInfo{
		Model:      visitor.MODEL_NAME_VISITOR,
		Collection: COLLECTION_NAME_VISITOR,
		Attribute:  utils.InterfaceToString(attributeName),
		Type:       "text",
		IsRequired: false,
		IsStatic:   false,
		Label:      utils.InterfaceToString(attributeLabel),
		Group:      "General",
		Editors:    "text",
		Options:    "",
		Default:    "",
		Validators: "",
		IsLayered:  false,
	}

	for key, value := range reqData {
		switch strings.ToLower(key) {
		case "type":
			attribute.Type = utils.InterfaceToString(value)
		case "group":
			attribute.Group = utils.InterfaceToString(value)
		case "editors":
			attribute.Editors = utils.InterfaceToString(value)
		case "options":
			attribute.Options = utils.InterfaceToString(value)
		case "default":
			attribute.Default = utils.InterfaceToString(value)
		case "validators":
			attribute.Validators = utils.InterfaceToString(value)
		case "isrequired", "required":
			attribute.IsRequired = utils.InterfaceToBool(value)
		case "islayered", "layered":
			attribute.IsLayered = utils.InterfaceToBool(value)
		}
	}

	err = visitorModel.AddNewAttribute(attribute)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return attribute, nil
}

// WEB REST API function used to remove custom attribute of visitor model
func restRemoveVisitorAttribute(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//--------------------
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attributeName, isSpecified := params.RequestURLParams["attribute"]
	if !isSpecified {
		return nil, env.ErrorNew("attribute name was not specified")
	}

	// remove attribute actions
	//-------------------------
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.RemoveAttribute(attributeName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// WEB REST API used to register new visitor (same as create but with email validation)
//   - visitor attributes must be included in POST form
//   - email attribute required
func restRegister(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(reqData, "email") {
		return nil, env.ErrorNew("email was not specified")
	}

	// register visitor operation
	//---------------------------
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		err := visitorModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	visitorModel.Set("created_at", time.Now())

	err = visitorModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.Invalidate()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorModel.ToHashMap(), nil
}

// WEB REST API used to validate e-mail address by key sent after registration
func restValidate(params *api.T_APIHandlerParams) (interface{}, error) {
	// check request params
	//---------------------
	validationKey, isKeySpecified := params.RequestURLParams["key"]
	if !isKeySpecified {
		return nil, env.ErrorNew("validation key was not specified")
	}

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.Validate(validationKey)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return api.T_RestRedirect{Location: app.GetStorefrontUrl("login"), DoRedirect: true}, nil
}

// WEB REST API used to sent new password to customer e-mail
func restForgotPassword(params *api.T_APIHandlerParams) (interface{}, error) {

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.LoadByEmail(params.RequestURLParams["email"])
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.GenerateNewPassword()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return api.T_RestRedirect{Result: "ok", Location: app.GetStorefrontUrl("login")}, nil
}

// WEB REST API function used to obtain visitor information
//   - visitor id must be specified in request URI
func restInfo(params *api.T_APIHandlerParams) (interface{}, error) {

	sessionValue := params.Session.Get(visitor.SESSION_KEY_VISITOR_ID)
	visitorID, ok := sessionValue.(string)
	if !ok {
		if api.ValidateAdminRights(params) == nil {
			return map[string]interface{}{"is_admin": true}, nil
		}

		return "you are not logined in", nil
	}

	// visitor info
	//--------------
	visitorModel, err := visitor.LoadVisitorById(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := visitorModel.ToHashMap()
	result["facebook_id"] = visitorModel.GetFacebookId()
	result["google_id"] = visitorModel.GetGoogleId()

	return result, nil
}

// WEB REST API function used to make visitor logout
func restLogout(params *api.T_APIHandlerParams) (interface{}, error) {

	params.Session.Close()

	return "ok", nil
}

// WEB REST API function used to make visitor login
//   - email and password information needed
func restLogin(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(reqData, "email", "password") {
		return nil, env.ErrorNew("email and/or password were not specified")
	}

	requestLogin := utils.InterfaceToString(reqData["email"])
	requestPassword := utils.InterfaceToString(reqData["password"])

	if !strings.Contains(requestLogin, "@") {
		rootLogin := utils.InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_STORE_ROOT_LOGIN))
		rootPassword := utils.InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_STORE_ROOT_PASSWORD))

		if requestLogin == rootLogin && requestPassword == rootPassword {
			params.Session.Set(api.SESSION_KEY_ADMIN_RIGHTS, true)

			return "ok", nil
		}

		return nil, env.ErrorNew("wrong login - should be email")
	}

	// visitor info
	//--------------
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.LoadByEmail(requestLogin)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	ok := visitorModel.CheckPassword(requestPassword)
	if !ok {
		return nil, env.ErrorNew("wrong password")
	}

	// api session updates
	if visitorModel.IsValidated() {
		params.Session.Set(visitor.SESSION_KEY_VISITOR_ID, visitorModel.GetId())
	} else {
		return nil, env.ErrorNew("visitor is not validated")
	}

	if visitorModel.IsAdmin() {
		params.Session.Set(api.SESSION_KEY_ADMIN_RIGHTS, true)
	}

	return "ok", nil
}

// WEB REST API function used to make login/registration via Facebook
//   - access_token and user_id params needed
//   - user needed information will be taken from Facebook
func restLoginFacebook(params *api.T_APIHandlerParams) (interface{}, error) {
	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(reqData, "access_token") {
		return nil, env.ErrorNew("access_token was not specified")
	}

	if !utils.KeysInMapAndNotBlank(reqData, "user_id") {
		return nil, env.ErrorNew("user_id was not specified")
	}

	// facebook login operation
	//-------------------------

	// using access token to get user information
	url := "https://graph.facebook.com/" + reqData["user_id"].(string) + "?access_token=" + reqData["access_token"].(string)
	facebookResponse, err := http.Get(url)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if facebookResponse.StatusCode != 200 {
		return nil, env.ErrorNew("Can't use google API: " + facebookResponse.Status)
	}

	var responseData []byte
	if facebookResponse.ContentLength > 0 {
		responseData = make([]byte, facebookResponse.ContentLength)
		_, err := facebookResponse.Body.Read(responseData)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	} else {
		responseData, err = ioutil.ReadAll(facebookResponse.Body)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	// response json workaround
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(responseData, &jsonMap)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.StrKeysInMap(jsonMap, "id", "email", "first_name", "last_name", "verified") {
		return nil, env.ErrorNew("unexpected facebook response")
	}

	// trying to load visitor from our DB
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorModel, ok := model.(visitor.I_Visitor)
	if !ok {
		return nil, env.ErrorNew("visitor model is not I_Visitor campatible")
	}

	// trying to load visitor by facebook_id
	err = visitorModel.LoadByFacebookId(reqData["user_id"].(string))
	if err != nil && strings.Contains(err.Error(), "not found") {

		// there is no such facebook_id in DB, trying to find by e-mail
		err = visitorModel.LoadByEmail(jsonMap["email"].(string))
		if err != nil && strings.Contains(err.Error(), "not found") {
			// visitor not exists in out DB - reating new one
			visitorModel.Set("email", jsonMap["email"])
			visitorModel.Set("first_name", jsonMap["first_name"])
			visitorModel.Set("last_name", jsonMap["last_name"])
			visitorModel.Set("facebook_id", jsonMap["id"])
			visitorModel.Set("created_at", time.Now())

			err := visitorModel.Save()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		} else {
			// we have visitor with that e-mail just updating token, if it verified
			if value, ok := jsonMap["verified"].(bool); !(ok && value) {
				return nil, env.ErrorNew("facebook account email unverified")
			}

			visitorModel.Set("facebook_id", jsonMap["id"])

			err := visitorModel.Save()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	// api session updates
	params.Session.Set(visitor.SESSION_KEY_VISITOR_ID, visitorModel.GetId())

	if visitorModel.IsAdmin() {
		params.Session.Set(api.SESSION_KEY_ADMIN_RIGHTS, true)
	}

	return "ok", nil
}

// WEB REST API function used to make login/registration via Google
//   - access_token param needed
//   - user needed information will be taken from Google
func restLoginGoogle(params *api.T_APIHandlerParams) (interface{}, error) {
	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorNew("unexpected request content")
	}

	if !utils.KeysInMapAndNotBlank(reqData, "access_token") {
		return nil, env.ErrorNew("access_token was not specified")
	}

	// google login operation
	//-------------------------

	// using access token to get user information
	url := "https://www.googleapis.com/oauth2/v1/userinfo?access_token=" + reqData["access_token"].(string)
	googleResponse, err := http.Get(url)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if googleResponse.StatusCode != 200 {
		return nil, env.ErrorNew("Can't use google API: " + googleResponse.Status)
	}

	var responseData []byte
	if googleResponse.ContentLength > 0 {
		responseData = make([]byte, googleResponse.ContentLength)
		_, err := googleResponse.Body.Read(responseData)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	} else {
		responseData, err = ioutil.ReadAll(googleResponse.Body)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	// response json workaround
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(responseData, &jsonMap)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.StrKeysInMap(jsonMap, "id", "email", "verified_email", "given_name", "family_name") {
		return nil, env.ErrorNew("unexpected google response")
	}

	// trying to load visitor from our DB
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorModel, ok := model.(visitor.I_Visitor)
	if !ok {
		return nil, env.ErrorNew("visitor model is not I_Visitor campatible")
	}

	// trying to load visitor by google_id
	err = visitorModel.LoadByGoogleId(jsonMap["email"].(string))
	if err != nil && strings.Contains(err.Error(), "not found") {

		// there is no such google_id in DB, trying to find by e-mail
		err = visitorModel.LoadByEmail(jsonMap["email"].(string))
		if err != nil && strings.Contains(err.Error(), "not found") {

			// visitor e-mail not exists in out DB - creating new one
			visitorModel.Set("email", jsonMap["email"])
			visitorModel.Set("first_name", jsonMap["given_name"])
			visitorModel.Set("last_name", jsonMap["family_name"])
			visitorModel.Set("google_id", jsonMap["id"])
			visitorModel.Set("created_at", time.Now())

			err := visitorModel.Save()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

		} else {
			// we have visitor with that e-mail just updating token, if it verified
			if value, ok := jsonMap["verified_email"].(bool); !(ok && value) {
				return nil, env.ErrorNew("google account email unverified")
			}

			visitorModel.Set("google_id", jsonMap["id"])

			err := visitorModel.Save()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	// api session updates
	params.Session.Set(visitor.SESSION_KEY_VISITOR_ID, visitorModel.GetId())

	if visitorModel.IsAdmin() {
		params.Session.Set(api.SESSION_KEY_ADMIN_RIGHTS, true)
	}

	return "ok", nil
}

// WEB REST API function used to get visitor order details information
func restVisitorOrderDetails(params *api.T_APIHandlerParams) (interface{}, error) {
	visitorID := visitor.GetCurrentVisitorId(params)
	if visitorID == "" {
		return "you are not logined in", nil
	}

	orderModel, err := order.LoadOrderById(params.RequestURLParams["id"])
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if utils.InterfaceToString(orderModel.Get("visitor_id")) != visitorID {
		return nil, env.ErrorNew("order is not belongs to logined user")
	}

	result := orderModel.ToHashMap()
	result["items"] = orderModel.GetItems()

	return result, nil
}

// WEB REST API function used to get visitor orders information
func restListVisitorOrders(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		reqData = make(map[string]interface{})
	}

	// list operation
	//---------------
	visitorID := visitor.GetCurrentVisitorId(params)
	if visitorID == "" {
		return "you are not logined in", nil
	}

	orderCollection, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = orderCollection.ListFilterAdd("visitor_id", "=", visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// filters handle
	api.ApplyFilters(params, orderCollection.GetDBCollection())

	// extra parameter handle
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := orderCollection.ListAddExtraAttribute(value)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	result, err := orderCollection.List()

	return result, env.ErrorDispatch(err)
}

// WEB REST API function used to get visitor orders information
func restVisitorSendMail(params *api.T_APIHandlerParams) (interface{}, error) {
	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		reqData = make(map[string]interface{})
	}

	if !utils.StrKeysInMap(reqData, "subject", "content", "visitor_ids") {
		return nil, env.ErrorNew("'visitor_ids', 'subject' or 'content' field was not set")
	}

	subject := utils.InterfaceToString(reqData["subject"])
	content := utils.InterfaceToString(reqData["content"])
	visitorIDs := utils.InterfaceToArray(reqData["visitor_ids"])

	for _, visitorID := range visitorIDs {
		visitorModel, err := visitor.LoadVisitorById(utils.InterfaceToString(visitorID))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		err = app.SendMail(visitorModel.GetEmail(), subject, content)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	return "ok", nil
}
