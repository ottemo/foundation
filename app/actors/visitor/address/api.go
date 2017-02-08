package address

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.POST("visitor/:visitorID/address", APICreateVisitorAddress)
	service.PUT("visitor/:visitorID/address/:addressID", APIUpdateVisitorAddress)
	service.DELETE("visitor/:visitorID/address/:addressID", APIDeleteVisitorAddress)

	service.GET("visitor/:visitorID/addresses", APIListVisitorAddresses)

	service.GET("visitors/addresses/attributes", api.IsAdmin(APIListVisitorAddressAttributes))
	service.DELETE("visitors/address/:addressID", APIDeleteVisitorAddress)
	service.PUT("visitors/address/:addressID", APIUpdateVisitorAddress)
	service.GET("visitors/address/:addressID", APIGetVisitorAddress)

	service.POST("visit/address", APICreateVisitorAddress)
	service.PUT("visit/address/:addressID", APIUpdateVisitorAddress)
	service.DELETE("visit/address/:addressID", APIDeleteVisitorAddress)
	service.GET("visit/addresses", APIListVisitorAddresses)
	service.GET("visit/address/:addressID", APIGetVisitorAddress)

	return nil
}

// APICreateVisitorAddress creates a new visitor address
//   - visitor address attributes should be specified in content
//   - "visitor_id" attribute required
func APICreateVisitorAddress(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if _, ok := requestData["visitor_id"]; !ok {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a9da4ac4-d073-48f3-b062-2ba536d2c577", "No Visitor ID found, unable to create address entry.  Please log in first.")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		if requestData["visitor_id"] != visitor.GetCurrentVisitorID(context) {
			return nil, env.ErrorDispatch(err)
		}
	}

	// create visitor address operation
	//---------------------------------
	visitorAddressModel, err := checkout.ValidateAddress(requestData)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorAddressModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorAddressModel.ToHashMap(), nil
}

// APIUpdateVisitorAddress updates existing visitor address
//   - visitor address id must be specified in "addressID" argument
//   - visitor address attributes should be specified in content
func APIUpdateVisitorAddress(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	addressID := context.GetRequestArgument("addressID")
	if addressID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fe7814c0-85fe-4d60-a134-415f7ac12075", "No visitor Address ID found, unable to process update request.")
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressModel, err := visitor.LoadVisitorAddressByID(addressID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		if visitorAddressModel.GetVisitorID() != visitor.GetCurrentVisitorID(context) {
			return nil, env.ErrorDispatch(err)
		}
	}

	_, err = checkout.ValidateAddress(requestData)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// update operation
	//-----------------
	for attribute, value := range requestData {
		err := visitorAddressModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = visitorAddressModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorAddressModel.ToHashMap(), nil
}

// APIDeleteVisitorAddress deletes existing visitor address
//   - visitor address id must be specified in "addressID" argument
func APIDeleteVisitorAddress(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//--------------------
	addressID := context.GetRequestArgument("addressID")
	if addressID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "eec1ef1b-25d9-4dbe-8bd2-b907a0897203", "No Visitor ID found, unable to process request.  Please log in first.")
	}

	visitorAddressModel, err := visitor.LoadVisitorAddressByID(addressID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		if visitorAddressModel.GetVisitorID() != visitor.GetCurrentVisitorID(context) {
			return nil, env.ErrorDispatch(err)
		}
	}

	// delete operation
	err = visitorAddressModel.Delete()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIListVisitorAddressAttributes returns a list of visitor address attributes
func APIListVisitorAddressAttributes(context api.InterfaceApplicationContext) (interface{}, error) {

	visitorAddressModel, err := visitor.GetVisitorAddressModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attrInfo := visitorAddressModel.GetAttributesInfo()
	return attrInfo, nil
}

// APIListVisitorAddresses returns visitor addresses list
//   - visitor id must be specified in "visitorID" argument
func APIListVisitorAddresses(context api.InterfaceApplicationContext) (interface{}, error) {

	// if visitorID was specified - using this otherwise, taking current visitor
	visitorID := context.GetRequestArgument("visitorID")
	if visitorID == "" {

		sessionVisitorID := visitor.GetCurrentVisitorID(context)
		if sessionVisitorID == "" {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2ac4c16b-9241-406e-b35a-399813bb6ca5", "No Visitor ID found, unable to retrieve associated addresses.  Please log in first.")
		}
		visitorID = sessionVisitorID
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		if visitorID != visitor.GetCurrentVisitorID(context) {
			return nil, env.ErrorDispatch(err)
		}
	}

	// list operation
	//---------------
	visitorAddressCollectionModel, err := visitor.GetVisitorAddressCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := visitorAddressCollectionModel.GetDBCollection()
	if err := dbCollection.AddStaticFilter("visitor_id", "=", visitorID); err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1a6fcea9-b640-4ad5-ad3d-550d55b3d99a", err.Error())
	}

	// filters handle
	if err := models.ApplyFilters(context, dbCollection); err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9fc25081-ecbb-4ac8-a5b5-b42de55afd07", err.Error())
	}

	// checking for a "count" request
	if context.GetRequestArgument("count") != "" {
		return visitorAddressCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	if err := visitorAddressCollectionModel.ListLimit(models.GetListLimit(context)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b7021bca-b95a-4e34-815b-92d70aa98abf", err.Error())
	}

	// extra parameter handle
	if err := models.ApplyExtraAttributes(context, visitorAddressCollectionModel); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fe8a4498-c21d-4492-a6dd-010fcfa52bec", err.Error())
	}

	return visitorAddressCollectionModel.List()
}

// APIGetVisitorAddress returns visitor address information
//   - visitor address id must be specified in "addressID" argument
func APIGetVisitorAddress(context api.InterfaceApplicationContext) (interface{}, error) {
	visitorAddressID := context.GetRequestArgument("addressID")
	if visitorAddressID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b94882c6-bbdd-428d-88b0-7ea5623d80f7", "No Visitor ID found, unable to retrieve associated address.  Please log in first.")
	}

	visitorAddressModel, err := visitor.LoadVisitorAddressByID(visitorAddressID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		if visitorAddressModel.GetVisitorID() != visitor.GetCurrentVisitorID(context) {
			return nil, env.ErrorDispatch(err)
		}
	}

	return visitorAddressModel.ToHashMap(), nil
}
