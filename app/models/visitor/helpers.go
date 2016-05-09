package visitor

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetVisitorAddressCollectionModel retrieves current InterfaceVisitorAddressCollection model implementation
func GetVisitorAddressCollectionModel() (InterfaceVisitorAddressCollection, error) {
	model, err := models.GetModel(ConstModelNameVisitorAddressCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressCollectionModel, ok := model.(InterfaceVisitorAddressCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "470605da-1e56-4bb2-a19d-43ed7610e3fd", "model "+model.GetImplementationName()+" is not 'InterfaceVisitorAddressCollection' capable")
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
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "26af56af-fd7e-4d32-a5c9-805873d7d03e", "model "+model.GetImplementationName()+" is not 'InterfaceVisitorAddress' capable")
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
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "65701c43-f2e1-4950-a100-1bf92104a701", "model "+model.GetImplementationName()+" is not 'InterfaceVisitorCollection' capable")
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
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "04e350af-4a73-4f39-b35a-95fda724fb93", "model "+model.GetImplementationName()+" is not 'InterfaceVisitor' capable")
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
func GetCurrentVisitorID(context api.InterfaceApplicationContext) string {
	sessionVisitorID, ok := context.GetSession().Get(ConstSessionKeyVisitorID).(string)
	if !ok {
		return ""
	}

	return sessionVisitorID
}

// GetCurrentVisitor returns visitor for current session if registered or nil - for guest visitor
func GetCurrentVisitor(context api.InterfaceApplicationContext) (InterfaceVisitor, error) {
	sessionVisitorID := utils.InterfaceToString(context.GetSession().Get(ConstSessionKeyVisitorID))
	if sessionVisitorID == "" {
		if app.ConstAllowGuest {
			visitorInstance, err := GetVisitorModel()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
			return visitorInstance, nil
		}

		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f5acb5ee-f689-4dd8-a85f-f2ec47425ba1", "not registered visitor")
	}

	visitorInstance, err := LoadVisitorByID(sessionVisitorID)

	return visitorInstance, env.ErrorDispatch(err)
}

// GetVisitorCardModel retrieves current InterfaceVisitorCard model implementation
func GetVisitorCardModel() (InterfaceVisitorCard, error) {
	model, err := models.GetModel(ConstModelNameVisitorCard)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorCardModel, ok := model.(InterfaceVisitorCard)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f1d8b09e-0936-46c6-a5e5-6a6df4e462f1", "model "+model.GetImplementationName()+" is not 'InterfaceVisitorCard' capable")
	}

	return visitorCardModel, nil
}

// LoadVisitorCardByID loads visitor address data into current InterfaceVisitorCard model implementation
func LoadVisitorCardByID(visitorCardID string) (InterfaceVisitorCard, error) {

	visitorCardModel, err := GetVisitorCardModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorCardModel.Load(visitorCardID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorCardModel, nil
}

// LoadVisitorCardByVID returns a list of cards belonging to the visitor
// includes the `customer_id` field, which is a customer token
// that can be used to charge the default card associated with the customer
// in stripe
func LoadVisitorCardByVID(vid string) []models.StructListItem {
	model, _ := GetVisitorCardCollectionModel()
	model.ListFilterAdd("visitor_id", "=", vid)

	// 3rd party customer identifier, used by stripe
	err := model.ListAddExtraAttribute("customer_id")
	if err != nil {
		env.ErrorDispatch(err)
	}

	resp, err := model.List()
	if err != nil {
		env.ErrorDispatch(err)
	}

	return resp
}

// GetVisitorCardModelAndSetID retrieves current InterfaceVisitorCard model implementation and sets its ID
func GetVisitorCardModelAndSetID(visitorCardID string) (InterfaceVisitorCard, error) {

	visitorCardModel, err := GetVisitorCardModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorCardModel.SetID(visitorCardID)
	if err != nil {
		return visitorCardModel, env.ErrorDispatch(err)
	}

	return visitorCardModel, nil
}

// GetVisitorCardCollectionModel retrieves current InterfaceVisitorAddressCollection model implementation
func GetVisitorCardCollectionModel() (InterfaceVisitorCardCollection, error) {
	model, err := models.GetModel(ConstModelNameVisitorCardCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorCardCollectionModel, ok := model.(InterfaceVisitorCardCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d978df3c-b908-41a4-b72c-da012ea93bad", "model "+model.GetImplementationName()+" is not 'InterfaceVisitorCardCollection' capable")
	}

	return visitorCardCollectionModel, nil
}
