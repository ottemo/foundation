// Package braintree is a "braintree payments" implementation of payment method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package braintree

import (
	"github.com/lionelbarrow/braintree-go"

	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	// --------------------------------------
	// Because of multiple payment modules supported by Braintree constant names and values are divided into
	// General - overall values
	// Method  - specific per method values

	// --------------------------------------
	// General

	constGeneralConfigPathGroup = "payment.braintree"
	constGeneralConfigPathEnvironment = "payment.braintree.environment"
	constGeneralConfigPathMerchantID = "payment.braintree.merchantID"
	constGeneralConfigPathPublicKey = "payment.braintree.publicKey"
	constGeneralConfigPathPrivateKey = "payment.braintree.privateKey"

	constEnvironmentSandbox    = braintree.Sandbox
	constEnvironmentProduction = braintree.Production

	constErrorModule = "payment/braintree"
	constErrorLevel  = env.ConstErrorLevelActor

	constLogStorage  = "braintree.log"

	// --------------------------------------
	// Credit Card Method

	constCCMethodConfigPathGroup   = "payment.braintree.cc"
	constCCMethodConfigPathEnabled = "payment.braintree.cc.enabled"
	constCCMethodConfigPathName    = "payment.braintree.cc.name" // User customized name of the payment method

	constCCMethodCode         = "braintreeCC"           // Method code used in business logic
	constCCMethodInternalName = "Braintree Credit Card" // Human readable name of payment method

)

// braintreeCCMethod is a implementer of InterfacePaymentMethod for a Credit Card payment method
type braintreeCCMethod struct{}
