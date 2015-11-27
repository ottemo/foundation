package jsonlogger

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultJSONLogger)

	api.RegisterJSONLogger(instance)
	app.OnAppStart(onAppStart)
}

// onAppStart makes module initialization on application startup
func onAppStart() error {
	// env.EventRegisterListener("api.request", requestHandler)
	env.EventRegisterListener("api.response", requestHandler)
	return nil
}
