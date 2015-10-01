package logger

import (
	"errors"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "6dee39ac-c930-420e-b777-b95f6cab8981", "can't obtain config")
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathError,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Error",
		Description: "error handling settings",
		Image:       "",
	}, nil)

	if err != nil {
		env.ErrorDispatch(err)
	}

	// Log level
	logLevelValidator := func(newValue interface{}) (interface{}, error) {
		newLevel := utils.InterfaceToInt(newValue)
		if newLevel > 10 || newLevel < 0 {
			return errorLogLevel, errors.New("'Log level' config value should be between 0 and 10")
		}
		errorLogLevel = newLevel

		return errorLogLevel, nil
	}
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathErrorLogLevel,
		Value:       5,
		Type:        env.ConstConfigTypeInteger,
		Editor:      "integer",
		Options:     nil,
		Label:       "Log level",
		Description: "errors below specified level will be send to logger service",
		Image:       "",
	}, logLevelValidator)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
