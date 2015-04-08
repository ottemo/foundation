package ini

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstIniGlobalSection  = ""  // ini file section name to be used as default section
	ConstAskForValuePrefix = "?" // prefix used before default ini value to be asked in console if not set

	ConstCmdArgStoreAllFlag = "--iniStoreAll"
	ConstCmdArgSectionName  = "--iniSection="
	ConstCmdArgTestFlag     = "--test"

	ConstEnvironmentIniFile    = "OTTEMO_INI"
	ConstEnvironmentIniSection = "OTTEMO_MODE"

	ConstTestSectionName = "test"
	ConstDefaultIniFile  = "ottemo.ini"

	ConstErrorModule = "env/ini"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// DefaultIniConfig is a default implementer of InterfaceIniConfig
type DefaultIniConfig struct {
	iniFilePath string

	iniFileValues  map[string]map[string]string
	currentSection string

	keysToStore map[string]bool
	storeAll    bool
}
