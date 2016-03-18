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
	modelCustomAttributes   = make(map[string]map[string]models.StructAttributeInfo)

	// modelExternalAttributes is a per model external attribute information (map[model][attribute] => info)
	modelExternalAttributes = make(map[string]map[string]models.StructAttributeInfo)

	// modelExternalDelegates is a per model attribute delegate mapping (map[model][attribute] => delegate)
	modelExternalDelegates  = make(map[string]map[string]interface{})

	// the mutexes to synchronize access on global variables
	modelCustomAttributesMutex   sync.Mutex
	modelExternalAttributesMutex sync.Mutex
	modelExternalDelegatesMutex  sync.Mutex
)

// ModelCustomAttributes type represents a set of attributes which could be modified (edited/added/removed) dynamically.
// The implementation relays on InterfaceCollection which is used to store values and have ability to add/remove
// columns on a fly.
type ModelCustomAttributes struct {
	model      string
	collection string
	instance   interface{}
	values     map[string]interface{}
}

// ModelExternalAttributes type represents a set of attributes managed by "external" package (outside of model package)
// which supposing at least InerfaceObject methods delegation, but also could delegate InterfaceStorable if the methods
// are implemented in delegate instance.
type ModelExternalAttributes struct {
	model     string
	instance  interface{}
}
