package emma

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"io"
)

const (
	constEmmaAPIURL = "https://api.e2ma.net/"
)

type emmaCredentialsType struct {
	AccountID     string
	PublicAPIKey  string
	PrivateAPIKey string
}

type emmaSubscribeInfoType struct {
	Email        string
	GroupIDsList []string
}

type emmaServiceType struct{}

type interfaceEmmaService interface {
	subscribe(subscribeInfo emmaSubscribeInfoType) (string, error)
}

// newEmmaService returns new emmaServiceType instance
func newEmmaService() *emmaServiceType {
	return &emmaServiceType{}
}

// subscribe
func (it *emmaServiceType) subscribe(credentials emmaCredentialsType, subscribeInfo emmaSubscribeInfoType) (string, error) {
	var result = "Error occurred"

	postData := map[string]interface{}{"email": subscribeInfo.Email}
	if len(subscribeInfo.GroupIDsList) > 0 {
		postData["group_ids"] = subscribeInfo.GroupIDsList
	}

	postDataJSON := utils.EncodeToJSONString(postData)

	var url = constEmmaAPIURL + credentials.AccountID + "/members/add"

	buf := bytes.NewBuffer([]byte(postDataJSON))
	request, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(credentials.PublicAPIKey, credentials.PrivateAPIKey)

	client := &http.Client{}
	response, err := client.Do(request)
	// require http response code of 200 or error out
	if response.StatusCode != http.StatusOK {

		var status string
		if response == nil {
			status = "nil"
		} else {
			status = response.Status
		}

		return result, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cad8ad77-dd4c-440c-ada2-1e315b706175", "Unable to subscribe visitor to Emma list, response code returned was "+status)
	}
	defer func (c io.ReadCloser){
		if err := c.Close(); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "039165b1-f795-4d3d-ac65-16a9087a3174", err.Error())
		}
	}(response.Body)

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	if isAdded, isset := jsonResponse["added"]; isset {
		result = "E-mail was added successfully"
		if isAdded == false {
			result = "E-mail already added"
		}
	}

	return result, nil
}
