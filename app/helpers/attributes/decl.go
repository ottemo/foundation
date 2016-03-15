// Package attributes represents an implementation of InterfaceCustomAttributes declared in
// "github.com/ottemo/foundation/app/models" package.
//
// In order to use it you should just embed CustomAttributes in your actor,
// you can found sample usage in "github.com/app/actors/product" package.
package attributes

import (
	"sync"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstCollectionNameCustomAttributes = "custom_attributes"

	ConstErrorModule = "attributes"
	ConstErrorLevel  = env.ConstErrorLevelHelper
)

// Package global variables
var (
	modelCustomAttributes      = map[string]map[string]models.StructAttributeInfo{}
	modelCustomAttributesMutex sync.Mutex
)

// CustomAttributes implements InterfaceCustomAttributes
//
// CustomAttributes type represents a set of attributes which could be modified (edited/added/removed) dynamically.
// The implementation relays on InterfaceCollection which is used to store values and have ability to add/remove
// columns on a fly.
type CustomAttributes struct {
	model      string
	collection string

	info   map[string]models.StructAttributeInfo
	values map[string]interface{}
}

// ExternalAttributes implements InterfaceExternalAttributes
//
// ExternalAttributes type represents a set of object attributes managed by "external" package (outside of implementor)
// which supposing setters/getters delegation routines handled by this type.
type ExternalAttributes struct {
	model  string
	info   map[string]models.StructAttributeInfo
	values map[string]interface{}
}
