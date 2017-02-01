package fsmedia

import (
	"bytes"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"
)

// GetName returns the media storage name
func (it *FilesystemMediaStorage) GetName() string {
	return "FilesystemMediaStorage"
}

// GetMediaPath returns the path needed to access the media file
func (it *FilesystemMediaStorage) GetMediaPath(model string, objID string, mediaType string) (string, error) {
	return mediaType + "/" + model + "/" + objID + "/", nil
}

// Load retrieves contents of the media entity for a given model object
func (it *FilesystemMediaStorage) Load(model string, objID string, mediaType string, mediaName string) ([]byte, error) {
	mediaPath, err := it.GetMediaPath(model, objID, mediaType)
	if err != nil {
		return nil, err
	}

	mediaFilePath := it.storageFolder + mediaPath + mediaName

	return ioutil.ReadFile(mediaFilePath)
}

// Save adds a media entity for the model object to the database and saves it on
// the filesystem.
func (it *FilesystemMediaStorage) Save(model string, objID string, mediaType string, mediaName string, mediaData []byte) error {
	mediaPath, err := it.GetMediaPath(model, objID, mediaType)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	mediaFolder := it.storageFolder + mediaPath
	mediaFilePath := mediaFolder + mediaName

	if _, err := os.Stat(mediaFolder); !os.IsExist(err) {
		err := os.MkdirAll(mediaFolder, os.ModePerm)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	// we have image associated media, so making special treatment
	if mediaType == media.ConstMediaTypeImage {

		// checking that image is png or jpeg, making it jpeg if not
		decodedImage, imageFormat, err := image.Decode(bytes.NewReader(mediaData))
		if err != nil {
			return env.ErrorDispatch(err)
		}

		// making sure file extension is right
		idx := strings.LastIndex(mediaName, ".")
		fileExt := ""
		if idx != -1 {
			fileExt = mediaName[idx:]
		}

		if imageFormat == "jpeg" && fileExt != ".jpeg" && fileExt != ".jpg" {
			mediaName += ".jpg"
			mediaFilePath += ".jpg"
		}

		if imageFormat == "png" && fileExt != ".png" {
			mediaName += ".png"
			mediaFilePath += ".png"
		}

		// converting image to known format (jpeg) if needed
		if imageFormat != "jpeg" && imageFormat != "png" {

			if fileExt != ".jpeg" && fileExt != ".jpg" {
				mediaName += ".jpg"
				mediaFilePath += ".jpg"
			}

			newFile, err := os.Create(mediaFilePath)
			defer newFile.Close()
			if err != nil {
				return env.ErrorDispatch(err)
			}

			err = jpeg.Encode(newFile, decodedImage, nil)
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}

		ioerr := ioutil.WriteFile(mediaFilePath, mediaData, os.ModePerm)
		if ioerr != nil {
			return env.ErrorDispatch(ioerr)
		}

		// ResizeMediaImage will check necessity of resize by it self
		for imageSize := range it.imageSizes {
			it.ResizeMediaImage(model, objID, mediaName, imageSize)
		}
		it.ResizeMediaImage(model, objID, mediaName, it.baseSize)

	} else {

		ioerr := ioutil.WriteFile(mediaFilePath, mediaData, os.ModePerm)
		if ioerr != nil {
			return env.ErrorDispatch(ioerr)
		}
	}

	// making database record
	//------------------------

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "92bd8bc9-8d04-43a0-9721-7b57aaabac7f", "Can't get database engine")
	}

	dbCollection, err := dbEngine.GetCollection(ConstMediaDBCollection)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbCollection.AddFilter("model", "=", model)
	dbCollection.AddFilter("object", "=", objID)
	dbCollection.AddFilter("type", "=", mediaType)
	dbCollection.AddFilter("media", "=", mediaName)

	count, err := dbCollection.Count()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if count == 0 {
		_, err = dbCollection.Save(map[string]interface{}{"model": model, "object": objID, "type": mediaType, "media": mediaName, "created_at": time.Now()})
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// Remove will delete the entity for the model object from the database and
// remove it from the filesystem.
func (it *FilesystemMediaStorage) Remove(model string, objID string, mediaType string, mediaName string) error {

	// preparing DB collection
	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel,
			"f3af959a-4d1e-4dd8-a827-7652fcb5402a",
			"Unable to find database engine to remove media entity.")
	}

	dbCollection, err := dbEngine.GetCollection(ConstMediaDBCollection)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.AddFilter("model", "=", model)
	err = dbCollection.AddFilter("object", "=", objID)
	err = dbCollection.AddFilter("type", "=", mediaType)
	err = dbCollection.AddFilter("media", "=", mediaName)

	// removing files
	records, err := dbCollection.Load()

	mediaFolder := it.storageFolder

	for _, record := range records {
		if mediaName, ok := record["media"].(string); ok {

			if path, err := it.GetMediaPath(model, objID, mediaType); err == nil {

				// looking for object image sizes to remove
				if mediaType == media.ConstMediaTypeImage {
					for imageSize := range it.imageSizes {
						mediaFilePath := mediaFolder + path + it.GetResizedMediaName(mediaName, imageSize)
						os.Remove(mediaFilePath)
					}
				}

				os.Remove(mediaFolder + path + mediaName)
			}
		}
	}

	// removing DB records
	_, err = dbCollection.Delete()

	return env.ErrorDispatch(err)
}

