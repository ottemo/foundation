package rest

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/api/context"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetName returns implementation name of our REST API service
func (it *DefaultRestService) GetName() string {
	return "httprouter"
}

// wrappedHandler Middleware around the handler you registered
// 1. Debug timers
// 1. Parses the request
// 1. Sets the ApplicationContext
// 1. Starts the Session
// 1. Handles the Referrer cookie
// 1. Calls handler on context
// 1. Handle redirects and response encoding (json/xml)
func (it *DefaultRestService) wrappedHandler(handler api.FuncAPIHandler) httprouter.Handle {
	// httprouter supposes other format of handler than we use, so we need wrapper
	wrappedHandler := func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {

		// catching API handler fails
		defer func() {
			if recoverResult := recover(); recoverResult != nil {
				err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "28d7ef2f-631f-4f38-a916-579bf822908b", "API call fail: "+fmt.Sprintf("%v", recoverResult))
				_ = env.ErrorDispatch(err)
			}
		}()

		// debug log related variables initialization
		var startTime time.Time
		var debugRequestIdentifier string

		if utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathAPILogEnable)) {
			startTime = time.Now()
			debugRequestIdentifier = startTime.Format("20060102150405")
		}

		// Request URL parameters detection
		//----------------------------------

		// getting URL arguments of request
		//   example: http://localhost:3000/category/:categoryID/product/:productID
		//   so, ":categoryID" and ":productID" are these arguments
		reqArguments := make(map[string]string)
		for _, param := range params {
			reqArguments[param.Key] = param.Value
		}

		// getting params from URL (after "?")
		//   example: http://localhost:3000/category/:categoryID/products?limit=10&extra=55
		//   so, "limit" and "extra" are these params
		urlParsedParams, err := url.ParseQuery(req.URL.RawQuery)
		if err == nil {
			for key, value := range urlParsedParams {
				reqArguments[key] = value[0]
			}
		}

		// Request content detection
		//----------------------------
		var content interface{}
		contentType := req.Header.Get("Content-Type")

		switch {

		// request contains JSON content
		case strings.Contains(contentType, "json"):
			newContent := map[string]interface{}{}

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				_ = env.ErrorDispatch(err)
			}
			if err := json.Unmarshal(body, &newContent); err != nil {
				_ = env.ErrorDispatch(err)
			}

			content = newContent

		// request contains POST form data
		case strings.Contains(contentType, "form-data"):
			newContent := map[string]interface{}{}

			if err := req.ParseForm(); err != nil {
				_ = env.ErrorDispatch(err)
			}
			for attribute, value := range req.PostForm {
				newContent[attribute], _ = url.QueryUnescape(value[0])
			}

			if err := req.ParseMultipartForm(32 << 20); err != nil {// 32 MB
				_ = env.ErrorDispatch(err)
			}
			if req.MultipartForm != nil {
				for attribute, value := range req.MultipartForm.Value {
					newContent[attribute], _ = url.QueryUnescape(value[0])
				}
			}

			content = newContent

		// request contains "x-www-form-urlencoded" data
		case strings.Contains(contentType, "urlencode"):
			newContent := map[string]interface{}{}

			if err := req.ParseForm(); err != nil {
				_ = env.ErrorDispatch(err)
			}
			for attribute, value := range req.PostForm {
				newContent[attribute], _ = url.QueryUnescape(value[0])
			}

			content = newContent

		// request contains POST text
		case strings.Contains(contentType, "text/plain"):
			fallthrough
		default:
			var body []byte

			body, err = ioutil.ReadAll(req.Body)
			if err != nil {
				_ = env.ErrorDispatch(err)
			}

			content = string(body)
		}

		// Handling request
		//------------------

		// preparing struct for API handler
		applicationContext := new(DefaultRestApplicationContext)
		applicationContext.Request = req
		applicationContext.RequestArguments = reqArguments
		applicationContext.RequestContent = content
		applicationContext.RequestFiles = make(map[string]io.Reader)
		applicationContext.ResponseWriter = resp
		applicationContext.ContextValues = make(map[string]interface{})

		// collecting request files
		if req.MultipartForm != nil && req.MultipartForm.File != nil {
			for _, fileInfoArray := range req.MultipartForm.File {
				for _, fileInfo := range fileInfoArray {
					attachedFile, err := fileInfo.Open()
					if err == nil {
						mediaFileName := fileInfo.Filename
						if contentDisposition, present := fileInfo.Header["Content-Disposition"]; present && len(contentDisposition) > 0 {
							if _, mediaParams, err := mime.ParseMediaType(contentDisposition[0]); err == nil {
								if value, present := mediaParams["name"]; present {
									if len(reqArguments[value]) != 0 {
										reqArguments[value] += ", "
									}
									reqArguments[value] += mediaFileName
								}
							}

						}
						applicationContext.RequestFiles[mediaFileName] = attachedFile
					}
				}
			}
		}

		// starting session for request
		currentSession, err := api.StartSession(applicationContext)
		if err != nil {
			err = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c8a3bbf8-215f-4dff-b0e7-3d0d102ad02d", "Session init fail: "+err.Error())
			_ = env.ErrorDispatch(err)
		}
		applicationContext.Session = currentSession

		if utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathAPILogEnable)) {
			allowLog := true
			apiExcludedURIs := env.ConfigGetValue(ConstConfigPathAPILogExclude)
			if apiExcludedURIs != nil {
				excludedURIs := utils.InterfaceToArray(apiExcludedURIs)
				needle := req.RequestURI + " {" + req.Method + "}"
				for i := 0; i < len(excludedURIs); i++ {
					if needle == excludedURIs[i] {
						allowLog = false
						break
					}
				}
			}
			if allowLog {
				env.Log(ConstDebugLogStorage, "REQUEST_"+debugRequestIdentifier, fmt.Sprintf("%s [%s]\n%#v\n", req.RequestURI, currentSession.GetID(), content))
				env.LogEvent(env.LogFields{
					"request_thread_id": debugRequestIdentifier,
					"session_id":        currentSession.GetID(),

					"uri":          req.RequestURI,
					"verb":         req.Method,
					"content":      content,
					"agent":        req.UserAgent(),
					"clientip":     req.RemoteAddr,
					"httpversion":  req.Proto,
					"host":         req.Host,
					"content_type": contentType,
				}, "request")
			}
		}

		// event for request
		eventData := map[string]interface{}{"session": currentSession, "context": applicationContext}
		env.Event("api.request", eventData)

		// store admin credentials for later in-call use
		var result interface{}
		context.MakeContext(func() {
			if callContext := context.GetContext(); callContext != nil {
				callContext["is_admin"] = false
				if api.ValidateAdminRights(applicationContext) == nil {
					callContext["is_admin"] = true
				}
			} else {
				err = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6b94a499-9d71-403e-9f67-06fd90d6250d", "can not get context for API handler")
			}

			if err == nil {
				// API handler processing
				result, err = handler(applicationContext)
			}
		})

		if err != nil {
			_ = env.ErrorDispatch(err)
			env.LogEvent(env.LogFields{
				"request_thread_id": debugRequestIdentifier,
				"session_id":        currentSession.GetID(),

				"uri":        req.RequestURI,
				"error_dump": err,
			}, "handler_error")
		}

		if err == nil {
			applicationContext.Result = result
		}

		// event for response
		eventData["response"] = result
		env.Event("api.response", eventData)
		result = eventData["response"]

		// result conversion before output
		redirectLocation := ""
		if redirect, ok := result.(api.StructRestRedirect); ok {
			if redirect.DoRedirect {
				resp.Header().Add("Location", redirect.Location)
				resp.WriteHeader(301)
				result = []byte("")
			} else {
				redirectLocation = redirect.Location
				result = redirect.Result
			}
		}

		// converting result to []byte if it is not already done
		if _, ok := result.([]byte); !ok {

			// JSON encode
			if resp.Header().Get("Content-Type") == "application/json" {
				var errorMsg map[string]interface{}
				if err != nil {
					if _, ok := err.(env.InterfaceOttemoError); !ok {
						err = env.ErrorDispatch(err)
					}

					if ottemoError, ok := err.(env.InterfaceOttemoError); ok {
						errorMsg = map[string]interface{}{
							"message": ottemoError.Error(),
							"level":   ottemoError.ErrorLevel(),
							"code":    ottemoError.ErrorCode(),
						}
					} else {
						_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bdbb8627-18e8-4969-a048-c8b482235f39", "can't convert error to ottemoError")
						errorMsg = map[string]interface{}{
							"message": err.Error(),
							"level":   env.ConstErrorLevelAPI,
							"code":    "896810b9-9b54-471a-830c-b77b33379adc",
						}
					}
				}

				response := map[string]interface{}{
					"result":   result,
					"error":    errorMsg,
					"redirect": redirectLocation,
				}

				if utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathAPILogEnable)) {
					responseTime := time.Now().Sub(startTime)
					env.Log(ConstDebugLogStorage, "RESPONSE_"+debugRequestIdentifier, fmt.Sprintf("%s (%dns)\n%s\n", req.RequestURI, responseTime, result))

					logFields := env.LogFields{
						"request_thread_id": debugRequestIdentifier,
						"session_id":        currentSession.GetID(),
						"uri":               req.RequestURI,
						"resp_time":         responseTime,
						"response":          response,
					}
					env.LogEvent(logFields, "response")
				}

				result, _ = json.Marshal(response)
			}

			// XML encode
			if resp.Header().Get("Content-Type") == "text/xml" && result != nil {
				xmlResult, _ := xml.MarshalIndent(result, "", "    ")
				result = []byte(xml.Header + string(xmlResult))
			}
		}

		if value, ok := result.([]byte); ok {
			if _, err := resp.Write(value); err != nil {
				_ = env.ErrorDispatch(err)
			}
		} else if result != nil {
			if _, err := resp.Write([]byte(fmt.Sprint(result))); err != nil {
				_ = env.ErrorDispatch(err)
			}
		}
	}

	return wrappedHandler
}

