// Package paypal is a PayPal implementation of payment method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package paypal

import (
	"sync"

	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	// PayPal general constants

	ConstLogStorage = "paypal.log"

	ConstErrorModule = "payment/paypal"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstPaymentActionSale          = "Sale"
	ConstPaymentActionAuthorization = "Authorization"

	// PayPal express constants

	ConstPaymentCode = "paypal_express"
	ConstPaymentName = "PayPal Express"

	ConstConfigPathGroup = "payment.paypal"

	ConstConfigPathEnabled = "payment.paypal.enabled"
	ConstConfigPathTitle   = "payment.paypal.title"

	ConstConfigPathNVP     = "payment.paypal.nvp"
	ConstConfigPathGateway = "payment.paypal.gateway"

	ConstConfigPathUser = "payment.paypal.user"
	ConstConfigPathPass = "payment.paypal.password"

	ConstConfigPathSignature = "payment.paypal.signature"
	ConstConfigPathAction    = "payment.paypal.action"

	ConstConfigStagingGateway = "https://www.sandbox.paypal.com/webscr?cmd=_express-checkout"
	ConstConfigProductionGateway = ""

	// PayPal PayFlow Pro API constants

	ConstPaymentPayPalPayflowCode = "paypal_payflow"
	ConstPaymentPayPalPayflowName = "PayPal Payflow"

	ConstConfigPathPayPalPayflowGroup = "payment.paypalpayflow"

	ConstConfigPathPayPalPayflowEnabled   = "payment.paypalpayflow.enabled"
	ConstConfigPathPayPalPayflowTokenable = "payment.paypalpayflow.tokanable"
	ConstConfigPathPayPalPayflowTitle     = "payment.paypalpayflow.title"

	ConstConfigPathPayPalPayflowURL  = "payment.paypalpayflow.url"
	ConstConfigPathPayPalPayflowHost = "payment.paypalpayflow.host"

	ConstConfigPathPayPalPayflowUser   = "payment.paypalpayflow.user"
	ConstConfigPathPayPalPayflowPass   = "payment.paypalpayflow.password"
	ConstConfigPathPayPalPayflowVendor = "payment.paypalpayflow.vendor"

	ConstConfigPayPalStagingHost = "https://pilot-payflowpro.paypal.com"
	ConstConfigPayPalProductionHost = ""
)

// Package global variables
var (
	waitingTokens      = make(map[string]interface{})
	waitingTokensMutex sync.RWMutex
)

// Express is a implementer of InterfacePaymentMethod for a PayPal Express method
type Express struct{}

// RestAPI is a implementer of InterfacePaymentMethod for a PayPal REST API method (currently not working)
type RestAPI struct{}

// PayFlowAPI is a implementer of PayPal Pro payflow API methods
type PayFlowAPI struct{}
