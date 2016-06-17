package rest

import (
	"io"
	"net/http"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/utils"
)

// GetRequest returns raw request object
func (it *DefaultRestApplicationContext) GetRequest() interface{} {
	return it.Request
}

// GetResponse returns raw response object
func (it *DefaultRestApplicationContext) GetResponse() interface{} {
	return it.ResponseWriter
}

// GetResponseWriter returns io.Writes for response (for DefaultRestApplicationContext it is clone to GetResponse)
func (it *DefaultRestApplicationContext) GetResponseWriter() io.Writer {
	return it.ResponseWriter
}

// GetRequestArguments returns all arguments provided to API function
//   - for REST API it is URI parameters "http://localhost/myfunc/:param1/get/:param2/:param3/"
func (it *DefaultRestApplicationContext) GetRequestArguments() map[string]string {
	return it.RequestArguments
}

// GetRequestArgument returns particular argument provided to API function or ""
func (it *DefaultRestApplicationContext) GetRequestArgument(name string) string {
	if value, present := it.RequestArguments[name]; present {
		return value
	}
	return ""
}

// GetRequestContent returns request contents of nil if not specified (HTTP request body)
func (it *DefaultRestApplicationContext) GetRequestContent() interface{} {
	return it.RequestContent
}

// GetRequestFiles returns files were attached to request
func (it *DefaultRestApplicationContext) GetRequestFiles() map[string]io.Reader {
	return it.RequestFiles
}

// GetRequestFile returns particular file attached to request or nil
func (it *DefaultRestApplicationContext) GetRequestFile(name string) io.Reader {
	if file, present := it.RequestFiles[name]; present {
		return file
	}
	return nil
}

// GetRequestSettings returns request related settings
//   - for REST API settings are HTTP headers and COOKIES
func (it *DefaultRestApplicationContext) GetRequestSettings() map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range it.Request.Header {
		result[key] = value
	}

	// cookie can overlap header information - it is supposed behaviour
	for _, cookie := range it.Request.Cookies() {
		result[cookie.Name] = cookie.Value
	}

	return result
}

// GetRequestSetting returns particular request related setting or nil
func (it *DefaultRestApplicationContext) GetRequestSetting(name string) interface{} {
	if value, err := it.Request.Cookie(name); err == nil && value != nil {
		return value.Value
	}

	if value, present := it.Request.Header[name]; present {
		if len(value) > 1 {
			return value
		}
		return value[0]
	}

	return nil
}

// GetRequestContentType returns MIME type of request content
func (it *DefaultRestApplicationContext) GetRequestContentType() string {
	return it.Request.Header.Get("Content-Type")
}

// GetResponseContentType returns MIME type of supposed response content
func (it *DefaultRestApplicationContext) GetResponseContentType() string {
	return it.ResponseWriter.Header().Get("Content-Type")
}

// SetResponseContentType changes response content type, returns error if not possible
func (it *DefaultRestApplicationContext) SetResponseContentType(mimeType string) error {
	it.ResponseWriter.Header().Set("Content-Type", mimeType)
	return nil
}

// GetResponseSetting returns specified setting value (for REST API returns header as settings)
func (it *DefaultRestApplicationContext) GetResponseSetting(name string) interface{} {
	return it.ResponseWriter.Header().Get(name)
}

// SetResponseSetting specifies response setting (for REST API it just sets additional header)
func (it *DefaultRestApplicationContext) SetResponseSetting(name string, value interface{}) error {
	it.ResponseWriter.Header().Set(name, utils.InterfaceToString(value))
	return nil
}

// SetResponseStatus will set an HTTP response code
//    - code is an integer correlating to HTTP response codes
func (it *DefaultRestApplicationContext) SetResponseStatus(code int) {
	it.ResponseWriter.WriteHeader(code)
}

// SetResponseStatusBadRequest will set the ResponseWriter to StatusBadRequest (400)
func (it *DefaultRestApplicationContext) SetResponseStatusBadRequest() {
	it.SetResponseStatus(http.StatusBadRequest)
}

// SetResponseStatusForbidden will set the ResponseWriter to StatusForbidden (403)
func (it *DefaultRestApplicationContext) SetResponseStatusForbidden() {
	it.SetResponseStatus(http.StatusForbidden)
}

// SetResponseStatusNotFound will set the ResponseWriter to StatusNotFound (404)
func (it *DefaultRestApplicationContext) SetResponseStatusNotFound() {
	it.SetResponseStatus(http.StatusNotFound)
}

// SetResponseStatusInternalServerError will set the ResponseWriter to StatusInternalServerError (500)
func (it *DefaultRestApplicationContext) SetResponseStatusInternalServerError() {
	it.SetResponseStatus(http.StatusInternalServerError)
}

// GetResponseResult returns result going to be written to response writer
func (it *DefaultRestApplicationContext) GetResponseResult() interface{} {
	return it.Result
}

// SetResponseResult changes result going to be written to response writer
func (it *DefaultRestApplicationContext) SetResponseResult(value interface{}) error {
	it.Result = value
	return nil
}

// GetContextValues returns current context related values map
func (it *DefaultRestApplicationContext) GetContextValues() map[string]interface{} {
	return it.ContextValues
}

// GetContextValue returns particular context related value or nil if not set
func (it *DefaultRestApplicationContext) GetContextValue(key string) interface{} {
	if value, present := it.ContextValues[key]; present {
		return value
	}
	return nil
}

// SetContextValue stores specified value in current context
func (it *DefaultRestApplicationContext) SetContextValue(key string, value interface{}) {
	it.ContextValues[key] = value
}

// SetSession assigns given session to current context
func (it *DefaultRestApplicationContext) SetSession(session api.InterfaceSession) error {
	it.Session = session
	return nil
}

// GetSession returns session assigned to current context or nil if nothing was assigned
func (it *DefaultRestApplicationContext) GetSession() api.InterfaceSession {
	return it.Session
}
