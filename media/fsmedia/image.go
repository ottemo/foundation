package fsmedia

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/ottemo/foundation/env"
)

// GetBiggestSize returns size with a image bounds limit 0x0 or with most resolution, "" - default size
func (it *FilesystemMediaStorage) GetBiggestSize() string {
	maxSize := ""

	maxWidth, maxHeight, _ := it.GetSizeDimensions(it.baseSize)
	if maxWidth == 0 && maxHeight == 0 {
		return maxSize
	}

	for imageSize := range it.imageSizes {
		width, height, _ := it.GetSizeDimensions(it.baseSize)
		if width == 0 && height == 0 {
			return imageSize
		}

		if width > maxWidth || height > maxHeight || width > maxHeight || height > maxWidth {
			maxWidth = width
			maxHeight = height
			maxSize = imageSize
		}
	}

	return maxSize
}

// UpdateBaseSize loads predefined sizes from config value
func (it *FilesystemMediaStorage) UpdateBaseSize(newValue string) error {
	newValue = strings.TrimSpace(newValue)
	_, _, err := it.GetSizeDimensions(newValue)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	it.baseSize = newValue
	it.biggestSize = it.GetBiggestSize()

	return nil
}

// UpdateSizeNames loads predefined sizes from config value
func (it *FilesystemMediaStorage) UpdateSizeNames(newValue string) error {
	result := make(map[string]string)
	for _, size := range strings.Split(newValue, ",") {
		var sizeName, sizeValue string
		sizeArray := strings.Split(size, ":")

		if len(sizeArray) > 1 {
			sizeName = strings.TrimSpace(sizeArray[0])
			sizeValue = strings.TrimSpace(sizeArray[1])
		} else {
			sizeValue = strings.TrimSpace(sizeArray[0])
		}

		_, _, err := it.GetSizeDimensions(sizeValue)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		if sizeValue != "" {
			if sizeName != "" {
				result[sizeName] = sizeValue
			} else {
				result[sizeValue] = sizeValue
			}
		}

	}
	it.imageSizes = result
	it.biggestSize = it.GetBiggestSize()
	return nil
}

// GetSizeDimensions returns width and height for specified size or error if size is ot valid
//   - supposed that 32-bit integer positive part (up to 2147483647) will be enough to reflect width/height dimension
func (it *FilesystemMediaStorage) GetSizeDimensions(size string) (int, int, error) {
	var height uint64

	if sizeValue, present := it.imageSizes[size]; present {
		size = sizeValue
	}

	if size == "" {
		if it.baseSize != "" {
			size = it.baseSize
		} else {
			return 0, 0, nil
		}
	}

	value := strings.Split(size, "x")
	width, err := strconv.ParseUint(value[0], 10, 0)
	if err != nil {
		return 0, 0, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f92ff65d-8645-4ee3-83e4-adea0fdb3588", "Invalid size")
	}

	if len(value) > 1 {
		height, _ = strconv.ParseUint(value[1], 10, 0)
	}

	return int(width), int(height), nil
}

// GetResizedMediaName returns media filename for a specified size, or error ir size is invalid
func (it *FilesystemMediaStorage) GetResizedMediaName(mediaName string, size string) string {

	// special case for base image size
	if size == "" || size == it.baseSize {
		return mediaName
	}

	// Get the extension and name
	var fileExtension string
	fileName := mediaName

	idx := strings.LastIndex(mediaName, ".")
	if idx != -1 {
		fileExtension = mediaName[idx:]
		fileName = mediaName[0:idx]
	}

	// Get the dimensions
	width, height, _ := it.GetSizeDimensions(size)

	return fmt.Sprintf("%s_%dx%d%s", fileName, width, height, fileExtension)
}

// TODO: need to refactor this method to several smaller methods to reduce complexity - jwv

