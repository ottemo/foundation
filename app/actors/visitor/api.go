package visitor

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"strconv"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	// Dashboard API
	err := api.GetRestService().RegisterAPI("visitor", api.ConstRESTOperationCreate, APICreateVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/:visitorID", api.ConstRESTOperationUpdate, APIUpdateVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/:visitorID", api.ConstRESTOperationDelete, APIDeleteVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/:visitorID", api.ConstRESTOperationGet, APIGetVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("visitors", api.ConstRESTOperationGet, APIListVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitors/attributes", api.ConstRESTOperationGet, APIListVisitorAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitors/attribute/:attribute", api.ConstRESTOperationDelete, APIDeleteVisitorAttribute)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitors/attribute/:attribute", api.ConstRESTOperationUpdate, APIUpdateVisitorAttribute)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitors/attribute", api.ConstRESTOperationCreate, APICreateVisitorAttribute)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Storefront API
	err = api.GetRestService().RegisterAPI("visitors/register", api.ConstRESTOperationCreate, APIRegisterVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitors/available", api.ConstRESTOperationCreate, APIVisitorEmailAvailable)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitors/validate/:key", api.ConstRESTOperationGet, APIValidateVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitors/invalidate/:email", api.ConstRESTOperationGet, APIInvalidateVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitors/forgot-password/:email", api.ConstRESTOperationGet, APIForgotPassword)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitors/mail", api.ConstRESTOperationCreate, APIMailToVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("visit", api.ConstRESTOperationGet, APIGetVisit)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visit", api.ConstRESTOperationUpdate, APIUpdateVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visit/logout", api.ConstRESTOperationGet, APILogout)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visit/login", api.ConstRESTOperationCreate, APILogin)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visit/login-facebook", api.ConstRESTOperationCreate, APIFacebookLogin)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visit/login-google", api.ConstRESTOperationCreate, APIGoogleLogin)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visit/orders", api.ConstRESTOperationGet, APIGetOrders)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visit/order/:orderID", api.ConstRESTOperationGet, APIGetOrder)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APICreateVisitor creates a new visitor
//   - visitor attributes should be specified in content
//   - "email" attribute required
func APICreateVisitor(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(requestData, "email") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a9610b78-add9-4ae5-b757-59462b646d2b", "No email address specified, please specify an email address.")
	}

	// create operation
	//-----------------
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		// always lowercase email address
		if attribute == "email" {
			value = strings.ToLower(utils.InterfaceToString(value))
		}
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

// APIUpdateVisitor updates existing visitor
//   - visitor id should be specified in "visitorID" argument
//   - visitor attributes should be specified in content
func APIUpdateVisitor(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	visitorID := context.GetRequestArgument("visitorID")
	if visitorID == "" {
		sessionValue := context.GetSession().Get(visitor.ConstSessionKeyVisitorID)
		sessionVisitorID, ok := sessionValue.(string)
		if !ok {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e7a97b45-a22a-48c5-96f8-f4ecd3f8380f", "Not logged in, please login.")
		}
		visitorID = sessionVisitorID
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorModel, err := visitor.LoadVisitorByID(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if err := api.ValidateAdminRights(context); err != nil {
		if visitor.GetCurrentVisitorID(context) != visitorID {
			return nil, env.ErrorDispatch(err)
		}
		if _, present := requestData["is_admin"]; present {
			return nil, env.ErrorDispatch(err)
		}
		// check when not admin try to change password, validate old password
		if _, present := requestData["password"]; present {
			if oldPass, present := requestData["old_password"]; present {
				if ok := visitorModel.CheckPassword(utils.InterfaceToString(oldPass)); !ok {
					return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "13a80ab1-d44e-4a90-979c-ea6914d9c012", "Password entered does not match stored password.")
				}
				delete(requestData, "old_password")
			} else {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "157df5fa-d775-4934-af94-b77ef8c826e9", "Please enter current password and try again.")
			}
		}
		// When admin user change password from storefront we will validate it
	} else if oldPass, present := requestData["old_password"]; present {
		if ok := visitorModel.CheckPassword(utils.InterfaceToString(oldPass)); !ok {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "13a80ab1-d44e-4a90-979c-ea6914d9c012", "Password entered does not match stored password.")
		}
		delete(requestData, "old_password")
	}

	// update operation
	//-----------------

	for attribute, value := range requestData {
		if err := visitorModel.Set(attribute, value); err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	if err = visitorModel.Save(); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorModel.ToHashMap(), nil
}

