package vantagepoint

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/db"
)

func init() {
	env.RegisterOnConfigStart(setupConfig)

	db.RegisterOnDatabaseStart(onDatabaseStart)
}

func onDatabaseStart() error {
	app.OnAppStart(onAppStart)

	return nil
}

func onAppStart() error {
	if err := CheckNewUploads(make(map[string]interface{})); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
