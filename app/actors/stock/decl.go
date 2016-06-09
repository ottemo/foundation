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
//
// Implementation details:
// 	Lets consider some product in the system. Any product in Ottemo system could have the options,
// 	and options combination could have quantity limit.
//
//		---------------------------------------------------------
//		| _id | ProductID | Options                       | Qty |
//		|-----|-----------|-------------------------------|-----|
//		|  1  |     1     | {}                            |  10 |
//		|  2  |     1     | {"color":"black"}             |   1 |
//		|  3  |     1     | {"color":"blue"}              |   5 |
//		|  4  |     1     | {"size":"l"}                  |   1 |
//		|  5  |     1     | {"size":"l", "color":"black"} |   6 |
//		|  6  |     1     | {"size":"l", "color":"red"}   |   1 |
//		---------------------------------------------------------
//
//	When the application asks for a product Qty, the currently applied options are used. The stock manager
//	searches for a minimal limits in the database which matching specified options, so if spe current product
//	options are set to nothing ({}), it will return the minimal amount which is 0 for a sample product,
//
//	When it searched for a Qty, it going through all the records for a specified product, and trying to match
//	each record option to a specified options.
//
//	So, if the options are set to {"color": "black"} it will take in account records 1, 2 and 5 for a provided
//	sample, as the minimal amount among these records is 1 - it would be the quantity of such product.
// 	If the options would be {"size":"l", "color":"black"}, the examined records would be 1 and 5, and the resulting
//	qty would be 6.
//
//	When the product going to be deleted, it removes all the records.

type DefaultStock struct{}
