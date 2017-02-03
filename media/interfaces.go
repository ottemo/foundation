package media

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstErrorModule = "media"
	ConstErrorLevel  = env.ConstErrorLevelService

	ConstMediaTypeImage    = "image"
	ConstMediaTypeLink     = "link"
	ConstMediaTypeDocument = "document"
)

// InterfaceMediaStorage is an interface to access media storage service
type InterfaceMediaStorage interface {
	GetName() string

	Load(model string, objID string, mediaType string, mediaName string) ([]byte, error)
	Save(model string, objID string, mediaType string, mediaName string, mediaData []byte) error

	Remove(model string, objID string, mediaType string, mediaName string) error

	ListMedia(model string, objID string, mediaType string) ([]string, error)
	ListMediaDetail(model string, objID string, mediaType string) ([]map[string]interface{}, error)

	GetMediaPath(model string, objID string, mediaType string) (string, error)

	GetAllSizes(model string, objID string, mediaType string) ([]map[string]string, error)
	GetSizes(model string, objID string, mediaType string, mediaName string) (map[string]string, error)

	ResizeAllMediaImages() error
}
