package default_visitor

import(
	"github.com/ottemo/foundation/models/visitor"
)

const (
	VISITOR_COLLECTION_NAME = "visitor"
)


type DefaultVisitor struct {
	id string

	Email     string
	FirstName string
	LastName  string

	BillingAddress  visitor.I_VisitorAddress
	ShippingAddress visitor.I_VisitorAddress
}
