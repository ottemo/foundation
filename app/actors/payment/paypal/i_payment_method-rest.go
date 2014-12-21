package paypal

/*
import (
	"bytes"

	"fmt"

	"text/template"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/utils"
)

func (it *RestAPI) GetName() string {
	return PAYMENT_NAME_REST
}

func (it *RestAPI) GetCode() string {
	return PAYMENT_CODE_REST
}

func (it *RestAPI) GetType() string {
	return checkout.ConstPaymentTypeCreditCard
}

func (it *RestAPI) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return true
}

func (it *RestAPI) Authorize(checkoutInstance checkout.InterfaceCheckout) error {

	ccInfo := utils.InterfaceToMap(checkoutInstance.GetInfo("cc"))
	if !utils.StrKeysInMap(ccInfo, "type", "number", "expire_month", "expire_year", "cvv") {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ce41da6b60514b3a883a4dd37c1c7a1a", "credit card info was not specified")
	}

	billingAddress := checkoutInstance.GetBillingAddress()
	if billingAddress == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9c57fbf9cb294472901ef77a9399dda6", "no billing address information")
	}

	order := checkoutInstance.GetOrder()
	if order == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6219944fd62b40d78495f844be0e562d", "no created order")
	}

	templateValues := map[string]interface{}{
		"intent":         "sale",
		"payment_method": "credit_card",
		"number":         utils.InterfaceToString(ccInfo["number"]),
		"type":           utils.InterfaceToString(ccInfo["type"]),
		"expire_month":   utils.InterfaceToString(ccInfo["expire_month"]),
		"expire_year":    utils.InterfaceToString(ccInfo["expire_year"]),
		"cvv2":           utils.InterfaceToString(ccInfo["cvv"]),
		"first_name":     billingAddress.GetFirstName(),
		"last_name":      billingAddress.GetLastName(),

		"line1":        billingAddress.GetAddressLine1(),
		"city":         billingAddress.GetCity(),
		"state":        billingAddress.GetState(),
		"postal_code":  billingAddress.GetZipCode(),
		"country_code": billingAddress.GetCountry(),

		"total":    fmt.Sprintf("%.2f", order.GetGrandTotal()),
		"currency": "USD",

		"description": "order id - " + order.GetID(),
	}

	bodyTemplate := `{
  "intent":"{{.intent}}",
  "payer":{
    "payment_method":"{{.payment_method}}",
    "funding_instruments":[
      {
        "credit_card":{
          "number":"{{.number}}",
          "type":"{{.type}}",
          "expire_month":{{.expire_month}},
          "expire_year":{{.expire_year}},
          "cvv2":"{{.cvv2}}",
          "first_name":"{{.first_name}}",
          "last_name":"{{.last_name}}",
          "billing_address":{
            "line1":"{{.line1}}",
            "city":"{{.city}}",
            "state":"{{.state}}",
            "postal_code":"{{.postal_code}}",
            "country_code":"{{.country_code}}"
          }
        }
      }
    ]
  },
  "transactions":[
    {
      "amount":{
        "total":"{{.total}}",
        "currency":"{{.currency}}"
      },
      "description":"{{.description}}"
    }
  ]
}`

	var body bytes.Buffer
	parsedTemplate, _ := template.New("paypal_payment").Parse(bodyTemplate)
	parsedTemplate.Execute(&body, templateValues)

	fmt.Println(body.String())

	request, err := http.NewRequest("POST", "https://api.sandbox.paypal.com/v1/payments/payment", &body)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	accessToken, err := it.GetAccessToken(checkoutInstance)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	fmt.Println(accessToken)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+accessToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	fmt.Println(response)
	fmt.Println(string(buf))

	result := make(map[string]interface{})
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if response.StatusCode != 201 {
		if response.StatusCode == 400 {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b022da8361ac4198a72d70f58f65acf0", utils.InterfaceToString(result["details"]))
		}
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "741e875d0ab54e8f81f98abbff10cf78", "payment was not processed")
	}

	//TODO: should store information to order

	return nil
}

func (it *RestAPI) Capture(checkoutInstance checkout.InterfaceCheckout) error {
	return nil
}

func (it *RestAPI) Refund(checkoutInstance checkout.InterfaceCheckout) error {
	return nil
}

func (it *RestAPI) Void(checkoutInstance checkout.InterfaceCheckout) error {
	return nil
}

// returns application access token needed for all other requests
func (it *RestAPI) GetAccessToken(checkoutInstance checkout.InterfaceCheckout) (string, error) {

	body := "grant_type=client_credentials"

	req, err := http.NewRequest("POST", "https://api.sandbox.paypal.com/v1/oauth2/token", bytes.NewBufferString(body))
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	req.SetBasicAuth("AbrcnhDi238ke9aG2NIQqVkW90oMJVg3B1QsjC68d2xRBLDq8boIrCaigPli", "EPcLWBCmfM_AwSOO1jC6TEDLCg-xZhFrUmXQnvTQ9yZV5_786xc5OkQ4Gx2-")

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en_US")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	if token, present := result["access_token"]; present {
		return utils.InterfaceToString(token), nil
	}

	return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "96d9554615954fe19156ef8d6e43f172", "unexpected response - without access_token")
}*/
