package api

import (
	"net/http"
)

// InterfaceSession is an interface to access private storage assigned to particular API request
type InterfaceSession interface {
	GetID() string

	Get(key string) interface{}
	Set(key string, value interface{})

	Close() error
}

// InterfaceRestService is an interface to interact with RESTFul API service
type InterfaceRestService interface {
	GetName() string

	Run() error
	RegisterAPI(service string, method string, uri string, handler FuncAPIHandler) error

	http.Handler
}
