package jsonlogger

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"net/http"
)

// requestHandler collect information about request and logs it in case of error
// eventData keys - session, context, referrer, response
// events - api.request and api.response
func requestHandler(event string, eventData map[string]interface{}) bool {
	go handleData(event, eventData)
	return true
}

// handleData
func handleData(event string, eventData map[string]interface{}) {
	var context api.InterfaceApplicationContext
	if eventItem, present := eventData["context"]; present {
		if typedItem, ok := eventItem.(api.InterfaceApplicationContext); ok {
			context = typedItem
		}
	}

	//var eventDataKeys = []string{"session", "context", "referrer", "response"}

	result := make(map[string]interface{})

	session, present := eventData["session"]
	if present {
		result["session"] = session
	}

	response, present := eventData["response"]
	if present {
		result["response"] = response
	}

	if context != nil {
		req := context.GetRequest()

		if req != nil {
			if typedItem, ok := req.(*http.Request); ok {
				result["request"] = typedItem.RequestURI
			}
		}
	}

	/*
		if errorData.Error != nil {
			if ottemoErr, ok := errorData.Error.(env.InterfaceOttemoError); ok {
				errorMap["stack_trace"] = ottemoErr.ErrorCallStack()
				errorMap["code"] = ottemoErr.ErrorCode()
				errorMap["level"] = ottemoErr.ErrorLevel()
				errorMap["message"] = ottemoErr.ErrorMessage()

			} else {
				errorMap["message"] = errorData.Error.Error()
			}
		}

		logInfo := map[string]interface {}{
			"error": errorMap,

		}
	*/

	jsonLogger := api.GetJSONLogger()
	err := jsonLogger.Log(ConstDebugLogStorage, result)
	if err != nil {
		env.ErrorDispatch(err)
	}
}
