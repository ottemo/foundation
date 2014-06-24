package app

import "github.com/ottemo/foundation/api"

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

// Serve starts and returns the REST Endpoint
func Serve() error {
	return api.GetEndPoint().Run()
}
