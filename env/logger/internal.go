package logger

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/api/context"
	"github.com/ottemo/foundation/env"
)

// logErrorToJSON collect required info from stack context to map for logging purpose
func (it *DefaultLogger) logErrorToJSON(err error) {
	stackContext := context.GetContext()
	result := make(map[string]interface{})

	if ottemoErr, ok := err.(env.InterfaceOttemoError); ok {
		result["stack_trace"] = ottemoErr.ErrorCallStack()
		result["code"] = ottemoErr.ErrorCode()
		result["level"] = ottemoErr.ErrorLevel()
		result["message"] = ottemoErr.ErrorMessage()
		result["type"] = "ottemo-error"

	} else {
		result["message"] = err.Error()
		result["type"] = "foundation-error"
	}

	if contextValue, present := stackContext["context"]; present {
		if applicationContext, ok := contextValue.(api.InterfaceApplicationContext); ok && applicationContext != nil {
			result["request"] = map[string]interface{}{
				"contentType": applicationContext.GetRequestContentType(),
				"arguments":   applicationContext.GetRequestArguments(),
				"content":     applicationContext.GetRequestContent(),
			}

			result["session"] = applicationContext.GetSession()
		}
	}

	if requestURL, present := stackContext["requestURL"]; present {
		result["request_url"] = requestURL
	}

	if response, present := stackContext["response"]; present {
		result["result"] = response
	}

	it.LogMap(defaultJSONErrorsFile, result)
}
