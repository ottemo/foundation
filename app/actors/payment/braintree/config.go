package braintree

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(constErrorModule, env.ConstErrorLevelStartStop, "f9aac1c6-781b-410f-916b-4c884c19bdfb", "internal error, unable to obtain environment configuration")
	}

	// --------------------------------------
	// General

	err := config.RegisterItem(env.StructConfigItem{
		Path:  constGeneralConfigPathGroup,
		Label: "Braintree General",
		Type:  env.ConstConfigTypeGroup,
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:   constGeneralConfigPathEnvironment,
		Value:  constEnvironmentSandbox,
		Type:   env.ConstConfigTypeVarchar,
		Editor: "select",
		Options: map[string]string{
			constEnvironmentSandbox:    "Sandbox",
			constEnvironmentProduction: "Production"},
		Label:       "Environment",
		Description: "Change Braintree environment according to the workflow mode",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        constGeneralConfigPathMerchantID,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Merchant ID",
		Description: "Environment merchant ID",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        constGeneralConfigPathPublicKey,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Public Key",
		Description: "Environment public key",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        constGeneralConfigPathPrivateKey,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "PRIVATE Key",
		Description: "Environment PRIVATE key",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// --------------------------------------
	// Credit Card

	err = config.RegisterItem(env.StructConfigItem{
		Path:  constCCMethodConfigPathGroup,
		Label: "Braintree Credit Card",
		Type:  env.ConstConfigTypeGroup,
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:   constCCMethodConfigPathEnabled,
		Label:  "Enabled",
		Type:   env.ConstConfigTypeBoolean,
		Editor: "boolean",
	}, func(value interface{}) (interface{}, error) {
		return utils.InterfaceToBool(value), nil
	})
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:   constCCMethodConfigPathName,
		Label:  "Name in checkout",
		Value:  constCCMethodInternalName,
		Type:   env.ConstConfigTypeVarchar,
		Editor: "line_text",
	}, func(value interface{}) (interface{}, error) {
		if utils.CheckIsBlank(value) {
			return nil, env.ErrorNew(constErrorModule, env.ConstErrorLevelStartStop, "cc1f027c-6337-497e-a158-7d5842d50eae", "name in checkout can't be blank")
		}
		return value, nil
	})
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
