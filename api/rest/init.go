package rest

import (
	"net/http"
	"sort"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"

	"github.com/julienschmidt/httprouter"
)

// init makes package self-initialization routine
func init() {
	var _ api.InterfaceApplicationContext = new(DefaultRestApplicationContext)

	instance := new(DefaultRestService)

	api.RegisterRestService(instance)
	env.RegisterOnConfigIniStart(instance.startup)
}

// service pre-initialization stuff
func (it *DefaultRestService) startup() error {

	it.ListenOn = ":3000"
	useXDomain := ""
	urlXDomain := ""
	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("rest.listenOn", it.ListenOn); iniValue != "" {
			it.ListenOn = iniValue
		}
		if iniValue := iniConfig.GetValue("xdomain.master", useXDomain); iniValue != "" {
			useXDomain = iniValue
		}
		if iniValue := iniConfig.GetValue("xdomain.min.js", useXDomain); iniValue != "" {
			urlXDomain = iniValue
		}
	}

	it.Router = httprouter.New()

	it.Router.PanicHandler = func(resp http.ResponseWriter, req *http.Request, params interface{}) {
		resp.WriteHeader(404)
		resp.Write([]byte("page not found"))
	}

	// handler for printing out all endpoints from active packages
	rootPageHandler := func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		newline := []byte("\n")

		resp.Header().Add("Content-Type", "text/plain")

		resp.Write([]byte("Ottemo REST API:"))
		resp.Write(newline)
		resp.Write([]byte("----"))
		resp.Write(newline)

		// sorting handlers before output
		handlers := make([]string, 0, len(it.Handlers))
		for handlerPath := range it.Handlers {
			handlers = append(handlers, handlerPath)
		}
		sort.Strings(handlers)

		for _, handlerPath := range handlers {
			resp.Write([]byte(handlerPath))
			resp.Write(newline)
		}
	}
	// show all registered API in text representation
	it.Router.GET("/", rootPageHandler)

	// support the use xdomain  - https://github.com/jpillora/xdomain
	if useXDomain != "" && urlXDomain != "" {
		xdomainHandler := func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
			newline := []byte("\n")

			resp.Header().Add("Content-Type", "text/html")

			resp.Write([]byte("<!DOCTYPE HTML>"))
			resp.Write(newline)
			resp.Write([]byte("<script src=\"" + urlXDomain + "\" master=\"" + useXDomain + "\"></script>"))
			resp.Write(newline)

		}
		it.Router.GET("/proxy.html", xdomainHandler)
	}

	it.Handlers = make(map[string]httprouter.Handle)

	api.OnRestServiceStart()

	return nil
}
