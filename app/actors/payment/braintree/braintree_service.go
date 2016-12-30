package braintree

import (
	"github.com/lionelbarrow/braintree-go"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/visitor"
)

func braintreeCardFormatExpirationDate(card braintree.CreditCard) (string, error) {
	var expirationDate = utils.InterfaceToString(card.ExpirationMonth)

	if len(card.ExpirationMonth) < 1 {
		return "", env.ErrorNew(constErrorModule, constErrorLevel, "f861d39f-516d-4bdd-8316-b7c4d34e3531", "unexpected month value coming back from braintree: "+card.ExpirationMonth)
	}

	// pad with a zero
	if len(card.ExpirationMonth) < 2 {
		expirationDate = "0" + expirationDate
	}

	// append the last two year digits
	year := utils.InterfaceToString(card.ExpirationYear)
	if len(year) == 4 {
		expirationDate = expirationDate + year[2:]
	} else {
		return "", env.ErrorNew(constErrorModule, constErrorLevel, "950aea13-16e8-4d20-9ad0-f5cee26c03c2", "unexpected year length coming back from braintree: "+year)
	}

	return expirationDate, nil
}

func braintreeAddressFromVisitorAddress(visitorAddress visitor.InterfaceVisitorAddress) *braintree.Address {
	return &braintree.Address{
		FirstName:       visitorAddress.GetFirstName(),
		LastName:        visitorAddress.GetLastName(),
		Company:         visitorAddress.GetCompany(),
		StreetAddress:   visitorAddress.GetAddressLine1(),
		ExtendedAddress: visitorAddress.GetAddressLine2(),

		CountryCodeAlpha2: visitorAddress.GetCountry(),
		Locality:          visitorAddress.GetCity(),
		Region:            visitorAddress.GetState(),
		PostalCode:        visitorAddress.GetZipCode(),
	}
}

func braintreeCardToAuthorizeResult(card braintree.CreditCard, customerID string) (map[string]interface{}, error) {
	expirationDate, err := braintreeCardFormatExpirationDate(card)
	if err != nil {
		return nil, env.ErrorNew(constErrorModule, constErrorLevel, "7aa0ea8e-679e-4ac3-b84b-40aad71ead5f", "unable to format expiration date: "+err.Error())
	}

	var result = map[string]interface{}{
		"transactionID":      card.Token,     // token_id
		"creditCardLastFour": card.Last4,     // number
		"creditCardType":     card.CardType,  // type
		"creditCardExp":      expirationDate, // expiration_date
		"customerID":         customerID,     // customer_id
	}

	return result, nil
}
