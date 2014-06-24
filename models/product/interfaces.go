package product

import (
	"github.com/ottemo/foundation/models"
)

type Product interface {
	GetSku() string
	GetName() string

	GetPrice() float64

	models.Model
	models.Object
	models.Storable
	models.Mapable

	models.Attribute
}
