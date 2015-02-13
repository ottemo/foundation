// Package visitor_test represents a ginko/gomega test for visitor's api
package product_test

import (
	"bytes"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/tests"
	"github.com/ottemo/foundation/utils"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Init settings for running the application in testing mode
var _ = BeforeSuite(func() {
	err := tests.StartAppInTestingMode()
	Expect(err).NotTo(HaveOccurred())

	go app.Serve()
	time.Sleep(1 * time.Second)
})

// General starting function for Ginkgo
var _ = Describe("Products test", func() {

	// Defining constants for testing
	const (
		ConstURLAppLogin     = "app/login"   // POST for admin login
		ConstURLVisitorLogin = "visit/login" // POST for log in as a visitor

		ConstURLProducts                   = "products"                      // {GET}
		ConstURLProductsShop               = "products/shop"                 // {GET}
		ConstURLProductsAttributes         = "products/attributes"           // {GET}
		ConstURLProductsAttribute          = "products/attribute"            // {POST}
		ConstURLProductsAttributeAttribute = "products/attribute/:attribute" // {PUT} {DELETE}
		ConstURLProductsShopLayers         = "products/shop/layers"          // {GET}
		ConstURLProduct                    = "product"                       // {POST}

		ConstURLProductID                 = "product/:productID"                    // {PUT} {GET} {DELETE}
		ConstURLProductIDRating           = "product/:productID/rating"             // {GET}
		ConstURLProductIDRelated          = "product/:productID/related"            // {GET}
		ConstURLProductIDReviews          = "product/:productID/reviews"            // {GET}
		ConstURLProductIDReview           = "product/:productID/review"             // {POST}
		ConstURLProductIDReviewRID        = "product/:productID/review/:reviewID"   // {DELETE}
		ConstURLProductIDRatedreviewStars = "product/:productID/ratedreview/:stars" // {POST}
		//ConstURLProductIDStock            = "product/:productID/stock"              // {POST}

		ConstAdminLogin      = "admin"
		ConstAdminPassword   = "admin"
		ConstVisitorEmail    = "a@i.ua"
		ConstVisitorPassword = "123"
	)

	// Defining variables for testing
	var (
		request            *http.Request
		client             = &http.Client{}
		loginAdminCookie   []*http.Cookie
		loginVisitorCookie []*http.Cookie
		productID          string
		attributeName      string
		urlString          string
		reviewID           string
		err                error
	)

	Describe("Prepare to work with products", func() {
		Context("to create admin requests", func() {
			It("create an admin session", func() {
				adminInfo := map[string]interface{}{
					"login":    ConstAdminLogin,
					"password": ConstAdminPassword}

				request, err = CreateRequest("POST", ConstURLAppLogin, adminInfo, nil)
				Expect(err).NotTo(HaveOccurred())

				// making app login request to become admin
				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				// getting sessionID and setting it to a loginAdminCookie
				loginAdminCookie = response.Cookies()
				By("We have logined and added Admin Cookies")
			})
		})
		Context("to create visitor requests", func() {
			It("create an visitor session", func() {
				visitorInfo := map[string]interface{}{
					"email":    ConstVisitorEmail,
					"password": ConstVisitorPassword}

				request, err = CreateRequest("POST", ConstURLVisitorLogin, visitorInfo, nil)
				Expect(err).NotTo(HaveOccurred())

				// making app login request to become admin
				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				// getting sessionID and setting it to a loginAdminCookie
				loginVisitorCookie = response.Cookies()
				By("We have logined and added Visitor Cookies")
			})
		})
	})

	Describe("Testing products general functions", func() {
		Context("Take a list of products", func() {
			It("Test of API: products  {GET}", func() {
				request, err = CreateRequest("GET", ConstURLProducts, nil, nil)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By(":" + utils.InterfaceToString(result))
			})
		})

		Context("Take a list of products in shop", func() {
			It("Testing: products/shop {GET}", func() {
				request, err = CreateRequest("GET", ConstURLProductsShop, nil, nil)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("list of products in shop: " + utils.InterfaceToString(result))
			})
		})

		Context("Creating products attriutes", func() {
			It("Testing: products/attribute {POST}", func() {
				attributeName = randomdata.SillyName()
				attributeInfo := map[string]interface{}{
					"Attribute":  attributeName,
					"Type":       "text",
					"Label":      "Feel free",
					"IsRequired": true,
					"IsStatic":   true,
					"Group":      "General",
					"Editors":    "multiline_text",
					"Options":    "",
					"Default":    "",
					"Validators": "",
					"IsLayered":  false,
					"IsPublic":   false}

				request, err = CreateRequest("POST", ConstURLProductsAttribute, attributeInfo, loginAdminCookie)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"].(map[string]interface{})
				Expect(ok).Should(BeTrue())
				By("created product attribute: " + utils.InterfaceToString(result))
				Expect(result).To(HaveKey("Attribute"))

				attributeName = utils.InterfaceToString(result["Attribute"])
				Expect(attributeName).ShouldNot(BeEmpty())
				By("created product attribute: " + attributeName)
			})
		})

		Context("Take a list of products attributes", func() {
			It("Testing: products/attributes  {GET}", func() {
				request, err = CreateRequest("GET", ConstURLProductsAttributes, nil, nil)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("list of products attributes: " + utils.InterfaceToString(result))
			})
		})

		Context("Take a list of products in shop by layers", func() {
			It("Testing: products/shop/layers {GET}", func() {
				request, err = CreateRequest("GET", ConstURLProductsShopLayers, nil, nil)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("list of products in shop by layers: " + utils.InterfaceToString(result))
			})
		})
	})

	Describe("Individual product tests", func() {
		Context("Creating a product", func() {
			It("Testing: product {POST}", func() {
				randomNumber := randomdata.StringNumber(4, "")
				randomName := randomdata.SillyName()
				randomSKU := randomName + randomNumber
				productInfo := map[string]interface{}{
					"enabled":           true,
					attributeName:       "hello world",
					"name":              randomName,
					"sku":               randomSKU,
					"price":             "10.7",
					"weight":            "1.3",
					"short_description": "some short description",
					"description":       "fuuuuuly description"}

				request, err = CreateRequest("POST", ConstURLProduct, productInfo, loginAdminCookie)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"].(map[string]interface{})
				Expect(ok).Should(BeTrue())
				Expect(result).To(HaveKey("_id"))

				productID = utils.InterfaceToString(result["_id"])
				Expect(productID).ShouldNot(BeEmpty())
				By("created product with id: " + productID)
			})
		})

		Context("Getting product info", func() {
			It("Testing: product/:productID {GET}", func() {
				urlString = strings.Replace(ConstURLProductID, ":productID", productID, 1)
				request, err = CreateRequest("GET", urlString, nil, nil)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("product info: " + utils.InterfaceToString(result))
			})
		})

		Context("Updating product info", func() {
			It("Testing: product/:productID {PUT}", func() {
				productInfo := map[string]interface{}{
					"price":             "17",
					"weight":            "4",
					"short_description": "another short description",
					"description":       "fuuuuuly second description"}

				urlString = strings.Replace(ConstURLProductID, ":productID", productID, 1)
				request, err = CreateRequest("PUT", urlString, productInfo, loginAdminCookie)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("Updated result: " + utils.InterfaceToString(result))
			})
		})

		Context("Getting product related products", func() {
			It("Testing: product/:productID/related {GET}", func() {
				urlString = strings.Replace(ConstURLProductIDRelated, ":productID", productID, 1)
				request, err = CreateRequest("GET", urlString, nil, nil)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				//Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("product related products: " + utils.InterfaceToString(result))
			})
		})

		Context("Seting product review", func() {
			It("Testing: product/:productID/review {POST}", func() {
				productInfo := map[string]interface{}{
					"it's fine": "wasssup"}

				urlString = strings.Replace(ConstURLProductIDReview, ":productID", productID, 1)
				request, err = CreateRequest("POST", urlString, productInfo, loginVisitorCookie)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				//Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"].(map[string]interface{})
				Expect(ok).Should(BeTrue())
				Expect(result).To(HaveKey("_id"))
				reviewID = utils.InterfaceToString(result["_id"])
				Expect(reviewID).ShouldNot(BeEmpty())
				By("- created review with ID:" + reviewID)
				By("product review add: " + utils.InterfaceToString(result))
			})
		})

		Context("Getting product reviews", func() {
			It("Testing: product/:productID/reviews {GET}", func() {
				urlString = strings.Replace(ConstURLProductIDReviews, ":productID", productID, 1)
				request, err = CreateRequest("GET", urlString, nil, nil)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("product reviews: " + utils.InterfaceToString(result))
			})
		})

		Context("Seting 'stars' on review", func() {
			It("Testing: product/:productID/ratedreview/:stars {POST}", func() {
				productInfo := map[string]interface{}{
					"stars": "5"}

				urlString = strings.Replace(ConstURLProductIDRatedreviewStars, ":productID", productID, 1)
				urlString = strings.Replace(urlString, ":stars", "4", 1)
				request, err = CreateRequest("POST", urlString, productInfo, loginVisitorCookie)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				//Expect(jsonResponse["error"]).Should(BeNil())
				//Expect(jsonResponse["result"]).ShouldNot(BeNil())
				eror, ok := jsonResponse["result"].(map[string]interface{})
				By("erors text: " + utils.InterfaceToString(eror))

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())

				By("product review stars: " + utils.InterfaceToString(result))
			})
		})

		Context("Getting product rating", func() {
			It("Testing: product/:productID/rating  {GET}", func() {
				urlString = strings.Replace(ConstURLProductIDRating, ":productID", productID, 1)
				request, err = CreateRequest("GET", urlString, nil, nil)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				//Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("product rating: " + utils.InterfaceToString(result))
			})
		})

		Context("Deleting review", func() {
			It("Testing: product/:productID/review/:reviewID {DELETE}", func() {
				urlString = strings.Replace(ConstURLProductIDReviewRID, ":productID", productID, 1)
				urlString = strings.Replace(urlString, ":reviewID", reviewID, 1)
				request, err = CreateRequest("DELETE", urlString, nil, loginVisitorCookie)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("review delete result: " + utils.InterfaceToString(result))
			})
		})

		Context("Getting product info", func() {
			It("Testing: product/:productID {GET}", func() {
				urlString = strings.Replace(ConstURLProductID, ":productID", productID, 1)
				request, err = CreateRequest("GET", urlString, nil, nil)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("product info: " + utils.InterfaceToString(result))
			})
		})

		Context("Deleting product", func() {
			It("Testing: product/:productID {DELETE}", func() {
				urlString = strings.Replace(ConstURLProductID, ":productID", productID, 1)
				request, err = CreateRequest("DELETE", urlString, nil, loginAdminCookie)
				Expect(err).NotTo(HaveOccurred())
				time.Sleep(1 * time.Second)
				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("product delete result: " + utils.InterfaceToString(result))
			})
		})

		Context("Deleting products attribute", func() {
			It("Testing: products/attribute/:attribute {DELETE}", func() {
				urlString = strings.Replace(ConstURLProductsAttributeAttribute, ":attribute", attributeName, 1)
				request, err = CreateRequest("DELETE", urlString, nil, loginAdminCookie)
				Expect(err).NotTo(HaveOccurred())

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("products attribute delete result: " + utils.InterfaceToString(result))
			})
		})

	})

})

func CreateRequest(typeAPI, requestURL string, contentMap map[string]interface{}, requestCookies []*http.Cookie) (request *http.Request, err error) {

	buffer := bytes.NewBuffer([]byte(""))

	if contentMap != nil {
		jsonString := utils.EncodeToJSONString(contentMap)
		buffer = bytes.NewBuffer([]byte(jsonString))
	}

	request, err = http.NewRequest(typeAPI, app.GetFoundationURL(requestURL), buffer)
	request.Header.Set("Content-Type", "application/json")

	if requestCookies != nil {
		for i := range requestCookies {
			request.AddCookie(requestCookies[i])
		}
	}

	return request, err
}
