// Package swatch is a default implementation of product swatches

package swatch

import (
	"github.com/ottemo/foundation/media"
)

// Package global constants
const (
	ConstErrorModule = "swatch"

	ConstStorageModel     = "swatch"
	ConstStorageObjectID  = "media"
	ConstStorageMediaType = "image"
)

var (
	mediaStorage media.InterfaceMediaStorage
)
