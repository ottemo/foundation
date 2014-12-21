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

// setupAPI setups package related API endpoint routines
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
	//TODO: why support POST only to push an error message?
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
	//TODO: are we missing a PUT to modify an attribute
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
	err = api.GetRestService().RegisterAPI("visitor", "GET", "invalidate/:email", restInvalidate)
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
func restCreateVisitor(params *api.StructAPIHandlerParams) (interface{}, error) {

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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a9610b78add94ae5b75759462b646d2b", "'email' was not specified")
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
func restUpdateVisitor(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	visitorID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {

		sessionValue := params.Session.Get(visitor.ConstSessionKeyVisitorID)
		sessionVisitorID, ok := sessionValue.(string)
		if !ok {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e7a97b45a22a48c596f8f4ecd3f8380f", "Not logged in, please login")
		}
		visitorID = sessionVisitorID
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if err := api.ValidateAdminRights(params); err != nil {
		if visitor.GetCurrentVisitorID(params) != visitorID {
			return nil, env.ErrorDispatch(err)
		}
		if _, present := reqData["is_admin"]; present {
			return nil, env.ErrorDispatch(err)
		}
	}

	// update operation
	//-----------------
	visitorModel, err := visitor.LoadVisitorByID(visitorID)
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
func restDeleteVisitor(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "157df5fad7754934af94b77ef8c826e9", "visitor id was not specified")
	}

	// delete operation
	//-----------------
	visitorModel, err := visitor.GetVisitorModelAndSetID(visitorID)
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
func restGetVisitor(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5863091923f544068676a7d1a629b35f", "visitor id was not specified")
	}

	// get operation
	//--------------
	visitorModel, err := visitor.LoadVisitorByID(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorModel.ToHashMap(), nil
}

// WEB REST API function used to obtain visitors count in model collection
func restCountVisitors(params *api.StructAPIHandlerParams) (interface{}, error) {

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
func restListVisitors(params *api.StructAPIHandlerParams) (interface{}, error) {
	// check request params
	//---------------------
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		if params.Request.Method == "POST" {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "be46cecfa4a14d8ba176a2658ae3b4d7", "Unexpected request for content")
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
func restListVisitorAttributes(params *api.StructAPIHandlerParams) (interface{}, error) {
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attrInfo := visitorModel.GetAttributesInfo()
	return attrInfo, nil
}

// WEB REST API function used to add new custom attribute to visitor model
func restAddVisitorAttribute(params *api.StructAPIHandlerParams) (interface{}, error) {

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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f91e2342d19d4d838c6df7dc4237a49f", "attribute name was not specified")
	}

	attributeLabel, isSpecified := reqData["Label"]
	if !isSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "23c7d427ab414183b446bf5ffef2e8fc", "attribute label was not specified")
	}

	// make product attribute operation
	//---------------------------------
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attribute := models.StructAttributeInfo{
		Model:      visitor.ConstModelNameVisitor,
		Collection: ConstCollectionNameVisitor,
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
func restRemoveVisitorAttribute(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//--------------------
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attributeName, isSpecified := params.RequestURLParams["attribute"]
	if !isSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4b0f0edf692642d8a6257d4d424ee819", "attribute name was not specified")
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
func restRegister(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(reqData, "email") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a37ff1e868e34201a7b4c9b4356dcbeb", "email was not specified")
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
func restValidate(params *api.StructAPIHandlerParams) (interface{}, error) {

	validationKey, isKeySpecified := params.RequestURLParams["key"]
	if !isKeySpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "0b1c841418bf4c448b29c462d38d5498", "validation key was not specified")
	}

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.Validate(validationKey)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return api.StructRestRedirect{Result: "ok", Location: app.GetStorefrontURL("login")}, nil
}

// WEB REST API used to invalidate customer e-mail
func restInvalidate(params *api.StructAPIHandlerParams) (interface{}, error) {

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorEmail := params.RequestURLParams["email"]

	err = visitorModel.LoadByEmail(visitorEmail)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// checking rights
	if err := api.ValidateAdminRights(params); err != nil {
		if visitorModel.IsValidated() {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = visitorModel.Invalidate()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// WEB REST API used to sent new password to customer e-mail
func restForgotPassword(params *api.StructAPIHandlerParams) (interface{}, error) {

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorEmail := params.RequestURLParams["email"]

	err = visitorModel.LoadByEmail(visitorEmail)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.GenerateNewPassword()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// WEB REST API function used to obtain visitor information
//   - visitor id must be specified in request URI
func restInfo(params *api.StructAPIHandlerParams) (interface{}, error) {

	// so if user was logged to app as admin, we want to reflect this
	isAdmin := false
	if api.ValidateAdminRights(params) == nil {
		isAdmin = true
	}

	sessionValue := params.Session.Get(visitor.ConstSessionKeyVisitorID)
	visitorID, ok := sessionValue.(string)
	if !ok {
		if isAdmin {
			return map[string]interface{}{"is_admin": true}, nil
		}
		return "you are not logined in", nil
	}

	// visitor info
	//--------------
	visitorModel, err := visitor.LoadVisitorByID(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := visitorModel.ToHashMap()
	result["facebook_id"] = visitorModel.GetFacebookID()
	result["google_id"] = visitorModel.GetGoogleID()

	// overriding DB value if currently logged as admin
	if isAdmin {
		result["is_admin"] = true
	}

	return result, nil
}

// WEB REST API function used to make visitor logout
func restLogout(params *api.StructAPIHandlerParams) (interface{}, error) {

	params.Session.Close()

	return "ok", nil
}

// WEB REST API function used to make visitor login
//   - email and password information needed
func restLogin(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(reqData, "email", "password") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "02b0583b28c34072afe2f392a163ed87", "email and/or password were not specified")
	}

	requestLogin := utils.InterfaceToString(reqData["email"])
	requestPassword := utils.InterfaceToString(reqData["password"])

	if !strings.Contains(requestLogin, "@") {
		rootLogin := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreRootLogin))
		rootPassword := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreRootPassword))

		if requestLogin == rootLogin && requestPassword == rootPassword {
			params.Session.Set(api.ConstSessionKeyAdminRights, true)

			return "ok", nil
		}
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3f10710a748442acaf49c69bce11ec13", "wrong login - should be email")
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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "13a80ab1d44e4a90979cea6914d9c012", "wrong password")
	}

	// api session updates
	if visitorModel.IsValidated() {
		params.Session.Set(visitor.ConstSessionKeyVisitorID, visitorModel.GetID())
	} else {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "29fba7a4bd85400e81c269189c50d0d0", "visitor is not validated, please check "+visitorModel.GetEmail()+" for a verify link we sent you")
	}

	if visitorModel.IsAdmin() {
		params.Session.Set(api.ConstSessionKeyAdminRights, true)
	}

	return "ok", nil
}

