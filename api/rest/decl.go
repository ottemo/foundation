// Package rest is a default implementation of InterfaceRestService
// declared in "github.com/ottemo/foundation/api" package
package rest

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstUseDebugLog     = true       // flag to use REST API logging
	ConstDebugLogStorage = "rest.txt" // log storage for debug log records

	ConstErrorModule = "api/rest"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// DefaultRestService is a default implementer of InterfaceRestService declared in
// "github.com/ottemo/foundation/api" package
type DefaultRestService struct {
	ListenOn string
	Router   *httprouter.Router
	Handlers map[string]httprouter.Handle
}
