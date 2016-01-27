package cron

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("cron/schedule", getSchedule)
	service.POST("cron/task", createTask)
	service.GET("cron/task", getTasks)
	service.GET("cron/task/enable/:taskIndex", enableTask)
	service.GET("cron/task/disable/:taskIndex", disableTask)
	service.PUT("cron/task/:taskIndex", updateTask)
	service.GET("cron/task/run/:taskIndex", runTask)

	return nil
}

// runTask - allows to execute task of schedule without updating of it
// taskIndex - need to be specified in request argument
func runTask(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reqTaskIndex := context.GetRequestArgument("taskIndex")
	if reqTaskIndex == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "41baa8b5-eea1-4a31-aad6-83aceb56ed2f", "task index should be specified")
	}

	taskIndex, err := utils.StringToInteger(reqTaskIndex)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	scheduler := env.GetScheduler()
	currentSchedules := scheduler.ListSchedules()

	if taskIndex > len(currentSchedules)-1 || taskIndex < 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "0a50509a-c638-49ab-bb52-a7b9750097b2", "task index is out of range for existing tasks")
	}

	useTaskParams := utils.InterfaceToBool(context.GetRequestArgument("useTaskParams"))

	var params map[string]interface{}
	if !useTaskParams {
		params = make(map[string]interface{})
	}

	for index, schedule := range currentSchedules {
		if index == taskIndex {
			err := schedule.RunTask(params)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
			break
		}
	}

	return "ok", nil
}

// getSchedule to get information about current schedules
func getSchedule(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	var result []interface{}

	scheduler := env.GetScheduler()

	for _, schedule := range scheduler.ListSchedules() {
		result = append(result, schedule.GetInfo())
	}

	return result, nil
}

// getTasks return scheduler registered tasks (functions that are available to execute)
func getTasks(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	scheduler := env.GetScheduler()

	return scheduler.ListTasks(), nil
}

// updateTask update scheduler task
//   - "taskIndex" should be specified as argument (task index can be obtained from getSchedules)
func updateTask(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reqTaskIndex := context.GetRequestArgument("taskIndex")
	if reqTaskIndex == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d4ee4c0c-124a-4098-aeef-23d868b0d682", "task index should be specified")
	}

	taskIndex, err := utils.StringToInteger(reqTaskIndex)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	scheduler := env.GetScheduler()
	scheduleParams := []string{"expr", "task", "repeat", "params", "time"}

	for index, schedule := range scheduler.ListSchedules() {
		if index == taskIndex {
			for _, param := range scheduleParams {
				if value, present := postValues[param]; present {
					err = schedule.Set(param, value)
					if err != nil {
						return nil, env.ErrorDispatch(err)
					}
				}
			}
		}
	}

	return "ok", nil
}

// createTask with request params
// in request params required are time or cronExpr for creating different type of tasks
func createTask(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cronExpression := utils.InterfaceToString(utils.GetFirstMapValue(postValues, "CronExpr", "cronExpr", "Expr", "expr"))
	scheduledTime := utils.InterfaceToTime(utils.GetFirstMapValue(postValues, "Time", "time"))

	if utils.IsZeroTime(scheduledTime) && cronExpression == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2250e0a5-e439-444d-b331-13d16f8e2401", "cronExpr or time were not specified")
	}

	scheduler := env.GetScheduler()

	taskName := utils.InterfaceToString(utils.GetFirstMapValue(postValues, "TaskName", "Task", "task"))

	if taskName == "" || !utils.IsInListStr(taskName, scheduler.ListTasks()) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "eafc6b15-e897-4d9f-a93a-f84cffa78497", "task not specified or not regisetered")
	}

	isRepeat := false
	if repeatValue, present := postValues["repeat"]; present {
		isRepeat = utils.InterfaceToBool(repeatValue)
	}

	taskParams := make(map[string]interface{})
	if params, present := postValues["params"]; present {
		taskParams = utils.InterfaceToMap(params)
	}

	var newSchedule env.InterfaceSchedule
	if !utils.IsZeroTime(scheduledTime) {
		newSchedule, err = scheduler.ScheduleAtTime(scheduledTime, taskName, taskParams)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	} else {
		if isRepeat {
			newSchedule, err = scheduler.ScheduleRepeat(cronExpression, taskName, taskParams)
		} else {
			newSchedule, err = scheduler.ScheduleOnce(cronExpression, taskName, taskParams)
		}

		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	return newSchedule, nil
}

// enableTask make schedule active
// taskIndex - need to be specified in request argument
func enableTask(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reqTaskIndex := context.GetRequestArgument("taskIndex")
	if reqTaskIndex == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d4ee4c0c-124a-4098-aeef-23d868b0d682", "task index should be specified")
	}

	taskIndex, err := utils.StringToInteger(reqTaskIndex)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	scheduler := env.GetScheduler()
	currentSchedules := scheduler.ListSchedules()

	if taskIndex > len(currentSchedules)-1 || taskIndex < 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5cf9ead0-d23d-4cb6-87b1-3578d54dd1ba", "task index is out of range for existing tasks")
	}

	for index, schedule := range currentSchedules {
		if index == taskIndex {
			err := schedule.Enable()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
			break
		}
	}

	return currentSchedules[taskIndex].GetInfo(), nil
}

// disableTask make schedule inactive
// taskIndex - need to be specified in request argument
func disableTask(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reqTaskIndex := context.GetRequestArgument("taskIndex")
	if reqTaskIndex == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "61285b1f-6c6c-4920-b1b1-5d4d31b58ad5", "task index should be specified")
	}

	taskIndex, err := utils.StringToInteger(reqTaskIndex)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	scheduler := env.GetScheduler()
	currentSchedules := scheduler.ListSchedules()

	if taskIndex > len(currentSchedules)-1 || taskIndex < 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6465b989-09f9-41a3-9752-9844fddfdc4a", "task index is out of range for existing tasks")
	}

	for index, schedule := range currentSchedules {
		if index == taskIndex {
			err := schedule.Disable()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
			break
		}
	}

	return currentSchedules[taskIndex].GetInfo(), nil
}
