package api

import (
	"net/http"
)

// InterfaceSessionService is an interface to access session managing service
type InterfaceSessionService interface {
	GetName() string

	New() (InterfaceSession, error)
	Get(sessionID string) (InterfaceSession, error)
}

// InterfaceSession is an interface represents private storage for particular API request
type InterfaceSession interface {
	GetID() string

	Get(key string) interface{}
	Set(key string, value interface{})

	SetModified()

	Close() error

	Load(id string) error
	Save() error
}

// InterfaceRestService is an interface to interact with RESTFul API service
type InterfaceRestService interface {
	GetName() string

	Run() error
	RegisterAPI(service string, method string, uri string, handler FuncAPIHandler) error

	http.Handler
}
