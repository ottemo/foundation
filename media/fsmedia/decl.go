package fsmedia

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstMediaDBCollection  = "media"    // database collection name to store media assignment information into
	ConstMediaDefaultFolder = "./media/" // filesystem folder path to store media files in there

	ConstResizeOnBackground = true

	ConstDefaultImageSize  = "1000x1000"      // "800x400"
	ConstDefaultImageSizes = "thumb: 280x350" // "small: 75x75, thumb: 260x300, big: 560x650"

	ConstConfigPathMediaImageSize  = "general.app.image_size"  // base image size
	ConstConfigPathMediaImageSizes = "general.app.image_sizes" // other image sizes required

	ConstConfigPathMediaBaseURL = "general.app.media_base_url"

	ConstMediaTypeImage    = "image"
	ConstMediaTypeLink     = "link"
	ConstMediaTypeDocument = "document"

	ConstErrorModule = "media/fsmedia"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// resizeImagesOnFly can be specified in ini config file by key "media.resize.images.onfly", false by default
var (
	resizeImagesOnFly bool
	mediaBasePath     = "media"
)

// FilesystemMediaStorage is a filesystem based implementer of InterfaceMediaStorage
type FilesystemMediaStorage struct {
	storageFolder string
	setupWaitCnt  int

	baseSize    string
	biggestSize string
	imageSizes  map[string]string
}
