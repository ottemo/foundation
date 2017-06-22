package vantagepoint

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app"
)

func init() {
	app.OnAppStart(onAppStart)
	env.RegisterOnConfigStart(setupConfig)
}

func onAppStart() error {
	if err := CheckNewUploads(make(map[string]interface{})); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}


