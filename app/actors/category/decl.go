// Package category is a default implementation of interfaces declared in
// "github.com/ottemo/foundation/app/models/category" package
package category

import (
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstCollectionNameCategory                = "category"
	ConstCollectionNameCategoryProductJunction = "category_product"

	ConstErrorModule = "category"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstCategoryMediaTypeImage = "image"
)

// DefaultCategory is a default implementer of InterfaceCategory
type DefaultCategory struct {
	id string

	Enabled     bool
	Name        string
	Description string
	Image       string
	Parent      category.InterfaceCategory
	Path        string
	ProductIds  []string
}

// DefaultCategoryCollection is a default implementer of InterfaceCategoryCollection
type DefaultCategoryCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}
