package xdomain

import (
	"fmt"

	"github.com/ottemo/foundation/api"
)

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

	if err := context.SetResponseContentType("text/html"); err != nil {
		fmt.Println("06a968c4-cd17-43cf-a93e-b7317b53b883", err)
	}

	if _, err := responseWriter.Write([]byte("<!DOCTYPE HTML>")); err != nil {
		fmt.Println("06a968c4-cd17-43cf-a93e-b7317b53b883", err)
	}
	if _, err := responseWriter.Write(newline); err != nil {
		fmt.Println("3968707c-c8e4-4e91-a4ed-f5869619e1bc", err)
	}
	if _, err := responseWriter.Write([]byte("<script src=\"//cdn.rawgit.com/jpillora/xdomain/0.7.4/dist/xdomain.min.js\" master=\"" + xdomainMasterURL + "\"></script>")); err != nil {
		fmt.Println("5b815f17-42a6-42a5-8591-82674f1a86bd", err)
	}
	if _, err := responseWriter.Write(newline); err != nil {
		fmt.Println("a6c11509-a68b-4920-bee2-3a8ab93ea30c", err)
	}

	return "", nil
}