// ListMedia returns the list of given type media entities for a given model object
func (it *FilesystemMediaStorage) ListMedia(model string, objID string, mediaType string) ([]string, error) {
	var result []string

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel,
			"d4c8dd6e-95be-4b8f-b606-5544b633cd7c",
			"Unable fo find database engine in order to return the current list of media.")
	}

	dbCollection, err := dbEngine.GetCollection(ConstMediaDBCollection)
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	dbCollection.AddFilter("model", "=", model)
	dbCollection.AddFilter("object", "=", objID)
	dbCollection.AddFilter("type", "=", mediaType)

	records, err := dbCollection.Load()

	for _, record := range records {
		if mediaName, ok := record["media"].(string); ok {
			result = append(result, mediaName)
		}
	}

	// checking that obj
	// ect have all image sizes
	if resizeImagesOnFly && mediaType == media.ConstMediaTypeImage {
		// ResizeMediaImage will check necessity of resize by it self
		for _, mediaName := range result {
			for imageSize := range it.imageSizes {
				it.ResizeMediaImage(model, objID, mediaName, imageSize)
			}
			it.ResizeMediaImage(model, objID, mediaName, it.baseSize)
		}
	}

	return result, env.ErrorDispatch(err)
}

// ListMediaDetail returns the list with media details of given type media entities for a given model object
// mediaType could be empty, single value or list of values separated by comma
func (it *FilesystemMediaStorage) ListMediaDetail(model string, objID string, mediaType string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel,
			"07ba2560-3d2d-4cbd-8fc7-a52cadc01b96",
			"Unable fo find database engine in order to return the current list of media.")
	}

	dbCollection, err := dbEngine.GetCollection(ConstMediaDBCollection)
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	if err := dbCollection.AddFilter("model", "=", model); err != nil {
		return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bdbc95b5-78ff-4d08-9945-c760fe0c8b02", "Internal error: unable to filter by model ["+model+"].")
	}

	if err := dbCollection.AddFilter("object", "=", objID); err != nil {
		return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f2a1b585-cd90-4efb-a12b-387fb09f4377", "Internal error: unable to filter by object ["+objID+"].")
	}

	var pathes = map[string]string{}

	if len(mediaType) > 0 {
		var mediaTypeValue = strings.Trim(mediaType, ",")
		var mediaTypeOptions = strings.Split(mediaTypeValue, ",")
		if err := dbCollection.AddFilter("type", "in", mediaTypeOptions); err != nil {
			return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d2031dbe-1f10-489b-ba95-ab5b64cf0a98",
				"Internal error: unable to filter by type ["+utils.InterfaceToString(mediaTypeOptions)+"].")
		}

		for _, mediaTypeOption := range mediaTypeOptions {
			path, _ := it.GetMediaPath(model, objID, mediaTypeOption)
			pathes[mediaTypeOption] = mediaBasePath + "/" + path
		}
	}

	records, err := dbCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, record := range records {
		name := utils.InterfaceToString(record["media"])
		recordType := utils.InterfaceToString(record["type"])

		path, present := pathes[recordType]
		if !present {
			path, _ = it.GetMediaPath(model, objID, recordType)
			path = mediaBasePath + "/" + path
			pathes[recordType] = path
		}

		mediaObject := map[string]interface{}{
			"id":         record["_id"],
			"name":       name,
			"url":        path + name,
			"created_at": record["created_at"],
			"type":       recordType,
		}

		result = append(result, mediaObject)
	}

	return result, env.ErrorDispatch(err)
}

