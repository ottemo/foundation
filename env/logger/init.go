package logger

import (
	"github.com/ottemo/foundation/env"
	"os"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultLogger)
	var _ env.InterfaceLogger = instance

	env.RegisterLogger(instance)
	env.RegisterOnConfigIniStart(startup)
}

// startup is a service pre-initialization stuff
func startup() error {
	if _, err := os.Stat(baseDirectory); !os.IsExist(err) {
		err := os.MkdirAll(baseDirectory, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
