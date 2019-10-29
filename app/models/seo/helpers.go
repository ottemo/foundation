package seo

import (
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/env"
)

// GetSEO shortcut to registered engine method
func GetSEO(seoType string, objectID string, urlPattern string) []InterfaceSEOItem {
	if seoEngine := GetRegisteredSEOEngine(); seoEngine != nil {
		return seoEngine.GetSEO(seoType, objectID, urlPattern)
	}
	return []InterfaceSEOItem{}
}

// GetSEOItemModel retrieves current InterfaceSEOItem model implementation
func GetSEOItemModel() (InterfaceSEOItem, error) {
	model, err := models.GetModel(ConstModelNameSEOItem)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	SEOItemModel, ok := model.(InterfaceSEOItem)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fa4bce5b-c500-4faf-81ba-9d28cfff72fb", "model "+model.GetImplementationName()+" is not 'InterfaceSEOItem' capable")
	}

	return SEOItemModel, nil
}

// GetSEOItemModelAndSetID retrieves current InterfaceSEOItem model implementation and sets its ID to some value
func GetSEOItemModelAndSetID(SEOItemID string) (InterfaceSEOItem, error) {

	SEOItemModel, err := GetSEOItemModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = SEOItemModel.SetID(SEOItemID)
	if err != nil {
		return SEOItemModel, env.ErrorDispatch(err)
	}

	return SEOItemModel, nil
}

// LoadSEOItemByID loads SEOItem data into current InterfaceSEOItem model implementation
func LoadSEOItemByID(SEOItemID string) (InterfaceSEOItem, error) {

	SEOItemModel, err := GetSEOItemModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = SEOItemModel.Load(SEOItemID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return SEOItemModel, nil
}

// GetProductCollectionModel retrieves current InterfaceProductCollection model implementation
func GetSEOItemCollectionModel() (InterfaceSEOCollection, error) {
	model, err := models.GetModel(ConstModelNameSEOItemCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stockModel, ok := model.(InterfaceSEOCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bc30efc8-fcad-4fe6-961d-b97208732376", "model "+model.GetImplementationName()+" is not 'InterfaceSEOCollection' capable")
	}

	return stockModel, nil
}
