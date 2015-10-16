package checkout

import (
	"strings"

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

// ValidateAddress makes checkout related address validation
func ValidateAddress(address interface{}) (visitor.InterfaceVisitorAddress, error) {
	var visitorAddress visitor.InterfaceVisitorAddress

	switch typedValue := address.(type) {
	case visitor.InterfaceVisitorAddress:
		visitorAddress = typedValue
	case map[string]interface{}:
		visitorAddressModel, err := visitor.GetVisitorAddressModel()
		if err != nil {
			return visitorAddressModel, err
		}
		visitorAddressModel.FromHashMap(typedValue)
	default:
		return visitorAddress, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f029c930-5f37-4999-9c76-3e739dd2f578", "Unknown address format")
	}

	if visitorAddress.GetAddressLine1() == "" {
		return visitorAddress, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a6039f78-ce9d-409f-bd74-ee00a9a54175", "Address is not specified")
	}

	if visitorAddress.GetCountry() == "" {
		return visitorAddress, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1da2910d-7735-47c8-904b-b4b56359ae4f", "Country is not specified")
	}

	if strings.ToUpper(visitorAddress.GetCountry()) == "US" && visitorAddress.GetState() == "" {
		return visitorAddress, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6dd79207-6d83-4f34-a310-b64bca723548", "State is not specified")
	}

	if visitorAddress.GetFirstName() == "" || visitorAddress.GetLastName() == "" {
		return visitorAddress, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "79ee118a-8663-4c6f-a33c-7234ff162e06", "Name is not specified")
	}

	if visitorAddress.GetZipCode() == "" {
		return visitorAddress, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "83217cce-e2da-4da6-999f-694a5a5aee97", "Zip code is not specified")
	}

	return visitorAddress, nil
}
