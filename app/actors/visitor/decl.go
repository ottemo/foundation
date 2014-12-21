// Package visitor is a default implementation of models/visitor package visitor related interfaces
package visitor

import (
	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"time"
)

// Package global constants
const (
	ConstCollectionNameVisitor = "visitor"

	ConstEmailValidateExpire = 60 * 60 * 24

	ConstErrorModule = "visitor"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultVisitor is a default implementer of InterfaceVisitor
type DefaultVisitor struct {
	id string

	Email      string
	FacebookID string
	GoogleID   string

	FirstName string
	LastName  string

	BillingAddress  visitor.InterfaceVisitorAddress
	ShippingAddress visitor.InterfaceVisitorAddress

	Password    string
	ValidateKey string

	Admin bool

	Birthday  time.Time
	CreatedAt time.Time

	*attributes.CustomAttributes
}

// DefaultVisitorCollection is a default implementer of InterfaceVisitorCollection
type DefaultVisitorCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}
