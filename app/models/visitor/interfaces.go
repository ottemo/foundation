// Package visitor represents abstraction of business layer visitor object
package visitor

import (
	"time"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstModelNameVisitor                  = "Visitor"
	ConstModelNameVisitorCollection        = "VisitorCollection"
	ConstModelNameVisitorAddress           = "VisitorAddress"
	ConstModelNameVisitorAddressCollection = "VisitorAddressCollection"
	ConstModelNameVisitorCard              = "VisitorCard"
	ConstModelNameVisitorCardCollection    = "VisitorCardCollection"

	ConstSessionKeyVisitorID = "visitor_id"

	ConstErrorModule = "visitor"
	ConstErrorLevel  = env.ConstErrorLevelModel
)

// InterfaceVisitor represents interface to access business layer implementation of visitor object
type InterfaceVisitor interface {
	GetEmail() string
	GetFacebookID() string
	GetGoogleID() string

	GetFullName() string
	GetFirstName() string
	GetLastName() string

	GetCreatedAt() time.Time

	GetShippingAddress() InterfaceVisitorAddress
	GetBillingAddress() InterfaceVisitorAddress

	SetShippingAddress(address InterfaceVisitorAddress) error
	SetBillingAddress(address InterfaceVisitorAddress) error

	SetPassword(passwd string) error
	CheckPassword(passwd string) bool
	GenerateNewPassword() error

	ResetPassword() error
	UpdateResetPassword(key string, passwd string) error

	IsAdmin() bool
	IsGuest() bool

	IsVerified() bool
	Invalidate() error
	Validate(key string) error

	LoadByEmail(email string) error
	LoadByFacebookID(facebookID string) error
	LoadByGoogleID(googleID string) error

	models.InterfaceModel
	models.InterfaceObject
	models.InterfaceStorable
	models.InterfaceListable
	models.InterfaceCustomAttributes
}

// InterfaceVisitorCollection represents interface to access business layer implementation of visitor collection
type InterfaceVisitorCollection interface {
	ListVisitors() []InterfaceVisitor

	models.InterfaceCollection
}

// InterfaceVisitorAddress represents interface to access business layer implementation of visitor address object
type InterfaceVisitorAddress interface {
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

	models.InterfaceModel
	models.InterfaceObject
	models.InterfaceStorable
}

// InterfaceVisitorAddressCollection represents interface to access business layer implementation of visitor address collection
type InterfaceVisitorAddressCollection interface {
	ListVisitorsAddresses() []InterfaceVisitorAddress

	models.InterfaceCollection
}

// InterfaceVisitorCard represents interface to access business layer implementation of visitor card object
type InterfaceVisitorCard interface {
	GetVisitorID() string

	GetHolderName() string
	GetPaymentMethod() string

	GetType() string
	GetNumber() string
	GetExpirationDate() string

	GetToken() string

	IsExpired() bool

	models.InterfaceModel
	models.InterfaceObject
	models.InterfaceStorable
}

// InterfaceVisitorCardCollection represents interface to access business layer implementation of visitor card collection
type InterfaceVisitorCardCollection interface {
	ListVisitorsCards() []InterfaceVisitorCard

	models.InterfaceCollection
}
