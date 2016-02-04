package xdomain

import "github.com/ottemo/foundation/api"

// endpoint configuration for xdomain package
func setupAPI() error {

	service := api.GetRestService()
	service.GET("proxy.html", xdomainHandler)

	return nil
}

// xdomainHandler will enable the usage of xdomain instead of CORS for legacy browsers
func xdomainHandler(context api.InterfaceApplicationContext) (interface{}, error) {

	responseWriter := context.GetResponseWriter()

	newline := []byte("\n")

	context.SetResponseContentType("text/html")

	responseWriter.Write([]byte("<!DOCTYPE HTML>"))
	responseWriter.Write(newline)
	responseWriter.Write([]byte("<script src=\"//cdn.rawgit.com/jpillora/xdomain/0.7.4/dist/xdomain.min.js\" master=\"" + xdomainMasterURL + "\"></script>"))
	responseWriter.Write(newline)

	return "", nil
}
