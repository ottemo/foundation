package api

import (
	"errors"
	"net/http"
)

type EndPoint interface {
	GetName() string

	Run() error
	RegisterJsonAPI(service string, uri string, handler func(req *http.Request) map[string]interface{}) error
}

var currentEndPoint EndPoint = nil
var callbacksOnEndPointStart = []func() error{}

func RegisterOnEndPointStart(callback func() error) {
	callbacksOnEndPointStart = append(callbacksOnEndPointStart, callback)
}

func OnEndPointStart() error {
	for _, callback := range callbacksOnEndPointStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

func RegisterEndPoint(ep EndPoint) error {
	if currentEndPoint == nil {
		currentEndPoint = ep
	} else {
		return errors.New("Sorry, '" + currentEndPoint.GetName() + "' already registered")
	}
	return nil
}

func GetEndPoint() EndPoint {
	return currentEndPoint
}
