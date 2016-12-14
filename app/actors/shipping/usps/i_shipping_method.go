package usps

import (
	"bytes"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"text/template"

	"gopkg.in/xmlpath.v1"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
)

// GetName returns name of shipping method
func (it *USPS) GetName() string {
	return ConstShippingName
}

// GetCode returns code of shipping method
func (it *USPS) GetCode() string {
	return ConstShippingCode
}

// IsAllowed checks for method applicability
func (it *USPS) IsAllowed(checkout checkout.InterfaceCheckout) bool {
	if utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEnabled)) == false {
		return false

	}

	if shippingAddress := checkout.GetShippingAddress(); shippingAddress != nil  {
		for _, countryCode := range utils.InterfaceToArray(env.ConfigGetValue(ConstConfigPathAllowCountries)) {
			if shippingAddress.GetCountry() == countryCode {
				return true
			}
		}
	}
	return false
}

// GetRates returns rates allowed by shipping method for a given checkout
func (it *USPS) GetRates(checkoutObject checkout.InterfaceCheckout) []checkout.StructShippingRate {

	var result []checkout.StructShippingRate

	useDebugLog := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathDebugLog))

	templateValues := make(map[string]interface{})

	templateValues["userid"] = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathUser))      // "133OTTEM1795",
	templateValues["origin"] = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathOriginZip)) // "44106",

	if templateValues["userid"] == "" || templateValues["origin"] == "" {
		return nil
	}

	if shippingAddress := checkoutObject.GetShippingAddress(); shippingAddress != nil && shippingAddress.GetZipCode() != "" {
		templateValues["destination"] = shippingAddress.GetZipCode()
	} else {
		return result
	}

	var pounds float64
	var ounces float64
	if checkoutCart := checkoutObject.GetCart(); checkoutCart != nil {

		cartItems := checkoutCart.GetItems()
		if len(cartItems) == 0 {
			return result
		}

		defaultWeight := utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathDefaultWeight))

		for _, cartItem := range cartItems {
			cartProduct := cartItem.GetProduct()
			if cartProduct == nil {
				continue
			}

			if cartProduct.GetWeight() == 0 {
				pounds += defaultWeight * utils.InterfaceToFloat64(cartItem.GetQty())
			} else {
				pounds += cartProduct.GetWeight() * utils.InterfaceToFloat64(cartItem.GetQty())
			}
		}
	}

	templateValues["pounds"] = pounds
	templateValues["ounces"] = ounces

	templateValues["container"] = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathContainer))
	templateValues["size"] = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathSize))

	templateValues["width"] = 0.1
	templateValues["height"] = 0.1
	templateValues["length"] = 0.1
	templateValues["girth"] = 0.1

	dimensions := strings.Split(utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathDefaultDimensions)), "x")
	for idx, dimensionValue := range dimensions {
		dimensionValue = strings.Trim(dimensionValue, " ")
		switch idx {
		case 0:
			templateValues["width"] = utils.InterfaceToFloat64(dimensionValue)
		case 1:
			templateValues["height"] = utils.InterfaceToFloat64(dimensionValue)
		case 2:
			templateValues["length"] = utils.InterfaceToFloat64(dimensionValue)
		case 3:
			templateValues["girth"] = utils.InterfaceToFloat64(dimensionValue)
		}

	}

	requestTemplate := `<RateV4Request USERID="{{.userid}}">
     <Revision/>
     <Package ID="1">
          <Service>ALL</Service>
          <ZipOrigination>{{.origin}}</ZipOrigination>
          <ZipDestination>{{.destination}}</ZipDestination>
          <Pounds>{{.pounds}}</Pounds>
          <Ounces>{{.ounces}}</Ounces>
          <Container>{{.container}}</Container>
          <Size>{{.size}}</Size>
          {{if eq .size "LARGE" }}
          <Width>{{.width}}</Width>
          <Length>{{.length}}</Length>
          <Height>{{.height}}</Height>
          <Girth>{{.girth}}</Girth>
          {{end}}
          <Machinable>True</Machinable>
    </Package>
    </RateV4Request>`

	var buff bytes.Buffer
	parsedTemplate, _ := template.New("usps").Parse(requestTemplate)
	parsedTemplate.Execute(&buff, templateValues)

	if useDebugLog {
		env.Log("usps.log", "REQUEST", buff.String())
	}

	query := ConstHTTPEndpoint + "?API=RateV4&XML=" + url.QueryEscape(buff.String())
	response, err := http.Get(query)
	if err != nil {
		return result
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result
	}

	if useDebugLog {
		env.Log("usps.log", "RESPONSE", string(responseData))
	}

	root, err := xmlpath.Parse(bytes.NewReader(responseData))
	if err != nil {
		return result
	}

	allowedMethodsArray := utils.InterfaceToArray(env.ConfigGetValue(ConstConfigPathAllowedMethods))

	postage, _ := xmlpath.Compile("//Postage")
	service, _ := xmlpath.Compile("./MailService")
	rate, _ := xmlpath.Compile("./Rate")
	code, _ := xmlpath.Compile("./@CLASSID")

	regexTags := regexp.MustCompile("<[^>]+>")

	for i := postage.Iter(root); i.Next(); {
		postageNode := i.Node()

		stringService, ok := service.String(postageNode)
		if !ok {
			continue
		}

		stringRate, ok := rate.String(postageNode)
		if !ok {
			continue
		}

		stringCode, ok := code.String(postageNode)
		if !ok {
			continue
		}

		if len(allowedMethodsArray) == 0 || utils.IsInArray(stringCode, allowedMethodsArray) {

			rateName := html.UnescapeString(stringService)
			if ConstRemoveRateNameTags {
				rateName = regexTags.ReplaceAllString(rateName, "")

			}

			result = append(result,
				checkout.StructShippingRate{
					Code:  stringCode,
					Name:  rateName,
					Price: utils.InterfaceToFloat64(stringRate),
				})
		}
	}

	return result
}

// GetAllRates returns all the shipping rates for the USPS Shipping method.
func (it USPS) GetAllRates() []checkout.StructShippingRate {
	result := []checkout.StructShippingRate{}

	for code, name := range ConstShippingMethods {
		resultItem := checkout.StructShippingRate{
			Code: code,
			Name: name,
		}

		result = append(result, resultItem)
	}

	return result
}
