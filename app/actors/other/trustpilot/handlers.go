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

const (
	requestLocale    = "en-US"
	accessTokenURL   = "https://api.trustpilot.com/v1/oauth/oauth-business-users-for-applications/accesstoken"
	serviceReviewURL = "https://invitations-api.trustpilot.com/v1/private/business-units/{businessUnitId}/invitation-links"
	productReviewURL = "https://api.trustpilot.com/v1/private/product-reviews/business-units/{businessUnitId}/invitation-links"
)

type tpCredentials struct {
	username  string
	password  string
	apiKey    string
	apiSecret string
}

// ProductReview is a container for TrustPilot product reviews
type ProductReview struct {
	ReferenceID string                 `json:"referenceId"`
	Locale      string                 `json:"locale"`
	Consumer    ProductReviewConsumer  `json:"consumer"`
	Products    []ProductReviewProduct `json:"products"`
}

// ProductReviewConsumer is a container for the TrustPilot reviewer's personal information
type ProductReviewConsumer struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ProductReviewProduct is a container object for a TrustPilot Product Review request based
// on items purchased by the consumer.
type ProductReviewProduct struct {
	ProductURL string `json:"productUrl"`
	ImageURL   string `json:"imageUrl"`
	Name       string `json:"name"`
	Sku        string `json:"sku"`
	Brand      string `json:"brand"`
}

