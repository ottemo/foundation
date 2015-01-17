// Package api is a set of interfaces representing API endpoint services.
//
// (currently only "InterfaceRestService" endpoint interface supported)
package api

import (
	"github.com/ottemo/foundation/env"
	"net/http"
)

// Package global constants
const (
	ConstSessionKeyAdminRights = "adminRights"   // session key used to flag that user have admin rights
	ConstSessionCookieName     = "OTTEMOSESSION" // cookie name which should contain sessionID

	ConstGETAuthParamName            = "auth"
	ConstConfigPathStoreRootLogin    = "general.store.root_login"
	ConstConfigPathStoreRootPassword = "general.store.root_password"

	ConstErrorModule = "api"
	ConstErrorLevel  = env.ConstErrorLevelHelper
)

// StructAPIHandlerParams is a structure to hold API request related information
type StructAPIHandlerParams struct {
	ResponseWriter   http.ResponseWriter
	Request          *http.Request
	RequestGETParams map[string]string
	RequestURLParams map[string]string
	RequestContent   interface{}
	Session          InterfaceSession
}

// StructRestRedirect is a structure you should return in API handler function if redirect needed
type StructRestRedirect struct {
	Result   interface{}
	Location string

	DoRedirect bool
}

// FuncAPIHandler is an API handler callback function
type FuncAPIHandler func(params *StructAPIHandlerParams) (interface{}, error)
