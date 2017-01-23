package authorizenet

import (
	"strconv"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/visitor"
)

type digits [6]int

// at returns the digits from the start to the given length
func (d *digits) at(i int) int {
	return d[i-1]
}

// getCardTypeByNumber get card type by card number
func getCardTypeByNumber(number string) (string, error) {
	ccLen := len(number)

	ccDigits := digits{}

	for i := 0; i < 6; i++ {
		if i < ccLen {
			ccDigits[i], _ = strconv.Atoi(number[:i+1])
		}
	}

	switch {
	case ccDigits.at(2) == 34 || ccDigits.at(2) == 37:
		return "American Express", nil
	case ccDigits.at(4) == 5610 || (ccDigits.at(6) >= 560221 && ccDigits.at(6) <= 560225):
		return "Bankcard", nil
	case ccDigits.at(2) == 62:
		return "China UnionPay", nil
	case ccDigits.at(3) >= 300 && ccDigits.at(3) <= 305 && ccLen == 15:
		return "Diners Club Carte Blanche", nil
	case ccDigits.at(4) == 2014 || ccDigits.at(4) == 2149:
		return "Diners Club enRoute", nil
	case ((ccDigits.at(3) >= 300 && ccDigits.at(3) <= 305) || ccDigits.at(3) == 309 || ccDigits.at(2) == 36 || ccDigits.at(2) == 38 || ccDigits.at(2) == 39) && ccLen <= 14:
		return "Diners Club International", nil
	case ccDigits.at(4) == 6011 || (ccDigits.at(6) >= 622126 && ccDigits.at(6) <= 622925) || (ccDigits.at(3) >= 644 && ccDigits.at(3) <= 649) || ccDigits.at(2) == 65:
		return "Discover", nil
	case ccDigits.at(3) == 636 && ccLen >= 16 && ccLen <= 19:
		return "InterPayment", nil
	case ccDigits.at(3) >= 637 && ccDigits.at(3) <= 639 && ccLen == 16:
		return "InstaPayment", nil
	case ccDigits.at(4) >= 3528 && ccDigits.at(4) <= 3589:
		return "JCB", nil
	case ccDigits.at(4) == 5018 || ccDigits.at(4) == 5020 || ccDigits.at(4) == 5038 || ccDigits.at(4) == 5612 || ccDigits.at(4) == 5893 || ccDigits.at(4) == 6304 || ccDigits.at(4) == 6759 || ccDigits.at(4) == 6761 || ccDigits.at(4) == 6762 || ccDigits.at(4) == 6763 || number[:3] == "0604" || ccDigits.at(4) == 6390:
		return "Maestro", nil
	case ccDigits.at(4) == 5019:
		return "Dankort", nil
	case ccDigits.at(2) >= 51 && ccDigits.at(2) <= 55:
		return "MasterCard", nil
	case ccDigits.at(4) == 4026 || ccDigits.at(6) == 417500 || ccDigits.at(4) == 4405 || ccDigits.at(4) == 4508 || ccDigits.at(4) == 4844 || ccDigits.at(4) == 4913 || ccDigits.at(4) == 4917:
		return "Visa Electron", nil
	case ccDigits.at(1) == 4:
		return "Visa", nil
	default:
		return "", env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d26eb21a-7ca9-47be-940e-986a0c443859", "Unknown credit card method.")
	}
}

// formatExpirationDate formats expiration date values to standard format
func formatExpirationDate(expireYear, expireMonth string) (string, error) {
	if len(expireMonth) < 1 {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9f1cdc5c-62c1-44b1-b82c-21c99aa0963a", "unexpected month value: "+expireMonth)
	}

	var expirationDate = utils.InterfaceToString(expireMonth)

	// pad with a zero
	if len(expireMonth) < 2 {
		expirationDate = "0" + expirationDate
	}

	// append the last two year digits
	year := utils.InterfaceToString(expireYear)
	if len(year) == 4 {
		expirationDate = expirationDate + year[2:]
	} else {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "16edcf8b-d183-4c64-a819-3a6405f27c38", "unexpected year length: "+year)
	}

	return expirationDate, nil
}

// getCustomerIDByVisitorID returns 3rd party customer ID by visitor registered ID
func getCustomerIDByVisitorID(visitorID string) string {
	var absentIDValue = ""
	var customerIDAttribute = "customer_id"

	if visitorID == "" {
		env.ErrorDispatch(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "16666b2f-1b40-44c5-8e0e-2cfcd868ba4c", "empty visitor ID passed to look up customer token"))
		return absentIDValue
	}

	model, _ := visitor.GetVisitorCardCollectionModel()
	model.ListFilterAdd("visitor_id", "=", visitorID)
	model.ListFilterAdd("payment", "=", ConstPaymentAuthorizeNetRestAPICode)

	// 3rd party customer identifier, used by braintree
	err := model.ListAddExtraAttribute(customerIDAttribute)
	if err != nil {
		env.ErrorDispatch(err)
	}

	visitorCards, err := model.List()
	if err != nil {
		env.ErrorDispatch(err)
	}

	for _, visitorCard := range visitorCards {
		return utils.InterfaceToString(visitorCard.Extra[customerIDAttribute])
	}

	return absentIDValue
}
