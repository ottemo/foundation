package visitor

import (
	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"time"
)

const (
	COLLECTION_NAME_VISITOR = "visitor"

	EMAIL_VALIDATE_EXPIRE = 60 * 60 * 24
)

type DefaultVisitor struct {
	id string

	Email      string
	FacebookId string
	GoogleId   string

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

type DefaultVisitorCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
