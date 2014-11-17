package address

import (
	"github.com/ottemo/foundation/db"
)

const (
	COLLECTION_NAME_VISITOR_ADDRESS = "visitor_address"
)

type DefaultVisitorAddress struct {
	id         string
	visitor_id string

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

type DefaultVisitorAddressCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
