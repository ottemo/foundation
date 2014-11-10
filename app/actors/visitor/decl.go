package visitor

import (
	"time"

	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
)

// Constants for Visitor
const (
	CollectionNameVisitor = "visitor"

	EmailValidateExpire = 60 * 60 * 24
)

// DefaultVisitor struct with default values
type DefaultVisitor struct {
	id string

	Email      string
	FacebookID string
	GoogleID   string

	FirstName string
	LastName  string

	BillingAddress  visitor.I_VisitorAddress
	ShippingAddress visitor.I_VisitorAddress

	Password    string
	ValidateKey string

	Admin bool

	Birthday  time.Time
	CreatedAt time.Time

	*attributes.CustomAttributes
}

// DefaultVisitorCollection struct holds the db collection and custom attributes
type DefaultVisitorCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
