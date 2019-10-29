package otto

import (
	"bytes"
	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/env"
	"github.com/robertkrimen/otto"
	"io"
	"sync"
)

// Package global constants
const (
	ConstSessionKey = "script_id"

	ConstErrorModule = "env/otto"
	ConstErrorLevel  = env.ConstErrorLevelService
)

var engine *ScriptEngine

// ScriptEngine is an implementer of InterfaceScriptEngine
type ScriptEngine struct {
	mutex     sync.RWMutex
	mappings  map[string]interface{}
	instances map[string]*Script
}

// Script is an implementer of InterfaceScriptEngine
type Script struct {
	id string
	vm *otto.Otto
}

// ApplicationContext is an implementor of api.InterfaceApplicationContext
type ApplicationContext struct {
	RequestParameters map[string]string
	RequestSettings   map[string]interface{}
	RequestArguments  map[string]string
	RequestContent    interface{}
	RequestFiles      map[string]io.Reader

	Session          api.InterfaceSession
	ContextValues    map[string]interface{}
	Result           interface{}
	Response         *bytes.Buffer
	ResponseSettings map[string]interface{}
}
