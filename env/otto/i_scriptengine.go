package otto

import (
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"
	"github.com/robertkrimen/otto"
)

// GetScriptName returns ScriptEngine instance name
func (it *ScriptEngine) GetScriptName() string {
	return "Otto"
}

// GetScriptInstance returns new instance for scripting
func (it *ScriptEngine) GetScriptInstance(id string) env.InterfaceScript {

	if instance, present := it.instances[id]; present {
		return instance
	}

	script := new(Script)
	script.vm = otto.New()
	script.id = id

	it.mutex.Lock()
	defer it.mutex.Unlock()

	for key, value := range it.mappings {
		script.vm.Set(key, value)
	}

	// TODO: implement instances cleanup (lifetime based)
	it.instances[script.id] = script

	return script
}

// Get returns value which is available for a new script instances
func (it *ScriptEngine) Get(path string) (interface{}, error) {
	it.mutex.Lock()
	defer it.mutex.Unlock()

	return utils.MapGetPathValue(it.mappings, path)
}

// Set specifies a value for all new script instances
func (it *ScriptEngine) Set(path string, value interface{}) error {
	it.mutex.Lock()
	defer it.mutex.Unlock()

	if value == nil {
		return utils.MapSetPathValue(it.mappings, path, value, true)
	}
	return utils.MapSetPathValue(it.mappings, path, value, false)
}
