package media

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("cms/media", APIListMediaImages)

	// Admin only
	service.POST("cms/media", api.IsAdmin(APIAddMediaImages))
	service.DELETE("cms/media/:mediaName", api.IsAdmin(APIRemoveMediaImage))

	return nil
}

// APIListMediaImages returns list of media files from media
func APIListMediaImages(context api.InterfaceApplicationContext) (interface{}, error) {

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return mediaStorage.ListMediaDetail(ConstStorageModel, ConstStorageObject, ConstStorageType)
}

// APIAddMediaImages uploads images to the media
//   - media file should be provided in "file" field with full name
func APIAddMediaImages(context api.InterfaceApplicationContext) (interface{}, error) {
	var result []interface{}

	files := context.GetRequestFiles()
	if len(files) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "75a2ddaf-b63d-4eed-b16d-4b32778f5fc1", "media file was not specified")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	for fileName, fileReader := range files {
		fileContent, err := ioutil.ReadAll(fileReader)
		if err != nil {
			return result, env.ErrorDispatch(err)
		}

		if !strings.Contains(fileName, ".") {
			result = append(result, "Image: '"+fileName+"', should contain extension")
			continue
		}

		// Handle image name, adding unique values to name
		fileName = strings.TrimSpace(fileName)
		mediaNameParts := strings.SplitN(fileName, ".", 2)
		imageName := mediaNameParts[0] + "_" + utils.InterfaceToString(time.Now().Nanosecond()) + "." + mediaNameParts[1]

		// save to media storage operation
		err = mediaStorage.Save(ConstStorageModel, ConstStorageObject, ConstStorageType, imageName, fileContent)
		if err != nil {
			env.ErrorDispatch(err)
			result = append(result, "Image: '"+fileName+"', returned error on save")
			continue
		}

		result = append(result, "ok")
	}

	return result, nil
}

// APIRemoveMediaImage removes image from media
//   - media name must be specified in "mediaName" argument
func APIRemoveMediaImage(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	imageName := context.GetRequestArgument("mediaName")
	if imageName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "63b37b08-3b21-48b7-9058-291bb7e635a1", "media name was not specified")
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

// APIGetMediaImage returns image from media
//   - media name must be specified in "mediaName" argument
//   - on success case not a JSON data returns, but media file
func APIGetMediaImage(context api.InterfaceApplicationContext) (interface{}, error) {

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

// APIGetMediaPath returns relative path to media library
//   - product id, media type must be specified in "productID" and "mediaType" arguments
func APIGetMediaPath(context api.InterfaceApplicationContext) (interface{}, error) {

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
