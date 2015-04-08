package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ottemo/foundation/app"

	// using standard set of packages
	_ "github.com/ottemo/foundation/basebuild"
)

// executable file start point
func main() {
	defer app.End() // application close event

	// we should intercept os signals to application as we should call app.End() before
	signalChain := make(chan os.Signal, 1)
	signal.Notify(signalChain, os.Interrupt, syscall.SIGTERM)
	go func() {
		for _ = range signalChain {
			err := app.End()
			if err != nil {
				fmt.Println(err.Error())
			}

			os.Exit(0)
		}
	}()

	// application start event
	if err := app.Start(); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	fmt.Println("Ottemo " + app.GetVerboseVersion())

	// starting HTTP server
	app.Serve()
}
