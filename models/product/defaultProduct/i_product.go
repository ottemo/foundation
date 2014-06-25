package defaultProduct

func (it *ProductModel) GetSku() string  { return it.Sku }
func (it *ProductModel) GetName() string { return it.Name }

func (it *ProductModel) GetPrice() float64 { return 10.5 }
