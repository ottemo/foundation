package models

import (
	"github.com/ottemo/foundation/db"
)

type I_Model interface {
	GetModelName() string
	GetImplementationName() string
	New() (I_Model, error)
}

type I_Storable interface {
	GetID() string
	SetID(string) error

	Save() error
	Load(id string) error
	Delete() error
}

type I_Object interface {
	Get(attribute string) interface{}
	Set(attribute string, value interface{}) error

	FromHashMap(hashMap map[string]interface{}) error
	ToHashMap() map[string]interface{}

	GetAttributesInfo() []T_AttributeInfo
}

type I_Listable interface {
	GetCollection() I_Collection
}

type I_CustomAttributes interface {
	GetCustomAttributeCollectionName() string

	AddNewAttribute(newAttribute T_AttributeInfo) error
	RemoveAttribute(attributeName string) error
}

type I_Media interface {
	AddMedia(mediaType string, mediaName string, content []byte) error
	RemoveMedia(mediaType string, mediaName string) error

	ListMedia(mediaType string) ([]string, error)

	GetMedia(mediaType string, mediaName string) ([]byte, error)
	GetMediaPath(mediaType string) (string, error)
}

type I_Collection interface {
	GetDBCollection() db.I_DBCollection

	List() ([]T_ListItem, error)

	ListAddExtraAttribute(attribute string) error

	ListFilterAdd(attribute string, operator string, value interface{}) error
	ListFilterReset() error

	ListLimit(offset int, limit int) error
}
