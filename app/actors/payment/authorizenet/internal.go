package authorizenet

import (
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
	"errors"
	"strconv"
)

type digits [6]int
// At returns the digits from the start to the given length
func (d *digits) At(i int) int {
	return d[i-1]
}
func getCardTypeByNumber(number string) (string, error) {
	ccLen := len(number)

	ccDigits := digits{}

	for i := 0; i < 6; i++ {
		if i < ccLen {
			ccDigits[i], _ = strconv.Atoi(number[:i+1])
		}
	}

	switch {
	case ccDigits.At(2) == 34 || ccDigits.At(2) == 37:
		return "American Express", nil
	case ccDigits.At(4) == 5610 || (ccDigits.At(6) >= 560221 && ccDigits.At(6) <= 560225):
		return "Bankcard", nil
	case ccDigits.At(2) == 62:
		return "China UnionPay", nil
	case ccDigits.At(3) >= 300 && ccDigits.At(3) <= 305 && ccLen == 15:
		return "Diners Club Carte Blanche", nil
	case ccDigits.At(4) == 2014 || ccDigits.At(4) == 2149:
		return "Diners Club enRoute", nil
	case ((ccDigits.At(3) >= 300 && ccDigits.At(3) <= 305) || ccDigits.At(3) == 309 || ccDigits.At(2) == 36 || ccDigits.At(2) == 38 || ccDigits.At(2) == 39) && ccLen <= 14:
		return "Diners Club International", nil
	case ccDigits.At(4) == 6011 || (ccDigits.At(6) >= 622126 && ccDigits.At(6) <= 622925) || (ccDigits.At(3) >= 644 && ccDigits.At(3) <= 649) || ccDigits.At(2) == 65:
		return "Discover", nil
	case ccDigits.At(3) == 636 && ccLen >= 16 && ccLen <= 19:
		return "InterPayment", nil
	case ccDigits.At(3) >= 637 && ccDigits.At(3) <= 639 && ccLen == 16:
		return "InstaPayment", nil
	case ccDigits.At(4) >= 3528 && ccDigits.At(4) <= 3589:
		return "JCB", nil
	case ccDigits.At(4) == 5018 || ccDigits.At(4) == 5020 || ccDigits.At(4) == 5038 || ccDigits.At(4) == 5612 || ccDigits.At(4) == 5893 || ccDigits.At(4) == 6304 || ccDigits.At(4) == 6759 || ccDigits.At(4) == 6761 || ccDigits.At(4) == 6762 || ccDigits.At(4) == 6763 || number[:3] == "0604" || ccDigits.At(4) == 6390:
		return "Maestro", nil
	case ccDigits.At(4) == 5019:
		return  "Dankort", nil
	case ccDigits.At(2) >= 51 && ccDigits.At(2) <= 55:
		return  "MasterCard", nil
	case ccDigits.At(4) == 4026 || ccDigits.At(6) == 417500 || ccDigits.At(4) == 4405 || ccDigits.At(4) == 4508 || ccDigits.At(4) == 4844 || ccDigits.At(4) == 4913 || ccDigits.At(4) == 4917:
		return "Visa Electron", nil
	case ccDigits.At(1) == 4:
		return "Visa", nil
	default:
		return "", errors.New("Unknown credit card method.")
	}
}

// getAuthorizenetCustomerToken We attach customer tokens to card tokens in the visitor_token table
// - the customer token is sensitive data because you can make a charge with it alone
// - if you are going to make a charge against a card that is attached to a customer though,
//   you must attach the customer id
func getAuthorizenetCustomerToken(vid string, paymentType string) string {
	const customerTokenPrefix = "cus"

	if vid == "" {
		env.ErrorDispatch(env.ErrorNew(ConstErrorModule, 1, "2ecfa3ec-7cfc-4783-9060-8467ca63beae", "empty vid passed to look up customer token"))
		return ""
	}

	model, _ := visitor.GetVisitorCardCollectionModel()
	model.ListFilterAdd("visitor_id", "=", vid)
	model.ListFilterAdd("payment", "=", paymentType)

	// 3rd party customer identifier, used by stripe
	err := model.ListAddExtraAttribute("customer_id")
	if err != nil {
		env.ErrorDispatch(err)
	}

	tokens, err := model.List()
	if err != nil {
		env.ErrorDispatch(err)
	}

	for _, t := range tokens {
		ts := utils.InterfaceToString(t.Extra["customer_id"])

		// Double check that this field is filled out
		if strings.HasPrefix(ts, customerTokenPrefix) {
			return ts
		}
	}

	return ""
}
