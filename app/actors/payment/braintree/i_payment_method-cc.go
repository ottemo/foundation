package braintree

import (
	"github.com/lionelbarrow/braintree-go"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
)

// GetCode returns payment method code for use in business logic
func (it *CreditCardMethod) GetCode() string {
	return constCCMethodCode
}

// GetInternalName returns the human readable name of the payment method
func (it *CreditCardMethod) GetInternalName() string {
	return constCCMethodInternalName
}

// GetName returns the user customized name of the payment method
func (it *CreditCardMethod) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstGeneralMethodConfigPathName))
}

// GetType returns type of payment method according to "github.com/ottemo/foundation/app/models/checkout"
func (it *CreditCardMethod) GetType() string {
	return checkout.ConstPaymentTypeCreditCard
}

// IsAllowed checks for payment method applicability
func (it *CreditCardMethod) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstGeneralConfigPathEnabled))
}

// IsTokenable returns possibility to save token for this payment method
func (it *CreditCardMethod) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return (true)
}

// Authorize makes payment method authorize operations
//  - just create token if set in paymentInfo
//  - otherwise create transaction
//  - `orderInstance = nil` when creating a token
func (it *CreditCardMethod) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	action, _ := paymentInfo[checkout.ConstPaymentActionTypeKey]
	creditCardInfo, present := paymentInfo["cc"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0e18570c-e35d-404f-a408-6c9fb4ecfabc", "credit card information has not been set")
	}

	creditCardInfoMap := utils.InterfaceToMap(creditCardInfo)

	if utils.InterfaceToString(action) == checkout.ConstPaymentActionTypeCreateToken {
		visitorInfo, present := paymentInfo["extra"]
		if !present {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "297dfae9-e7c2-41b1-bc54-4f9924d19ae1", "visitor information has not been set")
		}

		creditCardPtr, err := braintreeRegisterCardForVisitor(utils.InterfaceToMap(visitorInfo), creditCardInfoMap)
		if err != nil {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c11c22bf-05ed-432b-b094-0fb0606eb0f1", "unable to create credit card: "+err.Error())
		}

		return braintreeCardToAuthorizeResult(*creditCardPtr, (*creditCardPtr).CustomerId)
	}

	var transactionPtr *braintree.Transaction

	if creditCard, ok := creditCardInfo.(visitor.InterfaceVisitorCard); ok && creditCard != nil {
		var err error
		transactionPtr, err = chargeRegisteredVisitor(orderInstance, creditCard)
		if err != nil {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f7f75a44-1b59-41e7-92fd-d75371f46575", "unable to charge registered visitor: "+err.Error())
		}
	} else {
		var err error
		transactionPtr, err = chargeGuestVisitor(orderInstance, creditCardInfoMap)
		if err != nil {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "df4593b7-4bb0-46f2-a44e-b80488408dc2", "unable to charge guest visitor: "+err.Error())
		}
	}

	return braintreeCardToAuthorizeResult(*transactionPtr.CreditCard, transactionPtr.Customer.Id)
}

// Delete saved card from the payment system.  **This method is for future use**
func (it *CreditCardMethod) DeleteSavedCard(token visitor.InterfaceVisitorCard) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "035052da-95bf-4161-b068-07ebe78c93af", "Not implemented")
}

// Capture makes payment method capture operation
// - at time of implementation this method is not used anywhere
func (it *CreditCardMethod) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "772bc737-f025-4c81-a85a-c10efb67e1b3", " Capture method not implemented")
}

// Refund will return funds on the given order
// - at time of implementation this method is not used anywhere
func (it *CreditCardMethod) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "26febf8b-7e26-44d4-bfb4-e9b29126fe5a", "Refund method not implemented")
}

// Void will mark the order and capture as void
// - at time of implementation this method is not used anywhere
func (it *CreditCardMethod) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "561e0cc8-3bee-4ec4-bf80-585fa566abd4", "Void method not implemented")
}
