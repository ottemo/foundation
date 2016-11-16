package emma

import (
	"fmt"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/andelf/go-curl"
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

	//If emma is not enabled, ignore this request and do nothing
	if enabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEmmaEnabled)); !enabled {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b3548446-1453-4862-a649-393fc0aafda1", "emma does not active")
	}

	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	email := utils.InterfaceToString(requestData["email"])
	if email == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6372b9a3-29f3-4ea4-a19f-40051a8f330b", "email was not specified")
	}

	if !utils.ValidEmailAddress(email) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b54b0917-acc0-469f-925e-8f85a1feac7b", "The email address, " + email + ", is not in valid format.")
	}

	var account_id = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaAccountID))
	if account_id == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "88111f54-e8a1-4c43-bc38-0e660c4caa16", "account id was not specified")
	}

	var public_api_key = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaPublicAPIKey))
	if public_api_key == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1b5c42f5-d856-48c5-98a2-fd8b5929703c", "public api key was not specified")
	}

	var private_api_key = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaPrivateAPIKey))
	if private_api_key == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e0282f80-43b4-418e-a99b-60805e74c75d", "private api key was not specified")
	}

	var url = ConstEmmaApiUrl + account_id + "/members/add"

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

