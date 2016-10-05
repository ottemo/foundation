package swatch

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("swatch/media", APIListMediaImages)

	// Admin only
	service.POST("swatch/media", api.IsAdmin(APIAddMediaImages))
	service.DELETE("swatch/media/:mediaName", api.IsAdmin(APIRemoveMediaImage))

	return nil
}

// APIListMediaImages returns list of media files from media storage
func APIListMediaImages(context api.InterfaceApplicationContext) (interface{}, error) {

	// skip "unused parameter"
	_ = context

	return mediaStorage.ListMediaDetail(ConstStorageModel, ConstStorageObjectID, ConstStorageMediaType)
}

// APIAddMediaImages uploads images to the media
//   - media file should be provided in "file" field with full name
func APIAddMediaImages(context api.InterfaceApplicationContext) (interface{}, error) {

	var result []interface{}

	files := context.GetRequestFiles()
	if len(files) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c3f4695a-86d5-4269-8b4e-4885d324eb67", "media file was not specified")
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
		err = mediaStorage.Save(ConstStorageModel, ConstStorageObjectID, ConstStorageMediaType, imageName, fileContent)
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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f1b3d05b-9776-4354-a86f-59f5a6d09d50", "media name was not specified")
	}

	// remove media operation
	//---------------------
	err := mediaStorage.Remove(ConstStorageModel, ConstStorageObjectID, ConstStorageMediaType, imageName)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return "ok", nil
}
