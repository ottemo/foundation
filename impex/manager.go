package impex

import (
	"github.com/ottemo/foundation/env"
)

// RegisterImportCommand registers new command to import/export system
func RegisterImportCommand(commandName string, command InterfaceImpexImportCmd) error {
	if _, present := importCmd[commandName]; present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f4f50cfc3de64c76851802ad26790f44", commandName+" already registered in impex")
	}

	importCmd[commandName] = command

	return nil
}

// UnRegisterImportCommand un-registers command from import/export system
func UnRegisterImportCommand(commandName string) error {
	if _, present := importCmd[commandName]; !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a9df49e30d064afa8e3a7cde5a3b2a1d", "can't find registered command "+commandName)
	}

	delete(importCmd, commandName)

	return nil
}
