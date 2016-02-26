package rest

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstUseDebugLog         = true           // flag to use REST API logging
	ConstDebugLogStorage     = "rest.log"     // log storage for debug log records
	ConstDebugJSONLogStorage = "restJSON.log" // log storage for debug log records

	ConstErrorModule = "api/rest"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// DefaultRestService is a default implementer of InterfaceRestService
// declared in "github.com/ottemo/foundation/api" package
type DefaultRestService struct {
	ListenOn string
	Router   *httprouter.Router
	Handlers []string
}

// DefaultRestApplicationContext is a structure to hold API request related information
type DefaultRestApplicationContext struct {
	ResponseWriter    http.ResponseWriter
	Request           *http.Request
	RequestParameters map[string]string
	RequestArguments  map[string]string
	RequestContent    interface{}
	RequestFiles      map[string]io.Reader

	Session       api.InterfaceSession
	ContextValues map[string]interface{}
	Result        interface{}
}
