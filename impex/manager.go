package impex

import (
	"github.com/ottemo/foundation/env"
)

// RegisterImportCommand registers new command to import/export system
func RegisterImportCommand(commandName string, command InterfaceImpexImportCmd) error {
	if _, present := importCmd[commandName]; present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f4f50cfc-3de6-4c76-8518-02ad26790f44", commandName+" already registered in impex")
	}

	importCmd[commandName] = command

	return nil
}

// UnRegisterImportCommand un-registers command from import/export system
func UnRegisterImportCommand(commandName string) error {
	if _, present := importCmd[commandName]; !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a9df49e3-0d06-4afa-8e3a-7cde5a3b2a1d", "can't find registered command "+commandName)
	}

	delete(importCmd, commandName)

	return nil
}

// RegisterImpexModel registers model instance which supports InterfaceImpexModel interface to import/export system
func RegisterImpexModel(name string, instance InterfaceImpexModel) error {
	if _, present := impexModels[name]; present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8c5dba4f-6c54-4551-931a-a60d581252ab", name+" model already registered in impex")
	}

	impexModels[name] = instance

	return nil
}

// UnRegisterImpexModel un-registers InterfaceImpexModel capable model from import/export system
func UnRegisterImpexModel(name string) error {
	if _, present := impexModels[name]; !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0b0863dd-f8ba-4a61-8438-603e649189e1", "can't find registered model "+name)
	}

	delete(impexModels, name)

	return nil
}
