package testDiscount

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstSessionKeyAppliedTestDiscountCodes = "applied_test_discount_codes"

	ConstConfigPathTestDiscounts			= "general.testdiscounts"
	ConstConfigPathTestDiscountRule  		= "general.testdiscounts.testDiscount_rule"
	ConstConfigPathTestDiscountAction       = "general.testdiscounts.testDiscount_action"

	ConstErrorModule = "testDiscount"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultTestDiscount is a default implementer of InterfaceTestDiscount
type DefaultTestDiscount struct{}

