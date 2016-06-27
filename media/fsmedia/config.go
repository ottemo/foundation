package fsmedia

import (
	"github.com/ottemo/foundation/env"
)

// setupConfig setups package configuration values for a system
func (it *FilesystemMediaStorage) setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule,
			env.ConstErrorLevelStartStop,
			"7c6d9092-9a17-4d96-a796-fd0f7b64c00a",
			"Unable to find configuration media configuration variable, please specify one.")
		return env.ErrorDispatch(err)
	}

	imageSizesValidator := func(newValue interface{}) (interface{}, error) {
		if newValue, ok := newValue.(string); ok && newValue != "" {
			err := it.UpdateSizeNames(newValue)
			if err != nil {
				return ConstDefaultImageSizes, env.ErrorDispatch(err)
			}

			// resize all images instantly on change size if we not in "on fly" mode
			// and prevent from resizing while initial calls for validator
			if !resizeImagesOnFly && it.setupWaitCnt == 0 {
				err = it.ResizeAllMediaImages()
				if err != nil {
					env.ErrorDispatch(err)
				}
			}

			return newValue, nil
		}
		err := env.ErrorNew(ConstErrorModule,
			env.ConstErrorLevelStartStop,
			"165834bd-d8ca-4b1d-8cfa-b8a153288913",
			"Unexpected value for image size found.")
		return ConstDefaultImageSizes, env.ErrorDispatch(err)
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMediaImageSizes,
		Value:       ConstDefaultImageSizes,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Image sizes",
		Description: "Define image sizes in the format ([sizeName]: [maxWidth:maxHeight], ...)",
		Image:       "",
	}, imageSizesValidator)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	imageSizesValidator(env.ConfigGetValue(ConstConfigPathMediaImageSizes))

	imageDefaultSizeValidator := func(newValue interface{}) (interface{}, error) {
		if newValue, ok := newValue.(string); ok && newValue != "" {
			err := it.UpdateBaseSize(newValue)
			if err != nil {
				return ConstDefaultImageSize, env.ErrorDispatch(err)
			}
			return newValue, nil
		}
		err := env.ErrorNew(ConstErrorModule,
			env.ConstErrorLevelStartStop,
			"165834bd-d8ca-4b1d-8cfa-b8a153288913",
			"Unexpected value for image size found.")
		return ConstDefaultImageSize, env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMediaImageSize,
		Value:       ConstDefaultImageSize,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Base image size",
		Description: "Specify the size of primary image in the format [maxWidth:maxHeight], leave 0x0 of blank if resizing is not needed.",
		Image:       "",
	}, imageDefaultSizeValidator)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	imageDefaultSizeValidator(env.ConfigGetValue(ConstConfigPathMediaImageSize))

	imageBasePathValidator := func(newValue interface{}) (interface{}, error) {
		if newValue, ok := newValue.(string); ok {
			mediaPath := newValue
			if length := len(newValue); length > 1 && newValue[length-1:length] == "/" {
				mediaPath = newValue[0 : length-1]
			}
			mediaBasePath = mediaPath
			return mediaBasePath, nil
		}

		err := env.ErrorNew(ConstErrorModule,
			env.ConstErrorLevelStartStop,
			"eb97378b-a940-45b4-a653-3bdb47fe6b16",
			"Unexpected value found for the image base url.")
		return mediaBasePath, env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMediaBaseURL,
		Value:       "media",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Media base URL",
		Description: "URL application will use to generate media resources links",
		Image:       "",
	}, imageBasePathValidator)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	imageBasePathValidator(env.ConfigGetValue(ConstConfigPathMediaBaseURL))

	it.setupCheckDone()

	return nil
}
