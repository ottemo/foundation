// Package block is a default implementation of cms block related interfaces declared in
// "github.com/ottemo/foundation/app/models/csm" package
package block

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"time"
)

// Package global constants
const (
	ConstCmsBlockCollectionName = "cms_block"

	ConstErrorModule = "cms/block"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultCMSBlock is a default implementer of InterfaceCMSBlock
type DefaultCMSBlock struct {
	id string

	Identifier string
	Content    string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// DefaultCMSBlockCollection is a default implementer of InterfaceCMSBlockCollection
type DefaultCMSBlockCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}
