// Package gallery is a default implementation of cms page related interfaces declared in
// "github.com/ottemo/foundation/app/models/cms" package
package gallery

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
)

// Package global constants
const (
	ConstErrorModule = "cms/gallery"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstStorageModel  = "cms"
	ConstStorageType   = "image"
	ConstStorageObject = "gallery"
)

var (
	mediaStorage media.InterfaceMediaStorage
)
