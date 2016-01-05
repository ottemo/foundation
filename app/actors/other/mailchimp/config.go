package mailchimp

import (
	"github.com/ottemo/foundation/env"
)

func setupConfig() error {

	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "6b78d38a-35c5-4aa2-aec1-eaa16830ff61", "Error configuring Mailchimp module")
	}

	if err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailchimp,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Mailchimp",
		Description: "Mailchimp Settings",
		Image:       "",
	}, nil); err != nil {
		return env.ErrorDispatch(err)
	}

	if err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailchimpEnabled,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Mailchimp Enabled",
		Description: "Enable Mailchimp integration(defaults to false)",
		Image:       "",
	}, nil); err != nil {
		return env.ErrorDispatch(err)
	}

	if err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailchimpAPIKey,
		Value:       false,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Mailchimp API Key",
		Description: "Enter your Mailchimp API Key",
		Image:       "",
	}, nil); err != nil {
		return env.ErrorDispatch(err)
	}

	if err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailchimpBaseURL,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Mailchimp Base Url",
		Description: "Defines the base url for this account",
		Image:       "",
	}, nil); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
