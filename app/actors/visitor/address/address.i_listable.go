package address

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// GetCollection returns collection of current instance type
func (it *DefaultVisitorAddress) GetCollection() models.I_Collection {
	model, _ := models.GetModel(visitor.ModelNameVisitorAddressCollection)
	if result, ok := model.(visitor.I_VisitorAddressCollection); ok {
		return result
	}

	return nil
}
