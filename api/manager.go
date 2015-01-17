package api

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	currentRestService          InterfaceRestService    // currently registered RESTFul service in system
	currentSessionService       InterfaceSessionService // currently registered session service in system
	callbacksOnRestServiceStart = []func() error{}      // set of callback function on RESTFul service start
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
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9aecae03-4985-4682-b42d-145da67ba74b", "REST service '"+currentRestService.GetName()+"' was already registered")
	}
	return nil
}

// RegisterSessionService registers session managing service in the system
//   - will cause error if there are couple candidates for that role
func RegisterSessionService(newService InterfaceSessionService) error {
	if currentSessionService == nil {
		currentSessionService = newService
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2a0e046c-75bd-4460-b3c6-843ff761bd59", "Session service '"+currentSessionService.GetName()+"' was already registered")
	}
	return nil
}

// GetRestService returns currently using RESTFul service implementation
func GetRestService() InterfaceRestService {
	return currentRestService
}

// GetSessionService returns currently using session service implementation
func GetSessionService() InterfaceSessionService {
	return currentSessionService
}
