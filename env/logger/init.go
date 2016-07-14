package logger

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	log.SetFormatter(&logrus_logstash.LogstashFormatter{Type: "ottemo_api"})
	instance := new(DefaultLogger)
	var _ env.InterfaceLogger = instance

	env.RegisterLogger(instance)
	env.RegisterOnConfigIniStart(startup)
	env.RegisterOnConfigStart(setupConfig)
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
