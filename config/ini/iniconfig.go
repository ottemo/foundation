package ini

import (
	app "github.com/ottemo/foundation/app"
	config "github.com/ottemo/foundation/config"
	ini "github.com/vaughan0/go-ini"
)

type DefaultIniConfig struct {
	iniFileValues map[string]string
}

func init() {
	instance := new(DefaultIniConfig)

	app.OnAppStart(instance.startup)
	config.RegisterIniConfig(instance)
}

func (dic *DefaultIniConfig) startup() error {

	iniFile, _ := ini.LoadFile("ottemo.ini")
	dic.iniFileValues = iniFile.Section("")

	err := config.OnConfigIniStart()

	return err
}

func (dic *DefaultIniConfig) List() []string {
	result := make([]string, len(dic.iniFileValues))
	for itemName, _ := range dic.iniFileValues {
		result = append(result, itemName)
	}
	return result
}

func (dic *DefaultIniConfig) GetValue(Name string) string {
	if value, present := dic.iniFileValues[Name]; present {
		return value
	} else {
		return ""
	}
}
