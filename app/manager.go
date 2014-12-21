package app

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"

	"sync"
)

// Package global variables
var (
	// application status indicators to flag that system already in progress or done on kind of routine
	initFlag  bool
	startFlag bool
	endFlag   bool

	// synchronize locks to prevent simultaneous processing
	initMutex  sync.RWMutex
	startMutex sync.RWMutex
	endMutex   sync.RWMutex

	// registered callbacks for application events
	callbacksOnAppInit  = []func() error{}
	callbacksOnAppStart = []func() error{}
	callbacksOnAppEnd   = []func() error{}
)

// OnAppInit registers callback function on application init event
func OnAppInit(callback func() error) {
	callbacksOnAppInit = append(callbacksOnAppInit, callback)
}

// OnAppStart registers callback function on application start event
func OnAppStart(callback func() error) {
	callbacksOnAppStart = append(callbacksOnAppStart, callback)
}

// OnAppEnd registers callback function on application start event
func OnAppEnd(callback func() error) {
	callbacksOnAppEnd = append(callbacksOnAppEnd, callback)
}

// Init fires application init event for all registered modules
func Init() error {
	// prevents simultaneous execution
	initMutex.Lock()
	defer initMutex.Unlock()

	// runs all registered callbacks and doing it only once
	if !initFlag {
		for _, callback := range callbacksOnAppInit {
			if err := callback(); err != nil {
				return env.ErrorDispatch(err)
			}
		}
		initFlag = true
	}

	return nil
}

// Start fires application start event for all registered modules
func Start() error {
	// prevents simultaneous execution
	startMutex.Lock()
	defer startMutex.Unlock()

	// make sure we made init was made before start
	if !initFlag {
		err := Init()
		if err != nil {
			return err
		}
	}

	// runs all registered callbacks and doing it only once
	if !startFlag {
		for _, callback := range callbacksOnAppStart {
			if err := callback(); err != nil {
				return env.ErrorDispatch(err)
			}
		}
		startFlag = true
	}

	return nil
}

// End fires application end event for all registered modules
func End() error {
	endMutex.Lock()
	defer endMutex.Unlock()

	if !endFlag {
		endFlag = true

		for _, callback := range callbacksOnAppEnd {
			if err := callback(); err != nil {
				return env.ErrorDispatch(err)
			}
		}

		initFlag = false
		startFlag = false
		endFlag = false
	}

	return nil
}

// Serve runs HTTP server in current go routine
func Serve() error {
	return api.GetRestService().Run()
}
