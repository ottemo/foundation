// Package tax is a implementation of tax interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package tax

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstErrorModule = "tax"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstPriorityValue = 2.50

	ConstConfigPathGroup = "tax.default"
	ConstConfigPathEnabled = "tax.default.enabled"
)

var priority float64

// DefaultTax is a default implementer of InterfaceTax
type DefaultTax struct{}
