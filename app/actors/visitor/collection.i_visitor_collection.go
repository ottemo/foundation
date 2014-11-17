package visitor

import (
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
)

// returns database collection
func (it *DefaultVisitorCollection) GetDBCollection() db.I_DBCollection {
	return it.listCollection
}

// returns list of visitor model items
func (it *DefaultVisitorCollection) ListVisitors() []visitor.I_Visitor {
	result := make([]visitor.I_Visitor, 0)

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		visitorModel, err := visitor.GetVisitorModel()
		if err != nil {
			return result
		}
		visitorModel.FromHashMap(recordData)

		result = append(result, visitorModel)
	}

	return result
}
