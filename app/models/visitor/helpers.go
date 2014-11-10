package visitor

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// GetVisitorAddressCollectionModel retrieves current I_VisitorAddressCollection model implementation
func GetVisitorAddressCollectionModel() (I_VisitorAddressCollection, error) {
	model, err := models.GetModel(ModelNameVisitorAddressCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressCollectionModel, ok := model.(I_VisitorAddressCollection)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_VisitorAddressCollection' capable")
	}

	return visitorAddressCollectionModel, nil
}

// GetVisitorAddressModel retrieves current I_VisitorAddress model implementation
func GetVisitorAddressModel() (I_VisitorAddress, error) {
	model, err := models.GetModel(ModelNameVisitorAddress)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressModel, ok := model.(I_VisitorAddress)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_VisitorAddress' capable")
	}

	return visitorAddressModel, nil
}

// GetVisitorCollectionModel retrieves current I_VisitorCollection model implementation
func GetVisitorCollectionModel() (I_VisitorCollection, error) {
	model, err := models.GetModel(ModelNameVisitorCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorCollectionModel, ok := model.(I_VisitorCollection)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_VisitorCollection' capable")
	}

	return visitorCollectionModel, nil
}

// GetVisitorModel retrieves current I_Visitor model implementation
func GetVisitorModel() (I_Visitor, error) {
	model, err := models.GetModel(ModelNameVisitor)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorModel, ok := model.(I_Visitor)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_Visitor' capable")
	}

	return visitorModel, nil
}

// GetVisitorAddressModelAndSetID retrieves current I_VisitorAddress model implementation and sets its ID to some value
func GetVisitorAddressModelAndSetID(visitorAddressID string) (I_VisitorAddress, error) {

	visitorAddressModel, err := GetVisitorAddressModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorAddressModel.SetId(visitorAddressID)
	if err != nil {
		return visitorAddressModel, env.ErrorDispatch(err)
	}

	return visitorAddressModel, nil
}

// GetVisitorModelAndSetID retrieves current I_Visitor model implementation and sets its ID to some value
func GetVisitorModelAndSetID(visitorID string) (I_Visitor, error) {

	visitorModel, err := GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.SetId(visitorID)
	if err != nil {
		return visitorModel, env.ErrorDispatch(err)
	}

	return visitorModel, nil
}

// LoadVisitorAddressByID loads visitor address data into current I_VisitorAddress model implementation
func LoadVisitorAddressByID(visitorAddressID string) (I_VisitorAddress, error) {

	visitorAddressModel, err := GetVisitorAddressModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorAddressModel.Load(visitorAddressID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorAddressModel, nil
}

// LoadVisitorByID loads visitor data into current I_Visitor model implementation
func LoadVisitorByID(visitorID string) (I_Visitor, error) {

	visitorModel, err := GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.Load(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorModel, nil
}

// GetCurrentVisitorID returns visitor id for current session if registered or ""
func GetCurrentVisitorID(params *api.T_APIHandlerParams) string {
	sessionVisitorID, ok := params.Session.Get(SessionKeyVisitorID).(string)
	if !ok {
		return ""
	}

	return sessionVisitorID
}

// GetCurrentVisitor returns visitor for current session if registered or error
func GetCurrentVisitor(params *api.T_APIHandlerParams) (I_Visitor, error) {
	sessionVisitorID, ok := params.Session.Get(SessionKeyVisitorID).(string)
	if !ok {
		return nil, env.ErrorNew("not registered visitor")
	}

	visitorInstance, err := LoadVisitorByID(sessionVisitorID)

	return visitorInstance, env.ErrorDispatch(err)
}
