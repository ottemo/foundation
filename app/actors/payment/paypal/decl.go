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

	ConstConfigPathGroup = "payment.paypalExpress"

	ConstConfigPathEnabled = "payment.paypalExpress.enabled"
	ConstConfigPathTitle   = "payment.paypalExpress.title"

	ConstConfigPathUser = "payment.paypalExpress.user"
	ConstConfigPathPass = "payment.paypalExpress.password"

	ConstConfigPathSignature = "payment.paypalExpress.signature"
	ConstConfigPathAction    = "payment.paypalExpress.action"

	ConstConfigPathPayPalExpressGateway = "payment.paypalExpress.gateway"
	ConstConfigPathPayPalPayFlowGateway = "payment.paypalPayflowPro.gateway"
	ConstPaymentPayPalGatewaySandbox    = "sandbox"
	ConstPaymentPayPalGatewayProduction = "production"
	ConstPaymentPayPalGateway           = "gateway"
	ConstPaymentPayPalHost              = "host"
	ConstPaymentPayPalNvp               = "nvp"
	ConstPaymentPayPalUrl               = "url"

	// PayPal PayFlow Pro API constants

	ConstPaymentPayPalPayflowCode = "paypal_payflow"
	ConstPaymentPayPalPayflowName = "PayPal Payflow"

	ConstConfigPathPayPalPayflowGroup = "payment.paypalPayflowPro"

	ConstConfigPathPayPalPayflowEnabled   = "payment.paypalPayflowPro.enabled"
	ConstConfigPathPayPalPayflowTokenable = "payment.paypalPayflowPro.tokanable"
	ConstConfigPathPayPalPayflowTitle     = "payment.paypalPayflowPro.title"

	ConstConfigPathPayPalPayflowUser   = "payment.paypalPayflowPro.user"
	ConstConfigPathPayPalPayflowPass   = "payment.paypalPayflowPro.password"
	ConstConfigPathPayPalPayflowVendor = "payment.paypalPayflowPro.vendor"
)

// Package global variables
var (
	waitingTokens      = make(map[string]interface{})
	waitingTokensMutex sync.RWMutex

	paymentPayPalExpress = map[string]map[string]string{
		ConstPaymentPayPalNvp: {
			ConstPaymentPayPalGatewaySandbox:    "https://api-3t.sandbox.paypal.com/nvp",
			ConstPaymentPayPalGatewayProduction: "https://api-3t.paypal.com/nvp",
		},
		ConstPaymentPayPalGateway: {
			ConstPaymentPayPalGatewaySandbox:    "https://www.sandbox.paypal.com/webscr?cmd=_express-checkout",
			ConstPaymentPayPalGatewayProduction: "https://www.paypal.com/webscr?cmd=_express-checkout",
		},
	}

	paymentPayPalPayFlow = map[string]map[string]string{
		ConstPaymentPayPalUrl: {
			ConstPaymentPayPalGatewaySandbox:    "https://pilot-payflowpro.paypal.com",
			ConstPaymentPayPalGatewayProduction: "https://payflowpro.paypal.com",
		},
		ConstPaymentPayPalHost: {
			ConstPaymentPayPalGatewaySandbox:    "https://pilot-payflowpro.paypal.com",
			ConstPaymentPayPalGatewayProduction: "https://payflowpro.paypal.com",
		},
	}
)

// Express is a implementer of InterfacePaymentMethod for a PayPal Express method
type Express struct{}

// RestAPI is a implementer of InterfacePaymentMethod for a PayPal REST API method (currently not working)
type RestAPI struct{}

// PayFlowAPI is a implementer of PayPal Pro payflow API methods
type PayFlowAPI struct{}
