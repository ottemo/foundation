package visitor_test

import (
	"time"
	"strconv"
	"bytes"
	"math/rand"
	"net/http"
	"io/ioutil"


	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

)

var _ = Describe("Visitor", func() {
		var (
			request      *http.Request
		)

		rand.Seed(time.Now().UnixNano())

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
							num := rand.Intn(10000)
							strNum := strconv.Itoa(num)
							mailbox := `test` + strNum + `@ottemo.io`
							generalStr :=`"email": ` + `"` + mailbox + `"` + `"first_name": ` + `"test` + strNum + `"` + `"password": "123"` + `"is_admin": true`
							//for a current test
							generalStr = `{"email": "test@ottemo.ua", "password": "321", "first_name": "alex", "last_name": "bes", "is_admin": "true"}`
							jsonStartstr := []byte(generalStr)

							request, _ = http.NewRequest("POST", "http://192.168.56.101:3000/visitor/create", bytes.NewBuffer(jsonStartstr))
							request.Header.Set("Content-Type", "application/json")

							var jsonStr = []byte(`{"login": "admin", "password": "admin"}`)
							req, _ := http.NewRequest("POST", "http://192.168.56.101:3000/app/login", bytes.NewBuffer(jsonStr))
							req.Header.Set("Content-Type", "application/json")

							client := &http.Client{}
							resp, _ := client.Do(req)
							logcookies := resp.Cookies()
							for i:= range logcookies {
								request.AddCookie(logcookies[i])
							}

						})

						It("check the response", func() {

								client := &http.Client{}
								resp, err := client.Do(request)
								Expect(err).NotTo(HaveOccurred())

								defer resp.Body.Close()
								body, err := ioutil.ReadAll(resp.Body)
								Expect(err).NotTo(HaveOccurred())
								Expect(resp.Status).To(Equal("200 OK"))
								Expect((string(body))).NotTo(Equal(`{"error":{"code":"29be1531-cb6b-44cf-a78e-f1bf9aae1163","level":8,"message":"email already exists"},"redirect":"","result":null}`))
								Expect((string(body))).To(Equal("eror here"))
							})

					})

			})

})
