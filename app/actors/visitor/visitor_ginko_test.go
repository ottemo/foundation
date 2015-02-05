package visitor_test

import (
	"bytes"
	"fmt"
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

var _ = BeforeSuite(func() {
	err := tests.StartAppInTestingMode()
	Expect(err).NotTo(HaveOccurred())

	go app.Serve()
	time.Sleep(1 * time.Second)
})

var _ = Describe("Visitor", func() {

	const (
		ConstURLAppLogin         = "app/login"
		ConstURLVisitorInfo      = "visitor/info"
		ConstURLVisitorCreate    = "visitor/create"
		ConstURLVisitorLogin     = "visitor/login"
		ConstURLVisitorLogout    = "visitor/logout"
		ConstURLVisitorOrderList = "visitor/order/list"
		ConstURLVisitorUpdate    = "visitor/update"
		ConstURLVisitorDelete    = "visitor/delete"

		ConstVisitorPassword = "123"
	)

	var (
		request            *http.Request
		client             = &http.Client{}
		loginVisitorCookie []*http.Cookie
		loginAdminCookie   []*http.Cookie
		jsonString         string
		visitorID          string
		mailbox            string
		err                error
	)

	Describe("Visitor api testing", func() {
		Context("to create new visitor request", func() {
			It("nead to make a admin session", func() {
				jsonString = `{"login": "admin", "password": "admin"}`
				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLAppLogin), buffer)
				Expect(err).NotTo(HaveOccurred())

				request.Header.Set("Content-Type", "application/json")
				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())

				// getting sessionID and setting it to all following requests
				loginAdminCookie = response.Cookies()
			})
		})

		Context("POST /visitor/create as a admin", func() {
			It("Do request and check the response to "+ConstURLVisitorCreate+" url", func() {
				randomNumber := randomdata.StringNumber(3, "")
				firstName := randomdata.SillyName()
				mailbox = firstName + "_" + randomNumber + "@ottemo.io"

				// preparing request for a visitor creation
				visitorInfo := map[string]interface{}{
					"email":      mailbox,
					"first_name": firstName,
					"last_name":  firstName,
					"password":   ConstVisitorPassword,
					"is_admin":   true}
				jsonString = utils.EncodeToJSONString(visitorInfo)

				//By("- create request:" + jsonString)
				//fmt.Println("- create request:", jsonString)

				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLVisitorCreate), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				// making app login request to become admin

				for i := range loginAdminCookie {
					request.AddCookie(loginAdminCookie[i])
				}

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
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
				fmt.Println("- created visitor id:", visitorID)
			})
		})

		Context("POST /visitor/login", func() {
			It("check the response to "+ConstURLVisitorLogin+" url", func() {
				visitorInfo := map[string]interface{}{
					"email":    mailbox,
					"password": ConstVisitorPassword}
				jsonString = utils.EncodeToJSONString(visitorInfo)
				buffer := bytes.NewBuffer([]byte(jsonString))

				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLVisitorLogin), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				loginVisitorCookie = response.Cookies()
				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())

				By("Checking for not eror statment")
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

			})
		})

		Context("GET /visitor/info as a visitor", func() {
			It("check the response to "+ConstURLVisitorInfo+" url", func() {
				request, err = http.NewRequest("GET", app.GetFoundationURL(ConstURLVisitorInfo), nil)
				request.Header.Set("Content-Type", "application/json")
				Expect(err).NotTo(HaveOccurred())
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}
				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())

				result, ok := jsonResponse["result"].(map[string]interface{})
				Expect(ok).Should(BeTrue())
				fmt.Println("Take the info and get response result:" + utils.InterfaceToString(result))
			})
		})

		Context("GET /visitor/order/list as a visitor", func() {
			It("check the response to "+ConstURLVisitorOrderList+" url", func() {
				request, err = http.NewRequest("GET", app.GetFoundationURL(ConstURLVisitorOrderList), nil)
				request.Header.Set("Content-Type", "application/json")
				Expect(err).NotTo(HaveOccurred())
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}
				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
			})
		})

		Context("PUT /visitor/update as a visitor", func() {
			It("check the response to "+ConstURLVisitorUpdate+" url", func() {
				updateName := randomdata.SillyName()
				visitorInfo := map[string]interface{}{
					"last_name":  updateName,
					"first_name": updateName}
				jsonString = utils.EncodeToJSONString(visitorInfo)
				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("PUT", app.GetFoundationURL(ConstURLVisitorUpdate), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}
				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())

				By("Checking for not eror statment")
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())
				result, ok := jsonResponse["result"].(map[string]interface{})
				Expect(ok).Should(BeTrue())
				fmt.Println("Update visitor and get response result:" + utils.InterfaceToString(result))

			})
		})

		Context("POST /visitor/logout as a visitor", func() {
			It("check the response to "+ConstURLVisitorLogout+" url", func() {
				request, err = http.NewRequest("GET", app.GetFoundationURL(ConstURLVisitorLogout), nil)
				request.Header.Set("Content-Type", "application/json")
				Expect(err).NotTo(HaveOccurred())
				for i := range loginVisitorCookie {
					request.AddCookie(loginVisitorCookie[i])
				}
				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())

				result, ok := jsonResponse["result"]
				Expect(ok).Should(BeTrue())
				fmt.Println("Result from logout response:" + utils.InterfaceToString(result))
			})
		})

		Context("DELETE /visitor/delete as a admin", func() {
			It("check the response to "+ConstURLVisitorDelete+" url", func() {
				By("- delete a visitor with _ID: " + visitorID + " as admin")
				request, err = http.NewRequest("DELETE", app.GetFoundationURL(ConstURLVisitorDelete+"/"+visitorID), nil)
				request.Header.Set("Content-Type", "application/json")
				Expect(err).NotTo(HaveOccurred())
				for i := range loginAdminCookie {
					request.AddCookie(loginAdminCookie[i])
				}
				response, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
			})
		})

	})
})
