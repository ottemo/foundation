package defaultproduct

func (dpm *DefaultProductModel) GetSku() string  { return dpm.Sku }
func (dpm *DefaultProductModel) GetName() string { return dpm.Name }

func (dpm *DefaultProductModel) GetPrice() float64 { return 10.5 }