// APIDeleteVisitor deletes existing visitor
//   - visitor id should be specified in "visitorID" argument
func APIDeleteVisitor(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorID := context.GetRequestArgument("visitorID")
	if visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "157df5fa-d775-4934-af94-b77ef8c826e9", "No Visitor ID found, please specify a Visitor ID.")
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

// APIGetVisitor returns visitor information
//   - visitor id should be specified in "visitorID" argument
func APIGetVisitor(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorID := context.GetRequestArgument("visitorID")
	if visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "58630919-23f5-4406-8676-a7d1a629b35f", "No Visitor ID found, please specify a Visitor ID.")
	}

	// get operation
	//--------------
	visitorModel, err := visitor.LoadVisitorByID(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorModel.ToHashMap(), nil
}

// APIListVisitors returns a list of existing visitors
//   - if "action" parameter is set to "count" result value will be just a number of list items
func APIListVisitors(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorCollectionModel, err := visitor.GetVisitorCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// filters handle
	models.ApplyFilters(context, visitorCollectionModel.GetDBCollection())

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return visitorCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	visitorCollectionModel.ListLimit(models.GetListLimit(context))

	// extra parameter handle
	models.ApplyExtraAttributes(context, visitorCollectionModel)

	return visitorCollectionModel.List()
}

// APIListVisitorAttributes returns a list of visitor attributes
func APIListVisitorAttributes(context api.InterfaceApplicationContext) (interface{}, error) {
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attrInfo := visitorModel.GetAttributesInfo()
	return attrInfo, nil
}

// APICreateVisitorAttribute creates a new custom attribute for a visitor model
//   - attribute parameters "Attribute" and "Label" are required
func APICreateVisitorAttribute(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attributeName, isSpecified := requestData["Attribute"]
	if !isSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f91e2342-d19d-4d83-8c6d-f7dc4237a49f", "Attribute Name was not specified, this is a required field.")
	}

	attributeLabel, isSpecified := requestData["Label"]
	if !isSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "23c7d427-ab41-4183-b446-bf5ffef2e8fc", "Attribute Label was not specified, this is a required field.")
	}

	// create visitor by filling model
	//---------------------------------
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attribute := models.StructAttributeInfo{
		Model:      visitor.ConstModelNameVisitor,
		Collection: ConstCollectionNameVisitor,
		Attribute:  utils.InterfaceToString(attributeName),
		Type:       utils.ConstDataTypeText,
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

	for key, value := range requestData {
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

// APIUpdateVisitorAttribute updates existing custom attribute of visitor model
//   - attribute name/code should be provided in "attribute" argument
//   - attribute parameters should be provided in request content
//   - attribute parameters "id" and "name" will be ignored
//   - static attributes can not be changed
func APIUpdateVisitorAttribute(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attributeName := context.GetRequestArgument("attribute")
	if attributeName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cb8f7251-e22b-4605-97bb-e239df6c7aac", "Attribute Name was not specified, this is a required field.")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for _, attribute := range visitorModel.GetAttributesInfo() {
		if attribute.Attribute == attributeName {
			if attribute.IsStatic == true {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2893262f-a61a-42f8-9c75-e763e0a5c8ca", "Attribute is static, cannot edit static attributes.")
			}

			for key, value := range requestData {
				switch strings.ToLower(key) {
				case "label":
					attribute.Label = utils.InterfaceToString(value)
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
				case "ispublic", "public":
					attribute.IsPublic = utils.InterfaceToBool(value)
				}
			}
			err := visitorModel.EditAttribute(attributeName, attribute)
			if err != nil {
				return nil, err
			}
			return attribute, nil
		}
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2893262f-a61a-42f8-9c75-e763e0a5c8ca", "Unable to find specified attribute.")
}

// APIDeleteVisitorAttribute removes existing custom attribute of a visitor model
//   - attribute name/code should be provided in "attribute" argument
func APIDeleteVisitorAttribute(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//--------------------
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attributeName := context.GetRequestArgument("attribute")
	if attributeName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4b0f0edf-6926-42d8-a625-7d4d424ee819", "Attribute Name was not specified, this is required field.")
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

// APIRegisterVisitor creates a new visitor, sends verification link to visitor email
//   - visitor attributes should be included contents
//   - "email" attribute required
func APIRegisterVisitor(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(requestData, "email") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a37ff1e8-68e3-4201-a7b4-c9b4356dcbeb", "An email address was not specified, this is a required field.")
	}

	// register visitor operation
	//---------------------------
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
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

	// check if email verification is necessary to login visitor
	verifyEmail := utils.InterfaceToBool(env.ConfigGetValue(app.ConstConfigPathVerfifyEmail))
	if verifyEmail == true {
		// force the visitor to log in
		err = visitorModel.Invalidate()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	} else {
		// log visitor in, if site is not using verification emails
		context.GetSession().Set(visitor.ConstSessionKeyVisitorID, visitorModel.GetID())
	}

	return visitorModel.ToHashMap(), nil
}

// APIVisitorEmailAvailable creates a new visitor, sends verification link to visitor email
//   - "email" attribute required and email address in valid format
func APIVisitorEmailAvailable(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(requestData, "email") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fec90530-6441-49bf-aa55-187e18b1fc7c", "An email address was not specified, this is a required field.")
	}
	visitorEmail := strings.ToLower(utils.InterfaceToString(requestData["email"]))

	if !utils.ValidEmail(visitorEmail) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "362efeb5-331a-484c-b5f3-daef1e7b065b", "The email address specified is not in a valid format, "+visitorEmail+".")
	}
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.LoadByEmail(visitorEmail)
	if err != nil {
		// requested email address is available
		return true, nil
	}

	// requested email address is already in use
	return false, nil
}

