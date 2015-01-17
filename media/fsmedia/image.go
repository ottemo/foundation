package fsmedia

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"os"
	"strconv"
	"strings"

	"io/ioutil"

	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
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
func (it *FilesystemMediaStorage) GetSizeDimensions(size string) (uint, uint, error) {
	var height uint64

	if sizeValue, present := it.imageSizes[size]; present {
		size = sizeValue
	}

	if size == "" && it.baseSize != "" {
		size = it.baseSize
	} else {
		return 0, 0, nil
	}

	value := strings.Split(size, "x")
	width, err := strconv.ParseUint(value[0], 10, 0)
	if err != nil {
		return 0, 0, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f92ff65d-8645-4ee3-83e4-adea0fdb3588", "Invalid size")
	}

	if len(value) > 1 {
		height, _ = strconv.ParseUint(value[1], 10, 0)
	}

	return uint(width), uint(height), nil
}

// GetResizedMediaName returns media filename for a specified size, or error ir size is invalid
func (it *FilesystemMediaStorage) GetResizedMediaName(mediaName string, size string) string {

	// special case for base image size
	if size == "" || size == it.baseSize {
		return mediaName
	}

	// return size + "_" + mediaName

	// checking file extension
	var fileExtension string
	fileName := mediaName

	idx := strings.LastIndex(mediaName, ".")
	if idx != -1 {
		fileExtension = mediaName[idx:]
		fileName = mediaName[0:idx]
	}

	// if we have predefined size
	if _, present := it.imageSizes[size]; present {
		return fmt.Sprintf("%s_%s%s", fileName, size, fileExtension)
	}

	// otherwise
	width, height, _ := it.GetSizeDimensions(size)
	return fmt.Sprintf("%s_%dx%d", fileName, width, height)
}

// TODO: need to refactor this method to several smaller methods to reduce complexity - jwv

// ResizeMediaImage re-sizes specified media to given size, returns error if not possible
func (it *FilesystemMediaStorage) ResizeMediaImage(model string, objID string, mediaName string, size string) error {

	// initializing files path variables
	path, _ := it.GetMediaPath(model, objID, "image")
	path = it.storageFolder + path

	sourceFileName := path + it.GetResizedMediaName(mediaName, "")
	resizedFileName := path + it.GetResizedMediaName(mediaName, size)

	// using biggest resolution file if it exists, otherwise - default
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
	flagResizeNeeded := false
	if _, err := os.Stat(resizedFileName); !os.IsNotExist(err) {
		resizedFile, err := os.Open(resizedFileName)
		defer resizedFile.Close()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		resizedImage, _, err := image.Decode(resizedFile)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		resizedImageSize := resizedImage.Bounds().Size()
		width, height, err := it.GetSizeDimensions(size)
		if err == nil {
			if (int(width) != resizedImageSize.X && int(height) != resizedImageSize.Y) ||
				resizedImageSize.X > int(width) ||
				resizedImageSize.Y > int(height) {

				flagResizeNeeded = true
			}
		}

		if width == 0 && height == 0 {
			flagResizeNeeded = false
		}

	} else {
		flagResizeNeeded = true
	}

	// so, current file not exists or have wrong size
	if flagResizeNeeded {
		var sourceReader io.Reader

		// we can't read and write from same file
		if sourceFileName == resizedFileName {
			sourceContent, err := ioutil.ReadFile(sourceFileName)
			if err != nil {
				return env.ErrorDispatch(err)
			}
			sourceReader = bytes.NewReader(sourceContent)
		} else {
			sourceReader = sourceFile
		}

		// resizing stuff
		width, height, err := it.GetSizeDimensions(size)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		resizedFile, err := os.Create(resizedFileName)
		defer resizedFile.Close()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		sourceImage, imageFormat, err := image.Decode(sourceReader)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		originalSize := sourceImage.Bounds().Size()
		if width != 0 && height != 0 {
			if originalSize.X > originalSize.Y {
				height = 0
			} else {
				width = 0
			}
		}
		resizedImage := resize.Resize(width, height, sourceImage, resize.Bilinear)

		// encoding re-sized image
		switch imageFormat {
		case "jpeg":
			err = jpeg.Encode(resizedFile, resizedImage, nil)
		case "png":
			err = png.Encode(resizedFile, resizedImage)
		default:
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "42f0cbb3-9187-4e16-8953-5829ea8d2da8", "Unknown image format to encode")
		}

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
