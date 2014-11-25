// Package page is a default implementation of cms page related interfaces declared in
// "github.com/ottemo/foundation/app/models/csm" package
package page

import (
	"github.com/ottemo/foundation/db"
	"time"
)

// Package global constants
const (
	ConstCmsPageCollectionName = "cms_page"
)

// DefaultCMSPage is a default implementer of InterfaceCMSPage
type DefaultCMSPage struct {
	id string

	URL string

	Identifier string

	Title   string
	Content string

	MetaKeywords    string
	MetaDescription string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// DefaultCMSPageCollection is a default implementer of InterfaceCMSPageCollection
type DefaultCMSPageCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}
