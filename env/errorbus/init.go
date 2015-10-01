package errorbus

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine
func init() {
	instance := &DefaultErrorBus{listeners: make([]env.FuncErrorListener, 0)}
	var _ env.InterfaceErrorBus = instance

	var _ env.InterfaceOttemoError = new(OttemoError)

	env.RegisterErrorBus(instance)
	env.RegisterOnConfigIniStart(setupOnIniConfigStart)
	env.RegisterOnConfigStart(setupConfig)
}

// setupOnIniConfigStart is a initialization based on ini config service
func setupOnIniConfigStart() error {

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("error.instant.debug", "true"); iniValue != "" {
			debug = utils.InterfaceToBool(iniValue)
		}
	}

	return nil
}
