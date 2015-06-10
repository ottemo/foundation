package seo

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/seo"
)

// GetModelName returns model name
func (it *DefaultSEOItem) GetModelName() string {
	return seo.ConstModelNameSEOItem
}

// GetImplementationName returns model implementation name
func (it *DefaultSEOItem) GetImplementationName() string {
	return "DefaultSEOItem"
}

// New returns new instance of model implementation object
func (it *DefaultSEOItem) New() (models.InterfaceModel, error) {
	return &DefaultSEOItem{}, nil
}
