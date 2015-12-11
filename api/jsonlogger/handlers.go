package jsonlogger

import (
	"fmt"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"net/http"
)

// requestHandler collect information about request and logs it in case of error
// jsonErrorLogLevel is set to 5, so it will ignore most of simple client side mistakes
func requestHandler(event string, eventData map[string]interface{}) bool {
	go handleData(event, eventData)
	return true
}

// handleData collect data for logging of error in JSON format
func handleData(event string, eventData map[string]interface{}) {

	result := make(map[string]interface{})
	jsonLogger := api.GetJSONLogger()
	if jsonLogger == nil {
		return
	}

	responseError, present := eventData["responseError"]

	// handle only cases when we have an error in response
	if !present || responseError == nil {
		return
	}

	if ottemoErr, ok := responseError.(env.InterfaceOttemoError); ok {
		if ottemoErr.ErrorLevel() > jsonErrorLogLevel {
			return
		}
		result["stack_trace"] = ottemoErr.ErrorCallStack()
		result["code"] = ottemoErr.ErrorCode()
		result["level"] = ottemoErr.ErrorLevel()
		result["message"] = ottemoErr.ErrorMessage()
		result["type"] = "ottemo-error"

	} else {
		result["message"] = fmt.Sprintln(responseError)
		result["type"] = "foundation-error"
	}

	if eventItem, present := eventData["context"]; present {
		if context, ok := eventItem.(api.InterfaceApplicationContext); ok && context != nil {
			result["request"] = map[string]interface{}{
				"contentType": context.GetRequestContentType(),
				"arguments":   context.GetRequestArguments(),
				"content":     context.GetRequestContent(),
			}

			if req := context.GetRequest(); req != nil {
				if typedItem, ok := req.(*http.Request); ok {
					result["requestURL"] = typedItem.RequestURI
				}
			}
		}
	}

	session, present := eventData["session"]
	if present {
		result["session"] = session
	}

	response, present := eventData["response"]
	if present {
		result["result"] = response
	}

	err := jsonLogger.Log(defaultErrorsFile, result)
	if err != nil {
		env.ErrorDispatch(err)
	}
}
