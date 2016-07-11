package trustpilot

import (
	"time"
	"net/http"
	"io/ioutil"
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	service := api.GetRestService()
	service.GET("trustpilot/summaries", APIGetTrustpilotSummaries)
	return nil
}

// APIGetTrustpilotSummaries
// Sends a request to obtain review summaries from the Trustpilot
// Caches the response for 1 hour
func APIGetTrustpilotSummaries (context api.InterfaceApplicationContext) (interface{}, error) {
	isEnabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathTrustPilotEnabled))
	if !isEnabled {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d535ccc0-68ec-4249-8ec5-e6962d965ffc", "Trustpilot integration is disabled")
	}

	if summariesCache == nil || time.Since(lastTimeSummariesUpdate).Hours() > 1 {
		// Get configuration values
		apiKey := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAPIKey))
		apiSecret := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAPISecret))
		apiUsername := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotUsername))
		apiPassword := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotPassword))
		businessID := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotBusinessUnitID))

		// Verify the configuration values
		configs := []string{apiKey, apiSecret, apiUsername, apiPassword, businessID}
		if hasEmpty(configs) {
			return nil, env.ErrorDispatch(env.ErrorNew(ConstErrorModule, 1, "92485c24-66d4-4276-8978-88dabf2a47ac", "Some trustpilot settings are not configured"))
		}

		// Init credentials for the access token request
		credentials := tpCredentials{
			username:  apiUsername,
			password:  apiPassword,
			apiKey:    apiKey,
			apiSecret: apiSecret,
		}

		// Get the access token
		accessToken, err := getAccessToken(credentials)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// Send the request to obtain review summaries
		ratingUrl := strings.Replace(ConstRatingSummaryUrl, "{businessUnitId}", businessID, 1)
		request, err := http.NewRequest("GET", ratingUrl, nil)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer " + accessToken)

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

		if response.StatusCode >= 300 {
			errMsg := "Non 200 response while trying to get trustpilot reviews summaries: StatusCode:" + response.Status
			err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "198fd7b0-917a-4bdc-add8-2402876281ae", errMsg)
			fields := env.LogFields{
				"accessToken":        accessToken,
				"businessID":         businessID,
				"responseBody":       responseBody,
			}
			env.LogEvent(fields, "trustpilot-reviews-summary-error")
			return nil, env.ErrorDispatch(err)
		}

		// Retrieve the review summaries from the response
		jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		summaries, ok := jsonResponse["summaries"]
		if !ok {
			errorMessage := "Reviews summaries are empty"
			return nil, env.ErrorNew(ConstErrorModule, 1, "7329b79e-cf91-4663-a1cd-2776d56c648b", errorMessage)
		}

		return summaries, nil

		// Put the summaries in the cache and remember the request time
		summariesCache = summaries
		lastTimeSummariesUpdate = time.Now();
	}

	return summariesCache, nil
}
