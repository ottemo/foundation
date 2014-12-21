package visitor

import (
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
)

// GetDBCollection returns database collection for the Visitor
func (it *DefaultVisitorCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// ListVisitors returns list of visitor model items in the Visitor Collection
func (it *DefaultVisitorCollection) ListVisitors() []visitor.InterfaceVisitor {
	var result []visitor.InterfaceVisitor

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
