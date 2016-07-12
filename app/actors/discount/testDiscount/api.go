package testDiscount

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.POST("testDiscount/CreateTestRule", CreateTestRule)
	service.POST("testDiscount/CreateTestAction", CreateTestAction)
	service.GET("testDiscount/GetConfigForTestDiscount", GetConfigForTestDiscount)

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
	rule["type"] = "DiscountRule"
	rule["json"] = config.GetValue(ConstConfigPathTestDiscountRule)
	result["rule"] = rule

	action := make(map[string]interface{})
	action["type"] = "DiscountAction"
	action["json"] = config.GetValue(ConstConfigPathTestDiscountAction)
	result["action"] = action

	return result, nil
}
