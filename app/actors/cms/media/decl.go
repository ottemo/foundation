// Package media is a default implementation of cms page related interfaces declared in
// "github.com/ottemo/foundation/app/models/cms" package
package media

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
)

// Package global constants
const (
	ConstErrorModule = "cms/media"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstStorageModel  = "cms"
	ConstStorageType   = "image"
	ConstStorageObject = "media"
)

var (
	mediaStorage media.InterfaceMediaStorage
)
