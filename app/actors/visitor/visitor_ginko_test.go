package visitor_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/tests"

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
		ConstURLVisitorInfo   = "visitor/info"
		ConstURLVisitorCreate = "visitor/create"
		ConstURLAppLogin      = "app/login"

		ConstVisitorPassword = "123"
	)

	var (
		request *http.Request
		client = &http.Client{}

		err error
	)

	Describe("Requests /visitor/", func() {

		Context("GET /visitor/", func() {
			It("check the response", func() {
				response, err := http.Get(app.GetFoundationURL(ConstURLVisitorInfo))
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(body)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
			})
		})

		Context("POST /visitor/create", func() {

			BeforeEach(func() {

				randomNumber := randomdata.StringNumber(3, "")
				firstName := randomdata.SillyName()
				mailbox := firstName + "_" + randomNumber + "@ottemo.io"

				// preparing request for a visitor creation
				visitorInfo := map[string]interface{}{
					"email":      mailbox,
					"first_name": firstName,
					"last_name": firstName,
					"password":   ConstVisitorPassword,
					"is_admin":   true}
				jsonString := utils.EncodeToJSONString(visitorInfo)

				fmt.Println()
				fmt.Println("- create request:", jsonString)

				buffer := bytes.NewBuffer([]byte(jsonString))
				request, err = http.NewRequest("POST", app.GetFoundationURL(ConstURLVisitorCreate), buffer)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Set("Content-Type", "application/json")

				// making app login request to become admin
				jsonString = `{"login": "admin", "password": "admin"}`
				buffer = bytes.NewBuffer([]byte(jsonString))
				loginRequest, err := http.NewRequest("POST", app.GetFoundationURL(ConstURLAppLogin), buffer)
				Expect(err).NotTo(HaveOccurred())

				loginRequest.Header.Set("Content-Type", "application/json")
				loginResponse, err := client.Do(loginRequest)
				Expect(err).NotTo(HaveOccurred())

				// getting sessionID and setting it to all following requests
				logcookies := loginResponse.Cookies()
				for i := range logcookies {
					request.AddCookie(logcookies[i])
				}

			})

			It("check the response", func() {
				resp, err := client.Do(request)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				responseBody, err := ioutil.ReadAll(resp.Body)

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonResponse).To(HaveKey("error"))
				Expect(jsonResponse["error"]).Should(BeNil())
				Expect(jsonResponse["result"]).ShouldNot(BeNil())

				result, ok := jsonResponse["result"].(map[string]interface{})
				Expect(ok).Should(BeTrue())
				Expect(result).To(HaveKey("_id"))

				newVisitorId := utils.InterfaceToString(result["_id"])
				Expect(newVisitorId).ShouldNot(BeEmpty())
				fmt.Println("- created visitor id:", newVisitorId)
			})

		})

	})

})
