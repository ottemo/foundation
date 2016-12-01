package authorizenet

import (
	"strings"
	"errors"
	"strconv"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	
	"github.com/ottemo/foundation/app/models/visitor"
)

type digits [6]int
// At returns the digits from the start to the given length
func (d *digits) at(i int) int {
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
	case ccDigits.at(2) == 34 || ccDigits.At(2) == 37:
		return "American Express", nil
	case ccDigits.at(4) == 5610 || (ccDigits.At(6) >= 560221 && ccDigits.At(6) <= 560225):
		return "Bankcard", nil
	case ccDigits.at(2) == 62:
		return "China UnionPay", nil
	case ccDigits.at(3) >= 300 && ccDigits.At(3) <= 305 && ccLen == 15:
		return "Diners Club Carte Blanche", nil
	case ccDigits.at(4) == 2014 || ccDigits.At(4) == 2149:
		return "Diners Club enRoute", nil
	case ((ccDigits.at(3) >= 300 && ccDigits.At(3) <= 305) || ccDigits.At(3) == 309 || ccDigits.At(2) == 36 || ccDigits.At(2) == 38 || ccDigits.At(2) == 39) && ccLen <= 14:
		return "Diners Club International", nil
	case ccDigits.at(4) == 6011 || (ccDigits.At(6) >= 622126 && ccDigits.At(6) <= 622925) || (ccDigits.At(3) >= 644 && ccDigits.At(3) <= 649) || ccDigits.At(2) == 65:
		return "Discover", nil
	case ccDigits.at(3) == 636 && ccLen >= 16 && ccLen <= 19:
		return "InterPayment", nil
	case ccDigits.at(3) >= 637 && ccDigits.At(3) <= 639 && ccLen == 16:
		return "InstaPayment", nil
	case ccDigits.at(4) >= 3528 && ccDigits.At(4) <= 3589:
		return "JCB", nil
	case ccDigits.at(4) == 5018 || ccDigits.At(4) == 5020 || ccDigits.At(4) == 5038 || ccDigits.At(4) == 5612 || ccDigits.At(4) == 5893 || ccDigits.At(4) == 6304 || ccDigits.At(4) == 6759 || ccDigits.At(4) == 6761 || ccDigits.At(4) == 6762 || ccDigits.At(4) == 6763 || number[:3] == "0604" || ccDigits.At(4) == 6390:
		return "Maestro", nil
	case ccDigits.at(4) == 5019:
		return  "Dankort", nil
	case ccDigits.at(2) >= 51 && ccDigits.At(2) <= 55:
		return  "MasterCard", nil
	case ccDigits.at(4) == 4026 || ccDigits.At(6) == 417500 || ccDigits.At(4) == 4405 || ccDigits.At(4) == 4508 || ccDigits.At(4) == 4844 || ccDigits.At(4) == 4913 || ccDigits.At(4) == 4917:
		return "Visa Electron", nil
	case ccDigits.at(1) == 4:
		return "Visa", nil
	default:
		return "", errors.New("Unknown credit card method.")
	}
}
