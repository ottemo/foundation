package order

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

const (
	ConstConfigPathOrderGroup            = "general.order"
	ConstConfigPathShippingEmailSubject  = "general.order.shipping_status_email_subject"
	ConstConfigPathShippingEmailTemplate = "general.order.shipping_status_email_template"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "9028114d-cfee-46b5-bdf1-dac0db954d22", "Unable to obtain configuration for Order")
		return env.ErrorDispatch(err)
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathLastIncrementID,
		Value:       0,
		Type:        env.ConstConfigTypeInteger,
		Editor:      "integer",
		Options:     "",
		Label:       "Last Order Increment ID: ",
		Description: "Do not change this value unless you know what you doing",
		Image:       "",
	},
		func(value interface{}) (interface{}, error) {
			return utils.InterfaceToInt(value), nil
		})

	lastIncrementID = utils.InterfaceToInt(config.GetValue(ConstConfigPathLastIncrementID))

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathOrderGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Order",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathShippingEmailSubject,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Shipping Status Email Subject",
		Description: "",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathShippingEmailTemplate,
		Value:       "",
		Type:        env.ConstConfigTypeText,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Shipping Status Email Template",
		Description: "",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
