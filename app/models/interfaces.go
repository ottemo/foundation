// Package models represents abstraction of business layer object and basic access interfaces for it
package models

import (
	"github.com/ottemo/foundation/db"
)

// InterfaceModel represents interface for basic business layer implementation object
type InterfaceModel interface {
	GetModelName() string
	GetImplementationName() string
	New() (InterfaceModel, error)
}

// InterfaceStorable represents interface load/store business layer implementation object from database
type InterfaceStorable interface {
	GetID() string
	SetID(string) error

	Save() error
	Load(id string) error
	Delete() error
}

// InterfaceObject represents interface to access business layer implementation object via get/set functions
type InterfaceObject interface {
	Get(attribute string) interface{}
	Set(attribute string, value interface{}) error

	FromHashMap(hashMap map[string]interface{}) error
	ToHashMap() map[string]interface{}

	GetAttributesInfo() []StructAttributeInfo
}

// InterfaceListable represents interface to access business layer implementation collection via object instance
type InterfaceListable interface {
	GetCollection() InterfaceCollection
}

// InterfaceCustomAttributes represents interface to access business layer implementation object custom attributes
type InterfaceCustomAttributes interface {
	GetCustomAttributeCollectionName() string

	AddNewAttribute(newAttribute StructAttributeInfo) error
	RemoveAttribute(attributeName string) error
	EditAttribute(attributeName string, attributeValues StructAttributeInfo) error
}

// InterfaceExternalAttributes represents interface to access business layer implementation object external attributes
type InterfaceExternalAttributes interface {
	AddAttribute(newAttribute StructAttributeInfo, delegate interface{}) error
	RemoveAttribute(attributeName string) error
}

// InterfaceMedia represents interface to access business layer implementation object assigned media resources
type InterfaceMedia interface {
	AddMedia(mediaType string, mediaName string, content []byte) error
	RemoveMedia(mediaType string, mediaName string) error

	ListMedia(mediaType string) ([]string, error)

	GetMedia(mediaType string, mediaName string) ([]byte, error)
	GetMediaPath(mediaType string) (string, error)
}

// InterfaceCollection represents interface to access business layer implementation collection
type InterfaceCollection interface {
	GetDBCollection() db.InterfaceDBCollection

	List() ([]StructListItem, error)

	ListAddExtraAttribute(attribute string) error

	ListFilterAdd(attribute string, operator string, value interface{}) error
	ListFilterReset() error

	ListLimit(offset int, limit int) error
}
