package stripe

func getCCBrand(ccBrand string) string {
	switch ccBrand {
	case "Visa":
		return "VISA"

	case "American Express":
		return "AmericanExpress"

	case "MasterCard":
		return "MasterCard"

	case "Discover":
		return "Discover"

	case "JCB":
		return "JCB"

	case "Diners Club":
		return "DinersClub"

	case "Unknown":
		return "Unknown"
	}

	return ccBrand
}
