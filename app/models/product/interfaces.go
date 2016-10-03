// Package product represents abstraction of business layer product object
package product

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstModelNameProduct           = "Product"
	ConstModelNameProductCollection = "ProductCollection"

	ConstErrorModule = "product"
	ConstErrorLevel  = env.ConstErrorLevelModel

	// simple products for option key name
	ConstOptionSimpleIDsName = "_pids"
)

// InterfaceProduct represents interface to access business layer implementation of product object
type InterfaceProduct interface {
	GetEnabled() bool

	GetSku() string
	GetName() string

	GetShortDescription() string
	GetDescription() string

	GetDefaultImage() string

	GetPrice() float64
	GetWeight() float64

	GetAppliedOptions() map[string]interface{}
	GetOptions() map[string]interface{}

	ApplyOptions(map[string]interface{}) error

	models.InterfaceModel
	models.InterfaceObject
	models.InterfaceStorable
	models.InterfaceMedia
	models.InterfaceListable
	models.InterfaceCustomAttributes
	models.InterfaceExternalAttributes
}

// InterfaceProductCollection represents interface to access business layer implementation of product collection
type InterfaceProductCollection interface {
	ListProducts() []InterfaceProduct

	models.InterfaceCollection
}

// InterfaceStock represents interface to access business layer implementation of stock management
type InterfaceStock interface {
	SetProductQty(productID string, options map[string]interface{}, qty int) error
	GetProductQty(productID string, options map[string]interface{}) int
	GetProductOptions(productID string) []map[string]interface{}

	RemoveProductQty(productID string, options map[string]interface{}) error
	UpdateProductQty(productID string, options map[string]interface{}, deltaQty int) error
}
