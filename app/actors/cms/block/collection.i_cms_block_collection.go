package block

import (
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
)

// GetDBCollection returns database collection
func (it *DefaultCMSBlockCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// ListCMSBlocks returns list of cms block model items
func (it *DefaultCMSBlockCollection) ListCMSBlocks() []cms.InterfaceCMSBlock {
	var result []cms.InterfaceCMSBlock

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		cmsBlockModel, err := cms.GetCMSBlockModel()
		if err != nil {
			return result
		}
		cmsBlockModel.FromHashMap(recordData)

		result = append(result, cmsBlockModel)
	}

	return result
}
