package category_test

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

var _ = Describe("GinkgoCategory", func() {

	const (
		ConstURLAppLogin     = "app/login"   // POST for admin login
		ConstURLVisitorLogin = "visit/login" // POST for log in as a visitor
		ConstURLProducts     = "products"    // GET products to add to category

		ConstURLCategories                   = "categories"                              // {GET}
		ConstURLCategoriesAttributes         = "categories/attributes"                   // {GET}
		ConstURLCategoriesTree               = "categories/tree"                         // {GET}
		ConstURLCategory                     = "category"                                // {POST}
		ConstURLCategoryCategoryID           = "category/:categoryID"                    // {GET} {PUT} {DELETE}
		ConstURLCategoryCategoryIDProducts   = "category/:categoryID/products"           // {GET}
		ConstURLCategoryCategoryIDProductPID = "category/:categoryID/product/:productID" // {POST} {DELETE}
		ConstURLCategoryCategoryIDLayers     = "category/:categoryID/layers"             // {GET}

		ConstAdminLogin      = "admin"
		ConstAdminPassword   = "admin"
		ConstVisitorEmail    = "alex@ottemo.io"
		ConstVisitorPassword = "123"
	)

	// Defining variables for testing
	var (
		request            *http.Request
		client             = &http.Client{}
		loginAdminCookie   []*http.Cookie
		loginVisitorCookie []*http.Cookie
		productID          string
		productArray       []string
		urlString          string
		categoryID         string
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
		Context("Take a list of products", func() {
			It("to make possible add them", func() {
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
				resultString := utils.InterfaceToString(result)

				for {
					positionOfID := strings.Index(resultString, "ID")
					if positionOfID == -1 {
						break
					} else {
						positionOfID += 5
					}
					subStringID := resultString[positionOfID : positionOfID+24]
					resultString = resultString[positionOfID+24 : len(resultString)-1]
					By("current sub string: " + subStringID)
					productArray = append(productArray, subStringID)
				}

				productID = productArray[0]
				By("list of products in shop: " + (utils.InterfaceToString(productArray)))
				By("we will work with single product, it's ID: " + productID)
			})
		})
	})

	Describe("Testing categories general functions", func() {
		Context("Creating a category", func() {
			It("Testing: category {POST}", func() {
				randomName := randomdata.SillyName()
				newCategoryInfo := map[string]interface{}{
					"enabled": true,
					"name":    randomName}

				request, err = CreateRequest("POST", ConstURLCategory, newCategoryInfo, loginAdminCookie)
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

				categoryID = utils.InterfaceToString(result["_id"])
				Expect(categoryID).ShouldNot(BeEmpty())
				By("created category with id: " + categoryID)
			})
		})

		Context("Take a list of categories", func() {
			It("Test of API: categories  {GET}", func() {
				request, err = CreateRequest("GET", ConstURLCategories, nil, nil)
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
				By("Categories: " + utils.InterfaceToString(result))
			})
		})

		Context("Take attributes of categories", func() {
			It("Test of API: categories/attributes  {GET}", func() {
				request, err = CreateRequest("GET", ConstURLCategoriesAttributes, nil, nil)
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
				By("Attributes: " + utils.InterfaceToString(result))
			})
		})

		Context("Take categories tree", func() {
			It("Test of API: categories/tree  {GET}", func() {
				request, err = CreateRequest("GET", ConstURLCategoriesTree, nil, nil)
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
				By("Categories tree: " + utils.InterfaceToString(result))
			})
		})
	})

	Describe("Individual category tests", func() {
		Context("Adding products to category", func() {
			It("Testing: category/:categoryID {PUT}", func() {
				randomName := randomdata.SillyName()
				categoryNewInfo := map[string]interface{}{
					"products": productArray,
					"name":     randomName}
				urlString = strings.Replace(ConstURLCategoryCategoryID, ":categoryID", categoryID, 1)
				request, err = CreateRequest("PUT", urlString, categoryNewInfo, loginAdminCookie)
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
				Expect(result).To(HaveKey("name"))
				Expect(result).To(HaveKey("products"))

				categoryName := utils.InterfaceToString(result["name"])
				Expect(categoryName).ShouldNot(BeEmpty())
				By("category new name: " + categoryName)

				categoryProducts := utils.InterfaceToString(result["products"])
				Expect(categoryProducts).ShouldNot(BeEmpty())
				By("category products: " + categoryProducts)
			})
		})

		Context("Deleting a prodcut from category", func() {
			It("Testing: category/:categoryID/product/:productID {DELETE}", func() {
				urlString = strings.Replace(ConstURLCategoryCategoryIDProductPID, ":categoryID", categoryID, 1)
				urlString = strings.Replace(urlString, ":productID", productID, 1)
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
				By("Deleting product from category result: " + utils.InterfaceToString(result))
			})
		})

		Context("Take products from category", func() {
			It("Testing: category/:categoryID/products  {GET}", func() {
				urlString = strings.Replace(ConstURLCategoryCategoryIDProducts, ":categoryID", categoryID, 1)
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
				By("Category products: " + utils.InterfaceToString(result))
			})
		})

		Context("Adding product to the category", func() {
			It("Testing: category/:categoryID/product/:productID {POST}", func() {
				urlString = strings.Replace(ConstURLCategoryCategoryIDProductPID, ":categoryID", categoryID, 1)
				urlString = strings.Replace(urlString, ":productID", productID, 1)
				request, err = CreateRequest("POST", urlString, nil, loginAdminCookie)
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
				By("Adding product to the category result: " + utils.InterfaceToString(result))
			})
		})

		Context("Take category info", func() {
			It("Testing: category/:categoryID  {GET}", func() {
				urlString = strings.Replace(ConstURLCategoryCategoryID, ":categoryID", categoryID, 1)
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
				By("Category info: " + utils.InterfaceToString(result))
			})
		})
		//						ConstURLCategoryCategoryIDLayers
		Context("Take category layers", func() {
			It("Testing: category/:categoryID/layers  {GET}", func() {
				urlString = strings.Replace(ConstURLCategoryCategoryIDLayers, ":categoryID", categoryID, 1)
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
				By("Category layers info: " + utils.InterfaceToString(result))
			})
		})

		Context("Fake deleting of a category (permission)", func() {
			It("Error testing: category/:categoryID {DELETE}", func() {
				urlString = strings.Replace(ConstURLCategoryCategoryID, ":categoryID", categoryID, 1)
				request, err = CreateRequest("DELETE", urlString, nil, nil)
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
				Expect(jsonResponse["error"]).ShouldNot(BeNil())
				Expect(jsonResponse["result"]).Should(BeNil())

				result, ok := jsonResponse["error"]
				Expect(ok).Should(BeTrue())
				By("Deleting category error: " + utils.InterfaceToString(result))
			})
		})

		Context("Deleting a category", func() {
			It("Testing: category/:categoryID {DELETE}", func() {
				urlString = strings.Replace(ConstURLCategoryCategoryID, ":categoryID", categoryID, 1)
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
				By("Deleting category result: " + utils.InterfaceToString(result))
			})
		})
		Context("Creating a category with error (permission)", func() {
			It("Error testing: category {POST}", func() {
				randomName := randomdata.SillyName()
				newCategoryInfo := map[string]interface{}{
					"enabled": true,
					"name":    randomName}

				request, err = CreateRequest("POST", ConstURLCategory, newCategoryInfo, nil)
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
				Expect(jsonResponse["error"]).ShouldNot(BeNil())
				Expect(jsonResponse["result"]).Should(BeNil())

				result, ok := jsonResponse["error"].(map[string]interface{})
				Expect(ok).Should(BeTrue())
				By("Fake creating error: " + utils.InterfaceToString(result))
			})
		})
	})
})

func CreateRequest(typeAPI, requestURL string, contentMap interface{}, requestCookies []*http.Cookie) (request *http.Request, err error) {

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
