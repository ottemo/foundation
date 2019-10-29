package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ottemo/commerce/app"

	// using standard set of packages
	_ "github.com/ottemo/commerce/basebuild"
	"github.com/ottemo/commerce/env"
)

func init() {
	// time.Unix() should be in UTC (as it could be not by default)
	time.Local = time.UTC
}

// executable file start point
func main() {
	defer func() {
		if err := app.End(); err != nil { // application close event
			fmt.Println(err.Error())
		}
	}()

	// we should intercept os signals to application as we should call app.End() before
	signalChain := make(chan os.Signal, 1)
	signal.Notify(signalChain, os.Interrupt, syscall.SIGTERM)
	go func() {
		for _ = range signalChain {
			err := app.End()
			if err != nil {
				_ = env.ErrorDispatch(err)
				fmt.Println(err.Error())
			}

			os.Exit(0)
		}
	}()

	// application start event
	if err := app.Start(); err != nil {
		_ = env.ErrorDispatch(err)
		fmt.Println(err.Error())
		os.Exit(0)
	}

	fmt.Println("Ottemo " + app.GetVerboseVersion())

	go func() {
		for _, engine := range env.GetDeclaredScriptEngines() {
			engine.GetScriptInstance("main").Interact()
			break
		}
	}()

	// starting HTTP server
	if err := app.Serve(); err != nil {
		fmt.Println(err.Error())
	}
}
