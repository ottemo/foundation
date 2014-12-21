package block

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
)

// GetModelName returns model name
func (it *DefaultCMSBlock) GetModelName() string {
	return cms.ConstModelNameCMSBlock
}

// GetImplementationName returns model implementation name
func (it *DefaultCMSBlock) GetImplementationName() string {
	return "DefaultCMSBlock"
}

// New returns new instance of model implementation object
func (it *DefaultCMSBlock) New() (models.InterfaceModel, error) {
	return &DefaultCMSBlock{}, nil
}