// GET is a wrapper for the HTTP GET verb
func (it *DefaultRestService) GET(resource string, handler api.FuncAPIHandler) {
	path := "/" + resource
	it.Router.GET(path, it.wrappedHandler(handler))

	it.Handlers = append(it.Handlers, path+" {GET}")
}

// PUT is a wrapper for the HTTP PUT verb
func (it *DefaultRestService) PUT(resource string, handler api.FuncAPIHandler) {
	path := "/" + resource
	it.Router.PUT(path, it.wrappedHandler(handler))

	it.Handlers = append(it.Handlers, path+" {PUT}")
}

// POST is a wrapper for the HTTP POST verb
func (it *DefaultRestService) POST(resource string, handler api.FuncAPIHandler) {
	path := "/" + resource
	it.Router.POST(path, it.wrappedHandler(handler))

	it.Handlers = append(it.Handlers, path+" {POST}")
}

// DELETE is a wrapper for the HTTP DELETE verb
func (it *DefaultRestService) DELETE(resource string, handler api.FuncAPIHandler) {
	path := "/" + resource
	it.Router.DELETE(path, it.wrappedHandler(handler))

	it.Handlers = append(it.Handlers, path+" {DELETE}")
}

// ServeHTTP is an entry point for HTTP request, it takes control before request handled
// (go lang "http.server" package "Handler" interface implementation)
func (it DefaultRestService) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {

	// CORS fix-up
	responseWriter.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	responseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	responseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	responseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cookie, X-Referer, Content-Length, Accept-Encoding, X-CSRF-Token")

	responseWriter.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	responseWriter.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	responseWriter.Header().Set("Expires", "0")                                         // Proxies

	if request.Method == "GET" || request.Method == "POST" || request.Method == "PUT" || request.Method == "DELETE" {

		// default output format
		responseWriter.Header().Set("Content-Type", "application/json")

		request.URL.Path = strings.Replace(request.URL.Path, "/foundation", "", -1)

		it.Router.ServeHTTP(responseWriter, request)
	}
}

// Run is the Ottemo REST server startup function, analogous to "ListenAndServe"
func (it *DefaultRestService) Run() error {
	fmt.Println("REST API Service [HTTPRouter] starting to listen on " + it.ListenOn)
	err := http.ListenAndServe(it.ListenOn, it)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