// WEB REST API function used to make login/registration via Facebook
//   - access_token and user_id params needed
//   - user needed information will be taken from Facebook
func restLoginFacebook(params *api.StructAPIHandlerParams) (interface{}, error) {
	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(reqData, "access_token") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b4b0356d6c174b63bc720cfc943535e2", "access_token was not specified")
	}

	if !utils.KeysInMapAndNotBlank(reqData, "user_id") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "17b8481dcf784edbb866f067d496915d", "user_id was not specified")
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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9b304209107a4ba58329ebfaebb70ff5", "Can't use google API: "+facebookResponse.Status)
	}

	responseData, err := ioutil.ReadAll(facebookResponse.Body)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// response json workaround
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(responseData, &jsonMap)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.StrKeysInMap(jsonMap, "id", "email", "first_name", "last_name", "verified") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6258ffff833649ef9aafdd15f578a16f", "unexpected facebook response")
	}

	// trying to load visitor from our DB
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorModel, ok := model.(visitor.InterfaceVisitor)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "42144bb637a24dd087d8542270e3c5b7", "visitor model is not InterfaceVisitor campatible")
	}

	// trying to load visitor by facebook_id
	err = visitorModel.LoadByFacebookID(reqData["user_id"].(string))
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
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "904b3ad7435048238f8c6f278b032168", "facebook account email unverified")
			}

			visitorModel.Set("facebook_id", jsonMap["id"])

			err := visitorModel.Save()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	// api session updates
	params.Session.Set(visitor.ConstSessionKeyVisitorID, visitorModel.GetID())

	if visitorModel.IsAdmin() {
		params.Session.Set(api.ConstSessionKeyAdminRights, true)
	}

	return "ok", nil
}

