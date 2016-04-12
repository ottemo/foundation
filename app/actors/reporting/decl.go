package reporting

type AggrOrderItems struct {
	Name       string  `json:"name"`
	Sku        string  `json:"sku"`
	GrossSales float64 `json:"gross_sales"`
	UnitsSold  int     `json:"units_sold"`
}

type ByUnitsSold []AggrOrderItems

func (a ByUnitsSold) Len() int {
	return len(a)
}

func (a ByUnitsSold) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByUnitsSold) Less(i, j int) bool {
	if a[i].UnitsSold == a[j].UnitsSold {
		return a[i].Name < a[j].Name
	} else {
		return a[i].UnitsSold > a[j].UnitsSold
	}
}
