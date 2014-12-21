package api

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	currentRestService          InterfaceRestService // currently registered RESTFul service in system
	callbacksOnRestServiceStart = []func() error{}   // set of callback function on RESTFul service start
)

// RegisterOnRestServiceStart registers new callback on RESTFul service start
func RegisterOnRestServiceStart(callback func() error) {
	callbacksOnRestServiceStart = append(callbacksOnRestServiceStart, callback)
}

// OnRestServiceStart fires RESTFul service start event (callback handling)
func OnRestServiceStart() error {
	for _, callback := range callbacksOnRestServiceStart {
		if err := callback(); err != nil {
			return env.ErrorDispatch(err)
		}
	}
	return nil
}

// RegisterRestService registers RESTFul service in the system
//   - will cause error if there are couple candidates for that role
func RegisterRestService(newService InterfaceRestService) error {
	if currentRestService == nil {
		currentRestService = newService
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9aecae0349854682b42d145da67ba74b", "Sorry, '"+currentRestService.GetName()+"' already registered")
	}
	return nil
}

// GetRestService returns currently used RESTFul service implementation
func GetRestService() InterfaceRestService {
	return currentRestService
}
