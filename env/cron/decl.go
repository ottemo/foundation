package cron

import (
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstErrorModule = "env/cron"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// DefaultCronScheduler is a default implementer of InterfaceIniConfig
type DefaultCronScheduler struct {
	tasks     map[string]env.FuncCronTask
	schedules []*DefaultCronSchedule

	appStarted bool
}

// DefaultCronSchedule structure to hold schedule information (for internal usage)
type DefaultCronSchedule struct {
	CronExpr string
	TaskName string
	Params   map[string]interface{}
	Repeat   bool
	Time     time.Time
	active   bool

	task env.FuncCronTask
	expr *cronexpr.Expression

	scheduler *DefaultCronScheduler
}
