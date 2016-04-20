package trustpilot

import (
	"github.com/ottemo/foundation/env"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "b2c1c442-36b9-4994-b5d1-7c948a7552bd", "can't obtain config")
	}

	// Trust pilot config elements
	//----------------------------

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilot,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Trust Pilot",
		Description: "Trust Pilot settings",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotEnabled,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Enable Service Emails",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotBusinessUnitID,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Business Unit ID",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotUsername,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Username",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotPassword,
		Value:       "",
		Type:        env.ConstConfigTypeSecret,
		Editor:      "password",
		Options:     "",
		Label:       "Password",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotAPIKey,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "API Key",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotAPISecret,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "API Secret",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path: ConstConfigPathTrustPilotEmailTemplate,
		Value: `Dear {{.Visitor.name}},
please provide your feedback on recently purchase
<br />
{{.Visitor.link}}`,
		Type:        env.ConstConfigTypeHTML,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Service review e-mail template",
		Description: "Contents of this email will be sent to customers two weeks after purchase",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotProductBrand,
		Value:       "",
		Type:        env.ConstConfigTypeText,
		Editor:      "line_text",
		Options:     "",
		Label:       "Default Product Brand",
		Description: "If the product has a brand attribute that will be used, otherwise the default value will be sent to TrustPilot",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