// ResizeMediaImage re-sizes specified media to given size, returns error if not possible
func (it *FilesystemMediaStorage) ResizeMediaImage(model string, objID string, mediaName string, size string) error {

	// initializing files path variables
	path, _ := it.GetMediaPath(model, objID, "image")
	path = it.storageFolder + path

	sourceFileName := path + it.GetResizedMediaName(mediaName, "")
	destinationFileName := path + it.GetResizedMediaName(mediaName, size)

	// determination of source image for resizing
	//   - checking higher quality file existence (biggest resolution)
	//   - if not exists we will use default image as source
	biggestFileName := path + it.GetResizedMediaName(mediaName, it.biggestSize)
	if _, err := os.Stat(biggestFileName); !os.IsNotExist(err) {
		sourceFileName = biggestFileName
	}

	// opening source file
	sourceFile, err := os.Open(sourceFileName)
	defer sourceFile.Close()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// checking destination file
	//   - if file not exists then resizing is required
	//   - if file exists, opening it and checking if current dimensions fits requested
	flagResizeRequired := false
	if _, err := os.Stat(destinationFileName); !os.IsNotExist(err) {

		destinationImage, err := imaging.Open(destinationFileName)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		// reading and checking destination image sizes
		destinationImageSize := destinationImage.Bounds().Size()
		requiredWidth, requiredHeight, err := it.GetSizeDimensions(size)
		if err == nil {
			// checking that both image dimension are not equals to required or
			// that one of dimension equals required but other not fits required box
			if (requiredWidth != 0 && requiredWidth != destinationImageSize.X && requiredHeight != 0 && requiredHeight != destinationImageSize.Y) ||
				(destinationImageSize.Y == requiredHeight && requiredWidth != 0 && destinationImageSize.X != requiredWidth) ||
				(destinationImageSize.X == requiredWidth && requiredHeight != 0 && destinationImageSize.Y != requiredHeight) {

				flagResizeRequired = true
			}

			// images with background should exactly match for both dimensions
			if ConstResizeOnBackground &&
				((requiredWidth != 0 && requiredWidth != destinationImageSize.X) ||
					requiredHeight != 0 && requiredHeight != destinationImageSize.Y) {

				flagResizeRequired = true
			}
		}

	} else {
		flagResizeRequired = true
	}

	// making resize if it was requested
	if flagResizeRequired {

		// resize routine
		requiredWidth, requiredHeight, err := it.GetSizeDimensions(size)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		sourceImage, err := imaging.Open(sourceFileName)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		resizedImage := sourceImage

		if requiredWidth != 0 || requiredHeight != 0 {
			// imaging.Thumbnail works bad imaging.Fit - only reduces dimensions
			// so, using making imaging.Resize with own dimensions detect
			resizeWidth := requiredWidth
			resizeHeight := requiredHeight
			if requiredWidth != 0 && requiredHeight != 0 {
				bounds := sourceImage.Bounds()
				srcAspect := float64(bounds.Dx()) / float64(bounds.Dy())
				dstAspect := float64(requiredWidth) / float64(requiredHeight)

				if srcAspect > dstAspect {
					resizeHeight = 0
				} else {
					resizeWidth = 0
				}
			}

			resizedImage = imaging.Resize(sourceImage, resizeWidth, resizeHeight, imaging.Linear)

			if ConstResizeOnBackground {
				background := imaging.New(requiredWidth, requiredHeight, image.White)
				resizedImage = imaging.PasteCenter(background, resizedImage)
			}
		}

		err = imaging.Save(resizedImage, destinationFileName)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// GetResizedImage returns re-sized image contents or error if not possible
func (it *FilesystemMediaStorage) GetResizedImage(model string, objID string, mediaName string, size string) ([]byte, error) {
	path, _ := it.GetMediaPath(model, objID, "image")
	path = it.storageFolder + path

	resizedFileName := path + it.GetResizedMediaName(mediaName, size)
	resizedFileContents, err := ioutil.ReadFile(resizedFileName)
	if os.IsNotExist(err) {

		err = it.ResizeMediaImage(model, objID, mediaName, size)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		resizedFileContents, err = ioutil.ReadFile(resizedFileName)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

	}

	return resizedFileContents, env.ErrorDispatch(err)
}
