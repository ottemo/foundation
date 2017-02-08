package cart

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetCartModel retrieves current InterfaceCart model implementation
func GetCartModel() (InterfaceCart, error) {
	model, err := models.GetModel(ConstCartModelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cartModel, ok := model.(InterfaceCart)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "369c2954-4a20-4958-b929-036a4bdf944d", "model "+model.GetImplementationName()+" is not 'InterfaceCart' capable")
	}

	return cartModel, nil
}

// GetCartModelAndSetID retrieves current InterfaceCart model implementation and sets its ID to some value
func GetCartModelAndSetID(cartID string) (InterfaceCart, error) {

	cartModel, err := GetCartModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cartModel.SetID(cartID)
	if err != nil {
		return cartModel, env.ErrorDispatch(err)
	}

	return cartModel, nil
}

// LoadCartByID loads cart data into current InterfaceCart model implementation
func LoadCartByID(cartID string) (InterfaceCart, error) {

	cartModel, err := GetCartModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cartModel.Load(cartID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cartModel, nil
}

// GetCartForVisitor loads cart for visitor or creates new one
func GetCartForVisitor(visitorID string) (InterfaceCart, error) {
	cartModel, err := GetCartModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cartModel.MakeCartForVisitor(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cartModel, nil
}

// GetCurrentCart returns cart for current session or creates new one
func GetCurrentCart(context api.InterfaceApplicationContext, createNew bool) (InterfaceCart, error) {
	sessionCartID := context.GetSession().Get(ConstSessionKeyCurrentCart)
	visitorID := context.GetSession().Get(visitor.ConstSessionKeyVisitorID)

	// checking session for cart id
	if sessionCartID != nil && sessionCartID != "" {
		// cart id was found in session - loading cart by id
		sessionCart, err := LoadCartByID(utils.InterfaceToString(sessionCartID))
		if err == nil && sessionCart != nil {

			if visitorID != nil && sessionCart.GetVisitorID() == "" {
				visitorCart, err := GetCartForVisitor(utils.InterfaceToString(visitorID))
				if err == nil && visitorCart != nil {
					for _, item := range sessionCart.GetItems() {
						if _, err := visitorCart.AddItem(item.GetProductID(), item.GetQty(), item.GetOptions()); err != nil {
							return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "971c7278-0db5-4101-a5a7-b12c096dfa83", "unable to add item to visitor cart: "+err.Error())
						}
					}

					if err := visitorCart.SetSessionID(context.GetSession().GetID()); err != nil {
						return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "de31bc21-4073-47be-b773-e409f9f4ed55", "unable to set session for visitor cart: "+err.Error())
					}

					if err := visitorCart.Save(); err != nil {
						return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6a9da8ec-f1a8-401a-9dec-209220f642af", "unable to save visitor cart: "+err.Error())
					}

					if err := sessionCart.Delete(); err != nil {
						return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "8dfd01f9-2ffb-4347-84dc-ad3c3c67d0fd", "unable to delete session cart: "+err.Error())
					}

					context.GetSession().Set(ConstSessionKeyCurrentCart, visitorCart.GetID())

					return visitorCart, nil
				}
			}

			return sessionCart, nil
		}
	}

	// no cart id was in session, trying to get cart for visitor
	if visitorID != nil {
		currentCart, err := GetCartForVisitor(utils.InterfaceToString(visitorID))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		context.GetSession().Set(ConstSessionKeyCurrentCart, currentCart.GetID())

		return currentCart, nil
	}

	if createNew {

		// making new cart for guest if allowed
		if app.ConstAllowGuest {

			currentCart, err := GetCartModel()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
			if err := currentCart.SetSessionID(context.GetSession().GetID()); err != nil {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c5f6b9b8-8cf0-4c6e-bf9f-3900d84ab91c", "unable to set cart session: "+err.Error())
			}
			if err := currentCart.Activate(); err != nil {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1271856b-dd11-4b35-8d71-8676dec7e02f", "unable to activate cart: "+err.Error())
			}
			if err := currentCart.Save(); err != nil {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a8225952-291b-46e3-9812-c142552da2b8", "unable to save cart: "+err.Error())
			}

			context.GetSession().Set(ConstSessionKeyCurrentCart, currentCart.GetID())

			return currentCart, nil
		}

		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f5acb5ee-f689-4dd8-a85f-f2ec47425ba1", "not registered visitor")
	}

	return nil, nil
}
