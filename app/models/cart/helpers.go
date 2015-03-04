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
func GetCurrentCart(context api.InterfaceApplicationContext) (InterfaceCart, error) {
	sessionCartID := context.GetSession().Get(ConstSessionKeyCurrentCart)
	visitorID := context.GetSession().Get(visitor.ConstSessionKeyVisitorID)

	// checking session for cart id
	if sessionCartID != nil && sessionCartID != "" {
		// cart id was found in session - loading cart by id
		sessionCart, err := LoadCartByID(utils.InterfaceToString(sessionCartID))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		if visitorID != nil && sessionCart.GetVisitorID() == "" {
			visitorCart, err := GetCartForVisitor(utils.InterfaceToString(visitorID))
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			for _, item := range sessionCart.GetItems() {
				visitorCart.AddItem(item.GetProductID(), item.GetQty(), item.GetOptions())
			}

			visitorCart.SetSessionID(context.GetSession().GetID())

			err = visitorCart.Save()
			if err != nil {
				env.ErrorDispatch(err)
			}

			err = sessionCart.Delete()
			if err != nil {
				env.ErrorDispatch(err)
			}

			context.GetSession().Set(ConstSessionKeyCurrentCart, visitorCart.GetID())

			return visitorCart, nil
		}
		return sessionCart, nil

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

	// making new cart for guest if allowed
	if app.ConstAllowGuest {
		currentCart, err := GetCartModel()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		currentCart.SetSessionID(context.GetSession().GetID())
		currentCart.Activate()
		currentCart.Save()

		context.GetSession().Set(ConstSessionKeyCurrentCart, currentCart.GetID())

		return currentCart, nil
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f5acb5ee-f689-4dd8-a85f-f2ec47425ba1", "not registered visitor")
}
