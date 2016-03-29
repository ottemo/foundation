package trustpilot

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"

	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
)

// checkoutSuccessHandler is a handler for checkout success event which sends order information to TrustPilot
func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {

	var checkoutOrder order.InterfaceOrder
	if eventItem, present := eventData["order"]; present {
		if typedItem, ok := eventItem.(order.InterfaceOrder); ok {
			checkoutOrder = typedItem
		}
	}

	var checkoutCart cart.InterfaceCart
	if eventItem, present := eventData["cart"]; present {
		if typedItem, ok := eventItem.(cart.InterfaceCart); ok {
			checkoutCart = typedItem
		}
	}

	if checkoutOrder != nil && checkoutCart != nil {
		go sendOrderInfo(checkoutOrder, checkoutCart)
	}

	return true
}

// sendOrderInfo is a asynchronously calling request to TrustPilot
// 1. get a token from trustpilot
// 2. get a product review link
// 3. get a service review link, and set the product review url as the redirect once they complete the service review
// 4. set the service url on the order object
func sendOrderInfo(checkoutOrder order.InterfaceOrder, currentCart cart.InterfaceCart) error {
	if utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathTrustPilotEnabled)) {
		// taking TrustPilot settings into variables
		trustPilotAPIKey := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAPIKey))
		trustPilotAPISecret := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAPISecret))
		trustPilotBusinessUnitID := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotBusinessUnitID))
		trustPilotUsername := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotUsername))
		trustPilotPassword := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotPassword))
		trustPilotAccessTokenURL := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAccessTokenURL))
		trustPilotProductReviewURL := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotProductReviewURL))
		trustPilotServiceReviewURL := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotServiceReviewURL))

		trustPilotTestMode := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathTrustPilotTestMode))

		// verification of configuration values
		if trustPilotAPIKey != "" && trustPilotAPISecret != "" && trustPilotBusinessUnitID != "" && trustPilotUsername != "" &&
			trustPilotPassword != "" && trustPilotAccessTokenURL != "" && trustPilotProductReviewURL != "" && trustPilotServiceReviewURL != "" {

			/**
			 * 1. Get the access token
			 */

			bodyString := "grant_type=password&username=" + trustPilotUsername + "&password=" + trustPilotPassword
			buffer := bytes.NewBuffer([]byte(bodyString))

			valueAMIKeySecret := []byte(trustPilotAPIKey + ":" + trustPilotAPISecret)
			encodedString := base64.StdEncoding.EncodeToString(valueAMIKeySecret)

			// https://api.trustpilot.com/v1/oauth/oauth-business-users-for-applications/accesstoken
			request, err := http.NewRequest("POST", trustPilotAccessTokenURL, buffer)
			if err != nil {
				return env.ErrorDispatch(err)
			}

			request.Header.Set("Authorization", "Basic "+encodedString)
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				return env.ErrorDispatch(err)
			}
			defer response.Body.Close()

			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return env.ErrorDispatch(err)
			}

			if response.StatusCode >= 300 {
				errMsg := "Non 200 response while trying to get trustpilot access token: StatusCode:" + response.StatusCode + " Body:" + responseBody
				err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "376b178e-6cbf-4b4e-a3a8-fd65251d176b", errMsg)
				env.LogError(err)
				return err
			}

			jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
			if err != nil {
				return env.ErrorDispatch(err)
			}

			if accessToken, present := jsonResponse["access_token"]; present {
				/**
				 * 2. Create product review invitation link
				 *
				 * https://developers.trustpilot.com/product-reviews-api
				 *
				 * Given information about the consumer and the product(s) purchased, get a link that can be sent to
				 * the consumer to request reviews.
				 */

				cartItems := currentCart.GetItems()

				requestData := make(map[string]interface{})
				customerEmail := utils.InterfaceToString(checkoutOrder.Get("customer_email"))
				customerName := checkoutOrder.Get("customer_name")
				checkoutOrderID := checkoutOrder.GetID()

				if trustPilotTestMode {
					customerEmail = strings.Replace(customerEmail, "@", "_test@", 1)
				}

				requestData["consumer"] = map[string]interface{}{
					"email": customerEmail,
					"name":  customerName,
				}

				requestData["referenceId"] = checkoutOrderID
				requestData["locale"] = "en-US"

				mediaStorage, err := media.GetMediaStorage()
				if err != nil {
					return env.ErrorDispatch(err)
				}

				var productsOrdered []map[string]string

				// filling request with products information
				for _, productItem := range cartItems {
					currentProductID := productItem.GetProductID()
					currentProduct := productItem.GetProduct()

					mediaPath, err := mediaStorage.GetMediaPath("product", currentProductID, "image")
					if err != nil {
						return env.ErrorDispatch(err)
					}

					productOptions := productItem.GetOptions()
					productBrand := ConstProductBrand
					if brand, present := productOptions["brand"]; present {
						productBrand = utils.InterfaceToString(brand)
					}

					productInfo := map[string]string{
						"productUrl": app.GetStorefrontURL("product/" + currentProductID),
						"imageUrl":   app.GetStorefrontURL(mediaPath + currentProduct.GetDefaultImage()),
						"name":       currentProduct.GetName(),
						"sku":        currentProduct.GetSku(),
						"brand":      productBrand,
					}

					productsOrdered = append(productsOrdered, productInfo)
				}

				requestData["products"] = productsOrdered

				// https://api.trustpilot.com/v1/private/product-reviews/business-units/{businessUnitId}/invitation-links
				trustPilotProductReviewURL = strings.Replace(trustPilotProductReviewURL, "{businessUnitId}", trustPilotBusinessUnitID, 1)

				jsonString := utils.EncodeToJSONString(requestData)
				buffer := bytes.NewBuffer([]byte(jsonString))

				request, err := http.NewRequest("POST", trustPilotProductReviewURL, buffer)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				request.Header.Set("Content-Type", "application/json")
				request.Header.Set("Authorization", "Bearer "+utils.InterfaceToString(accessToken))

				client := &http.Client{}
				response, err := client.Do(request)
				if err != nil {
					return env.ErrorDispatch(err)
				}
				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				if response.StatusCode >= 300 {
					errMsg := "Non 200 response while trying to get trustpilot review link: StatusCode:" + response.StatusCode + " Body:" + responseBody
					err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e75b28c7-0da2-475b-8b65-b1a09f1f6926", errMsg)
					env.LogError(err)
					return err
				}

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				reviewLink, ok := jsonResponse["reviewUrl"]
				if !ok {
					errorMessage := "Review link empty, "
					if jsonMessage, present := jsonResponse["message"]; present {
						errorMessage += "error message: " + utils.InterfaceToString(jsonMessage)
					} else {
						errorMessage += "no error message provided"
					}
					env.LogError(env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "c53fd02f-2f5d-4111-8318-69a2cc2d2259", errorMessage))
					return err
				}

				/**
				 * 3. Generate service review invitation link
				 *
				 * https://developers.trustpilot.com/invitation-api#Generate service review invitation link
				 *
				 * Generate a unique invitation link that can be sent to a consumer by email or website. Use the request
				 * parameter called redirectURI to take the user to a product review link after the user has left a
				 * service review.
				 */

				// make service review link with the same token and product review link
				requestData = map[string]interface{}{
					"referenceId": checkoutOrderID,
					"email":       customerEmail,
					"name":        customerName,
					"locale":      "en-US",
					"redirectUri": reviewLink,
				}

				// https://invitations-api.trustpilot.com/v1/private/business-units/{businessUnitId}/invitation-links
				trustPilotServiceReviewURL = strings.Replace(trustPilotServiceReviewURL, "{businessUnitId}", trustPilotBusinessUnitID, 1)

				jsonString = utils.EncodeToJSONString(requestData)
				buffer = bytes.NewBuffer([]byte(jsonString))

				request, err = http.NewRequest("POST", trustPilotServiceReviewURL, buffer)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				request.Header.Set("Content-Type", "application/json")
				request.Header.Set("Authorization", "Bearer "+utils.InterfaceToString(accessToken))

				response, err = client.Do(request)
				if err != nil {
					return env.ErrorDispatch(err)
				}
				defer response.Body.Close()

				responseBody, err = ioutil.ReadAll(response.Body)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				if response.StatusCode >= 300 {
					errMsg := "Non 200 response while trying to get trustpilot review link: StatusCode:" + response.StatusCode + " Body:" + responseBody
					err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e75b28c7-0da2-475b-8b65-b1a09f1f6926", errMsg)
					env.LogError(err)
					return err
				}

				jsonResponse, err = utils.DecodeJSONToStringKeyMap(responseBody)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				serviceReviewLink, ok := jsonResponse["url"]
				if !ok {
					errorMessage := "Service review link empty, "
					if jsonMessage, present := jsonResponse["message"]; present {
						errorMessage += "error message: " + utils.InterfaceToString(jsonMessage)
					} else {
						errorMessage += "no error message provided"
					}
					env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e528633c-9413-41b0-bfe8-8cee581a616c", errorMessage))
					return nil
				}

				/**
				 * 4. Update order with the service review link
				 */

				orderCustomInfo := utils.InterfaceToMap(checkoutOrder.Get("custom_info"))
				orderCustomInfo[ConstOrderCustomInfoLinkKey] = serviceReviewLink
				orderCustomInfo[ConstOrderCustomInfoSentKey] = false

				err = checkoutOrder.Set("custom_info", orderCustomInfo)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				err = checkoutOrder.Save()
				if err != nil {
					return env.ErrorDispatch(err)
				}

			} else {
				env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "1293708d-9638-455a-8d49-3a387f086181", "Trustpilot didn't return an access token for our request")
				env.LogError(err)
				return err
			}
		} else {
			env.LogError(env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "22207d49-e001-4666-8501-26bf5ef0926b", "Some trustpilot settings are not configured"))
			return err
		}
	}
	return nil
}
