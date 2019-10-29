package otto

import (
	"github.com/robertkrimen/otto/repl"
)

// GetID returns the script identifier
func (it *Script) GetID() string {
	return it.id
}

// Execute executing given script
func (it *Script) Execute(code string) (interface{}, error) {
	value, err := it.vm.Run(code)
	if err != nil {
		return nil, err
	}
	return value.Export()
}

// Get returns the script context variable
func (it *Script) Get(name string) (interface{}, error) {
	return it.vm.Get(name)
}

// Set specifies script context variable
func (it *Script) Set(name string, value interface{}) error {
	return it.vm.Set(name, value)
}

// Interact run script instance in interaction mode
func (it *Script) Interact() error {
	return repl.RunWithOptions(it.vm, repl.Options{Prompt: "otto> ", Autocomplete: true})
}
