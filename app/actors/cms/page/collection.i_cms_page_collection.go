package page

import (
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetDBCollection returns database collection
func (it *DefaultCMSPageCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// ListCMSPages returns list of cms page model items
func (it *DefaultCMSPageCollection) ListCMSPages() []cms.InterfaceCMSPage {
	var result []cms.InterfaceCMSPage

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		cmsPageModel, err := cms.GetCMSPageModel()
		if err != nil {
			return result
		}
		if err := cmsPageModel.FromHashMap(recordData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8e74b4b4-08e9-44d5-805d-e158f5af518c", err.Error())
		}

		result = append(result, cmsPageModel)
	}

	return result
}