// APIValidateVisitors validates visitors by a key which sends after registration
//   - verification key should be provided in "key" argument
func APIValidateVisitors(context api.InterfaceApplicationContext) (interface{}, error) {

	verificationKey := context.GetRequestArgument("key")
	if verificationKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "0b1c8414-18bf-4c44-8b29-c462d38d5498", "Verification key was not specified, this is a required field.")
	}

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.Validate(verificationKey)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return api.StructRestRedirect{Result: "ok", Location: app.GetStorefrontURL("login")}, nil
}

// APIInvalidateVisitor invalidate visitor by email, sends a new verification key to email
//   - visitor email should be specified in "email" argument
func APIInvalidateVisitor(context api.InterfaceApplicationContext) (interface{}, error) {

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorEmail := strings.ToLower(context.GetRequestArgument("email"))

	err = visitorModel.LoadByEmail(visitorEmail)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// checking rights
	if err := api.ValidateAdminRights(context); err != nil {
		if visitorModel.IsVerfied() {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = visitorModel.Invalidate()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIForgotPassword changes and sends a new password to visitor e-mail
func APIForgotPassword(context api.InterfaceApplicationContext) (interface{}, error) {

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorEmail := strings.ToLower(context.GetRequestArgument("email"))

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

// APIGetVisit returns current visit related information
func APIGetVisit(context api.InterfaceApplicationContext) (interface{}, error) {

	// so if user was logged to app as admin, we want to reflect this
	isAdmin := false
	if api.ValidateAdminRights(context) == nil {
		isAdmin = true
	}

	sessionValue := context.GetSession().Get(visitor.ConstSessionKeyVisitorID)
	visitorID, ok := sessionValue.(string)
	if !ok {
		if isAdmin {
			return map[string]interface{}{"is_admin": true}, nil
		}
		return "No session found, please log in first.", nil
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

// APILogout makes logout for current visit
func APILogout(context api.InterfaceApplicationContext) (interface{}, error) {

	// if session cookie is set, expire it
	request := context.GetRequest()
	// use secure cookies by default
	var flagSecure = true
	var tmpSecure = ""
	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("secure_cookie", tmpSecure); iniValue != "" {
			tmpSecure = iniValue
			flagSecure, _ = strconv.ParseBool(tmpSecure)
		}
	}

	if request, ok := request.(*http.Request); ok {
		responseWriter := context.GetResponseWriter()
		if responseWriter, ok := responseWriter.(http.ResponseWriter); ok {

			// check for session cookie
			cookie, err := request.Cookie(api.ConstSessionCookieName)
			if err == nil {

				// expire the cookie
				cookieExpires := time.Now().AddDate(0, -1, 0)
				cookie = &http.Cookie{
					Name:     api.ConstSessionCookieName,
					Value:    "",
					Path:     "/",
					HttpOnly: true,
					MaxAge:   -1,
					Secure:   flagSecure,
					Expires:  cookieExpires,
				}

				http.SetCookie(responseWriter, cookie)
			}

		}
	}

	context.GetSession().Close()

	return "ok", nil
}

// APILogin makes login for a current visit
//   - "email" and "password" attributes required
func APILogin(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(requestData, "email", "password") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "02b0583b-28c3-4072-afe2-f392a163ed87", "Email Address and/or Password was not specified, these are required fields.")
	}

	requestLogin := strings.ToLower(utils.InterfaceToString(requestData["email"]))
	requestPassword := utils.InterfaceToString(requestData["password"])

	if !strings.Contains(requestLogin, "@") {
		rootLogin := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreRootLogin))
		rootPassword := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreRootPassword))

		if requestLogin == rootLogin && requestPassword == rootPassword {
			context.GetSession().Set(api.ConstSessionKeyAdminRights, true)

			return "ok", nil
		}
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3f10710a-7484-42ac-af49-c69bce11ec13", "Please enter a valid email address in the correct format.")
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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "13a80ab1-d44e-4a90-979c-ea6914d9c012", "The password entered does not match the stored password.")
	}

	// api session updates
	if visitorModel.IsVerfied() {
		context.GetSession().Set(visitor.ConstSessionKeyVisitorID, visitorModel.GetID())
	} else {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "29fba7a4-bd85-400e-81c2-69189c50d0d0", "This account has not been verfied, please check your email account: ,"+visitorModel.GetEmail()+" for a verification link sent to you.")
	}

	if visitorModel.IsAdmin() {
		context.GetSession().Set(api.ConstSessionKeyAdminRights, true)
	}

	return "ok", nil
}

