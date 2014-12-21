package checkout

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
)

// GetCheckoutModel retrieves current InterfaceCheckout model implementation
func GetCheckoutModel() (InterfaceCheckout, error) {
	model, err := models.GetModel(ConstCheckoutModelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	checkoutModel, ok := model.(InterfaceCheckout)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3e5976af4d4841d7a1cf72696907bce8", "model "+model.GetImplementationName()+" is not 'InterfaceCheckout' capable")
	}

	return checkoutModel, nil
}

// GetShippingMethodByCode retrieves shipping method for given unique code or nil if no shipping method with such code
func GetShippingMethodByCode(code string) InterfaceShippingMethod {

	for _, shippingMethod := range registeredShippingMethods {
		if shippingMethod.GetCode() == code {
			return shippingMethod
		}
	}

	return nil
}

// GetPaymentMethodByCode retrieves payment method for given unique code or nil if no payment method with such code
func GetPaymentMethodByCode(code string) InterfacePaymentMethod {

	for _, paymentMethod := range registeredPaymentMethods {
		if paymentMethod.GetCode() == code {
			return paymentMethod
		}
	}

	return nil
}

// GetCurrentCheckout returns checkout for current session or creates new one
func GetCurrentCheckout(params *api.StructAPIHandlerParams) (InterfaceCheckout, error) {
	sessionObject := params.Session.Get(ConstSessionKeyCurrentCheckout)

	var checkoutInstance InterfaceCheckout

	// trying to get checkout object from session, otherwise creating new one
	if sessionCheckout, ok := sessionObject.(InterfaceCheckout); ok {
		checkoutInstance = sessionCheckout

	} else {

		// making new checkout object
		newCheckoutInstance, err := GetCheckoutModel()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// storing checkout object to session
		params.Session.Set(ConstSessionKeyCurrentCheckout, newCheckoutInstance)

		//setting session
		newCheckoutInstance.SetSession(params.Session)

		checkoutInstance = newCheckoutInstance
	}

	// updating checkout object
	//-------------------------

	// setting cart
	currentCart, err := cart.GetCurrentCart(params)
	if err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}
	checkoutInstance.SetCart(currentCart)

	// setting visitor
	currentVisitor, err := visitor.GetCurrentVisitor(params)
	if err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}
	checkoutInstance.SetVisitor(currentVisitor)

	return checkoutInstance, nil
}
