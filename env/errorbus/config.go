package errorbus

import (
	"errors"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "6dee39ac-c930-420e-b777-b95f6cab8981", "can't obtain config")
		return env.ErrorDispatch(err)
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
		return env.ErrorDispatch(err)
	}

	// Hide level
	hideLevelValidator := func(newValue interface{}) (interface{}, error) {
		newLevel := utils.InterfaceToInt(newValue)
		if newLevel > 10 || newLevel < 0 {
			err := errors.New("'Hide level' config value should be between 0 and 10")
			return hideLevel, env.ErrorDispatch(err)
		}
		hideLevel = newLevel

		return hideLevel, nil
	}
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathErrorHideLevel,
		Value:       5,
		Type:        env.ConstConfigTypeInteger,
		Editor:      "integer",
		Options:     nil,
		Label:       "Hide level",
		Description: "errors above specified level will be replaced to stub system error message",
		Image:       "",
	}, hideLevelValidator)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Hide message
	hideMessageValidator := func(newValue interface{}) (interface{}, error) {
		if newMessage, ok := newValue.(string); ok {
			hideMessage = newMessage
		} else {
			err := errors.New("wrong type for 'Hide message' config value")
			return hideMessage, env.ErrorDispatch(err)
		}
		return hideMessage, nil
	}
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathErrorHideMessage,
		Value:       "System error has occured",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Hide message",
		Description: "system error message to substitute error message above hide level",
		Image:       "",
	}, hideMessageValidator)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
