package reporting

import (
	"time"
)

type ProductPerfItem struct {
	Name       string  `json:"name"`
	Sku        string  `json:"sku"`
	GrossSales float64 `json:"gross_sales"`
	UnitsSold  int     `json:"units_sold"`
}

type ProductPerf []ProductPerfItem

func (a ProductPerf) Len() int {
	return len(a)
}

func (a ProductPerf) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ProductPerf) Less(i, j int) bool {
	if a[i].UnitsSold == a[j].UnitsSold {
		return a[i].Name < a[j].Name
	}

	return a[i].UnitsSold > a[j].UnitsSold
}

type CustomerActivityItem struct {
	Email            string    `json:"email"`
	Name             string    `json:"name"`
	TotalSales       float64   `json:"total_sales"`
	TotalOrders      int       `json:"total_orders"`
	AverageSales     float64   `json:"avg_sales"`
	EarliestPurchase time.Time `json:"earliest_purchase"`
	LatestPurchase   time.Time `json:"latest_purchase"`
}

type CustomerActivityBySales []CustomerActivityItem

func (a CustomerActivityBySales) Len() int {
	return len(a)
}

func (a CustomerActivityBySales) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a CustomerActivityBySales) Less(i, j int) bool {
	if a[i].TotalSales == a[j].TotalSales {
		return a[i].Email < a[j].Email
	}

	return a[i].TotalSales > a[j].TotalSales
}

type CustomerActivityByOrders []CustomerActivityItem

func (a CustomerActivityByOrders) Len() int {
	return len(a)
}

func (a CustomerActivityByOrders) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a CustomerActivityByOrders) Less(i, j int) bool {
	if a[i].TotalOrders == a[j].TotalOrders {
		return a[i].Email < a[j].Email
	}

	return a[i].TotalOrders > a[j].TotalOrders

}
