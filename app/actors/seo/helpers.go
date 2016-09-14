package seo

import (
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/seo"
)

// GetSEOItemCollectionModel retrieves current InterfaceSEOCollection model implementation
func GetSEOItemCollectionModel() (seo.InterfaceSEOCollection, error) {
	model, err := models.GetModel(ConstCollectionNameURLRewrites)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	seoItemCollectionModel, ok := model.(seo.InterfaceSEOCollection)
	if !ok {
		return nil, env.ErrorNew(
			ConstErrorModule,
			ConstErrorLevel,
			"2198576e-2bf4-4631-a8b3-52b6f661f693",
			"model "+model.GetImplementationName()+" is not 'InterfaceSEOCollection' capable")
	}

	return seoItemCollectionModel, nil
}
