// Package visitor_test represents a ginko/gomega test for visitor's api
package address_test

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
var _ = Describe("Visitor address ginkgo test", func() {

	// Defining constants for testing
	const (
		ConstURLAppLogin      = "app/login"   // POST for admin login
		ConstURLVisitorCreate = "visitor"     // POST Create of new visitor
		ConstURLVisitorLogin  = "visit/login" // POST for log in as a visitor
		ConstURLVisitorDelete = "visitor/"    // DELETE of a visitor

		ConstURLVisitAddressCreate   = "visit/address"   // POST for create of address as a visitor
		ConstURLVisitAddressesList   = "visit/addresses" // GET list of addresses for current visitor
		ConstURLVisitAddressLoadID   = "visit/address/"  // GET address buy it ID for current visitor
		ConstURLVisitAddressUpdateID = "visit/address/"  // PUT in address buy it ID for current visitor
		ConstURLVisitAddressDeleteID = "visit/address/"  // DELETE address buy it ID for current visitor

		ConstURLVisitorVIDAddressCreate   = "visitor/:visitorID/address"   // POST for create of address as admin
		ConstURLVisitorVIDAddresses       = "visitor/:visitorID/addresses" // GET list of addresses as admin for certain visitor
		ConstURLVisitorVIDAddressUpdateID = "visitor/:visitorID/address/"  // PUT in address buy it ID for certain visitor
		ConstURLVisitorVIDAddressDeleteID = "visitor/:visitorID/address/"  // DELETE address buy it ID for certain visitor

		ConstURLVisitorsAddressLoadID   = "visitors/address/" // GET address buy it ID
		ConstURLVisitorsAddressUpdateID = "visitors/address/" // PUT in address buy it ID
		ConstURLVisitorsAddressDeleteID = "visitors/address/" // DELETE address buy it ID

		ConstVisitorPassword = "123"
	)

	// Defining variables for testing
	var (
		request            *http.Request
		client             = &http.Client{}
		loginVisitorCookie []*http.Cookie
		loginAdminCookie   []*http.Cookie
		jsonString         string
		visitorID          string
		addressID          string
		mailbox            string
		urlString          string
		err                error
	)

	Describe("Preparing for workin with addresses", func() {
		Context("to create requests", func() {
			It("making an admin session", func() {
				jsonString = `{"login": "admin", "password": "admin"}`
				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLAppLogin), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				// making app login request to become admin
				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				// getting sessionID and setting it to a loginAdminCookie
				loginAdminCookie = response.Cookies()
				By("We have logined and added Admin Cookies")
			})
		})

		Context("POST /visitor/create as a admin", func() {
			It("Creating the visitor. Testing "+ConstURLVisitorCreate, func() {
				randomNumber := randomdata.StringNumber(3, "")
				firstName := randomdata.SillyName()
				mailbox = firstName + "_" + randomNumber + "@ottemo.io"

				// preparing request for a visitor creation
				visitorInfo := map[string]interface{}{
					"email":      mailbox,
					"first_name": firstName,
					"last_name":  firstName,
					"password":   ConstVisitorPassword,
					"is_admin":   false}
				jsonString = utils.EncodeToJSONString(visitorInfo)

				By("Start to create visitor:" + jsonString)

				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLVisitorCreate), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding admins cookie to request
				for i := range loginAdminCookie {
					request.AddCookie(loginAdminCookie[i])
				}

				//taking response and checking of it
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

				visitorID = utils.InterfaceToString(result["_id"])
				Expect(visitorID).ShouldNot(BeEmpty())
				By("Finished creating visitor:" + utils.InterfaceToString(result))
				By("- created visitor id:" + visitorID)
			})
		})

		Context("POST /visitor/login", func() {
			It("Loggining as a visitor and taking it's cookie. Testing "+ConstURLVisitorLogin, func() {
				visitorInfo := map[string]interface{}{
					"email":    mailbox,
					"password": ConstVisitorPassword}
				jsonString = utils.EncodeToJSONString(visitorInfo)

				By("Start to login as a " + jsonString)
				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLVisitorLogin), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

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

				// getting visitor sessionID and setting it to a loginVisitorCookie
				loginVisitorCookie = response.Cookies()
				By("We have logined and added Visitor Cookies")
			})
		})
	})

	Describe("Testing an addresses as a visitor", func() {
		Context("POST /visit/address (create)", func() {
			It("Creating of visitor address as visitor. Testing "+ConstURLVisitAddressCreate, func() {
				randomNumber := randomdata.StringNumber(4, "")
				firstName := randomdata.SillyName()

				visitorAddressInfo := map[string]interface{}{
					"visitor_id":    visitorID,
					"address_line1": "My Street 17/8",
					"first_name":    firstName,
					"last_name":     firstName,
					"company":       "Speroteck",
					"city":          "Chikago",
					"phone":         randomNumber,
					"zip_code":      "07400"}
				jsonString = utils.EncodeToJSONString(visitorAddressInfo)
				By("Create address with params: " + jsonString)

				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLVisitAddressCreate), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding visitor cookie to request
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				By("We have response body " + string(responseBody))
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"].(map[string]interface{})
				Expect(ok).Should(BeTrue())
				Expect(result).To(HaveKey("_id"))

				addressID = utils.InterfaceToString(result["_id"])
				Expect(addressID).ShouldNot(BeEmpty())
				By("- created visitor address id:" + addressID)
			})
		})

		Context("GET /visit/address/  as a visitor", func() {
			It("Getting list of addresses for current visitor. Testing "+ConstURLVisitAddressesList, func() {
				request, err = http.NewRequest("GET", app.GetFoundationURL(ConstURLVisitAddressesList), nil)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding visitors cookie to request
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}

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

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				Expect(result).ShouldNot(BeNil())
				By("List of addresses for visitor: " + utils.InterfaceToString(result))
			})
		})

		Context("PUT /visit/address update as a visitor", func() {
			It("Updating an address. Testing "+ConstURLVisitAddressUpdateID, func() {
				By("Start updating of an address with _ID: " + addressID + " as visitor")
				updateName := randomdata.SillyName()
				visitorInfo := map[string]interface{}{
					"country":    "Ukraine",
					"last_name":  updateName,
					"first_name": updateName}
				jsonString = utils.EncodeToJSONString(visitorInfo)

				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("PUT", app.GetFoundationURL(ConstURLVisitAddressUpdateID+addressID), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding visitors cookie to request
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}

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
				By("Update visitor and get response result:" + utils.InterfaceToString(result))
			})
		})

		Context("GET /visit/address/ load as a visitor", func() {
			It("Make an address load from visitor. Testing "+ConstURLVisitAddressLoadID, func() {
				request, err = http.NewRequest("GET", app.GetFoundationURL(ConstURLVisitAddressLoadID+addressID), nil)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding visitor cookie to request
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}

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

				result, ok := jsonResponse["result"].(map[string]interface{})
				Expect(ok).Should(BeTrue())
				By("Address visitor load result: " + utils.InterfaceToString(result))
			})
		})

		Context("DELETE /visit/address/ delete as a visitor", func() {
			It("Deleting an address as a visitor. Testing "+ConstURLVisitAddressDeleteID, func() {
				By("Start deleting of an address with _ID: " + addressID + " as visitor")
				request, err = http.NewRequest("DELETE", app.GetFoundationURL(ConstURLVisitAddressDeleteID+addressID), nil)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding visitors cookie to request
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}

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
				By("Finished deleting of an address with _ID: " + addressID)
			})
		})
	})

	Describe("Testing an addresses as a admin", func() {
		Context("POST /visit/address (create)", func() {
			It("Creating of visitor address as admin. Testing "+ConstURLVisitorVIDAddressCreate, func() {
				randomNumber := randomdata.StringNumber(4, "")
				firstName := randomdata.SillyName()

				visitorAddressInfo := map[string]interface{}{
					"visitor_id":    visitorID,
					"address_line1": "My Street 17/8",
					"first_name":    firstName,
					"last_name":     firstName,
					"company":       "Speroteck",
					"city":          "Chikago",
					"phone":         randomNumber,
					"zip_code":      "07400"}
				jsonString = utils.EncodeToJSONString(visitorAddressInfo)
				By("Create address with params: " + jsonString)

				buffer := bytes.NewBuffer([]byte(jsonString))
				urlString = strings.Replace(ConstURLVisitorVIDAddressCreate, ":visitorID", visitorID, 1)
				request, err = http.NewRequest("POST", app.GetFoundationURL(urlString), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding admins cookie to request
				for i := range loginAdminCookie {
					request.AddCookie(loginAdminCookie[i])
				}

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				By("We have response body " + string(responseBody))
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"].(map[string]interface{})
				Expect(ok).Should(BeTrue())
				Expect(result).To(HaveKey("_id"))

				addressID = utils.InterfaceToString(result["_id"])
				Expect(addressID).ShouldNot(BeEmpty())
				By("- created visitor address id:" + addressID)
			})
		})

		Context("GET /visit/address/  as a admin", func() {
			It("Getting list of addresses for certain visitor. Testing "+ConstURLVisitorVIDAddresses, func() {
				urlString = strings.Replace(ConstURLVisitorVIDAddresses, ":visitorID", visitorID, 1)
				request, err = http.NewRequest("GET", app.GetFoundationURL(urlString), nil)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding admins cookie to request
				for i := range loginAdminCookie {
					request.AddCookie(loginAdminCookie[i])
				}

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

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				Expect(result).ShouldNot(BeNil())
				By("List of addresses for visitor: " + utils.InterfaceToString(result))
			})
		})

		Context("PUT /visit/address update as admin", func() {
			It("Updating an address. Testing "+ConstURLVisitorVIDAddressUpdateID, func() {
				By("Start updating of an address with _ID: " + addressID + " as admin")
				updateName := randomdata.SillyName()
				visitorInfo := map[string]interface{}{
					"country":    "Ukraine",
					"last_name":  updateName,
					"first_name": updateName}
				jsonString = utils.EncodeToJSONString(visitorInfo)

				buffer := bytes.NewBuffer([]byte(jsonString))
				urlString = strings.Replace(ConstURLVisitorVIDAddressUpdateID, ":visitorID", visitorID, 1)
				request, err = http.NewRequest("PUT", app.GetFoundationURL(urlString+addressID), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding admin cookie to request
				for i := range loginAdminCookie {
					request.AddCookie(loginAdminCookie[i])
				}

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
				By("Update visitor and get response result:" + utils.InterfaceToString(result))
			})
		})

		Context("DELETE /visit/address/ delete as a admin", func() {
			It("Deleting an address as admin. Testing "+ConstURLVisitorVIDAddressDeleteID, func() {
				By("Start deleting of an address with _ID: " + addressID + " as admin")
				urlString = strings.Replace(ConstURLVisitorVIDAddressDeleteID, ":visitorID", visitorID, 1)
				request, err = http.NewRequest("DELETE", app.GetFoundationURL(urlString+addressID), nil)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding admins cookie to request
				for i := range loginAdminCookie {
					request.AddCookie(loginAdminCookie[i])
				}

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
				By("Finished deleting of an address with _ID: " + addressID)
			})
		})
	})

	Describe("Testing an addresses combine", func() {
		Context("POST /visit/address (create)", func() {
			It("Creating of visitor address as a visitor. Testing "+ConstURLVisitAddressCreate, func() {
				randomNumber := randomdata.StringNumber(4, "")
				firstName := randomdata.SillyName()

				visitorAddressInfo := map[string]interface{}{
					"visitor_id":    visitorID,
					"address_line1": "My Street 17/8",
					"first_name":    firstName,
					"last_name":     firstName,
					"company":       "Speroteck",
					"city":          "Chikago",
					"phone":         randomNumber,
					"zip_code":      "07400"}
				jsonString = utils.EncodeToJSONString(visitorAddressInfo)
				By("Create address with params: " + jsonString)

				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLVisitAddressCreate), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding visitor cookie to request
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				By("We have response body " + string(responseBody))
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"].(map[string]interface{})
				Expect(ok).Should(BeTrue())
				Expect(result).To(HaveKey("_id"))

				addressID = utils.InterfaceToString(result["_id"])
				Expect(addressID).ShouldNot(BeEmpty())
				By("- created visitor address id:" + addressID)
			})
		})

		Context("GET /visitors/address/  as a visitor", func() {
			It("Getting Make an address load from visitor. Testing "+ConstURLVisitorsAddressLoadID, func() {
				request, err = http.NewRequest("GET", app.GetFoundationURL(ConstURLVisitorsAddressLoadID+addressID), nil)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding visitors cookie to request
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}

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

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				Expect(result).ShouldNot(BeNil())
				By("List of addresses for visitor: " + utils.InterfaceToString(result))
			})
		})
		Context("PUT /visitors/address update", func() {
			It("Updating an address. No cookies", func() {
				By("Start updating of an address with _ID: " + addressID + "as NONE")
				updateName := randomdata.SillyName()
				visitorInfo := map[string]interface{}{
					"country":    "Ukraine",
					"last_name":  updateName,
					"first_name": updateName}
				jsonString = utils.EncodeToJSONString(visitorInfo)

				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("PUT", app.GetFoundationURL(ConstURLVisitorsAddressUpdateID+addressID), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

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
				By("Fake update address and get response result:" + utils.InterfaceToString(result))
			})
		})

		Context("DELETE /visitors/address/ delete* ", func() {
			It("Deleting an address. No cookies", func() {
				By("Start deleting of an address with _ID: " + addressID + " as NONE")
				request, err = http.NewRequest("DELETE", app.GetFoundationURL(ConstURLVisitorsAddressDeleteID+addressID), nil)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

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

				result, ok := jsonResponse["error"].(map[string]interface{})
				Expect(ok).Should(BeTrue())
				By("Fake delete address and get response result:" + utils.InterfaceToString(result))
			})
		})

		Context("PUT /visitors/address update as admin", func() {
			It("Updating an address. Testing "+ConstURLVisitorsAddressUpdateID, func() {
				By("Start updating of an address with _ID: " + addressID + " as admin")
				updateName := randomdata.SillyName()
				visitorInfo := map[string]interface{}{
					"country":    "Ukraine",
					"last_name":  updateName,
					"first_name": updateName}
				jsonString = utils.EncodeToJSONString(visitorInfo)

				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("PUT", app.GetFoundationURL(ConstURLVisitorsAddressUpdateID+addressID), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding visitors cookie to request
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}

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
				By("Update visitor and get response result:" + utils.InterfaceToString(result))
			})
		})

		Context("DELETE /visitors/address/ delete ", func() {
			It("Deleting an address. Testing "+ConstURLVisitorsAddressDeleteID, func() {
				By("Start deleting of an address with _ID: " + addressID + " as admin")
				request, err = http.NewRequest("DELETE", app.GetFoundationURL(ConstURLVisitorsAddressDeleteID+addressID), nil)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding visitors cookie to request
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}

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
				By("Finished deleting of an address with _ID: " + addressID)
			})
		})

		Context("POST /visit/address (create)*", func() {
			It("Creating of visitor address as a visitor. Input wrong data", func() {
				randomNumber := randomdata.StringNumber(4, "")
				firstName := randomdata.SillyName()

				visitorAddressInfo := map[string]interface{}{
					"last_name": firstName,
					"company":   "Speroteck",
					"city":      "Chikago",
					"phone":     randomNumber,
					"zip_code":  "0000"}
				jsonString = utils.EncodeToJSONString(visitorAddressInfo)
				By("Create address with params: " + jsonString)

				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLVisitAddressCreate), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding visitor cookie to request
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				By("We have response body " + string(responseBody))
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).ShouldNot(BeNil())
				Expect(jsonResponse["result"]).Should(BeNil())
			})
		})

		Context("DELETE /visitor/delete as a admin", func() {
			It("Deleting a visitor as an admin. Testing "+ConstURLVisitorDelete, func() {
				By("Start deleting of a visitor with _ID: " + visitorID + " as admin")
				request, err = http.NewRequest("DELETE", app.GetFoundationURL(ConstURLVisitorDelete+visitorID), nil)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				//adding admins cookie to request
				for i := range loginAdminCookie {
					request.AddCookie(loginAdminCookie[i])
				}

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
				By("Finished deleting of a visitor with _ID: " + visitorID)
			})
		})

	})
})