// GetAllSizes returns a list of all image sizes in a []map[string], included is
// the path and type of media.
func (it *FilesystemMediaStorage) GetAllSizes(model string, objID string, mediaType string) ([]map[string]string, error) {

	var result []map[string]string

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel,
			"d4c8dd6e-95be-4b8f-b606-5544b633cd7c",
			"Unable to find a database engine to return all image sizes.")
	}

	dbCollection, err := dbEngine.GetCollection(ConstMediaDBCollection)
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	dbCollection.AddFilter("model", "=", model)
	dbCollection.AddFilter("object", "=", objID)
	dbCollection.AddFilter("type", "=", mediaType)

	records, err := dbCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, record := range records {

		if mediaName, ok := record["media"].(string); ok {
			mediaSet, err := it.GetSizes(model, objID, mediaType, mediaName)
			if err != nil {
				env.ErrorDispatch(err)
			}

			result = append(result, mediaSet)
		}
	}

	return result, nil
}

// GetSizes returns a list of all sizes for specificed image in a map[string],
// included is image path, model object and the image name.
func (it *FilesystemMediaStorage) GetSizes(model string, objID string, mediaType string, mediaName string) (map[string]string, error) {
	mediaSet := make(map[string]string)
	if mediaName == "" || model == "" || objID == "" || mediaType == "" {
		return mediaSet, nil
	}

	path, err := it.GetMediaPath(model, objID, mediaType)
	if err != nil {
		return mediaSet, env.ErrorDispatch(err)
	}

	path = mediaBasePath + "/" + path

	// Loop over the sizes we support
	for imageSize := range it.imageSizes {
		mediaSet[imageSize] = path + it.GetResizedMediaName(mediaName, imageSize)
		if resizeImagesOnFly {
			it.ResizeMediaImage(model, objID, mediaName, imageSize)
		}
	}

	if resizeImagesOnFly {
		it.ResizeMediaImage(model, objID, mediaName, it.baseSize)
	}

	return mediaSet, nil
}

// ResizeAllMediaImages will resize all images for currently specified sizes
func (it *FilesystemMediaStorage) ResizeAllMediaImages() error {
	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel,
			"92bd8bc9-8d04-43a0-9721-7b57aaabac7f",
			"Unable to find a database engine when attempting to resize all media.")
	}

	dbCollection, err := dbEngine.GetCollection(ConstMediaDBCollection)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbCollection.AddFilter("type", "=", media.ConstMediaTypeImage)

	records, err := dbCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var imagesResized int
	for _, record := range records {
		if !utils.KeysInMapAndNotBlank(record, "model", "object", "media") {
			env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel,
				"c85a9dc8-01eb-4f42-845b-34a0863dbd43",
				"Media associated with the following id, "+utils.InterfaceToString(record["id"])+", did not contain one of the required column values."))
			continue
		}

		mediaModel := utils.InterfaceToString(record["model"])
		mediaObject := utils.InterfaceToString(record["object"])
		mediaName := utils.InterfaceToString(record["media"])

		if err := it.ResizeMediaImage(mediaModel, mediaObject, mediaName, it.baseSize); err == nil {
			for _, size := range it.imageSizes {
				if it.ResizeMediaImage(mediaModel, mediaObject, mediaName, size) != nil {
					break
				}
			}
			imagesResized++
		} else {
			env.ErrorDispatch(err)
		}
	}

	if imagesResized != len(records) {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel,
			"73da10f7-8f40-4ba3-bcbf-ac9a505fa922",
			"Unable to resize all all images, result: "+utils.InterfaceToString(imagesResized)+" from "+utils.InterfaceToString(len(records)))
	}

	return nil
}
