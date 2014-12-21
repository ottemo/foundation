package cms

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// GetCMSPageCollectionModel retrieves current InterfaceCMSPageCollection model implementation
func GetCMSPageCollectionModel() (InterfaceCMSPageCollection, error) {
	model, err := models.GetModel(ConstModelNameCMSPageCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cmsPageCollectionModel, ok := model.(InterfaceCMSPageCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "36ee313f94f347a4a44cd487662dcc2a", "model "+model.GetImplementationName()+" is not 'InterfaceCMSPageCollection' capable")
	}

	return cmsPageCollectionModel, nil
}

// GetCMSPageModel retrieves current InterfaceCMSPage model implementation
func GetCMSPageModel() (InterfaceCMSPage, error) {
	model, err := models.GetModel(ConstModelNameCMSPage)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cmsPageModel, ok := model.(InterfaceCMSPage)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "36d47e0ae89e4975bf42a9971f1ea532", "model "+model.GetImplementationName()+" is not 'InterfaceCMSPage' capable")
	}

	return cmsPageModel, nil
}

// GetCMSPageModelAndSetID retrieves current InterfaceCMSPage model implementation and sets its ID to some value
func GetCMSPageModelAndSetID(cmsPageID string) (InterfaceCMSPage, error) {

	cmsPageModel, err := GetCMSPageModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cmsPageModel.SetID(cmsPageID)
	if err != nil {
		return cmsPageModel, env.ErrorDispatch(err)
	}

	return cmsPageModel, nil
}

// LoadCMSPageByID loads cmsPage data into current InterfaceCMSPage model implementation
func LoadCMSPageByID(cmsPageID string) (InterfaceCMSPage, error) {

	cmsPageModel, err := GetCMSPageModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cmsPageModel.Load(cmsPageID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cmsPageModel, nil
}

// GetCMSBlockCollectionModel retrieves current InterfaceCMSBlockCollection model implementation
func GetCMSBlockCollectionModel() (InterfaceCMSBlockCollection, error) {
	model, err := models.GetModel(ConstModelNameCMSBlockCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	csmBlockCollectionModel, ok := model.(InterfaceCMSBlockCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "81a29a453d7441eabc89753b0b1e48b7", "model "+model.GetImplementationName()+" is not 'InterfaceCMSBlockCollection' capable")
	}

	return csmBlockCollectionModel, nil
}

// CMS Block helpers
//------------------

// GetCMSBlockModel retrieves current InterfaceCMSBlock model implementation
func GetCMSBlockModel() (InterfaceCMSBlock, error) {
	model, err := models.GetModel(ConstModelNameCMSBlock)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	csmBlockModel, ok := model.(InterfaceCMSBlock)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6c8d19f969624a5bbc2ac25bfd20e3af", "model "+model.GetImplementationName()+" is not 'InterfaceCMSBlock' capable")
	}

	return csmBlockModel, nil
}

// GetCMSBlockModelAndSetID retrieves current InterfaceCMSBlock model implementation and sets its ID to some value
func GetCMSBlockModelAndSetID(csmBlockID string) (InterfaceCMSBlock, error) {

	csmBlockModel, err := GetCMSBlockModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = csmBlockModel.SetID(csmBlockID)
	if err != nil {
		return csmBlockModel, env.ErrorDispatch(err)
	}

	return csmBlockModel, nil
}

// LoadCMSBlockByID loads csmBlock data into current InterfaceCMSBlock model implementation
func LoadCMSBlockByID(csmBlockID string) (InterfaceCMSBlock, error) {

	csmBlockModel, err := GetCMSBlockModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = csmBlockModel.Load(csmBlockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return csmBlockModel, nil
}
