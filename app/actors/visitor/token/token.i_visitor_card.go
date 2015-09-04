package token

// GetVisitorID returns the Visitor ID for the Visitor Card
func (it *DefaultVisitorCard) GetVisitorID() string { return it.visitorID }

// GetHolderName returns the Holder of the Credit Card
func (it *DefaultVisitorCard) GetHolderName() string { return it.Holder }

// GetPaymentMethod returns the Payment method code of the Visitor Card
func (it *DefaultVisitorCard) GetPaymentMethod() string { return it.Payment }

// GetType will return the Type of the Visitor Card
func (it *DefaultVisitorCard) GetType() string { return it.Type }

// GetNumber will return the Number attribute of the Visitor Card
func (it *DefaultVisitorCard) GetNumber() string { return it.Number }

// GetExpirationDate will return the Expiration date  of the Visitor Card
func (it *DefaultVisitorCard) GetExpirationDate() string { return it.ExpirationDate }

// GetToken will return the Token of the Visitor Card
func (it *DefaultVisitorCard) GetToken() string { return it.Token }

// IsExpired will return Expired status of the Visitor Card
func (it *DefaultVisitorCard) IsExpired() bool {

	return it.ExpirationDate != ""
}