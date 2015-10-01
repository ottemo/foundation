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
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3e5976af-4d48-41d7-a1cf-72696907bce8", "model "+model.GetImplementationName()+" is not 'InterfaceCheckout' capable")
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
func GetCurrentCheckout(context api.InterfaceApplicationContext, bindToSession bool) (InterfaceCheckout, error) {
	sessionObject := context.GetSession().Get(ConstSessionKeyCurrentCheckout)

	var checkoutInstance InterfaceCheckout

	// trying to get checkout object from session, otherwise creating new one
	switch typedvalue := sessionObject.(type) {
	case InterfaceCheckout:
		checkoutInstance = typedvalue

	case map[string]interface{}:
		checkoutModel, err := GetCheckoutModel()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		err = checkoutModel.FromHashMap(typedvalue)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		checkoutInstance = checkoutModel

	default:

		// making new checkout object
		newCheckoutInstance, err := GetCheckoutModel()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		//setting session
		newCheckoutInstance.SetSession(context.GetSession())

		checkoutInstance = newCheckoutInstance
	}

	// storing checkout object to session
	if bindToSession {
		context.GetSession().Set(ConstSessionKeyCurrentCheckout, checkoutInstance)
	}

	// updating checkout object
	//-------------------------

	// setting cart
	currentCart, err := cart.GetCurrentCart(context, false)
	if err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}
	checkoutInstance.SetCart(currentCart)

	// setting visitor
	currentVisitor, err := visitor.GetCurrentVisitor(context)
	if err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}
	checkoutInstance.SetVisitor(currentVisitor)

	return checkoutInstance, nil
}

// SetCurrentCheckout assigns given checkout to current session
func SetCurrentCheckout(context api.InterfaceApplicationContext, checkout InterfaceCheckout) error {
	context.GetSession().Set(ConstSessionKeyCurrentCheckout, checkout)
	return nil
}