// APIFacebookLogin makes login and/or registration via Facebook
//   - "access_token" and "user_id" arguments required needed
//   - visitor attributes will be taken from Facebook
func APIFacebookLogin(context api.InterfaceApplicationContext) (interface{}, error) {
	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(requestData, "access_token") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b4b0356d-6c17-4b63-bc72-0cfc943535e2", "No access token found named 'access_token', this is a required field.")
	}

	if !utils.KeysInMapAndNotBlank(requestData, "user_id") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "17b8481d-cf78-4edb-b866-f067d496915d", "No ID found named 'user_id', this is a required field.")
	}

	// facebook login operation
	//-------------------------

	// using access token to get user information
	url := "https://graph.facebook.com/" + requestData["user_id"].(string) + "?access_token=" + requestData["access_token"].(string)
	facebookResponse, err := http.Get(url)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if facebookResponse.StatusCode != 200 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9b304209-107a-4ba5-8329-ebfaebb70ff5", "Facebook is not responding at this time, unable to use this login method. Current status code found: "+facebookResponse.Status)
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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6258ffff-8336-49ef-9aaf-dd15f578a16f", "The following fields are all required: id, email, first_name, last_name and verfied.")
	}

	// trying to load visitor from our DB
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorModel, ok := model.(visitor.InterfaceVisitor)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "42144bb6-37a2-4dd0-87d8-542270e3c5b7", "Unable to create visitor using Facebook credentials using using specified fields.")
	}

	// trying to load visitor by facebook_id
	err = visitorModel.LoadByFacebookID(requestData["user_id"].(string))
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
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "904b3ad7-4350-4823-8f8c-6f278b032168", "Account has not been verfied.  Please check your email account: "+visitorModel.GetEmail()+", for an email with a verification link.")
			}

			visitorModel.Set("facebook_id", jsonMap["id"])

			err := visitorModel.Save()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	// api session updates
	context.GetSession().Set(visitor.ConstSessionKeyVisitorID, visitorModel.GetID())

	if visitorModel.IsAdmin() {
		context.GetSession().Set(api.ConstSessionKeyAdminRights, true)
	}

	return "ok", nil
}

