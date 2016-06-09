// Package token allows to create and use tokens
package token

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"time"
)

// Package global constants
const (
	ConstCollectionNameVisitorToken = "visitor_token"

	ConstErrorModule = "visitor/token"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultVisitorCard is a default implementer of InterfaceVisitorCard
type DefaultVisitorCard struct {
	id        string
	visitorID string

	Holder  string
	Payment string

	Type   string
	Number string

	ExpirationDate  string
	ExpirationMonth int
	ExpirationYear  int

	CreatedAt    time.Time
	TokenUpdated time.Time

	tokenID    string
	customerID string
}

// DefaultVisitorCardCollection is a default implementer of InterfaceVisitorCardCollection
type DefaultVisitorCardCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}
