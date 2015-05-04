package seo

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

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
