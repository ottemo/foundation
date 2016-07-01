package cart

import (
	"github.com/ottemo/foundation/env"
)

func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "fed0dee4-3409-4533-a445-998d2290569a", "Cart module is unable to obtain configuration settings, using nil instead")
		return env.ErrorDispatch(err)
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:   ConstConfigPathCartAbandonEmailSendTime,
		Value:  "0",
		Type:   env.ConstConfigTypeVarchar,
		Editor: "select",
		Options: map[string]string{
			"0":   "Never",
			"-6":  "After 6 hours",
			"-24": "After 24 hours",
		},
		Label:       "Abandoned Cart Email - Send Time",
		Description: "If the customer abandons checkout, send them an email to complete their order.",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathCartAbandonEmailTemplate,
		Value:       "",
		Type:        env.ConstConfigTypeHTML,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Abandoned Cart Email - Template",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
