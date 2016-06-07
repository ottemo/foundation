package fedex

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	if config := env.GetConfig(); config != nil {
		err := config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathGroup,
			Value:       nil,
			Type:        env.ConstConfigTypeGroup,
			Editor:      "",
			Options:     nil,
			Label:       "FedEx",
			Description: "Federal express shipping method",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathEnabled,
			Value:       false,
			Type:        env.ConstConfigTypeBoolean,
			Editor:      "boolean",
			Options:     nil,
			Label:       "Enabled",
			Description: "enables/disables shipping method in checkout",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathTitle,
			Value:       "Federal Express",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "line_text",
			Options:     nil,
			Label:       "Title",
			Description: "shipping method name in checkout",
			Image:       "",
		}, func(value interface{}) (interface{}, error) {
			if utils.CheckIsBlank(value) {
				err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "e37e8ab2-ab87-40d0-b368-4e2f930d79b1", "can't be blank")
				return nil, env.ErrorDispatch(err)
			}
			return value, nil
		})

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathGateway,
			Value:       "https://wsbeta.fedex.com:443/web-services",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "line_text",
			Options:     nil,
			Label:       "Gateway",
			Description: "web services gateway",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathKey,
			Value:       "",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "line_text",
			Options:     nil,
			Label:       "Account Key",
			Description: "FedEx account key",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathPassword,
			Value:       "",
			Type:        env.ConstConfigTypeSecret,
			Editor:      "password",
			Options:     nil,
			Label:       "Account Password",
			Description: "FedEx account password",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathNumber,
			Value:       "",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "line_text",
			Options:     nil,
			Label:       "Account Number",
			Description: "FedEx account number",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathMeter,
			Value:       "",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "line_text",
			Options:     nil,
			Label:       "Account Meter",
			Description: "FedEx account meter value",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDefaultWeight,
			Value:       0.1,
			Type:        env.ConstConfigTypeDecimal,
			Editor:      "decimal",
			Options:     nil,
			Label:       "Default weight",
			Description: "Will be used if product do not have this value (in pounds)",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathAllowedMethods,
			Value:       "",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "multi_select",
			Options:     ConstShippingMethods,
			Label:       "Allowed methods",
			Description: "To customer will be proposed only allowed methods",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDropoff,
			Value:       "REGULAR_PICKUP",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "select",
			Options:     ConstShippingDropoff,
			Label:       "Dropoff",
			Description: "dropoff method",
			Image:       "",
		}, func(value interface{}) (interface{}, error) {
			stringValue := utils.InterfaceToString(value)
			if _, present := ConstShippingDropoff[stringValue]; !present {
				err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "c2129328-0477-4741-a09f-fc96863e936d", "wrong value")
				return nil, env.ErrorDispatch(err)
			}
			return value, nil
		})

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathPackaging,
			Value:       "FEDEX_PAK",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "select",
			Options:     ConstShippingPackaging,
			Label:       "Packing",
			Description: "packing method",
			Image:       "",
		}, func(value interface{}) (interface{}, error) {
			stringValue := utils.InterfaceToString(value)
			if _, present := ConstShippingPackaging[stringValue]; !present {
				err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "890ae42a-b890-4481-81ba-5165a78acd18", "wrong value")
				return nil, env.ErrorDispatch(err)
			}
			return value, nil
		})

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDebugLog,
			Value:       false,
			Type:        env.ConstConfigTypeBoolean,
			Editor:      "boolean",
			Options:     nil,
			Label:       "Debug log",
			Description: "enables/disables shipping method debug log",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return env.ErrorDispatch(err)
		}
	} else {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "ba621a64-5ab2-4f65-8953-1165e4a401ac", "Unable to obtain configuration for FedEx")
		return env.ErrorDispatch(err)
	}

	return nil
}
