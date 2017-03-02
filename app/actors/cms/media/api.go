package media

import (
	"io/ioutil"
	"net/http"
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

	// By default "type" is image
	service.GET("cms/media/:mediaType", APIListMedia)
	// Deprecated: Method with explicit mediaType should be used
	service.GET("cms/media", APIListMedia)

	// Admin only
	service.POST("cms/media", api.IsAdminHandler(APIAddMedia))

	// By default "type" is image
	service.DELETE("cms/media/:mediaName/:mediaType", api.IsAdminHandler(APIRemoveMedia))
	// Deprecated: Method with explicit mediaType should be used
	service.DELETE("cms/media/:mediaName", api.IsAdminHandler(APIRemoveMedia))

	return nil
}

// APIListMedia returns list of media files from media
//  - if mediaType is empty - all types will be used
//  - if mediaType is explicit value - only this value will be used
//  - if mediaType is list of types separated by comma - only these types will be shown
func APIListMedia(context api.InterfaceApplicationContext) (interface{}, error) {

	var mediaType = context.GetRequestArgument("mediaType")

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return mediaStorage.ListMediaDetail(ConstStorageModel, ConstStorageObject, mediaType)
}

// APIAddMedia uploads images to the media
//   - media file should be provided in "file" field with full name
//   - the mediaType of files are detected automatically
func APIAddMedia(context api.InterfaceApplicationContext) (interface{}, error) {
	var result []interface{}

	files := context.GetRequestFiles()
	if len(files) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a5a24434-67e8-49c8-aabd-44315ddc9d61", "media files has not been specified")
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
			result = append(result, "Media: '"+fileName+"', should contain extension")
			continue
		}

		var fileContentType = http.DetectContentType(fileContent)

		// Handle image name, adding unique values to name
		var mediaName = strings.TrimSpace(fileName)
		var mediaType = media.ConstMediaTypeImage

		switch fileContentType {
		case "image/gif", "image/png", "image/jpeg":
			{
				mediaNameParts := strings.SplitN(fileName, ".", 2)
				mediaName = mediaNameParts[0] + "_" + utils.InterfaceToString(time.Now().Nanosecond()) + "." + mediaNameParts[1]
			}
		case "application/pdf":
			{
				mediaType = media.ConstMediaTypeDocument
			}
		}

		// save to media storage operation
		err = mediaStorage.Save(ConstStorageModel, ConstStorageObject, mediaType, mediaName, fileContent)
		if err != nil {
			_ = env.ErrorDispatch(err)
			result = append(result, "Media: '"+fileName+"', returned error on save")
			continue
		}

		result = append(result, "ok")
	}

	return result, nil
}

// APIRemoveMedia removes media from storage
//   - media name must be specified in "mediaName" argument
func APIRemoveMedia(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	mediaName := context.GetRequestArgument("mediaName")
	if mediaName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b2f9f0ee-894a-4d8c-8e7a-d150cd9827e2", "media name has not been specified")
	}

	var mediaType = correctMediaType(context.GetRequestArgument("mediaType"))

	// remove media operation
	//---------------------
	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	err = mediaStorage.Remove(ConstStorageModel, ConstStorageObject, mediaType, mediaName)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIGetMedia returns media from storage
//   - media name must be specified in "mediaName" argument
//   - on success case not a JSON data returns, but media file
func APIGetMedia(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	mediaName := context.GetRequestArgument("mediaName")
	if mediaName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "be4cbede-b87e-4384-ab7d-2bb3c05317f5", "media name has not been specified")
	}

	var mediaType = correctMediaType(context.GetRequestArgument("mediaType"))

	// list media operation
	//---------------------
	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	imageFile, err := mediaStorage.Load(ConstStorageModel, ConstStorageObject, mediaType, mediaName)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}
	return imageFile, nil
}

// APIGetMediaPath returns relative path to media library
//   - product id, media type must be specified in "productID" and "mediaType" arguments
func APIGetMediaPath(context api.InterfaceApplicationContext) (interface{}, error) {

	var mediaType = correctMediaType(context.GetRequestArgument("mediaType"))

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	mediaPath, err := mediaStorage.GetMediaPath(ConstStorageModel, ConstStorageObject, mediaType)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return mediaPath, nil
}
