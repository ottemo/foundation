package address

import (
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
)

// returns database collection
func (it *DefaultVisitorAddressCollection) GetDBCollection() db.I_DBCollection {
	return it.listCollection
}

// returns list of visitor model items
func (it *DefaultVisitorAddressCollection) ListVisitorsAddresses() []visitor.I_VisitorAddress {
	result := make([]visitor.I_VisitorAddress, 0)

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
