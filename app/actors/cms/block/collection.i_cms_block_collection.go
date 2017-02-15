package block

import (
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
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
		if err := cmsBlockModel.FromHashMap(recordData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "143f447e-4bb7-46e5-9481-523ccf48fc70", err.Error())
		}

		result = append(result, cmsBlockModel)
	}

	return result
}
