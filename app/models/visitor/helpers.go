package visitor

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// retrieves current I_VisitorAddressCollection model implementation
func GetVisitorAddressCollectionModel() (I_VisitorAddressCollection, error) {
	model, err := models.GetModel(MODEL_NAME_VISITOR_ADDRESS_COLLECTION)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressCollectionModel, ok := model.(I_VisitorAddressCollection)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_VisitorAddressCollection' capable")
	}

	return visitorAddressCollectionModel, nil
}

// retrieves current I_VisitorAddress model implementation
func GetVisitorAddressModel() (I_VisitorAddress, error) {
	model, err := models.GetModel(MODEL_NAME_VISITOR_ADDRESS)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressModel, ok := model.(I_VisitorAddress)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_VisitorAddress' capable")
	}

	return visitorAddressModel, nil
}

// retrieves current I_VisitorCollection model implementation
func GetVisitorCollectionModel() (I_VisitorCollection, error) {
	model, err := models.GetModel(MODEL_NAME_VISITOR_COLLECTION)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorCollectionModel, ok := model.(I_VisitorCollection)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_VisitorCollection' capable")
	}

	return visitorCollectionModel, nil
}

// retrieves current I_Visitor model implementation
func GetVisitorModel() (I_Visitor, error) {
	model, err := models.GetModel(MODEL_NAME_VISITOR)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorModel, ok := model.(I_Visitor)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_Visitor' capable")
	}

	return visitorModel, nil
}

// retrieves current I_VisitorAddress model implementation and sets its ID to some value
func GetVisitorAddressModelAndSetId(visitorAddressId string) (I_VisitorAddress, error) {

	visitorAddressModel, err := GetVisitorAddressModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorAddressModel.SetId(visitorAddressId)
	if err != nil {
		return visitorAddressModel, env.ErrorDispatch(err)
	}

	return visitorAddressModel, nil
}

// retrieves current I_Visitor model implementation and sets its ID to some value
func GetVisitorModelAndSetId(visitorId string) (I_Visitor, error) {

	visitorModel, err := GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.SetId(visitorId)
	if err != nil {
		return visitorModel, env.ErrorDispatch(err)
	}

	return visitorModel, nil
}

// loads visitor address data into current I_VisitorAddress model implementation
func LoadVisitorAddressById(visitorAddressId string) (I_VisitorAddress, error) {

	visitorAddressModel, err := GetVisitorAddressModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorAddressModel.Load(visitorAddressId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorAddressModel, nil
}

// loads visitor data into current I_Visitor model implementation
func LoadVisitorById(visitorId string) (I_Visitor, error) {

	visitorModel, err := GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.Load(visitorId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorModel, nil
}

// returns visitor id for current session if registered or ""
func GetCurrentVisitorId(params *api.T_APIHandlerParams) string {
	sessionVisitorId, ok := params.Session.Get(SESSION_KEY_VISITOR_ID).(string)
	if !ok {
		return ""
	}

	return sessionVisitorId
}

// returns visitor for current session if registered or error
func GetCurrentVisitor(params *api.T_APIHandlerParams) (I_Visitor, error) {
	sessionVisitorId, ok := params.Session.Get(SESSION_KEY_VISITOR_ID).(string)
	if !ok {
		return nil, env.ErrorNew("not registered visitor")
	}

	visitorInstance, err := LoadVisitorById(sessionVisitorId)

	return visitorInstance, env.ErrorDispatch(err)
}
