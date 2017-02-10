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
	"github.com/ottemo/foundation/app/models/visitor"

	"github.com/ottemo/foundation/app/models/order"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	// Dashboard API
	service.POST("visitor", api.IsAdmin(APICreateVisitor))
	service.PUT("visitor/:visitorID", APIUpdateVisitor)
	service.DELETE("visitor/:visitorID", api.IsAdmin(APIDeleteVisitor))
	service.GET("visitor/:visitorID", api.IsAdmin(APIGetVisitor))

	service.GET("visitors", api.IsAdmin(APIListVisitors))
	service.GET("visitors/attributes", APIListVisitorAttributes)
	service.DELETE("visitors/attribute/:attribute", api.IsAdmin(APIDeleteVisitorAttribute))
	service.PUT("visitors/attribute/:attribute", api.IsAdmin(APIUpdateVisitorAttribute))
	service.POST("visitors/attribute", api.IsAdmin(APICreateVisitorAttribute))
	service.GET("visitors/guests", api.IsAdmin(APIGetGuestsList))

	// Storefront API
	service.POST("visitors/register", APIRegisterVisitor)
	service.POST("visitors/available", APIVisitorEmailAvailable)
	service.GET("visitors/validate/:key", APIValidateVisitors)
	service.GET("visitors/invalidate/:email", APIInvalidateVisitor)
	service.GET("visitors/forgot-password/:email", APIForgotPassword)
	service.POST("visitors/mail", APIMailToVisitor)

	service.POST("visitors/reset-password", APIResetPassword)

	service.GET("visit", APIGetVisit)
	service.PUT("visit", APIUpdateVisitor)
	service.GET("visit/logout", APILogout)
	service.POST("visit/login", APILogin)
	service.POST("visit/login-facebook", APIFacebookLogin)
	service.POST("visit/login-google", APIGoogleLogin)

	return nil
}

// APIGetGuestsList returns list of guests GROUP(ed) BY customer_email
func APIGetGuestsList(context api.InterfaceApplicationContext) (interface{}, error) {

	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dbCollection := orderCollectionModel.GetDBCollection()
	if err := dbCollection.AddFilter("visitor_id", "=", ""); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "388a3235-8be4-4d89-b705-720233451b29", err.Error())
	}

	// handle filters
	if err := models.ApplyFilters(context, dbCollection); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c76e230f-43b7-456b-8017-9219847f5374", err.Error())
	}

	// remove limits to get right GROUP BY result
	if err := dbCollection.SetLimit(0, 0); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6383e2e6-caf1-4c9f-b668-3230f2ce7c59", err.Error())
	}

	orders := orderCollectionModel.ListOrders()

	guests := make([]map[string]interface{}, 0)
	uniqEmails := make([]string, 0)

	for _, order := range orders {
		orderItem := order.ToHashMap()
		email := utils.InterfaceToString(orderItem["customer_email"])
		name := utils.InterfaceToString(orderItem["customer_name"])

		// GROUP BY customer_email
		if !utils.IsInListStr(email, uniqEmails) && name != "" {
			item := map[string]interface{}{
				"customer_email": email,
				"customer_name":  name,
			}

			uniqEmails = append(uniqEmails, email)
			guests = append(guests, item)
		}
	}

	count := len(guests)

	// limit handle
	from, howMany := models.GetListLimit(context)
	to := from + howMany
	if to > count {
		to = count
	}

	result := map[string]interface{}{
		"guests": guests[from:to],
		"count":  count,
	}

	return result, nil
}

