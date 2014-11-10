package address

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/visitor"
)

// REST API registration function
func setupAPI() error {
	err := api.GetRestService().RegisterAPI("visitor/address", "POST", "create", restCreateVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/address", "PUT", "update/:id", restUpdateVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/address", "DELETE", "delete/:id", restDeleteVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("visitor/address", "GET", "attribute/list", restListVisitorAddressAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/address", "GET", "list", restListVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/address", "POST", "list", restListVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/address", "GET", "count", restCountVisitorAddresses)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/address", "GET", "list/:visitorId", restListVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/address", "POST", "list/:visitorId", restListVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visitor/address", "GET", "load/:id", restGetVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API used to create new visitor address
//   - visitor address attributes must be included in POST form
//   - visitor id required
func restCreateVisitorAddress(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if _, ok := reqData["visitor_id"]; !ok {
		return nil, env.ErrorNew("visitor id was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		if reqData["visitor_id"] != visitor.GetCurrentVisitorID(params) {
			return nil, env.ErrorDispatch(err)
		}
	}

	// create visitor address operation
	//---------------------------------
	visitorAddressModel, err := visitor.GetVisitorAddressModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
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

// WEB REST API used to update existing visitor address
//   - visitor address id must be specified in request URI
//   - visitor address attributes must be included in POST form
func restUpdateVisitorAddress(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	addressID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew("visitor address 'id' was not specified")
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressModel, err := visitor.LoadVisitorAddressByID(addressID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		if visitorAddressModel.GetVisitorID() != visitor.GetCurrentVisitorID(params) {
			return nil, env.ErrorDispatch(err)
		}
	}

	// update operation
	//-----------------
	for attribute, value := range reqData {
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

// WEB REST API used to delete visitor address
//   - visitor address attributes must be included in POST form
func restDeleteVisitorAddress(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//--------------------
	addressID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew("visitor address id was not specified")
	}

	visitorAddressModel, err := visitor.LoadVisitorAddressByID(addressID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		if visitorAddressModel.GetVisitorID() != visitor.GetCurrentVisitorID(params) {
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

// WEB REST API function used to obtain visitor address attributes information
func restListVisitorAddressAttributes(params *api.T_APIHandlerParams) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressModel, err := visitor.GetVisitorAddressModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attrInfo := visitorAddressModel.GetAttributesInfo()
	return attrInfo, nil
}

// WEB REST API function used to obtain visitors addresses count in model collection
func restCountVisitorAddresses(params *api.T_APIHandlerParams) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressCollectionModel, err := visitor.GetVisitorAddressCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := visitorAddressCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(params, dbCollection)

	return dbCollection.Count()
}

// WEB REST API function used to obtain visitor addresses list
//   - visitor id must be specified in request URI
func restListVisitorAddress(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		if params.Request.Method == "POST" {
			return nil, env.ErrorNew("unexpected request content")
		}

		reqData = make(map[string]interface{})
	}

	visitorID, isSpecifiedID := params.RequestURLParams["visitorId"]
	if !isSpecifiedID {

		sessionVisitorID := visitor.GetCurrentVisitorID(params)
		if sessionVisitorID == "" {
			return nil, env.ErrorNew("you are not logined in")
		}
		visitorID = sessionVisitorID
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		if visitorID != visitor.GetCurrentVisitorID(params) {
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

	// limit parameter handle
	visitorAddressCollectionModel.ListLimit(api.GetListLimit(params))

	// filters handle
	api.ApplyFilters(params, dbCollection)

	// extra parameter handle
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := visitorAddressCollectionModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	return visitorAddressCollectionModel.List()
}

// WEB REST API used to get visitor address object
//   - visitor address id must be specified in request URI
func restGetVisitorAddress(params *api.T_APIHandlerParams) (interface{}, error) {
	visitorAddressID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew("visitor 'id' was not specified")
	}

	visitorAddressModel, err := visitor.LoadVisitorAddressByID(visitorAddressID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		if visitorAddressModel.GetVisitorID() != visitor.GetCurrentVisitorID(params) {
			return nil, env.ErrorDispatch(err)
		}
	}

	return visitorAddressModel.ToHashMap(), nil
}
