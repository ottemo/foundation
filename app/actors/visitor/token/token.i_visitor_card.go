package token

import (
	"github.com/ottemo/foundation/utils"
	"time"
)

// GetVisitorID returns the Visitor ID for the Visitor Card
func (it *DefaultVisitorCard) GetVisitorID() string {
	return it.visitorID
}

// GetHolderName returns the Holder of the Credit Card
func (it *DefaultVisitorCard) GetHolderName() string {
	return it.Holder
}

// GetPaymentMethodCode returns the Payment method code of the Visitor Card
func (it *DefaultVisitorCard) GetPaymentMethodCode() string {
	return it.Payment
}

// GetType will return the Type of the Visitor Card
func (it *DefaultVisitorCard) GetType() string {
	return it.Type
}

// GetNumber will return the Number attribute of the Visitor Card
func (it *DefaultVisitorCard) GetNumber() string {
	return it.Number
}

// GetExpirationDate will return the Expiration date  of the Visitor Card
func (it *DefaultVisitorCard) GetExpirationDate() string {

	if it.ExpirationDate == "" {
		it.ExpirationDate = utils.InterfaceToString(it.ExpirationMonth) + "/" + utils.InterfaceToString(it.ExpirationYear)
	}

	return it.ExpirationDate
}

// GetToken will return the Token of the Visitor Card
func (it *DefaultVisitorCard) GetToken() string {
	return it.tokenID
}

// IsExpired will return Expired status of the Visitor Card
func (it *DefaultVisitorCard) IsExpired() bool {
	current := time.Now()
	return it.ExpirationYear < utils.InterfaceToInt(current.Year()) || it.ExpirationMonth < utils.InterfaceToInt(current.Month())
}
