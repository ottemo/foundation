package visitor

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/visitor"
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
		if err := visitorModel.FromHashMap(recordData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7c4664a6-52eb-419f-adec-df7b7fd146a1", err.Error())
		}

		result = append(result, visitorModel)
	}

	return result
}
