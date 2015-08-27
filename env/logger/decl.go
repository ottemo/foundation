package logger

// Package global variables
var (
	baseDirectory     = "./var/log/" // folder location where to store logs
	defaultLogFile    = "system.log" // filename for default log
	defaultErrorsFile = "errors.log" // filename for errors log
)

// DefaultLogger is a default implementer of InterfaceLogger
type DefaultLogger struct{}
