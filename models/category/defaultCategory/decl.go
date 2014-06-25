package defaultCategory

import "github.com/ottemo/foundation/models/product"

const (
	CATEGORY_COLLECTION_NAME                  = "category"
	CATEGORY_PRODUCT_JUNCTION_COLLECTION_NAME = "category_product"
)

type DefaultCategory struct {
	id string

	Name     string
	Parent   ICategory
	Products []product.IProduct
}
