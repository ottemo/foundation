// Package visitor_test represents a ginko/gomega test for visitor's api
package address_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/tests"
	"github.com/ottemo/foundation/utils"

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

// General starting function for Ginko
var _ = Describe("Visitor", func() {

	// Defining constants for testing
	const (
		ConstURLAppLogin                    = "app/login"
		ConstURLVisitorCreate               = "visitor/create"
		ConstURLVisitorLogin                = "visitor/login"
		ConstURLVisitorDelete               = "visitor/delete/"
		ConstURLVisitorAddressAttributeList = "visitor/address/attribute/list"
		ConstURLVisitorAddressCount         = "visitor/address/count"
		ConstURLVisitorAddressCreate        = "visitor/address/create"
		ConstURLVisitorAddressDeleteID      = "visitor/address/delete/"
		ConstURLVisitorAddressList          = "visitor/address/list"
		ConstURLVisitorAddressListVID       = "visitor/address/list/"
		ConstURLVisitorAddressLoadID        = "visitor/address/load/"
		ConstURLVisitorAddressUpdateID      = "visitor/address/update/"

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
		err                error
	)

	Describe("Visitor address api testing", func() {
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
				//fmt.Println("- create request:", jsonString)

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

		Context("POST /visitor/address/create", func() {
			It("Creating of visitor address as visitor. Testing "+ConstURLVisitorAddressCreate, func() {
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
				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLVisitorAddressCreate), buffer)
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

		Context("GET /visitor/address/attribute/list as a admin", func() {
			It("Getting address attribute list as admin. Testing "+ConstURLVisitorAddressAttributeList, func() {
				request, err = http.NewRequest("GET", app.GetFoundationURL(ConstURLVisitorAddressAttributeList), nil)
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
				By("Address attribute list get result: " + utils.InterfaceToString(result))
			})
		})
		Context("GET /visitor/address/count as a admin", func() {
			It("Getting address count as admin. Testing "+ConstURLVisitorAddressCount, func() {
				request, err = http.NewRequest("GET", app.GetFoundationURL(ConstURLVisitorAddressCount), nil)
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
				By("Address admin count get result: " + utils.InterfaceToString(result))
			})
		})

		Context("GET /visitor/address/list as a admin", func() {
			It("Get an address list from admin. Testing "+ConstURLVisitorAddressListVID, func() {
				request, err = http.NewRequest("GET", app.GetFoundationURL(ConstURLVisitorAddressListVID+visitorID), nil)
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
				By("Address admin list get result: " + utils.InterfaceToString(result))
			})
		})

		Context("POST /visitor/address/list as a admin", func() {
			It("Post an address list from admin. Testing "+ConstURLVisitorAddressListVID, func() {
				//randomNumber := randomdata.StringNumber(5, "")
				//firstName := randomdata.SillyName()

				visitorAddressInfo := map[string]interface{}{
				/*"visitor_id":    visitorID,
				"address_line1": "My Street 0/1",
				"first_name":    firstName,
				"last_name":     firstName,
				"company":       "Speroteck",
				"country":       `"US":"United States"`,
				"city":          "Chikago",
				"phone":         randomNumber,
				"state":         `"IL":"Illinois"`,
				"zip_code":      "07400"
				*/}
				jsonString = utils.EncodeToJSONString(visitorAddressInfo)
				buffer := bytes.NewBuffer([]byte(jsonString))

				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLVisitorAddressListVID+visitorID), buffer)
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
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("Address admin list get result: " + utils.InterfaceToString(result))
			})
		})

		Context("GET /visitor/address/list as a visitor", func() {
			It("Get an address list from visitor. Testing "+ConstURLVisitorAddressList, func() {
				request, err = http.NewRequest("GET", app.GetFoundationURL(ConstURLVisitorAddressList), nil)
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

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("Address visitor list get result: " + utils.InterfaceToString(result))
			})
		})

		Context("POST /visitor/address/list as a visitor", func() {
			It("Post an address list from visitor. Testing "+ConstURLVisitorAddressList, func() {
				//randomNumber := randomdata.StringNumber(5, "")
				//firstName := randomdata.SillyName()

				visitorAddressInfo := map[string]interface{}{
				/*					"visitor_id":    visitorID,
									"address_line1": "My Street 10/20",
									"first_name":    firstName,
									"last_name":     firstName,
									"company":       "Speroteck",
									"city":          "Chikago",
									"phone":         randomNumber,
									"zip_code":      "07400"
				*/}
				jsonString = utils.EncodeToJSONString(visitorAddressInfo)
				buffer := bytes.NewBuffer([]byte(jsonString))

				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLVisitorAddressList), buffer)
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
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				By("Address visitor list post result: " + utils.InterfaceToString(result))
			})
		})

		Context("GET /visitor/address/load as a visitor", func() {
			It("Make an address load from visitor. Testing "+ConstURLVisitorAddressLoadID, func() {
				request, err = http.NewRequest("GET", app.GetFoundationURL(ConstURLVisitorAddressLoadID+addressID), nil)
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

		Context("PUT /visitor/address/update as an admin", func() {
			It("Updating an address. Testing "+ConstURLVisitorAddressUpdateID, func() {
				By("Start updating of an address with _ID: " + addressID + " as admin")
				updateName := randomdata.SillyName()
				visitorInfo := map[string]interface{}{
					"last_name":  updateName,
					"first_name": updateName}
				jsonString = utils.EncodeToJSONString(visitorInfo)
				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("PUT", app.GetFoundationURL(ConstURLVisitorAddressUpdateID+addressID), buffer)
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

		Context("DELETE /visitor/address/delete as a admin", func() {
			It("Deleting an address as an admin. Testing "+ConstURLVisitorAddressDeleteID, func() {
				By("Start deleting of an address with _ID: " + addressID + " as admin")
				request, err = http.NewRequest("DELETE", app.GetFoundationURL(ConstURLVisitorAddressDeleteID+addressID), nil)
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
