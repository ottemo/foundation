package emma

import (
	"fmt"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/andelf/go-curl"
	"github.com/ottemo/foundation/utils"
)

const (
	EMMA_API_URL = "https://api.e2ma.net/"
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

	//If emma is not enabled, ignore this request and do nothing
	if enabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEmmaEnabled)); !enabled {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "feb3a463-622b-477e-a22d-c0a3fd1972dc", "emma does not active")
	}

	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	email := utils.InterfaceToString(requestData["email"])
	if email == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "feb3a463-622b-477e-a22d-c0a3fd1972dc", "email was not specified")
	}

	var account_id = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaAccountID))
	if account_id == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "feb3a463-622b-477e-a22d-c0a3fd1972dc", "account id was not specified")
	}

	var public_api_key = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaPublicAPIKey))
	if public_api_key == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "feb3a463-622b-477e-a22d-c0a3fd1972dc", "public api key was not specified")
	}

	var private_api_key = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaPrivateAPIKey))
	if private_api_key == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "feb3a463-622b-477e-a22d-c0a3fd1972dc", "private api key was not specified")
	}

	var url = EMMA_API_URL + account_id + "/members/add"

	postData := map[string]interface{}{"email": email}
	postDataJson := utils.EncodeToJSONString(postData)

	easy := curl.EasyInit()
	defer easy.Cleanup()
	easy.Setopt(curl.OPT_URL, url)
	easy.Setopt(curl.OPT_USERPWD, public_api_key + ":" + private_api_key)
	easy.Setopt(curl.OPT_POSTFIELDS, postDataJson)
	easy.Setopt(curl.OPT_HTTPHEADER, []string{"Content-type: application/json"})
	easy.Setopt(curl.OPT_SSL_VERIFYPEER, false)
	easy.Setopt(curl.OPT_POST, 1)

	responseBody := ""
	easy.Setopt(curl.OPT_WRITEFUNCTION, func(buf []byte, userdata interface{}) bool {
		responseBody += string(buf)
		return true
	})

	if err := easy.Perform(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	var result = "Error occurred";
	if responseCode, err := easy.Getinfo(curl.INFO_RESPONSE_CODE); responseCode == 200 && err == nil {
		jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if isAdded, isset := jsonResponse["added"]; isset {
			result = "E-mail was added successfully"
			if isAdded == false {
				result = "E-mail already added"
			}
		}
	}

	return result, nil
}

