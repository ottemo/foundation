package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
)

// GetModelName returns model name
func (it *DefaultCategory) GetModelName() string {
	return category.ConstModelNameCategory
}

// GetImplementationName returns model implementation name
func (it *DefaultCategory) GetImplementationName() string {
	return "Default" + category.ConstModelNameCategory
}

// New returns new instance of model implementation object
func (it *DefaultCategory) New() (models.InterfaceModel, error) {
	return &DefaultCategory{ProductIds: make([]string, 0)}, nil
}
