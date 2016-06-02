package product

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetEnabled returns enabled flag for the given product
func (it *DefaultProduct) GetEnabled() bool {
	return it.Enabled
}

// GetSku returns requested sku for the given product
func (it *DefaultProduct) GetSku() string {
	return it.Sku
}

// GetName returns the name of the given product
func (it *DefaultProduct) GetName() string {
	return it.Name
}

// GetShortDescription returns the short description of the requested product
func (it *DefaultProduct) GetShortDescription() string {
	return it.ShortDescription
}

// GetDescription returns the long description of the requested product
func (it *DefaultProduct) GetDescription() string {
	return it.Description
}

// GetDefaultImage returns the imaged identified as defult for the given product
func (it *DefaultProduct) GetDefaultImage() string {
	return it.DefaultImage
}

// GetPrice returns the price as a float64 for the given product
func (it *DefaultProduct) GetPrice() float64 {
	return it.Price
}

// GetWeight returns the weight for the given product
func (it *DefaultProduct) GetWeight() float64 {
	return it.Weight
}

// GetOptions returns current products possible options as a map[string]interface{}
func (it *DefaultProduct) GetOptions() map[string]interface{} {
	options := it.Options
	eventData := map[string]interface{}{"id": it.GetID(), "product": it, "options": options}
	env.Event("product.getOptions", eventData)
	return options
}

// GetRelatedProductIds returns the related product id list
func (it *DefaultProduct) GetRelatedProductIds() []string {
	return it.RelatedProductIds
}

// GetRelatedProducts returns related products instances list
func (it *DefaultProduct) GetRelatedProducts() []product.InterfaceProduct {
	var result []product.InterfaceProduct

	for _, productID := range it.RelatedProductIds {
		if productID == "" {
			continue
		}

		productModel, err := product.LoadProductByID(productID)
		if err == nil {
			result = append(result, productModel)
		}
	}

	return result
}

// GetQty updates and returns the product qty if stock manager is available and Qty was not set to instance
// otherwise returns qty was set
func (it *DefaultProduct) GetQty() int {
	if stockManager := product.GetRegisteredStock(); it.Qty == 0 && stockManager != nil {
		it.Qty = stockManager.GetProductQty(it.GetID(), it.GetAppliedOptions())
	}

	return it.Qty
}

// GetAppliedOptions returns applied options for current product instance
func (it *DefaultProduct) GetAppliedOptions() map[string]interface{} {
	if it.appliedOptions != nil {
		return it.appliedOptions
	}
	return make(map[string]interface{})
}

func (it *DefaultProduct) GetInventory() []map[string]interface{} {
	if stockManager := product.GetRegisteredStock(); stockManager != nil {
		it.Inventory = stockManager.GetProductOptions(it.GetID())
	}

	return it.Inventory
}

