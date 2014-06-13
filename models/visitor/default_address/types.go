package default_address

const (
	VISITOR_ADDRESS_COLLECTION_NAME = "visitor_address"
)


type DefaultVisitorAddress struct {
	id string

	Street string
	City string
	State string
	Phone string
	ZipCode string
}
