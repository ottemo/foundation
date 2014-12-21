// Package address is a default implementation of models/visitor package visitor address related  interfaces
package address

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstCollectionNameVisitorAddress = "visitor_address"

	ConstErrorModule = "visitor/address"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultVisitorAddress is a default implementer of InterfaceVisitorAddress
type DefaultVisitorAddress struct {
	id        string
	visitorID string

	FirstName string
	LastName  string

	Company string

	Country string
	State   string
	City    string

	AddressLine1 string
	AddressLine2 string

	Phone   string
	ZipCode string
}

// DefaultVisitorAddressCollection is a default implementer of InterfaceVisitorAddressCollection
type DefaultVisitorAddressCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}