// APICreateVisitor creates a new visitor
//   - visitor attributes should be specified in content
//   - "email" attribute required
func APICreateVisitor(context api.InterfaceApplicationContext) (interface{}, error) {

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
	if err := visitorModel.Set("created_at", time.Now()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "38f3fe0f-4ab6-43f0-9677-8b4e4a3ba45a", err.Error())
	}

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
		// Visitor. Not admin.
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
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2ba2d5ac-84f6-421e-a7a6-da10c16e85f3", "Please enter current password and try again.")
			}
		}
	} else if oldPass, present := requestData["old_password"]; present {
		// Admin
		// When admin user change password from storefront we will validate it
		if ok := visitorModel.CheckPassword(utils.InterfaceToString(oldPass)); !ok {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9c4c53c3-6b8d-4719-8515-7a31df1d08bb", "Password entered does not match stored password.")
		}
		delete(requestData, "old_password")
	}

	// update operation
	//-----------------

	for attribute, value := range requestData {
		if attribute == "email" {
			visitorEmail := strings.ToLower(utils.InterfaceToString(value))

			if !utils.ValidEmailAddress(visitorEmail) {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "10dce5f7-ff43-41b5-b05f-48731fe22115", "The email address specified is not in a valid format, "+visitorEmail+".")
			}

			visitorCollectionModel, err := visitor.GetVisitorCollectionModel()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			err = visitorCollectionModel.ListFilterAdd("email", "=", visitorEmail)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			visitorItems, err := visitorCollectionModel.List()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			if len(visitorItems) > 0 {
				if len(visitorItems) > 1 {
					// Do not say to visitor that email is in use - security reason
					return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b9f147cd-ab4f-489a-81de-34d568331fc1", "Wrong email. Please, ask site administrator.")
				} else if visitorID != visitorItems[0].ID { // we have only one record
					// Do not say to visitor that email is in use - security reason
					return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "65ab4441-ce5e-47fb-b21b-9bf8a94ee075", "Wrong email. Please, ask site administrator.")
				} // else Ok: visitor with visitorID is the owner of email
			} // else it's ok - no visitor with this email found
		}

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

	visitorCollectionModel, err := visitor.GetVisitorCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// filters handle
	if err := models.ApplyFilters(context, visitorCollectionModel.GetDBCollection()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a77781e4-1c24-44c6-9dbd-307b7d2a1713", err.Error())
	}

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return visitorCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	if err := visitorCollectionModel.ListLimit(models.GetListLimit(context)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "578bd204-2e56-4b86-a72f-e5475d20a69c", err.Error())
	}

	// extra parameter handle
	if err := models.ApplyExtraAttributes(context, visitorCollectionModel); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "57962242-bf9c-4c11-847a-a21656a378b1", err.Error())
	}

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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e7d28bb9-5969-40b8-ae59-9134baa4d8d1", "Attribute Name was not specified, this is a required field.")
	}

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for _, attribute := range visitorModel.GetAttributesInfo() {
		if attribute.Attribute == attributeName {
			if attribute.IsStatic == true {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "8a4e1d22-cd3f-4da9-a9b0-d6950245b349", "Attribute is static, cannot edit static attributes.")
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

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4137313b-f2ae-4545-923e-eac0324d923f", "Unable to find specified attribute.")
}

// APIDeleteVisitorAttribute removes existing custom attribute of a visitor model
//   - attribute name/code should be provided in "attribute" argument
func APIDeleteVisitorAttribute(context api.InterfaceApplicationContext) (interface{}, error) {

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

	if err := visitorModel.Set("created_at", time.Now()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "0a12b807-6998-4d68-91bb-ce85ad804b08", err.Error())
	}

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

	if !utils.ValidEmailAddress(visitorEmail) {
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
		if visitorModel.IsVerified() {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = visitorModel.Invalidate()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIResetPassword update visitor password by using verification key from email
func APIResetPassword(context api.InterfaceApplicationContext) (interface{}, error) {

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(requestData, "key", "password") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "991eceb5-37a6-4044-b0e9-5039055044a0", "Verification key and new password are required.")
	}

	verificationKey := utils.InterfaceToString(requestData["key"])
	newPassword := utils.InterfaceToString(requestData["password"])

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.UpdateResetPassword(verificationKey, newPassword)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorModel.GetEmail(), nil
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
		if strings.Contains(err.Error(), "Unable to find") {
			return "ok", nil
		}

		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.ResetPassword()
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

	if err := context.GetSession().Close(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "8fc6c2fa-7f0e-475b-856f-d4af161812fd", err.Error())
	}

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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9fd0a895-4b42-4d89-9aa7-9104ba23f96a", "The password entered does not match the stored password.")
	}

	// api session updates
	if visitorModel.IsVerified() {
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

	// using access token to get user information,
	// NOTE: API Version supported at least until April 2018
	url := "https://graph.facebook.com/v2.5/" + utils.InterfaceToString(requestData["user_id"]) + "?access_token=" + utils.InterfaceToString(requestData["access_token"]) + "&fields=id,email,first_name,last_name"
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

	if !utils.StrKeysInMap(jsonMap, "id", "email", "first_name", "last_name") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6258ffff-8336-49ef-9aaf-dd15f578a16f", "The following fields are all required: id, email, first_name, last_name.")
	}

	if responseError, present := jsonMap["error"]; present && responseError != nil {
		errorMap := utils.InterfaceToMap(responseError)
		if errorMessage, present := errorMap["message"]; present {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9f003c22-65e2-4c17-abb7-cd2a74b7eb44", utils.InterfaceToString(errorMessage))
		}
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "df03de97-e779-4779-b6c9-e8d331bf590d", "Unexpected error during retrieving information from facebook service")
	}

	visitorFacebookAccount := map[string]string{
		"id":         utils.InterfaceToString(jsonMap["id"]),
		"email":      utils.InterfaceToString(jsonMap["email"]),
		"first_name": utils.InterfaceToString(jsonMap["first_name"]),
		"last_name":  utils.InterfaceToString(jsonMap["last_name"]),
	}

	// trying to load visitor from our DB
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// trying to load visitor by facebook_id
	err = visitorModel.LoadByFacebookID(visitorFacebookAccount["id"])
	if err != nil && strings.Contains(err.Error(), "Unable to find") {

		// there is no such facebook_id in DB, trying to find by e-mail
		err = visitorModel.LoadByEmail(visitorFacebookAccount["email"])

		// visitor not exists in out DB - creating new one
		if err != nil && strings.Contains(err.Error(), "Unable to find") {
			if err := visitorModel.Set("email", visitorFacebookAccount["email"]); err != nil {
				_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "bdc1cd48-51e1-4c6d-a2b9-8b1e32ce813a", err.Error())
			}
			if err := visitorModel.Set("first_name", visitorFacebookAccount["first_name"]); err != nil {
				_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4416efb3-af73-443b-a741-a36e437b77d1", err.Error())
			}
			if err := visitorModel.Set("last_name", visitorFacebookAccount["last_name"]); err != nil {
				_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "0ce238fe-189c-48dd-af30-3a12c6b209b2", err.Error())
			}
			if err := visitorModel.Set("created_at", time.Now()); err != nil {
				_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "dc2dbf8e-7754-4128-ae9b-d8915f00de47", err.Error())
			}
		}

		if err := visitorModel.Set("facebook_id", visitorFacebookAccount["id"]); err != nil {
			_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2f40a1c5-8a80-44a0-8317-cb1cf1a64484", err.Error())
		}

		err := visitorModel.Save()
		if err != nil {
			return nil, env.ErrorDispatch(err)
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

	visitorGoogleAccount := map[string]string{
		"id":          utils.InterfaceToString(jsonMap["id"]),
		"email":       utils.InterfaceToString(jsonMap["email"]),
		"given_name":  utils.InterfaceToString(jsonMap["given_name"]),
		"family_name": utils.InterfaceToString(jsonMap["family_name"]),
	}

	// trying to load visitor from our DB
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// trying to load visitor by google_id
	err = visitorModel.LoadByGoogleID(visitorGoogleAccount["id"])
	if err != nil && strings.Contains(err.Error(), "Unable to find") {

		// there is no such google_id in DB, trying to find by e-mail
		err = visitorModel.LoadByEmail(visitorGoogleAccount["email"])
		if err != nil && strings.Contains(err.Error(), "Unable to find") {

			// As the visitor does not exist, create a new one
			if err := visitorModel.Set("email", visitorGoogleAccount["email"]); err != nil {
				_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1d871e0a-4e16-423c-8a11-e551ec771e4f", err.Error())
			}
			if err := visitorModel.Set("first_name", visitorGoogleAccount["given_name"]); err != nil {
				_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "168d491d-e432-4fd6-8cd4-f650857cab82", err.Error())
			}
			if err := visitorModel.Set("last_name", visitorGoogleAccount["family_name"]); err != nil {
				_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9a4607bb-d439-48a6-a625-b8b5a6ae04de", err.Error())
			}
			if err := visitorModel.Set("created_at", time.Now()); err != nil {
				_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "de3bc9ec-51f7-44b0-85a3-d8b955e7485d", err.Error())
			}
		}

		if err := visitorModel.Set("google_id", visitorGoogleAccount["id"]); err != nil {
			_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "8f1f5929-0f72-4157-b341-af669cb82e37", err.Error())
		}

		err := visitorModel.Save()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	// api session updates
	context.GetSession().Set(visitor.ConstSessionKeyVisitorID, visitorModel.GetID())

	if visitorModel.IsAdmin() {
		context.GetSession().Set(api.ConstSessionKeyAdminRights, true)
	}

	return "ok", nil
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
