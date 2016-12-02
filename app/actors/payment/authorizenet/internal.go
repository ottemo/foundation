package authorizenet

import (
	"strconv"

	"github.com/ottemo/foundation/env"
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
		return  "Dankort", nil
	case ccDigits.at(2) >= 51 && ccDigits.at(2) <= 55:
		return  "MasterCard", nil
	case ccDigits.at(4) == 4026 || ccDigits.at(6) == 417500 || ccDigits.at(4) == 4405 || ccDigits.at(4) == 4508 || ccDigits.at(4) == 4844 || ccDigits.at(4) == 4913 || ccDigits.at(4) == 4917:
		return "Visa Electron", nil
	case ccDigits.at(1) == 4:
		return "Visa", nil
	default:
		return "", env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d26eb21a-7ca9-47be-940e-986a0c443859", "Unknown credit card method.")
	}
}
