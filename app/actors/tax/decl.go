// Package tax is a implementation of tax interface declared in
// "github.com/ottemo/commerce/app/models/checkout" package
package tax

import (
	"github.com/ottemo/commerce/env"
)

// Package global constants
const (
	ConstErrorModule = "tax"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstPriorityValue = 2.50

	ConstProductTaxableAttribute = "taxable"
)

var priority float64

// DefaultTax is a default implementer of InterfaceTax
type DefaultTax struct{}
