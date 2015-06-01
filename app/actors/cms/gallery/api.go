package gallery

import (
	"io/ioutil"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("cms/gallery/image/:mediaName", api.ConstRESTOperationGet, APIGetGalleryImage)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("cms/gallery/images", api.ConstRESTOperationGet, APIListGalleryImages)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("cms/gallery/image/:mediaName", api.ConstRESTOperationCreate, APIAddGalleryImage)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("cms/gallery/image/:mediaName", api.ConstRESTOperationDelete, APIRemoveGalleryImage)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("cms/gallery/path", api.ConstRESTOperationGet, APIGetGalleryPath)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIGetGalleryPath returns relative path to gallery library
//   - product id, media type must be specified in "productID" and "mediaType" arguments
func APIGetGalleryPath(context api.InterfaceApplicationContext) (interface{}, error) {

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	mediaPath, err := mediaStorage.GetMediaPath(ConstStorageModel, ConstStorageObject, ConstStorageType)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return mediaPath, nil
}

// APIListGalleryImages returns list of media files from gallery
func APIListGalleryImages(context api.InterfaceApplicationContext) (interface{}, error) {

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	mediaList, err := mediaStorage.ListMedia(ConstStorageModel, ConstStorageObject, ConstStorageType)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return mediaList, nil
}

// APIAddGalleryImage uploads image to the gallery
//   - media name should be specified in "mediaName" arguments
//   - media file should be provided in "file" field
func APIAddGalleryImage(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	imageName := context.GetRequestArgument("mediaName")
	if imageName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "23fb7617-f19a-4505-b706-10f7898fd980", "media name was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// income file processing
	//-----------------------
	files := context.GetRequestFiles()
	if len(files) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "75a2ddaf-b63d-4eed-b16d-4b32778f5fc1", "media file was not specified")
	}

	var fileContents []byte
	for _, fileReader := range files {
		contents, err := ioutil.ReadAll(fileReader)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		fileContents = contents
		break
	}

	// add media operation
	//--------------------
	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	err = mediaStorage.Save(ConstStorageModel, ConstStorageObject, ConstStorageType, imageName, fileContents)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIRemoveGalleryImage removes image from gallery
//   - media name must be specified in "mediaName" argument
func APIRemoveGalleryImage(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	imageName := context.GetRequestArgument("mediaName")
	if imageName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "63b37b08-3b21-48b7-9058-291bb7e635a1", "media name was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// remove media operation
	//---------------------
	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	err = mediaStorage.Remove(ConstStorageModel, ConstStorageObject, ConstStorageType, imageName)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIGetGalleryImage returns image from gallery
//   - media name must be specified in "mediaName" argument
//   - on success case not a JSON data returns, but media file
func APIGetGalleryImage(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	imageName := context.GetRequestArgument("mediaName")
	if imageName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "124c8b9d-1a6b-491c-97ba-a03e8c828337", "media name was not specified")
	}

	// list media operation
	//---------------------
	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	imageFile, err := mediaStorage.Load(ConstStorageModel, ConstStorageObject, ConstStorageType, imageName)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}
	return imageFile, nil
}
