// Package zeropay is a "Zero Payment" implementation of payment method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package zeropay

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstPaymentZeroPaymentCode = "zeropay"
	ConstPaymentName            = "Zero Pay"

	ConstConfigPathGroup   = "payment.zeropay"
	ConstConfigPathEnabled = "payment.zeropay.enabled"
	ConstConfigPathName    = "payment.zeropay.name"

	ConstErrorModule = "payment/zeropay"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// ZeroAmountPayment is a implementer of InterfacePaymentMethod for zero amount payments
type ZeroAmountPayment struct{}
