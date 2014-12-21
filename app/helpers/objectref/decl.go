// Package objectref intended to unify and simplify a way of model instance changes tracking (currently not implemented)
package objectref

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstErrorModule = "objectref"
	ConstErrorLevel  = env.ConstErrorLevelHelper
)

// DBObjectRef is a object state tracking helper and implementer of InterfaceObject and InterfaceStorable
type DBObjectRef struct {
	id string

	loaded   bool
	modified bool

	origData map[string]interface{}
	currData map[string]interface{}
}

// MarkAsLoaded marks object instance as loaded from DB
func (it *DBObjectRef) MarkAsLoaded() {
	it.loaded = true
}

// MarkAsModified marks object instance as modified
func (it *DBObjectRef) MarkAsModified() {
	it.modified = true
}

// IsModified returns value of modification flag
func (it *DBObjectRef) IsModified() bool {
	return it.modified
}

// IsLoaded returns value of load from DB flag
func (it *DBObjectRef) IsLoaded() bool {
	return it.loaded
}
