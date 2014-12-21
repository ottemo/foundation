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
	globalCustomAttributes      = map[string]map[string]models.StructAttributeInfo{}
	globalCustomAttributesMutex sync.RWMutex
)

// CustomAttributes is a implementer of InterfaceCustomAttributes
type CustomAttributes struct {
	model      string
	collection string

	attributes map[string]models.StructAttributeInfo

	values map[string]interface{}
}
