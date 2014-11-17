package visitor

import (
	"github.com/ottemo/foundation/app/models"
	"time"
)

const (
	MODEL_NAME_VISITOR                    = "Visitor"
	MODEL_NAME_VISITOR_COLLECTION         = "VisitorCollection"
	MODEL_NAME_VISITOR_ADDRESS            = "VisitorAddress"
	MODEL_NAME_VISITOR_ADDRESS_COLLECTION = "VisitorAddressCollection"

	SESSION_KEY_VISITOR_ID = "visitor_id"
)

type I_Visitor interface {
	GetEmail() string
	GetFacebookId() string
	GetGoogleId() string

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
	LoadByFacebookId(facebookId string) error
	LoadByGoogleId(googleId string) error

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Listable
	models.I_CustomAttributes
}

type I_VisitorCollection interface {
	ListVisitors() []I_Visitor

	models.I_Collection
}

type I_VisitorAddress interface {
	GetVisitorId() string

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

type I_VisitorAddressCollection interface {
	ListVisitorsAddresses() []I_VisitorAddress

	models.I_Collection
}
