// Package fsmedia is a default implementation of InterfaceMediaStorage declared in
// "github.com/ottemo/foundation/media" package
package fsmedia

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstMediaDBCollection  = "media"    // database collection name to store media assignment information into
	ConstMediaDefaultFolder = "./media/" // filesystem folder path to store media files in there

	ConstDefaultImageSize  = "800x400"    // "800x400"
	ConstDefaultImageSizes = "large: 0x0" // "small: 75x75, thumb: 260x300, big: 560x650"

	ConstConfigPathMediaImageSize  = "general.app.image_size"  // base image size
	ConstConfigPathMediaImageSizes = "general.app.image_sizes" // other image sizes required

	ConstMediaTypeImage    = "image"
	ConstMediaTypeLink     = "link"
	ConstMediaTypeDocument = "document"

	ConstErrorModule = "media/fsmedia"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// FilesystemMediaStorage is a filesystem based implementer of InterfaceMediaStorage
type FilesystemMediaStorage struct {
	storageFolder string
	setupWaitCnt  int

	baseSize    string
	biggestSize string
	imageSizes  map[string]string
}
