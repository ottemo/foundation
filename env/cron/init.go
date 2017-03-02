package cron

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultCronScheduler)
	var _ env.InterfaceScheduler = instance
	var _ env.InterfaceSchedule = new(DefaultCronSchedule)

	instance.tasks = make(map[string]env.FuncCronTask)
	instance.schedules = make([]*DefaultCronSchedule, 0)

	app.OnAppInit(instance.appInitEvent)
	app.OnAppEnd(instance.appEndEvent)
	api.RegisterOnRestServiceStart(setupAPI)

	if err := env.RegisterScheduler(instance); err != nil {
		_ = env.ErrorDispatch(err)
	}
}

// routines before application end
func (it *DefaultCronScheduler) appEndEvent() error {
	return nil
}

// routines before application start (on init phase)
func (it *DefaultCronScheduler) appInitEvent() error {

	// TODO: load manually specified tasks from DB

	it.appStarted = true

	return nil
}
