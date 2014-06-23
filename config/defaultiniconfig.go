package config

import (
	app "github.com/ottemo/foundation/app"
	ini "github.com/vaughan0/go-ini"
)

// DefaultIniConfig contains all values defined in the INI configuration file
type DefaultIniConfig struct {
	iniFileValues map[string]string
}

func init() {
	instance := new(DefaultIniConfig)

	app.OnAppStart(instance.startup)
	RegisterIniConfig(instance)
}

func (dic *DefaultIniConfig) startup() error {

	iniFile, _ := ini.LoadFile("ottemo.ini")
	dic.iniFileValues = iniFile.Section("")

	err := OnConfigIniStart()

	return err
}

// List contains all INI values
func (dic *DefaultIniConfig) List() []string {
	result := make([]string, len(dic.iniFileValues))
	for itemName := range dic.iniFileValues {
		result = append(result, itemName)
	}
	return result
}

// GetValue retrieves the given value of the INI entity.
func (dic *DefaultIniConfig) GetValue(Name string) string {
	if value, present := dic.iniFileValues[Name]; present {
		return value
	}

	return ""
}
