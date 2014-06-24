package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/ottemo/foundation/config"
)

type HTTPHandler func(resp http.ResponseWriter, req *http.Request)
type JSONHandler func(req *http.Request) map[string]interface{}

type RestService struct {
	Negroni  *negroni.Negroni
	Mux      *http.ServeMux
	ListenOn string

	Handlers map[string]HTTPHandler
}

func init() {
	instance := new(RestService)

	RegisterEndPoint(instance)
	config.RegisterOnConfigIniStart(instance.Startup)
}

func (rs *RestService) Startup() error {
	rs.Mux = http.NewServeMux()

	rs.ListenOn = ":9000"
	if iniConfig := config.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("negroni.port"); iniValue != "" {
			rs.ListenOn = iniValue
		}
	}

	rs.Mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		newline := []byte("\n")
		for path, _ := range rs.Handlers {
			resp.Header().Add("Content-Type", "text")
			resp.Write([]byte(path))
			resp.Write(newline)
		}
	})

	rs.Negroni = negroni.Classic()
	rs.Negroni.UseHandler(rs.Mux)

	rs.Handlers = make(map[string]HTTPHandler)

	OnEndPointStart()

	return nil
}

func (rs *RestService) GetName() string {
	return "Negroni"
}

func (rs *RestService) RegisterJsonAPI(service string, uri string, handler func(req *http.Request) map[string]interface{}) error {

	jsonHandler := func(resp http.ResponseWriter, req *http.Request) {
		result, _ := json.Marshal(handler(req))

		resp.Header().Add("Content-Type", "application/json")
		resp.Write(result)
	}

	path := "/" + service + "/" + uri

	if _, present := rs.Handlers[path]; present {
		return errors.New("There is already registered handler for " + path)
	}

	rs.Mux.HandleFunc(path, jsonHandler)
	rs.Handlers[path] = jsonHandler

	return nil
}

func (rs *RestService) Run() error {
	rs.Negroni.Run(rs.ListenOn)

	return nil
}
