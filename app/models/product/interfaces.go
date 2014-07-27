package product

import (
	"github.com/ottemo/foundation/app/models"
)

const (
	PRODUCT_MODEL_NAME = "Product"
)

type I_Product interface {
	GetSku() string
	GetName() string

	GetShortDescription() string
	GetDescription() string

	GetDefaultImage() string

	GetPrice() float64

	GetWeight() float64
	GetSize() float64

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Listable
	models.I_Media

	models.I_CustomAttributes
}
