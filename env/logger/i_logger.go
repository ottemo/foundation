package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/ottemo/foundation/env"
)

// Log is a general case logging function
func (it *DefaultLogger) Log(storage string, prefix string, msg string) {
	message := time.Now().Format(time.RFC3339) + " [" + prefix + "]: " + msg + "\n"

	logFile, err := os.OpenFile(baseDirectory+storage, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(message)
		return
	}

	logFile.Write([]byte(message))

	logFile.Close()
}

// LogError makes error log
func (it *DefaultLogger) LogError(err error) {
	if err != nil {
		if ottemoErr, ok := err.(env.InterfaceOttemoError); ok {
			if ottemoErr.ErrorLevel() <= errorLogLevel && !ottemoErr.IsLogged() {
				it.Log(defaultErrorsFile, env.ConstLogPrefixError, ottemoErr.ErrorFull())
				ottemoErr.MarkLogged()
			}
		} else {
			it.Log(defaultErrorsFile, env.ConstLogPrefixError, err.Error())
		}
	}
}

// LogToStorage logs info type message to specific storage
func (it *DefaultLogger) LogToStorage(storage string, msg string) {
	it.Log(storage, env.ConstLogPrefixInfo, msg)
}

// LogWithPrefix logs prefixed message to default storage
func (it *DefaultLogger) LogWithPrefix(prefix string, msg string) {
	it.Log(defaultLogFile, prefix, msg)
}

// LogMessage logs info message to default storage
func (it *DefaultLogger) LogMessage(msg string) {
	it.Log(defaultLogFile, env.ConstLogPrefixInfo, msg)
}
