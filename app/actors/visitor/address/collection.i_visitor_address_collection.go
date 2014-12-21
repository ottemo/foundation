package address

import (
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
)

// GetDBCollection returns the database collection of the Visitor Address
func (it *DefaultVisitorAddressCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// ListVisitorsAddresses returns list of visitor model items for the Visitor Address
func (it *DefaultVisitorAddressCollection) ListVisitorsAddresses() []visitor.InterfaceVisitorAddress {
	var result []visitor.InterfaceVisitorAddress

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		visitorAddressModel, err := visitor.GetVisitorAddressModel()
		if err != nil {
			return result
		}
		visitorAddressModel.FromHashMap(recordData)

		result = append(result, visitorAddressModel)
	}

	return result
}
