package order

import (
	"github.com/ottemo/foundation/db"
)

// GetDBCollection returns database collection
func (it *DefaultOrderItemCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}