// ServiceReview is a container object for the TrustPilot Service Review.  This is the first
// link the consumer will receive in an email request to review a TrustPilot registered service.
type ServiceReview struct {
	ReferenceID string `json:"referenceId"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Locale      string `json:"locale"`
	RedirectURI string `json:"redirectUri"`
}

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
		go SendOrderInfo(checkoutOrder, checkoutCart)
	}

	return true
}

// SendOrderInfo Makes requests to the trustpilot api to obtain an access token, then a product review url, then a
// service review url. This last url is then saved to the order so that we can setup a cronjob to email customers
// at any given time.
func SendOrderInfo(checkoutOrder order.InterfaceOrder, currentCart cart.InterfaceCart) error {

	isEnabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathTrustPilotEnabled))
	if !isEnabled {
		return nil
	}

	apiKey := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAPIKey))
	apiSecret := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAPISecret))
	apiUsername := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotUsername))
	apiPassword := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotPassword))
	businessID := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotBusinessUnitID))
	defaultProductBrand := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotProductBrand))

	// verification of configuration values
	configs := []string{apiKey, apiSecret, apiUsername, apiPassword, businessID, defaultProductBrand}
	if hasEmpty(configs) {
		return env.ErrorDispatch(env.ErrorNew(ConstErrorModule, 1, "22207d49-e001-4666-8501-26bf5ef0926b", "Some trustpilot settings are not configured"))
	}

	// Init some variables
	credentials := tpCredentials{
		username:  apiUsername,
		password:  apiPassword,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	customerEmail := utils.InterfaceToString(checkoutOrder.Get("customer_email"))
	customerName := utils.InterfaceToString(checkoutOrder.Get("customer_name"))
	orderID := checkoutOrder.GetID()

	// 1. Get the access token
	accessToken, err := getAccessToken(credentials)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// 2. Create product review invitation link
	productReviewData := ProductReview{
		Consumer: ProductReviewConsumer{
			Email: customerEmail,
			Name:  customerName,
		},
		Products:    buildProductInfo(currentCart, defaultProductBrand),
		ReferenceID: orderID,
		Locale:      requestLocale,
	}

	productReviewLink, err := getProductReviewLink(productReviewData, businessID, accessToken)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// 3. Generate service review invitation link, which will then redirect to the product review link
	serviceReviewData := ServiceReview{
		ReferenceID: orderID,
		Email:       customerEmail,
		Name:        customerName,
		Locale:      requestLocale,
		RedirectURI: productReviewLink, // product review link is daisy chained
	}

	serviceReviewLink, err := getServiceReviewLink(serviceReviewData, businessID, accessToken)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// 4. Update order with the service review link
	customInfo := utils.InterfaceToMap(checkoutOrder.Get("custom_info"))
	customInfo[ConstOrderCustomInfoLinkKey] = serviceReviewLink
	customInfo[ConstOrderCustomInfoSentKey] = false

	checkoutOrder.Set("custom_info", customInfo)

	err = checkoutOrder.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

func hasEmpty(testStrings []string) bool {
	for _, test := range testStrings {
		if test == "" {
			return true
		}
	}

	return false
}

func getAccessToken(cred tpCredentials) (string, error) {
	bodyString := "grant_type=password&username=" + cred.username + "&password=" + cred.password
	buffer := bytes.NewBuffer([]byte(bodyString))

	valueAMIKeySecret := []byte(cred.apiKey + ":" + cred.apiSecret)
	encodedString := base64.StdEncoding.EncodeToString(valueAMIKeySecret)

	request, err := http.NewRequest("POST", accessTokenURL, buffer)
	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", "Basic "+encodedString)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode >= 300 {
		errMsg := "Non 200 response while trying to get trustpilot access token: StatusCode:" + response.Status
		err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "376b178e-6cbf-4b4e-a3a8-fd65251d176b", errMsg)
		return "", err
	}

	jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
	if err != nil {
		return "", err
	}

	token := utils.InterfaceToString(jsonResponse["access_token"])
	if token == "" {
		return "", env.ErrorNew(ConstErrorModule, 1, "1293708d-9638-455a-8d49-3a387f086181", "Trustpilot didn't return an access token for our request")
	}

	return token, nil
}

/**
 * 2. Create product review invitation link
 *
 * https://developers.trustpilot.com/product-reviews-api
 *
 * Given information about the consumer and the product(s) purchased, get a link that can be sent to
 * the consumer to request reviews.
 */
func getProductReviewLink(requestData ProductReview, businessID string, accessToken string) (string, error) {
	reviewURL := strings.Replace(productReviewURL, "{businessUnitId}", businessID, 1)

	jsonString := utils.EncodeToJSONString(requestData)
	buffer := bytes.NewBuffer([]byte(jsonString))

	request, err := http.NewRequest("POST", reviewURL, buffer)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	if response.StatusCode >= 300 {
		errMsg := "Non 200 response while trying to get trustpilot review link: StatusCode:" + response.Status
		err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e75b28c7-0da2-475b-8b65-b1a09f1f6926", errMsg)
		return "", env.ErrorDispatch(err)
	}

	jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	reviewLinkI, ok := jsonResponse["reviewUrl"]
	if !ok {
		errorMessage := "Review link empty"
		return "", env.ErrorNew(ConstErrorModule, 1, "c53fd02f-2f5d-4111-8318-69a2cc2d2259", errorMessage)
	}
	reviewLink := utils.InterfaceToString(reviewLinkI)

	return reviewLink, nil
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
func getServiceReviewLink(requestData ServiceReview, businessUnitID string, accessToken string) (string, error) {

	reviewURL := strings.Replace(serviceReviewURL, "{businessUnitId}", businessUnitID, 1)

	jsonString := utils.EncodeToJSONString(requestData)
	buffer := bytes.NewBuffer([]byte(jsonString))

	request, err := http.NewRequest("POST", reviewURL, buffer)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode >= 300 {
		errMsg := "Non 200 response while trying to get trustpilot review link: StatusCode:" + response.Status
		err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e75b28c7-0da2-475b-8b65-b1a09f1f6926", errMsg)
		return "", err
	}

	jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
	if err != nil {
		return "", err
	}

	serviceReviewLinkI, ok := jsonResponse["url"]
	if !ok {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e528633c-9413-41b0-bfe8-8cee581a616c", "Service review link empty")
	}
	serviceReviewLink := utils.InterfaceToString(serviceReviewLinkI)

	return serviceReviewLink, nil
}

func buildProductInfo(cCart cart.InterfaceCart, defaultBrand string) []ProductReviewProduct {
	var productsOrdered []ProductReviewProduct
	mediaStorage, _ := media.GetMediaStorage()

	cartItems := cCart.GetItems()
	for _, productItem := range cartItems {
		pid := productItem.GetProductID()
		p := productItem.GetProduct()

		mediaPath, _ := mediaStorage.GetMediaPath("product", pid, "image")

		productOptions := productItem.GetOptions()

		productBrand := defaultBrand
		if brand, present := productOptions["brand"]; present {
			productBrand = utils.InterfaceToString(brand)
		}

		productInfo := ProductReviewProduct{
			ProductURL: app.GetStorefrontURL("product/" + pid),
			ImageURL:   app.GetStorefrontURL(mediaPath + p.GetDefaultImage()),
			Name:       p.GetName(),
			Sku:        p.GetSku(),
			Brand:      productBrand,
		}

		productsOrdered = append(productsOrdered, productInfo)
	}

	return productsOrdered
}
