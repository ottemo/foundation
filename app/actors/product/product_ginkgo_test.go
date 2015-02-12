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

		ConstURLProducts           = "products"             // {GET}
		ConstURLProductsShop       = "products/shop"        // {GET}
		ConstURLProductsAttributes = "products/attributes"  // {GET}
		ConstURLProductsShopLayers = "products/shop/layers" // {GET}
		ConstURLProduct            = "product"              // {POST}

		ConstURLProductID        = "product/:productID"         // {PUT} {GET} {DELETE}
		ConstURLProductIDRating  = "product/:productID/rating"  // {GET}
		ConstURLProductIDRelated = "product/:productID/related" // {GET}
		ConstURLProductIDReviews = "product/:productID/reviews" // {GET}
		ConstURLProductIDReview  = "product/:productID/review"  // {POST}

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
		urlString          string
		err                error
		//urlString = strings.Replace(ConstURLProducts, ":productID", productID, 1)
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

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
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

		Context("Deleting product", func() {
			It("Testing: product/:productID {DELETE}", func() {
				urlString = strings.Replace(ConstURLProductID, ":productID", productID, 1)
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
				By("product delete result: " + utils.InterfaceToString(result))
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
