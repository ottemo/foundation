package api

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstRESTOperationGet    = "GET"
	ConstRESTOperationUpdate = "PUT"
	ConstRESTOperationCreate = "POST"
	ConstRESTOperationDelete = "DELETE"
	ConstRESTActionParameter = "action"

	ConstSessionKeyAdminRights = "adminRights"   // session key used to flag that user have admin rights
	ConstSessionCookieName     = "OTTEMOSESSION" // cookie name which should contain sessionID
	ConstSessionKeyTimeZone    = "timeZone"      // session key for setting time zone

	ConstGETAuthParamName            = "auth"
	ConstConfigPathStoreRootLogin    = "general.store.root_login"
	ConstConfigPathStoreRootPassword = "general.store.root_password"

	ConstErrorModule = "api"
	ConstErrorLevel  = env.ConstErrorLevelHelper
)

// ApplicationContext is a type you can embed in your model for application context support
type ApplicationContext struct{ InterfaceApplicationContext }

// GetApplicationContext returns current application context or nil
func (it *ApplicationContext) GetApplicationContext() InterfaceApplicationContext {
	return it.InterfaceApplicationContext
}

// SetApplicationContext assigns given application context to type
func (it *ApplicationContext) SetApplicationContext(context InterfaceApplicationContext) error {
	it.InterfaceApplicationContext = context
	return nil
}

// StructRestRedirect is a structure you should return in API handler function if redirect needed
type StructRestRedirect struct {
	Result   interface{}
	Location string

	DoRedirect bool
}

// FuncAPIHandler is an API handler callback function
type FuncAPIHandler func(context InterfaceApplicationContext) (interface{}, error)
