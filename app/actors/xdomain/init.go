package xdomain

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
)

// init performs self-initialization routine before app start
func init() {

	env.RegisterOnConfigIniStart(setupIniConfig)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupIniConfig reads the setting from the ottemo.ini file
func setupIniConfig() error {

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("xdomain.master", xdomainMasterURL); iniValue != "" {
			xdomainMasterURL = iniValue
		}
	}

	return nil
}
