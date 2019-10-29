package otto

import (
	"bytes"
	"github.com/ottemo/commerce/env"
	"io"
	"net/http"

	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/utils"
)

// makeApplicationContext creates new InterfaceApplicationContext instance
func makeApplicationContext() api.InterfaceApplicationContext {
	applicationContext := new(ApplicationContext)
	applicationContext.Response = bytes.NewBufferString("")
	applicationContext.Result = applicationContext.Response
	applicationContext.ResponseSettings = make(map[string]interface{})

	applicationContext.RequestArguments = make(map[string]string)
	applicationContext.RequestSettings = make(map[string]interface{})
	applicationContext.RequestParameters = make(map[string]string)
	applicationContext.RequestContent = nil
	applicationContext.RequestFiles = make(map[string]io.Reader)
	applicationContext.ContextValues = make(map[string]interface{})

	if session, err := api.NewSession(); err == nil {
		applicationContext.Session = session
	}

	return applicationContext
}

// apiHandler returns API handler for a giver resource
func apiHandler(method string, resource string) (api.FuncAPIHandler, error) {
	if service := api.GetRestService(); service != nil {
		if handler := service.GetHandler(method, resource); handler != nil {
			return handler, nil
		}
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f0e17e01-b0a7-43a0-a7e1-90646cd5c309", "Handler not found")
	}
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e5c20766-e3bc-4020-a918-c5d396b321df", "API service is not available")
}

// apiCall performs API call for a given resource
func apiCall(method string, resource string, context api.InterfaceApplicationContext) (interface{}, error) {
	if context == nil {
		context = makeApplicationContext()
	}

	handler, err := apiHandler(method, resource)
	if err != nil {
		return nil, err
	}
	return handler(context)
}

// GetRequest returns raw request object
func (it *ApplicationContext) GetRequest() interface{} {
	return nil
}

// GetResponse returns raw response object
func (it *ApplicationContext) GetResponse() interface{} {
	return it.Result
}

// GetResponseWriter returns io.Writes for response (for ApplicationContext it is clone to GetResponse)
func (it *ApplicationContext) GetResponseWriter() io.Writer {
	return it.Response
}

// GetRequestArguments returns all arguments provided to API function
//   - for REST API it is URI parameters "http://localhost/myfunc/:param1/get/:param2/:param3/"
func (it *ApplicationContext) GetRequestArguments() map[string]string {
	return it.RequestArguments
}

// GetRequestArgument returns particular argument provided to API function or ""
func (it *ApplicationContext) GetRequestArgument(name string) string {
	if value, present := it.RequestArguments[name]; present {
		return value
	}
	return ""
}

// GetRequestContent returns request contents of nil if not specified (HTTP request body)
func (it *ApplicationContext) GetRequestContent() interface{} {
	return it.RequestContent
}

// GetRequestFiles returns files were attached to request
func (it *ApplicationContext) GetRequestFiles() map[string]io.Reader {
	return it.RequestFiles
}

// GetRequestFile returns particular file attached to request or nil
func (it *ApplicationContext) GetRequestFile(name string) io.Reader {
	if file, present := it.RequestFiles[name]; present {
		return file
	}
	return nil
}

// GetRequestSettings returns request related settings
//   - for REST API settings are HTTP headers and COOKIES
func (it *ApplicationContext) GetRequestSettings() map[string]interface{} {
	return it.RequestSettings
}

// GetRequestSetting returns particular request related setting or nil
func (it *ApplicationContext) GetRequestSetting(name string) interface{} {
	if value, present := it.RequestSettings[name]; present {
		return value
	}
	return nil
}

// GetRequestContentType returns MIME type of request content
func (it *ApplicationContext) GetRequestContentType() string {
	if value, present := it.RequestSettings["Content-Type"]; present {
		return utils.InterfaceToString(value)
	}
	return ""
}

// GetResponseContentType returns MIME type of supposed response content
func (it *ApplicationContext) GetResponseContentType() string {
	if value, present := it.ResponseSettings["Content-Type"]; present {
		return utils.InterfaceToString(value)
	}
	return ""
}

// SetResponseContentType changes response content type, returns error if not possible
func (it *ApplicationContext) SetResponseContentType(mimeType string) error {
	it.ResponseSettings["Content-Type"] = mimeType
	return nil
}

// GetResponseSetting returns specified setting value (for REST API returns header as settings)
func (it *ApplicationContext) GetResponseSetting(name string) interface{} {
	if value, present := it.RequestSettings[name]; present {
		return value
	}
	return nil
}

// SetResponseSetting specifies response setting (for REST API it just sets additional header)
func (it *ApplicationContext) SetResponseSetting(name string, value interface{}) error {
	it.ResponseSettings[name] = value
	return nil
}

// SetResponseStatus will set an HTTP response code
//    - code is an integer correlating to HTTP response codes
func (it *ApplicationContext) SetResponseStatus(code int) {
	it.ResponseSettings["status"] = code
}

// SetResponseStatusBadRequest will set the ResponseWriter to StatusBadRequest (400)
func (it *ApplicationContext) SetResponseStatusBadRequest() {
	it.ResponseSettings["status"] = http.StatusBadRequest
}

// SetResponseStatusForbidden will set the ResponseWriter to StatusForbidden (403)
func (it *ApplicationContext) SetResponseStatusForbidden() {
	it.ResponseSettings["status"] = http.StatusForbidden
}

// SetResponseStatusNotFound will set the ResponseWriter to StatusNotFound (404)
func (it *ApplicationContext) SetResponseStatusNotFound() {
	it.ResponseSettings["status"] = http.StatusNotFound
}

// SetResponseStatusInternalServerError will set the ResponseWriter to StatusInternalServerError (500)
func (it *ApplicationContext) SetResponseStatusInternalServerError() {
	it.ResponseSettings["status"] = http.StatusInternalServerError
}

// GetResponseResult returns result going to be written to response writer
func (it *ApplicationContext) GetResponseResult() interface{} {
	return it.Result
}

// SetResponseResult changes result going to be written to response writer
func (it *ApplicationContext) SetResponseResult(value interface{}) error {
	it.Result = value
	return nil
}

// GetContextValues returns current context related values map
func (it *ApplicationContext) GetContextValues() map[string]interface{} {
	return it.ContextValues
}

// GetContextValue returns particular context related value or nil if not set
func (it *ApplicationContext) GetContextValue(key string) interface{} {
	if value, present := it.ContextValues[key]; present {
		return value
	}
	return nil
}

// SetContextValue stores specified value in current context
func (it *ApplicationContext) SetContextValue(key string, value interface{}) {
	it.ContextValues[key] = value
}

// SetSession assigns given session to current context
func (it *ApplicationContext) SetSession(session api.InterfaceSession) error {
	it.Session = session
	return nil
}

// GetSession returns session assigned to current context or nil if nothing was assigned
func (it *ApplicationContext) GetSession() api.InterfaceSession {
	return it.Session
}
