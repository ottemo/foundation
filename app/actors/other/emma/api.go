package emma

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)


// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	// Public
	service.POST("emma/contact", APIEmmaAddContact)

	return nil
}

// APIEmmaAddContact - return message, after add contact
// - email should be specified in "email" argument
func APIEmmaAddContact(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	email := utils.InterfaceToString(requestData["email"])
	if email == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6372b9a3-29f3-4ea4-a19f-40051a8f330b", "email was not specified")
	}

	if !utils.ValidEmailAddress(email) {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b54b0917-acc0-469f-925e-8f85a1feac7b", "The email address, " + email + ", is not in valid format.")
	}

	result, err := subscribe(email)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return result, nil
}

