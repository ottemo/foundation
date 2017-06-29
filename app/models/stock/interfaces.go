// Package product represents abstraction of business layer product object
package stock

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstModelNameStock           = "Stock"
	ConstModelNameStockCollection = "StockCollection"

	ConstErrorModule = "stock"
	ConstErrorLevel  = env.ConstErrorLevelModel

	//ConstOptionProductIDs = "_ids"
)

// InterfaceStock represents interface to access business layer implementation of stock management
type InterfaceStock interface {

	GetProductID() string
	GetOptions() string
	GetQty() int

	SetProductID(product_id string) error
	SetOptions(options string) error
	SetQty(qty int) error

	SetProductQty(productID string, options map[string]interface{}, qty int) error
	GetProductQty(productID string, options map[string]interface{}) int
	GetProductOptions(productID string) []map[string]interface{}

	RemoveProductQty(productID string, options map[string]interface{}) error
	UpdateProductQty(productID string, options map[string]interface{}, deltaQty int) error

	models.InterfaceModel
	models.InterfaceObject
	models.InterfaceStorable
	models.InterfaceListable
}

// InterfaceStockCollection represents interface to access business layer implementation of stock collection
type InterfaceStockCollection interface {
	ListStocks() []InterfaceStock

	models.InterfaceCollection
}
