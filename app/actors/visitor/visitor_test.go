package visitor_test

import (
	"time"
	"strconv"
	"bytes"
	"math/rand"
	"net/http"
	//"net/http/httptest"
	//"encoding/json"
	"io/ioutil"

	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/tests"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

)

var _ = Describe("Visitor", func() {
		var (
			adminVisitor visitor.InterfaceVisitor
			//request      *http.Request
			//recorder     *httptest.ResponseRecorder
		)

		rand.Seed(time.Now().UnixNano())

		BeforeEach(func() {
			By("Starting a new instance for testing")
			err := tests.StartAppInTestingMode()
			Expect(err).NotTo(HaveOccurred())
			By("Creating a new Admin Visitor")
			visitorModel, err := visitor.GetVisitorModel()

			Expect(err).NotTo(HaveOccurred())

			num := rand.Intn(10000)
			strNum := strconv.Itoa(num)

			visitorModel.Set("email", "alex+"+strNum+"@ottemo.io")
			visitorModel.Set("first_name", "alex")
			visitorModel.Set("last_name", "bes")
			visitorModel.Set("passwd", "123")
			visitorModel.Set("birthday", "30-03-1994")
			visitorModel.Set("is_admin", true)
			visitorModel.Set("created_at", time.Now())

			By("Persisting the Visitor to the database")
			err = visitorModel.Save()
			Expect(err).NotTo(HaveOccurred())

			adminVisitor = visitorModel

		})

		Describe("Validating all attributes of a Registration against regex", func() {
				It("should have a valid email address", func() {
						By("Using conforming characters")
						Expect(adminVisitor.GetEmail()).Should(MatchRegexp("[a-zA-Z0-9_\x2E\x2B-]+@[a-zA-Z0-9-]+\x2E[a-zA-Z0-9-\x2E]+"))

					})

				It("should have a valid first name", func() {
						By("Using a conforming first name")
						Expect(adminVisitor.GetFirstName()).Should(MatchRegexp("[a-zA-Z0-9-]"))
					})

				It("should have a valid last name", func() {
						By("Using a conforming last name")
						Expect(adminVisitor.GetLastName()).Should(MatchRegexp("[a-zA-Z0-9-]"))
					})

			})

		Describe("Requests /visitor/", func() {

				Context("GET /visitor/", func() {
						It("check the response", func() {
								resp, _ := http.Get("http://192.168.56.101:3000/visitor/info")

								defer resp.Body.Close()
								body, err := ioutil.ReadAll(resp.Body)
								Expect(err).NotTo(HaveOccurred())
								Expect(resp.Status).To(Equal("200 OK"))
								Expect((string(body))).To(Equal(`{"error":null,"redirect":"","result":"you are not logined in"}`))
							})
					})

				Context("POST /visitor/customer/create", func() {

						BeforeEach(func() {
							//body, _ := json.Marshal(adminVisitor.ToHashMap())
							//request, _ = http.NewRequest("POST", "http://192.168.56.101:3000/visitor/create", bytes.NewReader(body))
							var jsonStr = []byte(`{"login": "admin", "password": "admin"}`)
							req, _ := http.NewRequest("POST", "http://192.168.56.101:3000/app/login", bytes.NewBuffer(jsonStr))
							req.Header.Set("Content-Type", "application/json")
							client := &http.Client{}
							client.Do(req)
						})

						It("check the response", func() {

								url := "http://192.168.56.101:3000/visitor/create?auth=admin:admin"

								var jsonStr = []byte(`{"email": "alex1a0@i.ua", "password": "321", "first_name": "alex", "last_name": "bes"}`)
								req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
								Expect(err).NotTo(HaveOccurred())
								req.Header.Set("Content-Type", "application/json")

								client := &http.Client{}
								resp, err := client.Do(req)
								Expect(err).NotTo(HaveOccurred())

								defer resp.Body.Close()
								body, err := ioutil.ReadAll(resp.Body)
								Expect(err).NotTo(HaveOccurred())
								Expect(resp.Status).To(Equal("200 OK"))
								Expect((string(body))).NotTo(Equal(`{"error":{"code":"29be1531-cb6b-44cf-a78e-f1bf9aae1163","level":8,"message":"email already exists"},"redirect":"","result":null}`))
								Expect((string(body))).To(Equal(
									`{"error":null,
									"redirect":"",
									"result":{
									"_id":"54d0bdb73f4c0d62e2e78c72",
									"billing_address":null,
									"created_at":"2015-02-03T06:23:19.141622944-06:00",
									"email":"alex1@i.ua",
									"facebook_id":"",
									"first_name":"alex",
									"google_id":"",
									"is_admin":false,
									"last_name":"bes",
									"password":"4bb4ac294ac120494dc771576e61a2c2",
									"shipping_address":null,
									"validate":""}}`))
							})

					})

			})

})
