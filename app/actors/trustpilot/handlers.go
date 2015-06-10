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

		// config values validation
		if trustPilotAPIKey != "" && trustPilotAPISecret != "" && trustPilotBusinessUnitID != "" && trustPilotUsername != "" &&
			trustPilotPassword != "" && trustPilotAccessTokenURL != "" && trustPilotProductReviewURL != "" {

			// making request to get authentication token required for following requests
			bodyString := "grant_type=password&username=" + trustPilotUsername + "&password=" + trustPilotPassword
			buffer := bytes.NewBuffer([]byte(bodyString))

			valueAMIKeySecret := []byte(trustPilotAPIKey + ":" + trustPilotAPISecret)
			encodedString := base64.StdEncoding.EncodeToString(valueAMIKeySecret)

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

			jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
			if err != nil {
				return env.ErrorDispatch(err)
			}

			if accessToken, present := jsonResponse["access_token"]; present {
				// access token received - making requests
				cartItems := currentCart.GetItems()

				requestData := make(map[string]interface{})
				customerEmail := utils.InterfaceToString(checkoutOrder.Get("customer_email"))

				if ConstTestMode {
					customerEmail = strings.Replace(customerEmail, "@", "_test@", 1)
				}

				requestData["consumer"] = map[string]interface{}{
					"email": customerEmail,
					"name":  checkoutOrder.Get("customer_name"),
				}

				requestData["referenceId"] = checkoutOrder.GetID()
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
					productBrand := ""
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

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				if _, ok := jsonResponse["reviewUrl"]; !ok {
					errorMessage := "Review link empty, "
					if jsonMessage, present := jsonResponse["message"]; present {
						errorMessage += "error message: " + utils.InterfaceToString(jsonMessage)
					} else {
						errorMessage += "no error message provided"
					}
					return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "c53fd02f-2f5d-4111-8318-69a2cc2d2259", errorMessage)
				}

			} else {
				return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "1293708d-9638-455a-8d49-3a387f086181", "access token is empty")
			}
		} else {
			return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "22207d49-e001-4666-8501-26bf5ef0926b", "some of trust pilot settings are blank")
		}
	}
	return nil
}
