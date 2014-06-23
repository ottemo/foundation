package default_product


func (it *DefaultProductModel) GetSku() string  { return it.Sku }
func (it *DefaultProductModel) GetName() string { return it.Name }

func (it *DefaultProductModel) GetPrice() float64 { return 10.5 }
