// Package stock is a default implementation of stock interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package stock

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstCollectionNameStock = "stock"

	ConstConfigPathGroup   = "general.stock"
	ConstConfigPathEnabled = "general.stock.enabled"

	ConstErrorModule = "stock"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultStock is a default implementer of InterfaceStock
type DefaultStock struct{}
