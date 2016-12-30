package braintree

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/visitor"
)

func getTokenByVisitorID(visitorID string) string {
	var absendToken = ""

	if visitorID == "" {
		env.ErrorDispatch(env.ErrorNew(constErrorModule, constErrorLevel, "0f6c678f-66a3-470e-8a80-5cc2ff619058", "empty visitor ID passed to look up customer token"))
		return absendToken
	}

	model, _ := visitor.GetVisitorCardCollectionModel()
	model.ListFilterAdd("visitor_id", "=", visitorID)
	model.ListFilterAdd("payment", "=", constCCMethodCode)

	// 3rd party customer identifier, used by braintree
	// TODO: separate function to add 3rd party identifier
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
