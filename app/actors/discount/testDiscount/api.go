package testDiscount

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("testDiscount/CreateTestRule", api.ConstRESTOperationCreate, CreateTestRule)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("testDiscount/CreateTestAction", api.ConstRESTOperationCreate, CreateTestAction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("testDiscount/GetConfigForTestDiscount", api.ConstRESTOperationGet, GetConfigForTestDiscount)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

func CreateTestRule(context api.InterfaceApplicationContext) (interface{}, error) {
	config := env.GetConfig()

	var setValue interface{}

	setValue = context.GetRequestContent()
	configPath := ConstConfigPathTestDiscountRule

	err := config.SetValue(configPath, setValue)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "rule was saved successfully", nil
}

func CreateTestAction(context api.InterfaceApplicationContext) (interface{}, error) {
	config := env.GetConfig()

	var setValue interface{}

	setValue = context.GetRequestContent()
	configPath := ConstConfigPathTestDiscountAction

	err := config.SetValue(configPath, setValue)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "action was saved successfully", nil
}

func GetConfigForTestDiscount(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})

	config := env.GetConfig()

	rule := make(map[string]interface{})
	rule["type"]   = "DiscountRule"
	rule["json"]   = config.GetValue(ConstConfigPathTestDiscountRule)
	result["rule"] = rule

	action := make(map[string]interface{})
	action["type"]   = "DiscountAction"
	action["json"]   = config.GetValue(ConstConfigPathTestDiscountAction)
	result["action"] = action

	return result, nil
}

