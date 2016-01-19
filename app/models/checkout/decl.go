package checkout

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstCheckoutModelName = "Checkout"

	ConstConfigPathGroup             = "general.checkout"
	ConstConfigPathConfirmationEmail = "general.checkout.order_confirmation_email"
	ConstConfigPathOversell          = "general.checkout.oversell"

	ConstConfigPathShippingGroup              = "shipping"
	ConstConfigPathShippingOriginGroup        = "shipping.origin"
	ConstConfigPathShippingOriginCountry      = "shipping.origin.country"
	ConstConfigPathShippingOriginState        = "shipping.origin.state"
	ConstConfigPathShippingOriginCity         = "shipping.origin.city"
	ConstConfigPathShippingOriginAddressline1 = "shipping.origin.addressline1"
	ConstConfigPathShippingOriginAddressline2 = "shipping.origin.addressline2"
	ConstConfigPathShippingOriginZip          = "shipping.origin.zip"

	ConstConfigPathPaymentGroup              = "payment"
	ConstConfigPathPaymentOriginGroup        = "payment.origin"
	ConstConfigPathPaymentOriginCountry      = "payment.origin.country"
	ConstConfigPathPaymentOriginState        = "payment.origin.state"
	ConstConfigPathPaymentOriginCity         = "payment.origin.city"
	ConstConfigPathPaymentOriginAddressline1 = "payment.origin.addressline1"
	ConstConfigPathPaymentOriginAddressline2 = "payment.origin.addressline2"
	ConstConfigPathPaymentOriginZip          = "payment.origin.zip"

	ConstPaymentActionTypeKey         = "actionType"
	ConstPaymentActionTypeCreateToken = "createToken"
	ConstPaymentActionTypeUseToken    = "useToken"

	ConstPaymentTypeSimple     = "simple"
	ConstPaymentTypeCreditCard = "cc"
	ConstPaymentTypeRemote     = "remote"
	ConstPaymentTypePost       = "post"
	ConstPaymentTypePostCC     = "post_cc"

	ConstPaymentErrorDeclined  = "Payment declined."
	ConstPaymentErrorTechnical = "Technical error."

	ConstCalculateTargetSubtotal   = 1.0
	ConstCalculateTargetShipping   = 2.0
	ConstCalculateTargetGrandTotal = 3.0

	ConstDiscountObjectCart = "cart"

	ConstSessionKeyCurrentCheckout = "Checkout"

	ConstErrorModule = "checkout"
	ConstErrorLevel  = env.ConstErrorLevelModel
)
