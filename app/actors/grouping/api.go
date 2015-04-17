package grouping

import (
	"io/ioutil"
	"bytes"
"os"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("grouper/rules", api.ConstRESTOperationGet, APIGetRules)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("grouper/rules", api.ConstRESTOperationUpdate, APIUpdateRules)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	return nil
}

// APIGetRules get rules from saved in file
func APIGetRules(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}
	var result
	// read whole the file
	b, err := ioutil.ReadFile("group_rules")
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	str := string(b)
	result[0] = str


	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return result, nil
}

// APIUpdateRules save a posted rule to the file
func APIUpdateRules(context api.InterfaceApplicationContext) (interface{}, error) {

	var mode os.FileMode()
	mode = 0644
	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	rule := context.GetRequestArgument("rule")

	buff := []byte(rule)

	// write whole the body
	err := ioutil.WriteFile("group_rules", buff, mode)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}


	return "ok", nil
}