// APIGoogleLogin associates the specified email address with a Google account
//   - "access_token" attribute needed
//   - visitor attributes will be taken from Google
func APIGoogleLogin(context api.InterfaceApplicationContext) (interface{}, error) {
	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "03a6782e-e34f-4621-8e1a-0d5d5d0dc0c3", "Unexpected error occurred when parsing request.")
	}

	if !utils.KeysInMapAndNotBlank(requestData, "access_token") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6838cf4f-c1bc-41fb-b73e-426be0ee6f17", "No key-pair found named 'access_token' in request, this is a required field.")
	}

	// google login operation
	//-------------------------

	// using access token to get user information
	url := "https://www.googleapis.com/oauth2/v1/userinfo?access_token=" + requestData["access_token"].(string)
	googleResponse, err := http.Get(url)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if googleResponse.StatusCode != 200 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6bc2fd24-8ca3-442c-91d8-341428c3d45d", "Unable to access Google API at this time, current response code: "+googleResponse.Status)
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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d793e872-98b4-44c8-ac74-351a177dab68", "One of the following required key-pair fields are missing: id, email, verfified_email, given_name and family_name.")
	}

	// trying to load visitor from our DB
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorModel, ok := model.(visitor.InterfaceVisitor)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "ec3f04a0-a9fb-4096-abcf-38062355e0a6", "Unable to create visitor using Google credentials using using specified fields.")
	}

	// trying to load visitor by google_id
	err = visitorModel.LoadByGoogleID(jsonMap["email"].(string))
	if err != nil && strings.Contains(err.Error(), "Unable to find the email address, "+jsonMap["email"].(string)+", associated with the provided Google ID.") {

		// there is no such google_id in DB, trying to find by e-mail
		err = visitorModel.LoadByEmail(jsonMap["email"].(string))
		if err != nil && strings.Contains(err.Error(), "Unable to find the visitor account with the following email address: "+jsonMap["email"].(string)+".") {

			// As the visitor does not exist, create a new one
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
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b7b65d5c-0372-4c9f-816a-c5157110894d", "Account has not been verfied.  Please check your email account: "+visitorModel.GetEmail()+", for an email with a verification link.")
			}

			visitorModel.Set("google_id", jsonMap["id"])

			err := visitorModel.Save()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	// api session updates
	context.GetSession().Set(visitor.ConstSessionKeyVisitorID, visitorModel.GetID())

	if visitorModel.IsAdmin() {
		context.GetSession().Set(api.ConstSessionKeyAdminRights, true)
	}

	return "ok", nil
}

// APIGetOrder returns current visitor order details for specified order
//   - orderID should be specified in arguments
func APIGetOrder(context api.InterfaceApplicationContext) (interface{}, error) {
	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return "No Visitor ID found, unable to process order request. Please log in first.", nil
	}

	orderModel, err := order.LoadOrderByID(context.GetRequestArgument("orderID"))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if utils.InterfaceToString(orderModel.Get("visitor_id")) != visitorID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c5ca1fdb-7008-4a1c-a168-9df544df9825", "There is a mis-match between the current Visitor ID and the Visitor ID on the order.")
	}

	result := orderModel.ToHashMap()
	result["items"] = orderModel.GetItems()

	return result, nil
}

// APIGetOrders returns list of orders related to current vivitor
func APIGetOrders(context api.InterfaceApplicationContext) (interface{}, error) {

	// list operation
	//---------------
	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return "No Visitor ID found, unable to process request.  Please log in first.", nil
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
	models.ApplyFilters(context, orderCollection.GetDBCollection())

	// extra parameter handle
	models.ApplyExtraAttributes(context, orderCollection)

	result, err := orderCollection.List()

	return result, env.ErrorDispatch(err)
}

// APIMailToVisitor sends email to specified visitors
//   - "subject", "content", "visitor_ids" arguments required
func APIMailToVisitor(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, err
	}

	if !utils.StrKeysInMap(requestData, "subject", "content", "visitor_ids") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "21ac9f9d-956f-4963-b37d-fbd0469b8890", "One of the following fields was not set: visitor_ids, content, or subject.")
	}

	subject := utils.InterfaceToString(requestData["subject"])
	content := utils.InterfaceToString(requestData["content"])
	visitorIDs := utils.InterfaceToArray(requestData["visitor_ids"])

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
