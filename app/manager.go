package app

import (
	"github.com/ottemo/foundation/rest_service"
)

var callbacksOnAppStart = []func() error{}

func OnAppStart(callback func() error) {
	callbacksOnAppStart = append(callbacksOnAppStart, callback)
}

func Start() error {
	for _, callback := range callbacksOnAppStart {
		if err := callback(); err != nil {
			return err
		}
	}

	return nil
}

func Serve() error {
	return rest_service.GetRestService().Run()
}
