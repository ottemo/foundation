package visitor

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// GetVisitorAddressCollectionModel retrieves current InterfaceVisitorAddressCollection model implementation
func GetVisitorAddressCollectionModel() (InterfaceVisitorAddressCollection, error) {
	model, err := models.GetModel(ConstModelNameVisitorAddressCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressCollectionModel, ok := model.(InterfaceVisitorAddressCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "470605da1e564bb2a19d43ed7610e3fd", "model "+model.GetImplementationName()+" is not 'InterfaceVisitorAddressCollection' capable")
	}

	return visitorAddressCollectionModel, nil
}

// GetVisitorAddressModel retrieves current InterfaceVisitorAddress model implementation
func GetVisitorAddressModel() (InterfaceVisitorAddress, error) {
	model, err := models.GetModel(ConstModelNameVisitorAddress)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressModel, ok := model.(InterfaceVisitorAddress)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "26af56affd7e4d32a5c9805873d7d03e", "model "+model.GetImplementationName()+" is not 'InterfaceVisitorAddress' capable")
	}

	return visitorAddressModel, nil
}

// GetVisitorCollectionModel retrieves current InterfaceVisitorCollection model implementation
func GetVisitorCollectionModel() (InterfaceVisitorCollection, error) {
	model, err := models.GetModel(ConstModelNameVisitorCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorCollectionModel, ok := model.(InterfaceVisitorCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "65701c43f2e14950a1001bf92104a701", "model "+model.GetImplementationName()+" is not 'InterfaceVisitorCollection' capable")
	}

	return visitorCollectionModel, nil
}

// GetVisitorModel retrieves current InterfaceVisitor model implementation
func GetVisitorModel() (InterfaceVisitor, error) {
	model, err := models.GetModel(ConstModelNameVisitor)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorModel, ok := model.(InterfaceVisitor)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "04e350af4a734f39b35a95fda724fb93", "model "+model.GetImplementationName()+" is not 'InterfaceVisitor' capable")
	}

	return visitorModel, nil
}

// GetVisitorAddressModelAndSetID retrieves current InterfaceVisitorAddress model implementation and sets its ID to some value
func GetVisitorAddressModelAndSetID(visitorAddressID string) (InterfaceVisitorAddress, error) {

	visitorAddressModel, err := GetVisitorAddressModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorAddressModel.SetID(visitorAddressID)
	if err != nil {
		return visitorAddressModel, env.ErrorDispatch(err)
	}

	return visitorAddressModel, nil
}

// GetVisitorModelAndSetID retrieves current InterfaceVisitor model implementation and sets its ID to some value
func GetVisitorModelAndSetID(visitorID string) (InterfaceVisitor, error) {

	visitorModel, err := GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.SetID(visitorID)
	if err != nil {
		return visitorModel, env.ErrorDispatch(err)
	}

	return visitorModel, nil
}

// LoadVisitorAddressByID loads visitor address data into current InterfaceVisitorAddress model implementation
func LoadVisitorAddressByID(visitorAddressID string) (InterfaceVisitorAddress, error) {

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

// LoadVisitorByID loads visitor data into current InterfaceVisitor model implementation
func LoadVisitorByID(visitorID string) (InterfaceVisitor, error) {

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
func GetCurrentVisitorID(params *api.StructAPIHandlerParams) string {
	sessionVisitorID, ok := params.Session.Get(ConstSessionKeyVisitorID).(string)
	if !ok {
		return ""
	}

	return sessionVisitorID
}

// GetCurrentVisitor returns visitor for current session if registered or nil - for guest visitor
func GetCurrentVisitor(params *api.StructAPIHandlerParams) (InterfaceVisitor, error) {
	sessionVisitorID, ok := params.Session.Get(ConstSessionKeyVisitorID).(string)
	if !ok {
		if app.ConstAllowGuest {
			visitorInstance, err := GetVisitorModel()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
			return visitorInstance, nil
		}

		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f5acb5eef6894dd8a85ff2ec47425ba1", "not registered visitor")
	}

	visitorInstance, err := LoadVisitorByID(sessionVisitorID)

	return visitorInstance, env.ErrorDispatch(err)
}
