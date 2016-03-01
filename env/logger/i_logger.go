package logger

import (
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
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

// LogEvent Saves log details out to a file for logstash consumption
func (it *DefaultLogger) LogEvent(fields env.LogFields, eventName string) {
	f, err := os.OpenFile(baseDirectory+"events.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		it.LogError(err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.WithFields(log.Fields(fields)).Info(eventName)
}
