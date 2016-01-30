package rest

import (
	"net/http"
	"sort"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"

	"github.com/julienschmidt/httprouter"
)

// init performs the package self-initialization routine
func init() {
	var _ api.InterfaceApplicationContext = new(DefaultRestApplicationContext)

	instance := new(DefaultRestService)

	api.RegisterRestService(instance)
	env.RegisterOnConfigIniStart(instance.startup)
}

// service pre-initialization stuff
func (it *DefaultRestService) startup() error {

	it.ListenOn = ":3000"
	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("rest.listenOn", it.ListenOn); iniValue != "" {
			it.ListenOn = iniValue
		}
	}

	it.Router = httprouter.New()

	it.Router.PanicHandler = func(w http.ResponseWriter, r *http.Request, params interface{}) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
	}

	it.Router.GET("/", it.rootPageHandler)
	it.Router.GET("/proxy.html", it.xdomainHandler)

	api.OnRestServiceStart()

	return nil
}

// rootPageHandler Display a list of the registered endpoints
func (it *DefaultRestService) rootPageHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	newline := []byte("\n")

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Ottemo REST API:"))
	w.Write(newline)
	w.Write([]byte("----"))
	w.Write(newline)

	// sorting handlers before output
	sort.Strings(it.Handlers)

	for _, handlerPath := range it.Handlers {
		w.Write([]byte(handlerPath))
		w.Write(newline)
	}
}

// xdomainHandler will enale the usage of xdomain instead of CORS for legacy browsers
func (it *DefaultRestService) xdomainHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	useXDomain := ""
	urlXDomain := ""

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("xdomain.master", useXDomain); iniValue != "" {
			useXDomain = iniValue
		}
		if iniValue := iniConfig.GetValue("xdomain.min.js", useXDomain); iniValue != "" {
			urlXDomain = iniValue
		}
	}

	if useXDomain != "" && urlXDomain != "" {
		newline := []byte("\n")

		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("<!DOCTYPE HTML>"))
		w.Write(newline)
		w.Write([]byte("<script src=\"" + urlXDomain + "\" master=\"" + useXDomain + "\"></script>"))
		w.Write(newline)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
	}
}
