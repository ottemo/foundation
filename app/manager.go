package app

import (
	"github.com/ottemo/foundation/rest_service"
)

var callbacksOnAppStart = []func() error{}

// OnAppStart is a place to register callbacks upon application initialization
func OnAppStart(callback func() error) {
	callbacksOnAppStart = append(callbacksOnAppStart, callback)
}

// Start executes the registered callback chain when Foundation Server is first started.
func Start() error {
	for _, callback := range callbacksOnAppStart {
		if err := callback(); err != nil {
			return err
		}
	}

	return nil
}

// Serve starts and returns the REST service
func Serve() error {
	return rest_service.GetRestService().Run()
}
