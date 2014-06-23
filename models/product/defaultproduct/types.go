package defaultproduct

type DefaultProductModel struct {
	id string

	Sku  string
	Name string

	*custom_attributes.CustomAttributes
}
