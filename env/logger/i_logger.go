package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ottemo/foundation/api/context"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
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

// LogMap allows to write mapped data into log file
func (it *DefaultLogger) LogMap(storage string, data map[string]interface{}) {

	if storage == "" {
		storage = defaultLogFile
	}

	logData := map[string]interface{}{
		"@version":   "1",
		"@timestamp": time.Now().Format(time.RFC3339),
	}

	for key, value := range data {
		logData[key] = value
	}

	message := "\n"
	jsonMessage, err := json.Marshal(logData)
	if err == nil {
		message = string(jsonMessage) + message
	} else {
		message = utils.HardEncodeToJSONString(logData) + message
	}

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
		logToJSON := false
		if ottemoErr, ok := err.(env.InterfaceOttemoError); ok {
			if ottemoErr.ErrorLevel() <= errorLogLevel && !ottemoErr.IsLogged() {
				it.Log(defaultErrorsFile, env.ConstLogPrefixError, ottemoErr.ErrorFull())
				ottemoErr.MarkLogged()
				logToJSON = true
			}
		} else {
			it.Log(defaultErrorsFile, env.ConstLogPrefixError, err.Error())
			logToJSON = true
		}

		if logToJSON {
			stackContext := make(map[string]interface{})
			for key, value := range context.GetContext() {
				stackContext[key] = value
			}
			it.LogMap(defaultJSONErrorsFile, handleErrorToJSONLog(err, stackContext))
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
