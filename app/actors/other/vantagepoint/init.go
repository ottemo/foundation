package vantagepoint

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

func init() {
	initHoursList()

	db.RegisterOnDatabaseStart(onDatabaseStart)

	env.RegisterOnConfigStart(setupConfig)
}

func onDatabaseStart() error {
	app.OnAppStart(onAppStart)

	return nil
}

func onAppStart() error {
	if err := scheduleCheckNewUploads(); err != nil {
		_ = env.ErrorDispatch(err)
	}

	return nil
}
