package order

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()

	config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathLastIncrementID,
		Value:       0,
		Type:        "int",
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

	return nil
}
