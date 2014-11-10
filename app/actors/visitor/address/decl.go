package address

import (
	"github.com/ottemo/foundation/db"
)

// Constants for working with the Visitor Adddress collections
const (
	CollectionNameVisitorAddress = "visitor_address"
)

// DefaultVisitorAddress is the base struct for holding a Visitor address
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

// DefaultVisitorAddressCollection holds the customm attributes for addresses
type DefaultVisitorAddressCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
