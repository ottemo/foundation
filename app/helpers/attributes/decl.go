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
	modelCustomAttributes = make(map[string]map[string]models.StructAttributeInfo)

	// modelExternalAttributes is a per model external attribute information (map[model][attribute] => delegate)
	modelExternalAttributes = make(map[string]map[string]models.InterfaceAttributesDelegate)

	// the mutexes to synchronize access on global variables
	modelCustomAttributesMutex   sync.Mutex
	modelExternalAttributesMutex sync.Mutex
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
// which supposing at least InterfaceObject methods delegation, but also could delegate InterfaceStorable if the methods
// are implemented in delegate instance.
//
// Workflow diagram:
//                    helper makes proxy for Object/Storable interface methods:
//                                 Get(), Set(), ListAttributes(),
//           +---------------+       FromHashMap(), ToHashMap(),       +------------------+
//           | Model Package |    GetId(), Load(), Sate(), Delete()    | External package |
//           |               |                                         |                  |
//           |   +-------+   |                                         |   +----------+   |
//     +---------> Model <---------------------+ +-------------------------> Delegate |   |
//     |     |   +---+---+   |   model-helper  | |   helper-delegate   |   +----^-----+   |
//     |     |       |       |    proxy call   | |      proxy call     |        |         |
//     |     +-------|-------+                 | |                     +--------|---+-----+
//     |             |                         | |                              |   |
//     |             |   instance    +---------v-v--------------+               |   |
//     |             +---------------> *ModelExternalAttributes ----------------+   |
//     |   Model.New() instantiates  +-- -+---------------------+   Delegate.New()  |  Registering delegate
//     |   ExternalAttributes helper      |                         instantiates    |  for a specific model
//     |   which transfers instance       |                           delegate      |
//     +----------------------------------+-GetInstance()                           |
//        helper have reference to model  |                                         |
//                                        +-AddExternalAttributes() <---------------+
//                                        +-RemoveExternalAttributes()
//                                        |
//                                        +-ListExternalAttributes()
type ModelExternalAttributes struct {
	model     string
	instance  interface{}
	delegates map[string]models.InterfaceAttributesDelegate
}
