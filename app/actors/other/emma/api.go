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
	service.POST("emma/", APIEmmaSubscribeEmail)

	return nil
}

// APIEmmaSubscribeEmail - return message, after subscribe
// - email should be specified in "email" argument
func APIEmmaSubscribeEmail(context api.InterfaceApplicationContext) (interface{}, error) {

	//If mailchimp is not enabled, ignore this handler and do nothing
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

	easy := curl.EasyInit()
	defer easy.Cleanup()

	// @todo
	account_id := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaAccountID))
	if account_id == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "feb3a463-622b-477e-a22d-c0a3fd1972dc", "account id was not specified")
	}

	public_api_key := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaPublicAPIKey))
	if public_api_key == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "feb3a463-622b-477e-a22d-c0a3fd1972dc", "public api key was not specified")
	}

	private_api_key := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaPrivateAPIKey))
	if private_api_key == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "feb3a463-622b-477e-a22d-c0a3fd1972dc", "private api key was not specified")
	}

	//default_group_id := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmmaDefaultGroupID))
	url := "https://api.e2ma.net/" + account_id + "/members/add"
	//var postData = [2]string
	//postData["email"] = email
	//postData["group_ids"] = default_group_id

	easy.Setopt(curl.OPT_URL, url)
	easy.Setopt(curl.OPT_USERPWD, public_api_key + ":" + private_api_key)
	//easy.Setopt(curl.OPT_POSTFIELDS, )
	//easy.Setopt(curl.OPT_HTTPHEADER, ["Content-type: application/json"])
	easy.Setopt(curl.OPT_SSL_VERIFYPEER, false)
//	curl_setopt($ch, CURLOPT_USERPWD, $public_api_key . ":" . $private_api_key);
//curl_setopt($ch, CURLOPT_URL, $url);

//curl_setopt($ch, CURLOPT_POST, count($member_data));
//curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($member_data));
//curl_setopt($ch, CURLOPT_HTTPHEADER, array('Content-type: application/json'));
//curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
//curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);

	// make a callback function
	fooTest := func(buf []byte, userdata interface{}) bool {
		println("DEBUG: size=>", len(buf))
		println("DEBUG: content=>", string(buf))
		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

	if err := easy.Perform(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	//result = array
	//
	//return result, nil
	return nil, nil
}

