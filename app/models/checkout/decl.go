package checkout

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstCheckoutModelName = "Checkout"

	ConstConfigPathGroup             = "general.checkout"
	ConstConfigPathConfirmationEmail = "general.checkout.order_confirmation_email"
	ConstConfigPathCheckoutType      = "general.checkout.checkout_type"
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

	ConstPaymentTypeSimple     = "simple"
	ConstPaymentTypeCreditCard = "cc"
	ConstPaymentTypeRemote     = "remote"
	ConstPaymentTypePost       = "post"
	ConstPaymentTypePostCC     = "post_cc"

	ConstSessionKeyCurrentCheckout = "Checkout"

	ConstErrorModule = "checkout"
	ConstErrorLevel  = env.ConstErrorLevelModel
)
