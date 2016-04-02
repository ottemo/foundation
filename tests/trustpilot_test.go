package tests

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/ottemo/foundation/app/actors/other/trustpilot"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

type TPToken struct {
	Token      string
	StatusCode int
}

// TestTPToken
func TestTPToken(t *testing.T) {
	tpToken, err := obtainToken()
	if err != nil {
		t.Errorf("Error retrieving oauth token: '%s' ", err)
	}

	// check if token exists
	if tpToken != nil {
		if tpToken.Token == "" {
			t.Error("Trustpilot oauth token cannot be empty")
		}
		// error on any status code but 200
		if tpToken.StatusCode != 200 {
			t.Errorf("Invalid http status received from Trustpilot.  StatusCode: '%s'.", tpToken.StatusCode)
		}
	}
}

// TestTPOrders
func TestTPOrders(t *testing.T) {

	tpToken, err := obtainToken()
	if err != nil {
		t.Errorf("Error retrieving oauth token: '%s' ", err)
	}

	addReviewLinkToOrders(tpToken.Token)
}

// Obtain OAUTH token
func obtainToken() (*TPToken, error) {
	// starting application in test mode
	err := StartAppInTestingMode()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	apiKey := "8N3P9QOvhrC4ej2PawPRoRD6O4BjjsWM"
	apiSecret := "LeiZXj9t7kf16dDG"
	userID := "engineering@ottemo.io"
	userPass := "Ottemo2016!"
	tokenURL := "https://api.trustpilot.com/v1/oauth/oauth-business-users-for-applications/accesstoken"
	// unitID := "54f078980000ff00057db649"
	// productReview := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotProductReviewURL))
	// serviceReview := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotServiceReviewURL))

	body := "grant_type=password&username=" + userID + "&password=" + userPass
	buffer := bytes.NewBuffer([]byte(body))

	apiCreds := []byte(apiKey + ":" + apiSecret)
	encodedCreds := base64.StdEncoding.EncodeToString(apiCreds)

	// https://api.trustpilot.com/v1/oauth/oauth-business-users-for-applications/accesstoken
	fmt.Printf("URL: %s\n", tokenURL)
	fmt.Printf("Buffer: %s\n", buffer)
	request, err := http.NewRequest("POST", tokenURL, buffer)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	request.Header.Set("Authorization", "Basic "+encodedCreds)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	jsonEntity, err := utils.DecodeJSONToStringKeyMap(responseBody)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	accessToken, present := jsonEntity["access_token"]
	fmt.Printf("Found Token: '%s'\n", utils.InterfaceToString(accessToken))
	if !present {
		return nil, env.ErrorDispatch(env.ErrorNew(trustpilot.ConstErrorModule, 1, "c4fae732-d524-44ea-9a20-d0115105a89a", "Trustpilot did not return an access token."))
	}

	token := &TPToken{
		Token:      utils.InterfaceToString(accessToken),
		StatusCode: response.StatusCode,
	}

	return token, nil
}

func addReviewLinkToOrders(token string) {

	fmt.Printf("Using Token to create links: %s\n", token)
	// starting application in test mode
	err := StartAppInTestingMode()
	if err != nil {
		panic(err)
	}

	startDate, err := time.Parse(time.UnixDate, "Tue Mar 1 00:00:00 PST 2016")
	fmt.Printf("Start time: %s\n", startDate)
	if err != nil {
		panic(err)
	}
	endDate, err := time.Parse(time.UnixDate, "Wed Mar 5 11:59:59 PST 2016")
	fmt.Printf("Start time: %s\n", endDate)
	if err != nil {
		panic(err)
	}
	orderList := order.GetOrdersCreatedBetween(startDate, endDate)
	fmt.Printf("Orders: %s\n", orderList)
}

//populate customer data

// populate order item data

// make call to TrustPilot to get unique link

// populate order and save
