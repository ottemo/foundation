package stripe

import (
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// getCCBrand Standardize the cc brand
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

// getStripeCustomerToken We attach customer tokens to card tokens in the visitor_token table
// - the customer token is sensitive data because you can make a charge with it alone
// - if you are going to make a charge against a card that is attached to a customer though,
//   you must attach the customer id
func getStripeCustomerToken(vid string) string {
	const customerTokenPrefix = "cus"

	if vid == "" {
		_ = env.ErrorDispatch(env.ErrorNew(ConstErrorModule, 1, "2ecfa3ec-7cfc-4783-9060-8467ca63beae", "empty vid passed to look up customer token"))
		return ""
	}

	model, _ := visitor.GetVisitorCardCollectionModel()
	if err := model.ListFilterAdd("visitor_id", "=", vid); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "ec585a6e-2046-47da-be9e-d010d5149838", err.Error())
	}
	if err := model.ListFilterAdd("payment", "=", ConstPaymentCode); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "dc4bc550-705b-497e-b771-3c49ecd15d3a", err.Error())
	}

	// 3rd party customer identifier, used by stripe
	err := model.ListAddExtraAttribute("customer_id")
	if err != nil {
		_ = env.ErrorDispatch(err)
	}

	tokens, err := model.List()
	if err != nil {
		_ = env.ErrorDispatch(err)
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
