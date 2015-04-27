package grouping

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine before app start
func init() {
	app.OnAppStart(initListners)
	env.RegisterOnConfigStart(setupConfig)
}

// init Listeners for current model
func initListners() error {

	env.EventRegisterListener("api.cart.updatedCart", updateCartHandler)

	return nil
}
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "701e85e4-b63c-48f4-a990-673ba0ed6a2a", "can't obtain config")
	}

	// Grouping rules config setup
	//---------
	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstGroupingConfigPath,
		Value:       nil,
		Type:        env.ConstConfigTypeText,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Rules for grouping items",
		Description: `decribe products that will be grouped; type [][][]map[string]interface{}, example: [[[{"pid":"pid1","qty":"n"},...],...],[[{"options":{},"pid":"resultpid1","qty":"n"}],...]] `,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}
	return nil
}
