// Package attributes represents an implementation of InterfaceCustomAttributes declared in
// "github.com/ottemo/foundation/app/models" package.
//
// In order to use it you should just embed ModelCustomAttributes in your actor,
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
	// modelCustomAttributes is a per model attribute information storage (map[model][attribute])
	modelCustomAttributes      = make(map[string]map[string]models.StructAttributeInfo)
	// modelCustomAttributesMutex is a synchronization for modelCustomAttributes
	modelCustomAttributesMutex sync.Mutex

	// modelExternalAttributes is a per model attribute information storage (map[model][attribute])
	modelExternalAttributes      = make(map[string]*ModelExternalAttributes)
	// modelExternalAttributesMutex is a synchronization for modelCustomAttributes
	modelExternalAttributesMutex sync.Mutex
)

// ModelCustomAttributes type represents a set of attributes which could be modified (edited/added/removed) dynamically.
// The implementation relays on InterfaceCollection which is used to store values and have ability to add/remove
// columns on a fly.
type ModelCustomAttributes struct {
	model      string
	collection string

	info      map[string]models.StructAttributeInfo
	infoMutex sync.Mutex

	instance   interface{}
	values map[string]interface{}
}

// ModelExternalAttributes type represents a set of attributes managed by "external" package (outside of model package)
// which supposing at least InerfaceObject methods delegation, but also could delegate InterfaceStorable if the methods
// are implemented in delegate instance.
type ModelExternalAttributes struct {
	model     string

	// the info is a shared field, so it must by synchronized
	info  map[string]models.StructAttributeInfo
	mutex sync.Mutex

	instance  interface{}
	delegates map[string]interface{}
}
