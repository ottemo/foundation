package braintree

import (
	"github.com/lionelbarrow/braintree-go"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/visitor"
)

func getBraintreeCustomerToken(visitorID string) string {
	var absendToken = ""

	if visitorID == "" {
		env.ErrorDispatch(env.ErrorNew(constErrorModule, constErrorLevel, "0f6c678f-66a3-470e-8a80-5cc2ff619058", "empty visitor ID passed to look up customer token"))
		return absendToken
	}

	model, _ := visitor.GetVisitorCardCollectionModel()
	model.ListFilterAdd("visitor_id", "=", visitorID)
	model.ListFilterAdd("payment", "=", constCCMethodCode)

	// 3rd party customer identifier, used by braintree
	err := model.ListAddExtraAttribute("customer_id")
	if err != nil {
		env.ErrorDispatch(err)
	}

	tokens, err := model.List()
	if err != nil {
		env.ErrorDispatch(err)
	}

	for _, token := range tokens {
		return utils.InterfaceToString(token.Extra["customer_id"])
	}

	return absendToken
}

func formatCardExpirationDate(card braintree.CreditCard) string {
	var expirationDate = utils.InterfaceToString(card.ExpirationMonth)

	// pad with a zero
	if len(card.ExpirationMonth) < 2 {
		expirationDate = "0" + expirationDate
	}

	// append the last two year digits
	year := utils.InterfaceToString(card.ExpirationYear)
	if len(year) == 4 {
		expirationDate = expirationDate + year[2:]
	} else {
		env.ErrorDispatch(env.ErrorNew(constErrorModule, constErrorLevel, "950aea13-16e8-4d20-9ad0-f5cee26c03c2", "unexpected year length coming back from braintree "+year))
	}

	return expirationDate
}
