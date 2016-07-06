// Package visitor is a default implementation of models/visitor package visitor related interfaces
package visitor

import (
	"time"

	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstCollectionNameVisitor = "visitor"

	ConstEmailVerifyExpire = 60 * 60 * 24

	ConstEmailPasswordResetExpire = 30 * 60

	ConstErrorModule = "visitor"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstConfigPathLostPasswordEmailSubject  = "general.mail.lost_password_email_subject"
	ConstConfigPathLostPasswordEmailTemplate = "general.mail.lost_password_email_template"
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

	Password        string
	VerificationKey string

	Admin bool

	CreatedAt time.Time

	*attributes.CustomAttributes
}

// DefaultVisitorCollection is a default implementer of InterfaceVisitorCollection
type DefaultVisitorCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}
