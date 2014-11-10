package cart

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetCartModel retrieves current I_Cart model implementation
func GetCartModel() (I_Cart, error) {
	model, err := models.GetModel(CART_MODEL_NAME)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cartModel, ok := model.(I_Cart)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_Cart' capable")
	}

	return cartModel, nil
}

// GetCartModelAndSetId retrieves current I_Cart model implementation and sets its ID to some value
func GetCartModelAndSetId(cartId string) (I_Cart, error) {

	cartModel, err := GetCartModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cartModel.SetId(cartId)
	if err != nil {
		return cartModel, env.ErrorDispatch(err)
	}

	return cartModel, nil
}

// LoadCartById loads cart data into current I_Cart model implementation
func LoadCartById(cartId string) (I_Cart, error) {

	cartModel, err := GetCartModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cartModel.Load(cartId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cartModel, nil
}

// GetCartForVisitor loads cart for visitor or creates new one
func GetCartForVisitor(visitorId string) (I_Cart, error) {
	cartModel, err := GetCartModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cartModel.MakeCartForVisitor(visitorId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cartModel, nil
}

// returns cart for current session or creates new one
func GetCurrentCart(params *api.T_APIHandlerParams) (I_Cart, error) {
	sessionCartId := params.Session.Get(SESSION_KEY_CURRENT_CART)

	if sessionCartId != nil && sessionCartId != "" {

		// cart id was found in session - loading cart by id
		currentCart, err := LoadCartById(utils.InterfaceToString(sessionCartId))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		return currentCart, nil

	} else {

		// no cart id was in session, trying to get cart for visitor
		visitorId := params.Session.Get(visitor.SessionKeyVisitorID)
		if visitorId != nil {
			currentCart, err := GetCartForVisitor(utils.InterfaceToString(visitorId))
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			params.Session.Set(SESSION_KEY_CURRENT_CART, currentCart.GetId())

			return currentCart, nil
		} else {
			return nil, env.ErrorNew("you are not registered")
		}

	}
}
