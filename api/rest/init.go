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
	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("rest.listenOn", it.ListenOn); iniValue != "" {
			it.ListenOn = iniValue
		}
	}

	it.Router = httprouter.New()

	it.Router.PanicHandler = func(resp http.ResponseWriter, req *http.Request, params interface{}) {
		resp.WriteHeader(404)
		resp.Write([]byte("page not found"))
	}

	// our homepage - shows all registered API in text representation
	it.Router.GET("/", it.rootPageHandler)

	api.OnRestServiceStart()

	return nil
}

func (it *DefaultRestService) rootPageHandler(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	newline := []byte("\n")

	resp.Header().Add("Content-Type", "text/plain")

	resp.Write([]byte("Ottemo REST API:"))
	resp.Write(newline)
	resp.Write([]byte("----"))
	resp.Write(newline)

	// sorting handlers before output
	sort.Strings(it.Handlers)

	for _, handlerPath := range it.Handlers {
		resp.Write([]byte(handlerPath))
		resp.Write(newline)
	}
}
