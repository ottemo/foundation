package visitor

import (
	"github.com/ottemo/foundation/env"
)

// configuration constants
const (
	ConstConfigPathLostPasswordEmailSubject  = "general.mail.lost_password_email_subject"
	ConstConfigPathLostPasswordEmailTemplate = "general.mail.lost_password_email_template"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "b196712b-99c3-4510-bd58-935c5eedfc9f", "Unable to obtain configuration")
		return env.ErrorDispatch(err)
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathLostPasswordEmailSubject,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Lost Password Email Subject",
		Description: "",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathLostPasswordEmailTemplate,
		Value:       "",
		Type:        env.ConstConfigTypeText,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Lost Password Email Template",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
