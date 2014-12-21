package page

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetModelName returns model name
func (it *DefaultCMSPageCollection) GetModelName() string {
	return cms.ConstModelNameCMSPageCollection
}

// GetImplementationName returns model implementation name
func (it *DefaultCMSPageCollection) GetImplementationName() string {
	return "Default" + cms.ConstModelNameCMSPageCollection
}

// New returns new instance of model implementation object
func (it *DefaultCMSPageCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCmsPageCollectionName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultCMSPageCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
