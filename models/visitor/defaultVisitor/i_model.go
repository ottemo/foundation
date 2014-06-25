package defaultVisitor

import (
	"github.com/ottemo/foundation/models"
)

func (it *DefaultVisitor) GetModelName() string {
	return "Visitor"
}

func (it *DefaultVisitor) GetImplementationName() string {
	return "DefaultVisitor"
}

func (it *DefaultVisitor) New() (models.IModel, error) {
	return &DefaultVisitor{}, nil
}
