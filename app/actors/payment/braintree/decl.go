// Package braintree is a "braintreepayments" implementation of payment method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package braintree

import "github.com/ottemo/foundation/env"

// Package global constants
const (
	// Payment method code used in business logic
	ConstPaymentCode = "braintree"

	// Human readable name of payment method
	ConstPaymentInternalName = "Braintree"

	// Config attribute for User customized name of the payment method
	ConstConfigPathName = "payment.braintree.name"

	ConstConfigPathGroup   = "payment.braintree"
	ConstConfigPathEnabled = "payment.braintree.enabled"
	ConstConfigPathBraintreeEnabled = "payment.braintree.enabled"

	ConstErrorModule = "payment/braintree"

	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstLogStorage = "braintree.log"
)

// Braintree is a implementer of InterfacePaymentMethod for a "braintreepayments" payment method
type BraintreePaymentMethod struct{}


