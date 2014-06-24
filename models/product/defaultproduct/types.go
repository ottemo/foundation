package defaultproduct

import "github.com/ottemo/foundation/models"

type DefaultProductModel struct {
	id string

	Sku  string
	Name string

	*models.CustomAttribute
}
