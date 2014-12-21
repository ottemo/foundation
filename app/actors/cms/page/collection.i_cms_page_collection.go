package page

import (
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
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
		cmsPageModel.FromHashMap(recordData)

		result = append(result, cmsPageModel)
	}

	return result
}
