// Package authorizenet is a Authorize.Net implementation of payment method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package authorizenet

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	// ConstTransactionApproved constants for transaction state
	ConstTransactionApproved      = "1"
	ConstTransactionDeclined      = "2"
	ConstTransactionError         = "3"
	ConstTransactionWaitingReview = "4"

	ConstLogStorage = "authorizenet.log"

	ConstPaymentCodeDPM = "authorizeNetDPM"
	ConstPaymentNameDPM = "Authorize.Net (Direct Post)"

	ConstDPMActionAuthorizeOnly       = "AUTH_ONLY"
	ConstDPMActionAuthorizeAndCapture = "AUTH_CAPTURE"

	ConstConfigPathDPMGroup = "payment.authorizeNetDPM"

	ConstConfigPathDPMEnabled = "payment.authorizeNetDPM.enabled"
	ConstConfigPathDPMAction  = "payment.authorizeNetDPM.action"
	ConstConfigPathDPMTitle   = "payment.authorizeNetDPM.title"

	ConstConfigPathDPMLogin   = "payment.authorizeNetDPM.login"
	ConstConfigPathDPMKey     = "payment.authorizeNetDPM.key"
	ConstConfigPathDPMGateway = "payment.authorizeNetDPM.gateway"

	ConstConfigPathDPMTest     = "payment.authorizeNetDPM.test"
	ConstConfigPathDPMDebug    = "payment.authorizeNetDPM.debug"
	ConstConfigPathDPMCheckout = "payment.authorizeNetDPM.checkout"

	ConstErrorModule = "payment/authorizenet"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DirectPostMethod is a implementer of InterfacePaymentMethod for a Authorize.Net Direct Post method
type DirectPostMethod struct{}
