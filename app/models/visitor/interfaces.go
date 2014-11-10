package visitor

import (
	"time"

	"github.com/ottemo/foundation/app/models"
)

// Constants related to working with Visitors and their attributes
const (
	ModelNameVisitor                  = "Visitor"
	ModelNameVisitorCollection        = "VisitorCollection"
	ModelNameVisitorAddress           = "VisitorAddress"
	ModelNameVisitorAddressCollection = "VisitorAddressCollection"

	SessionKeyVisitorID = "visitor_id"
)

// I_Visitor is the default Visitor interface
type I_Visitor interface {
	GetEmail() string
	GetFacebookID() string
	GetGoogleID() string

	GetFullName() string
	GetFirstName() string
	GetLastName() string

	GetBirthday() time.Time
	GetCreatedAt() time.Time

	GetShippingAddress() I_VisitorAddress
	GetBillingAddress() I_VisitorAddress

	SetShippingAddress(address I_VisitorAddress) error
	SetBillingAddress(address I_VisitorAddress) error

	SetPassword(passwd string) error
	CheckPassword(passwd string) bool
	GenerateNewPassword() error

	IsAdmin() bool

	IsValidated() bool
	Invalidate() error
	Validate(key string) error

	LoadByEmail(email string) error
	LoadByFacebookID(facebookID string) error
	LoadByGoogleID(googleID string) error

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Listable
	models.I_CustomAttributes
}

// I_VisitorCollection is the interface for working with groups of Visitors
type I_VisitorCollection interface {
	ListVisitors() []I_Visitor

	models.I_Collection
}

// I_VisitorAddress is the interface which holds address information for Visitors
type I_VisitorAddress interface {
	GetVisitorID() string

	GetFirstName() string
	GetLastName() string

	GetCompany() string

	GetCountry() string
	GetState() string
	GetCity() string

	GetAddress() string
	GetAddressLine1() string
	GetAddressLine2() string

	GetPhone() string
	GetZipCode() string

	models.I_Model
	models.I_Object
	models.I_Storable
}

// I_VisitorAddressCollection is the interfac for working wiht groups of addresses
type I_VisitorAddressCollection interface {
	ListVisitorsAddresses() []I_VisitorAddress

	models.I_Collection
}
