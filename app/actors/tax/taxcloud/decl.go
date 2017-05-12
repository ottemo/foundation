package taxcloud

import (
	"github.com/ottemo/foundation/app/actors/tax/taxcloud/model"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
)

const (
	ConstErrorModule   = "taxcloud"
	ConstErrorLevel    = env.ConstErrorLevelActor
	ConstPriorityValue = checkout.ConstCalculateTargetGrandTotal - 0.0001

	ConstConfigPathGroup      = "tax.taxCloud"
	ConstConfigPathAPILoginID = "tax.taxCloud.apiLoginID"
	ConstConfigPathAPIKey     = "tax.taxCloud.apiKey"
	ConstConfigPathEnabled    = "tax.taxCloud.enabled"

	ConstTicIdAttribute = "tic_id"

	ConstDefaultTicID = 0
)

// TaxCloudPriceAdjustment is a tax calculation helper based on TaxCloud service
type TaxCloudPriceAdjustment struct{}

// TicDelegate type implements InterfaceAttributesDelegate to extend product model
type TicDelegate struct {
	productInstance product.InterfaceProduct

	productTicPtr *model.InterfaceProductTic
}

var ticDelegate models.InterfaceAttributesDelegate

// ticsCachePtr contains taxability information codes from TaxCloud
var ticsCachePtr *map[int]string

// DefaultProductTic is a default model for storing tic attribute per product
type DefaultProductTic struct {
	id        string
	productID string
	ticID     int
}
