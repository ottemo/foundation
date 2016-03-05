package shareasale

import "github.com/ottemo/foundation/env"

// package constants for Share A Sale module
const (
	ConstErrorModule = "shareasale"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstShareASaleURL = "https://shareasale.com/sale.cfm"

	ConstConfigPathShareASale           = "general.shareasale"
	ConstConfigPathShareASaleEnabled    = "general.shareasale.enabled"
	ConstConfigPathShareASaleMerchantID = "general.shareasale.merchant_id"
)

// AffiliateSale is a struct to hold a single affiliate order for Share A Sale promotions
type AffiliateSale struct {
	SubTotal   float64
	OrderNo    string
	MerchantID string
}