// ApplyOptions updates current product attributes according to given product options,
// returns error if specified option are not possible for the product
func (it *DefaultProduct) ApplyOptions(options map[string]interface{}) error {
	// taking item specified options and product options
	productOptions := it.GetOptions()

	// storing start price for a case of percentage price modifier
	startPrice := it.GetPrice()

	// sorting applicable product attributes according to "order" field
	// optionsApplyOrder := make([]string, 0)
	var optionsApplyOrder []string

	for itemOptionName := range options {

		// looking only for options that customer set for item
		if productOption, present := productOptions[itemOptionName]; present {
			if productOption, ok := productOption.(map[string]interface{}); ok {

				orderValue := int(^uint(0) >> 1) // default order - max integer
				if optionValue, present := productOption["order"]; present {
					orderValue = utils.InterfaceToInt(optionValue)
				}

				// encoding key order to string "000000000000001 [attribute name]"
				// for future sort as string (16 digits - max for js integer)
				key := fmt.Sprintf("%.16d %s", orderValue, itemOptionName)
				optionsApplyOrder = append(optionsApplyOrder, key)
			}
		}
	}
	sort.Strings(optionsApplyOrder)

	// function to modify orderItem according to option values
	applyOptionModifiers := func(optionToApply map[string]interface{}) {

		// price modifier
		if optionValue, present := optionToApply["price"]; present {
			if stringValue, ok := optionValue.(string); ok {
				if stringValue != "" && strings.Trim(stringValue, "1234567890+-%.") == "" {
					isDelta := false
					isPercent := false

					stringValue = strings.TrimSpace(stringValue)
					if strings.HasSuffix(stringValue, "%") {
						isPercent = true
						stringValue = strings.TrimSuffix(stringValue, "%")
					}

					var priceValue float64
					switch {
					case strings.HasPrefix(stringValue, "+"):
						stringValue = strings.TrimPrefix(stringValue, "+")
						isDelta = true
						priceValue = utils.InterfaceToFloat64(stringValue)
					case strings.HasPrefix(stringValue, "-"):
						stringValue = strings.TrimPrefix(stringValue, "-")
						isDelta = true
						priceValue = -1 * utils.InterfaceToFloat64(stringValue)
					default:
						priceValue = utils.InterfaceToFloat64(stringValue)
					}

					if isPercent {
						it.Price += startPrice * priceValue / 100
					} else if isDelta {
						it.Price += priceValue
					} else {
						it.Price = priceValue
					}
				}
			} else {
				it.Set("price", optionValue)
			}
		}

		// sku modifier
		if optionValue, present := optionToApply["sku"]; present {
			skuModifier := utils.InterfaceToString(optionValue)
			if strings.HasPrefix(skuModifier, "-") || strings.HasPrefix(skuModifier, "_") {
				it.Sku += skuModifier
			} else {
				it.Sku = skuModifier
			}
		}
	}

	// loop over item applied option in right order
	for _, itemOptionName := range optionsApplyOrder {
		// decoding key order after sort
		itemOptionName := itemOptionName[strings.Index(itemOptionName, " ")+1:]
		itemOptionValue := options[itemOptionName]

		if productOption, present := productOptions[itemOptionName]; present {
			if productOptions, ok := productOption.(map[string]interface{}); ok {

				// product option itself can contain price, sku modifiers
				applyOptionModifiers(productOptions)

				// if product option value have predefined option values, then checking their modifiers
				if productOptionValues, present := productOptions["options"]; present {
					if productOptionValues, ok := productOptionValues.(map[string]interface{}); ok {

						// option user set can be single on multi-value
						// making it uniform
						// itemOptionValueSet := make([]string, 0)
						var itemOptionValueSet []string
						switch typedOptionValue := itemOptionValue.(type) {
						case string:
							itemOptionValueSet = append(itemOptionValueSet, typedOptionValue)
						case []string:
							itemOptionValueSet = typedOptionValue
						case []interface{}:
							for _, value := range typedOptionValue {
								if value, ok := value.(string); ok {
									itemOptionValueSet = append(itemOptionValueSet, value)
								}
							}
						default:
							return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6d02be30-ca5e-46e4-94c7-2f01782f30b2", "unexpected option value for "+itemOptionName+" option")
						}

						// loop through option values customer set for product
						for _, itemOptionValue := range itemOptionValueSet {

							if productOptionValue, present := productOptionValues[itemOptionValue]; present {
								if productOptionValue, ok := productOptionValue.(map[string]interface{}); ok {
									applyOptionModifiers(productOptionValue)
								}
							} else {
								return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8f2baf8a-af91-44f3-8364-b42099959ec4", "invalid '"+itemOptionName+"' option value: '"+itemOptionValue)
							}

						}

						// cleaning option values were not used by customer
						for productOptionValueName := range productOptionValues {
							if !utils.IsInArray(productOptionValueName, itemOptionValueSet) {
								delete(productOptionValues, productOptionValueName)
							}
						}
					}
				}
			}
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "96246e83-fb80-4781-b671-2d3d75a65e56", "unknown option '"+itemOptionName+"'")
		}
	}

	// cleaning options were not used by customer
	for productOptionName, productOption := range productOptions {
		if _, present := options[productOptionName]; present {
			if productOption, ok := productOption.(map[string]interface{}); ok {
				productOption["value"] = options[productOptionName]
			}
		} else {
			delete(productOptions, productOptionName)
		}
	}

	it.Price = utils.RoundPrice(it.Price)

	it.appliedOptions = options

	// stock management stuff
	if stockManager := product.GetRegisteredStock(); stockManager != nil {
		it.Qty = stockManager.GetProductQty(it.GetID(), it.GetAppliedOptions())
	}

	return nil
}
