package seo

import (
	"github.com/ottemo/foundation/app/models/seo"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetSEO returns records for a given filter based on function arguments
//   - use blank string to exclude filter field
func (it *DefaultSEOEngine) GetSEO(seoType string, objectID string, urlPattern string) []seo.InterfaceSEOItem {

	result := []seo.InterfaceSEOItem{}

	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		env.ErrorDispatch(err)
		return result
	}

	if seoType != "" {
		collection.AddFilter("type", "=", seoType)
	}
	if objectID != "" {
		collection.AddFilter("rewrite", "=", objectID)
	}
	if urlPattern != "" {
		collection.AddFilter("url", " like ", urlPattern)
	}

	records, err := collection.Load()

	for _, record := range records {
		seoItem, err := seo.LoadSEOItemByID(utils.InterfaceToString(record["_id"]))
		if err == nil {
			result = append(result, seoItem)
		}
	}

	return result
}
