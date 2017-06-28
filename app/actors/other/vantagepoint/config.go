package vantagepoint

import (
	"github.com/ottemo/foundation/env"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0a7d5958-395e-4691-b943-237054a0b561", "can't obtain config")
		return env.ErrorDispatch(err)
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathVantagePoint,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Vantage Point",
		Description: "Vantage Point settings",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathVantagePointScheduleEnabled,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Enable Vantage Point scheduled update",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathVantagePointUploadPath,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Upload Path",
		Description: "Path to uploaded files",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	validateScheduleHour := func(value interface{}) (interface{}, error) {
		_ = setScheduleHour(value)

		return value, nil
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathVantagePointScheduleHour,
		Value:       "0",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     hoursList,
		Label:       "Schedule hour",
		Description: "The hour for update inventory task execution",
		Image:       "",
	}, env.FuncConfigValueValidator(validateScheduleHour))

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
