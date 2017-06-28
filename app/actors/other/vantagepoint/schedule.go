package vantagepoint

import (
	"fmt"
	"time"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)


var checkNewUploadsSchedule env.InterfaceSchedule
var hoursList = map[string]string{}

type postponedFunc func() error
var postponedScheduleChanges = map[string]postponedFunc{}

// initHoursList prepares hour selector values
func initHoursList() {
	hoursList = map[string]string{}

	for hour := 0; hour < 24; hour++ {
		hourStr := utils.InterfaceToString(hour)
		hoursList[hourStr] = hourStr
	}
}

// scheduleCheckNewUploads registers and schedules CheckNewUploads execution
func scheduleCheckNewUploads() error {
	var err error

	scheduler := env.GetScheduler()
	if scheduler == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "905f6a45-e733-42a9-8699-00a27762c044", "scheduler is not registered")
	}

	if err := scheduler.RegisterTask(ConstSchedulerTaskName, runCheckNewUploadsSchedule); err != nil {
		return env.ErrorDispatch(err)
	}

	if checkNewUploadsSchedule == nil {
		checkNewUploadsSchedule, err = scheduler.ScheduleRepeat("0 * * * *", ConstSchedulerTaskName, nil)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		for _, worker := range(postponedScheduleChanges) {
			if err := worker(); err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}

	return nil
}

// runCheckNewUploadsSchedule checks if task could be executed and run it.
// The check is outside of real executor to allow call CheckNewUploads independently of schedule.
func runCheckNewUploadsSchedule(params map[string]interface{}) error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e53354b3-1728-4433-a8f5-aa32b6013b24", "can't obtain config")
	}

	if !utils.InterfaceToBool(config.GetValue(ConstConfigPathVantagePointScheduleEnabled)) {
		env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "scheduled update disabled")
		return nil
	}

	return CheckNewUploads(params)
}

// setScheduleHour updates scheduled task execution time
func setScheduleHour(hour interface{}) error {
	var worker = func() error {
		currentExpr := utils.InterfaceToString(checkNewUploadsSchedule.Get("expr"))
		newExpr := "0 " + utils.InterfaceToString(hour) + " * * *"
		if currentExpr != newExpr {
			env.Log(ConstLogStorage, env.ConstLogPrefixInfo, fmt.Sprintf("change schedule cronExpr from [%s] to [%s]", currentExpr, newExpr))

			if err := checkNewUploadsSchedule.Set("expr", newExpr); err != nil {
				return env.ErrorDispatch(err)
			}

			if err := checkNewUploadsSchedule.Set("time", time.Now()); err != nil {
				return env.ErrorDispatch(err)
			}
		}

		return nil
	}

	if checkNewUploadsSchedule != nil {
		return worker()
	}

	postponedScheduleChanges["hour"] = worker

	return nil
}
