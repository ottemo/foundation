package trustpilot

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	service := api.GetRestService()
	service.GET("trustpilot/products/summaries", APIGetTrustpilotProductsSummaries)
	return nil
}

// APIGetTrustpilotProductsSummaries Makes a request to Trustpilot api to obtain a list of reviews summaries for every product,
// caches the response
// https://developers.trustpilot.com/product-reviews-api#Get product reviews summaries list
func APIGetTrustpilotProductsSummaries(context api.InterfaceApplicationContext) (interface{}, error) {
	isEnabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathTrustPilotEnabled))
	if !isEnabled {
		context.SetResponseStatusForbidden()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d535ccc0-68ec-4249-8ec5-e6962d965ffc", "Trustpilot integration is disabled")
	}

	// TODO: we should use some caching module instead of just global variables
	if summariesCache == nil || time.Since(lastTimeSummariesUpdate).Hours() >= 24 {
		// Get configuration values
		apiKey := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAPIKey))
		apiSecret := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAPISecret))
		apiUsername := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotUsername))
		apiPassword := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotPassword))
		businessID := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotBusinessUnitID))

		// Verify the configuration values
		configs := []string{apiKey, apiSecret, apiUsername, apiPassword, businessID}
		if hasEmpty(configs) {
			context.SetResponseStatusInternalServerError()
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
			context.SetResponseStatusInternalServerError()
			return nil, env.ErrorDispatch(err)
		}

		// Send the request to obtain review summaries
		ratingURL := strings.Replace(ConstRatingSummaryURL, "{businessUnitId}", businessID, 1)
		request, err := http.NewRequest("GET", ratingURL, nil)
		if err != nil {
			context.SetResponseStatusInternalServerError()
			return nil, env.ErrorDispatch(err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+accessToken)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			context.SetResponseStatusInternalServerError()
			return nil, env.ErrorDispatch(err)
		}
		defer func (c io.ReadCloser){
			if err := c.Close(); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "05ffa9ae-3b84-4835-a732-af4aa7f6bc2a", err.Error())
			}
		}(response.Body)

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			context.SetResponseStatusInternalServerError()
			return nil, env.ErrorDispatch(err)
		}

		if response.StatusCode >= 300 {
			errMsg := "Non 200 response while trying to get trustpilot reviews summaries: StatusCode:" + response.Status
			err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "198fd7b0-917a-4bdc-add8-2402876281ae", errMsg)
			fields := env.LogFields{
				"accessToken":  accessToken,
				"businessID":   businessID,
				"responseBody": responseBody,
			}
			env.LogEvent(fields, "trustpilot-reviews-summary-error")
			context.SetResponseStatusInternalServerError()
			return nil, env.ErrorDispatch(err)
		}

		// Retrieve the review summaries from the response
		jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
		if err != nil {
			context.SetResponseStatusInternalServerError()
			return nil, env.ErrorDispatch(err)
		}

		summaries, ok := jsonResponse["summaries"]
		if !ok {
			errorMessage := "Reviews summaries are empty"
			context.SetResponseStatusInternalServerError()
			return nil, env.ErrorNew(ConstErrorModule, 1, "7329b79e-cf91-4663-a1cd-2776d56c648b", errorMessage)
		}

		// Put the summaries in the cache and remember the request time
		summariesCache = summaries
		lastTimeSummariesUpdate = time.Now()
	}

	return summariesCache, nil
}
