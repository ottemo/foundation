// Package braintree is a "braintreepayments" implementation of payment method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package braintree

import "github.com/ottemo/foundation/env"

// Package global constants
const (
	// Payment method code used in business logic
	ConstPaymentCode = "braintree"
	ConstPaypalPaymentCode = "braintreePaypal"

	// Human readable name of payment method
	ConstPaymentInternalName = "Braintree"
	ConstPaypalPaymentInternalName = "Braintree Paypal"

	// Config attribute for User customized name of the payment method
	ConstConfigPathName = "payment.braintree.name"
	ConstConfigPathPaypalName = "payment.braintree.paypal.name"

	ConstConfigPathGroup   = "payment.braintree"
	ConstConfigPathEnabled = "payment.braintree.enabled"
	ConstConfigPathBraintreeEnabled = "payment.braintree.enabled"
	ConstConfigPathBraintreePaypalEnabled = "payment.braintree.paypal.enabled"

	ConstErrorModule = "payment/braintree"

	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstLogStorage = "braintree.log"
)

// Braintree is a implementer of InterfacePaymentMethod for a "braintreepayments" payment method
type BraintreePaymentMethod struct{}

type BraintreePaypalPaymentMethod struct{}