// WEB REST API function used to make login/registration via Google
//   - access_token param needed
//   - user needed information will be taken from Google
func restLoginGoogle(params *api.StructAPIHandlerParams) (interface{}, error) {
	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "03a6782ee34f46218e1a0d5d5d0dc0c3", "unexpected request content")
	}

	if !utils.KeysInMapAndNotBlank(reqData, "access_token") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6838cf4fc1bc41fbb73e426be0ee6f17", "access_token was not specified")
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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6bc2fd248ca3442c91d8341428c3d45d", "Can't use google API: "+googleResponse.Status)
	}

	responseData, err := ioutil.ReadAll(googleResponse.Body)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// response json workaround
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(responseData, &jsonMap)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.StrKeysInMap(jsonMap, "id", "email", "verified_email", "given_name", "family_name") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d793e87298b444c8ac74351a177dab68", "unexpected google response")
	}

	// trying to load visitor from our DB
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorModel, ok := model.(visitor.InterfaceVisitor)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "ec3f04a0a9fb4096abcf38062355e0a6", "visitor model is not InterfaceVisitor campatible")
	}

	// trying to load visitor by google_id
	err = visitorModel.LoadByGoogleID(jsonMap["email"].(string))
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
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b7b65d5c03724c9f816ac5157110894d", "google account email unverified")
			}

			visitorModel.Set("google_id", jsonMap["id"])

			err := visitorModel.Save()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	// api session updates
	params.Session.Set(visitor.ConstSessionKeyVisitorID, visitorModel.GetID())

	if visitorModel.IsAdmin() {
		params.Session.Set(api.ConstSessionKeyAdminRights, true)
	}

	return "ok", nil
}

// WEB REST API function used to get visitor order details information
func restVisitorOrderDetails(params *api.StructAPIHandlerParams) (interface{}, error) {
	visitorID := visitor.GetCurrentVisitorID(params)
	if visitorID == "" {
		return "you are not logined in", nil
	}

	orderModel, err := order.LoadOrderByID(params.RequestURLParams["id"])
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if utils.InterfaceToString(orderModel.Get("visitor_id")) != visitorID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c5ca1fdb70084a1ca1689df544df9825", "order is not belongs to logined user")
	}

	result := orderModel.ToHashMap()
	result["items"] = orderModel.GetItems()

	return result, nil
}

// WEB REST API function used to get visitor orders information
func restListVisitorOrders(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		reqData = make(map[string]interface{})
	}

	// list operation
	//---------------
	visitorID := visitor.GetCurrentVisitorID(params)
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

// WEB REST API function used to send email to a Visitor
func restVisitorSendMail(params *api.StructAPIHandlerParams) (interface{}, error) {
	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		reqData = make(map[string]interface{})
	}

	if !utils.StrKeysInMap(reqData, "subject", "content", "visitor_ids") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "21ac9f9d956f4963b37dfbd0469b8890", "'visitor_ids', 'subject' or 'content' field was not set")
	}

	subject := utils.InterfaceToString(reqData["subject"])
	content := utils.InterfaceToString(reqData["content"])
	visitorIDs := utils.InterfaceToArray(reqData["visitor_ids"])

	for _, visitorID := range visitorIDs {
		visitorModel, err := visitor.LoadVisitorByID(utils.InterfaceToString(visitorID))
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
