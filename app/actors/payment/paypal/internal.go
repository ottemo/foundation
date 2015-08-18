package paypal

import (
	"github.com/ottemo/foundation/app/models"
)

// getCreditCardName returns credit card valid name
func getCreditCardName(creditCardType string) string {

	switch creditCardType {
	case "0", "VI", "Visa", "VISA":
		return "VISA"

	case "1", "MC", "Master Card", "MasterCard":
		return "MasterCard"

	case "2", "DS", "Discover":
		return "Discover"

	case "3", "AX", "AE", "AmericanExpress", "American Express":
		return "AmericanExpress"

	case "4", "DC", "Diner’s Club", "Diner’sClub", "DinersClub":
		return "DinersClub"

	case "5", "JC", "JCB":
		return "JCB"
	}

	return creditCardType
}

// getCountryCode return valid country code from it's name
func getCountryCode(country string) string {

	switch country {
	case "USA", "US", "United States", "UnitedStates":
		return "US"

	default:
		for code, name := range models.ConstCountriesList {
			if name == country {
				return code
			}
		}
		return "US"
	}
}
