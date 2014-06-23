package product

import (
	"github.com/ottemo/foundation/models"
)

type I_Product interface {
	GetSku() string
	GetName() string

	GetPrice() float64

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Mapable

	models.I_CustomAttributes
}
