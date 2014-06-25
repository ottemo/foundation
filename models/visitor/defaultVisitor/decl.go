package defaultVisitor

import "github.com/ottemo/foundation/models/visitor"

const (
	VISITOR_COLLECTION_NAME = "visitor"
)

type DefaultVisitor struct {
	id string

	Email     string
	FirstName string
	LastName  string

	BillingAddress  visitor.IVisitorAddress
	ShippingAddress visitor.IVisitorAddress
}
