package taxcloud

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/checkout"
)

const (
	ConstErrorModule = "taxcloud"
	ConstErrorLevel  = env.ConstErrorLevelActor
	ConstPriorityValue = checkout.ConstCalculateTargetGrandTotal - 0.0001

	ConstConfigPathGroup = "general.tax.taxCloud"
	ConstConfigPathAPILoginID = "general.tax.taxCloud.apiLoginID"
	ConstConfigPathAPIKey = "general.tax.taxCloud.apiKey"
	ConstConfigPathEnabled = "general.tax.taxCloud.enabled"
)

// DefaultTax is a default implementer of InterfaceTax
type DefaultTaxCloud struct {}

