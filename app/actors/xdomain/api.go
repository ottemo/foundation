package xdomain

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
)

// endpoint configuration for xdomain package
func setupAPI() error {

	service := api.GetRestService()
	service.GET("/proxy.html", xdomainHandler)

	return nil
}

// xdomainHandler will enable the usage of xdomain instead of CORS for legacy browsers
func (it *DefaultRestService) xdomainHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	xdomainMasterUrl := ""

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("xdomain.master", xdomainMasterUrl); iniValue != "" {
			xdomainMasterUrl = iniValue
		}
	}

	if xdomainMasterUrl != "" {
		newline := []byte("\n")

		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("<!DOCTYPE HTML>"))
		w.Write(newline)
		w.Write([]byte("<script src=\"" + urlXDomain + "\" master=\"" + xdomainMasterUrl + "\"></script>"))
		w.Write(newline)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
	}
}
