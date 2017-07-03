package fedex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"text/template"

	"gopkg.in/xmlpath.v1"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetName returns name of shipping method
func (it *FedEx) GetName() string {
	return ConstShippingName
}

// GetCode returns code of shipping method
func (it *FedEx) GetCode() string {
	return ConstShippingCode
}

// IsAllowed checks for method applicability
func (it *FedEx) IsAllowed(checkout checkout.InterfaceCheckout) bool {
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
func (it *FedEx) GetRates(checkoutObject checkout.InterfaceCheckout) []checkout.StructShippingRate {

	var result []checkout.StructShippingRate

	useDebugLog := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathDebugLog))

	templateValues := map[string]interface{}{
		"key":           	utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathKey)),
		"password":      	utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPassword)),
		"number":        	utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathNumber)),
		"meter":         	utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMeter)),
		"originZip":     	utils.InterfaceToString(env.ConfigGetValue(checkout.ConstConfigPathShippingOriginZip)),
		"originCountry": 	utils.InterfaceToString(env.ConfigGetValue(checkout.ConstConfigPathShippingOriginCountry)),
		"dropoff":       	utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathDropoff)),
		"packaging":     	utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPackaging)),
		"residential":     	utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathResidential)),
		"destinationZip":     	nil,
		"destinationCountry": 	nil,
		"weight":             	"0.01",
		"amount":             	checkoutObject.GetSubtotal(),
	}

	// getting destination zip code
	//-----------------------------
	if shippingAddress := checkoutObject.GetShippingAddress();
		shippingAddress != nil && shippingAddress.GetZipCode() != "" && shippingAddress.GetCountry() != "" {
		templateValues["destinationZip"] = shippingAddress.GetZipCode()
		templateValues["destinationCountry"] = shippingAddress.GetCountry()
	} else {
		return result
	}

	// calculating weight
	//-------------------
	var pounds float64
	if checkoutCart := checkoutObject.GetCart(); checkoutCart != nil {

		cartItems := checkoutCart.GetItems()
		if len(cartItems) == 0 {
			return result
		}

		defaultWeight := utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathDefaultWeight))

		if defaultWeight == 0.0 {
			defaultWeight = 0.01
		}

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
	templateValues["weight"] = pounds

	// preparing SOAP request
	//------------------------
	requestTemplate := `<?xml version="1.0" encoding="utf-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns="http://fedex.com/ws/rate/v16">
   <SOAP-ENV:Body>
      <RateRequest>
         <WebAuthenticationDetail>
            <UserCredential>
               <Key>{{.key}}</Key>
               <Password>{{.password}}</Password>
            </UserCredential>
         </WebAuthenticationDetail>
         <ClientDetail>
            <AccountNumber>{{.number}}</AccountNumber>
            <MeterNumber>{{.meter}}</MeterNumber>
         </ClientDetail>

         <Version>
            <ServiceId>crs</ServiceId>
            <Major>16</Major>
            <Intermediate>0</Intermediate>
            <Minor>0</Minor>
         </Version>

         <RequestedShipment>
            <DropoffType>{{.dropoff}}</DropoffType>
            <PackagingType>{{.packaging}}</PackagingType>

            <Shipper>
               <Address>
                  <PostalCode>{{.originZip}}</PostalCode>
                  <CountryCode>{{.originCountry}}</CountryCode>
               </Address>
            </Shipper>

            <Recipient>
               <Address>
                  <PostalCode>{{.destinationZip}}</PostalCode>
                  <CountryCode>{{.destinationCountry}}</CountryCode>
                  <Residential>{{.residential}}</Residential>
               </Address>
            </Recipient>

            <RateRequestTypes>LIST</RateRequestTypes>
            <PackageCount>1</PackageCount>

            <RequestedPackageLineItems>
               <GroupPackageCount>1</GroupPackageCount>
               <InsuredValue>
	          <Currency>USD</Currency>
	          <Amount>{{.amount}}</Amount>
               </InsuredValue>
               <Weight>
                  <Units>LB</Units>
                  <Value>{{.weight}}</Value>
               </Weight>
            </RequestedPackageLineItems>
         </RequestedShipment>
      </RateRequest>
   </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	var body bytes.Buffer
	parsedTemplate, _ := template.New("fedex").Parse(requestTemplate)
	if err := parsedTemplate.Execute(&body, templateValues); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e61049c9-4728-4f76-bf32-2ee577281542", err.Error())
	}

	if useDebugLog {
		env.Log("fedex.log", "REQUEST", string(body.Bytes()))
	}

	url := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathGateway)) + "/rate"
	request, err := http.NewRequest("POST", url, &body)
	if err != nil {
		_ = env.ErrorDispatch(err)
		return result
	}

	request.Header.Add("SOAPAction", "http://fedex.com/ws/rate/v16/getRates")
	request.Header.Add("Content-Type", "text/xml")

	// doing SOAP request and getting result
	//--------------------------------------
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		_ = env.ErrorDispatch(err)
		return result
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = env.ErrorDispatch(err)
		return result
	}

	if useDebugLog {
		env.Log("fedex.log", "RESPONSE", string(responseData))
	}

	// parsing xml response to rates
	//------------------------------
	xmlRoot, err := xmlpath.Parse(bytes.NewReader(responseData))
	if err != nil {
		_ = env.ErrorDispatch(err)
		return result
	}

	xmlErrorFlag, _ := xmlpath.Compile("//Severity")
	if value, ok := xmlErrorFlag.String(xmlRoot); ok && value == "ERROR" {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3bf4bb50-0ba8-4aaf-90a9-695fed9d497a", string(responseData))
	}

	xmlPostage, _ := xmlpath.Compile("//RateReplyDetails")
	xmlMethod, _ := xmlpath.Compile("./ServiceType")
	xmlRate, _ := xmlpath.Compile("./RatedShipmentDetails[1]/ShipmentRateDetail/TotalNetCharge/Amount")

	allowedMethodsArray := utils.InterfaceToArray(env.ConfigGetValue(ConstConfigPathAllowedMethods))

	for i := xmlPostage.Iter(xmlRoot); i.Next(); {
		postageNode := i.Node()

		stringMethod, ok := xmlMethod.String(postageNode)
		if !ok {
			continue
		}

		stringRate, ok := xmlRate.String(postageNode)
		if !ok {
			continue
		}

		if len(allowedMethodsArray) == 0 || utils.IsInArray(stringMethod, allowedMethodsArray) {
			rateName, present := ConstShippingMethods[stringMethod]
			if !present {
				rateName = stringMethod
			}

			result = append(result,
				checkout.StructShippingRate{
					Code:  stringMethod,
					Name:  rateName,
					Price: utils.InterfaceToFloat64(stringRate),
				})
		}

	}

	return result
}

// GetAllRates will return all the rates for the FedEx shipping method.
func (it FedEx) GetAllRates() []checkout.StructShippingRate {
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
