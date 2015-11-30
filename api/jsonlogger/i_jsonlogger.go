package jsonlogger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// GetName returns implementation name of JSON logger service
func (it *DefaultJSONLogger) GetName() string {
	return "jsonlogger"
}

// Log allows to write mapped data into log file
func (it *DefaultJSONLogger) Log(storage string, data map[string]interface{}) error {

	if storage == "" {
		storage = defaultLogFile
	}

	logData := map[string]interface{}{
		"@version":   "1",
		"@timestamp": time.Now().Format(time.RFC3339),
		"type":       "foundation-error",
	}

	for key, value := range data {
		logData[key] = value
	}

	jsonMessage, err := json.Marshal(logData)
	if err != nil {
		return err
	}

	message := string(jsonMessage) + "\n"

	logFile, err := os.OpenFile(baseDirectory+storage, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(message)
		return err
	}

	logFile.Write([]byte(message))

	logFile.Close()

	return nil
}
