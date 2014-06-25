package defaultProduct

import "github.com/ottemo/foundation/models/attribute"

type ProductModel struct {
	id string

	Sku  string
	Name string

	*attribute.CustomAttributes
}
