package logger

import "github.com/ottemo/foundation/env"

// Package global constants
const (
	ConstCollectCallStack = true // flag to indicate that call stack information within error is required

	ConstConfigPathError         = "general.error"
	ConstConfigPathErrorLogLevel = "general.error.log_level"

	ConstErrorModule = "env/logger"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// Package global variables
var (
	baseDirectory     = "./var/log/" // folder location where to store logs
	defaultLogFile    = "system.log" // filename for default log
	defaultErrorsFile = "errors.log" // filename for errors log

	errorLogLevel = 5
)

// DefaultLogger is a default implementer of InterfaceLogger
type DefaultLogger struct{}
