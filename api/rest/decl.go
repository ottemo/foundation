package rest

import (
	"io"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/env"
)

// Package global constants
const (
	ConstDebugLogStorage = "rest.log" // log storage for debug log records

	ConstErrorModule = "api/rest"
	ConstErrorLevel  = env.ConstErrorLevelService

	ConstConfigPathAPI           = "api"
	ConstConfigPathAPILog        = "api.log"
	ConstConfigPathAPILogEnable  = "api.log.enable"
	ConstConfigPathAPILogExclude = "api.log.exclude"
)

// DefaultRestService is a default implementer of InterfaceRestService
// declared in "github.com/ottemo/commerce/api" package
type DefaultRestService struct {
	ListenOn string
	Router   *httprouter.Router
	Handlers []string

	RawHandler api.FuncAPIHandler
	RawHandlerMutex sync.RWMutex
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
