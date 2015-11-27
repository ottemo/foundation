package jsonlogger

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstUseDebugLog     = true       // flag to use full data logging
	ConstDebugLogStorage = "json.log" // log storage for debug log records

	ConstErrorModule = "api/jsonlogger"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// Package global variables
var (
	baseDirectory = "./var/log/" // folder location where to store logs

	errorLogLevel = 5
)

// DefaultJSONLogger is a structure to hold related information
type DefaultJSONLogger struct {
	Data interface{}
}
