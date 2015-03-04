package address

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	err := api.GetRestService().RegisterAPI("visitor/:visitorID/address", api.ConstRESTOperationCreate, APICreateVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/:visitorID/address/:addressID", api.ConstRESTOperationUpdate, APIUpdateVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/:visitorID/address/:addressID", api.ConstRESTOperationDelete, APIDeleteVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("visitor/:visitorID/addresses", api.ConstRESTOperationGet, APIListVisitorAddresses)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("visitors/addresses/attributes", api.ConstRESTOperationGet, APIListVisitorAddressAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitors/address/:addressID", api.ConstRESTOperationDelete, APIDeleteVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitors/address/:addressID", api.ConstRESTOperationUpdate, APIUpdateVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitors/address/:addressID", api.ConstRESTOperationGet, APIGetVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("visit/address", api.ConstRESTOperationCreate, APICreateVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visit/address/:addressID", api.ConstRESTOperationUpdate, APIUpdateVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visit/address/:addressID", api.ConstRESTOperationDelete, APIDeleteVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visit/addresses", api.ConstRESTOperationGet, APIListVisitorAddresses)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visit/address/:addressID", api.ConstRESTOperationGet, APIGetVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}

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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a9da4ac4-d073-48f3-b062-2ba536d2c577", "visitor id was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		if requestData["visitor_id"] != visitor.GetCurrentVisitorID(context) {
			return nil, env.ErrorDispatch(err)
		}
	}

	// create visitor address operation
	//---------------------------------
	visitorAddressModel, err := visitor.GetVisitorAddressModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

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

// APIUpdateVisitorAddress updates existing visitor address
//   - visitor address id must be specified in "addressID" argument
//   - visitor address attributes should be specified in content
func APIUpdateVisitorAddress(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	addressID := context.GetRequestArgument("addressID")
	if addressID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fe7814c0-85fe-4d60-a134-415f7ac12075", "visitor address 'id' was not specified")
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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "eec1ef1b-25d9-4dbe-8bd2-b907a0897203", "visitor address id was not specified")
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

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

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
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2ac4c16b-9241-406e-b35a-399813bb6ca5", "you are not logined in")
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
	dbCollection.AddStaticFilter("visitor_id", "=", visitorID)

	// filters handle
	models.ApplyFilters(context, dbCollection)

	// checking for a "count" request
	if context.GetRequestArgument("count") != "" {
		return visitorAddressCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	visitorAddressCollectionModel.ListLimit(models.GetListLimit(context))

	// extra parameter handle
	models.ApplyExtraAttributes(context, visitorAddressCollectionModel)

	return visitorAddressCollectionModel.List()
}

// APIGetVisitorAddress returns visitor address information
//   - visitor address id must be specified in "addressID" argument
func APIGetVisitorAddress(context api.InterfaceApplicationContext) (interface{}, error) {
	visitorAddressID := context.GetRequestArgument("addressID")
	if visitorAddressID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b94882c6-bbdd-428d-88b0-7ea5623d80f7", "visitor 'id' was not specified")
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
