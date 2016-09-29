package product

// DefaultProduct type implements:
// 	- InterfaceProduct
// 	- InterfaceModel
// 	- InterfaceObject
// 	- InterfaceStorable
// 	- InterfaceListable
// 	- InterfaceMedia

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

// ---------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models")
// ---------------------------------------------------------------------------------

// GetModelName returns model name
func (it *DefaultProduct) GetModelName() string {
	return product.ConstModelNameProduct
}

// GetImplementationName returns model implementation name
func (it *DefaultProduct) GetImplementationName() string {
	return "Default" + product.ConstModelNameProduct
}

// New returns new instance of model implementation object
func (it *DefaultProduct) New() (models.InterfaceModel, error) {
	newInstance := new(DefaultProduct)

	customAttributes, err := attributes.CustomAttributes(product.ConstModelNameProduct, ConstCollectionNameProduct)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	newInstance.customAttributes = customAttributes

	externalAttributes, err := attributes.ExternalAttributes(newInstance)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	newInstance.externalAttributes = externalAttributes

	return newInstance, nil
}

// -------------------------------------------------------------------------------------------
// InterfaceProduct implementation (package "github.com/ottemo/foundation/app/models/product")
// -------------------------------------------------------------------------------------------

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

// GetAppliedOptions returns applied options for current product instance
func (it *DefaultProduct) GetAppliedOptions() map[string]interface{} {
	if it.appliedOptions != nil {
		return it.appliedOptions
	}
	return make(map[string]interface{})
}

// ApplyOptions updates current product attributes according to given product options,
// returns error if specified option are not possible for the product
func (it *DefaultProduct) ApplyOptions(options map[string]interface{}) error {
	env.Log("errors.log", env.ConstLogPrefixDebug, "========================\n\nRequested options: "+utils.InterfaceToString(options))
	// taking item specified options and product options
	productOptions := it.GetOptions()
	env.Log("errors.log", env.ConstLogPrefixDebug, "========================\n\nproductOptions: "+utils.InterfaceToString(productOptions))

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

	var selectedProductIDs []string
	var foundOptionsCount int
	for itemOptionName := range options {
		// get product option (color, size, etc)
		productOption, ok := productOptions[itemOptionName]
		env.Log("errors.log", env.ConstLogPrefixDebug, "productOption: "+utils.InterfaceToString(productOption))
		if !ok {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5b7d5166-f687-4c65-b109-5c2cee27e3f5", "unknown option '"+itemOptionName+"'")
		}

		productOptionHashMap, ok := productOption.(map[string]interface{})
		if !ok {
			continue
		}

		// check it have options (color:{options:...})
		optionOptions, ok := productOptionHashMap["options"]
		env.Log("errors.log", env.ConstLogPrefixDebug, "optionOptions: "+utils.InterfaceToString(optionOptions))
		if !ok {
			continue
		}

		optionOptionsHashMap, ok := optionOptions.(map[string]interface{})
		env.Log("errors.log", env.ConstLogPrefixDebug, "optionOptionsHashMap: "+utils.InterfaceToString(optionOptionsHashMap))
		if !ok {
			continue
		}

		// check single value (color:{options:{red:{...}}})
		itemOptionValue := utils.InterfaceToString(options[itemOptionName])
		env.Log("errors.log", env.ConstLogPrefixDebug, "looking for "+utils.InterfaceToString(itemOptionName)+"="+utils.InterfaceToString(itemOptionValue))
		optionOptionsItem, present := optionOptionsHashMap[itemOptionValue]
		if !present {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "11b4cf1b-2324-4ed1-b592-32bd66b751e4", "invalid '"+itemOptionName+"' option value: '"+itemOptionValue)
		}

		optionOptionsItemHashMap, ok := optionOptionsItem.(map[string]interface{})
		if !ok {
			continue
		}
		env.Log("errors.log", env.ConstLogPrefixDebug, "optionOptionsItemHashMap: "+utils.InterfaceToString(optionOptionsItemHashMap))
		foundOptionsCount++

		// check option product IDs (color:{options:{red:{product_ids:"1,2,3", ...}}})
		productIDsHashMap, present := optionOptionsItemHashMap["product_ids"]
		if !present {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e9ed53c4-1400-4ffb-b2a4-5d78b9ee2904", "no product specified for '"+itemOptionName+"' option value: '"+itemOptionValue+"'")
		}

		// filter selected product ids by current option product ids
		productIDs := utils.InterfaceToStringArray(productIDsHashMap)
		env.Log("errors.log", env.ConstLogPrefixDebug, "productIDs: "+utils.InterfaceToString(productIDs));
		if len(selectedProductIDs) == 0 {
			// initialization
			selectedProductIDs = productIDs
		} else {
			// filtering
			var newProductIDs []string
			for _, productID := range productIDs {
				if utils.IsInArray(productID, selectedProductIDs) {
					newProductIDs = append(newProductIDs, productID)
				}
			}
			selectedProductIDs = newProductIDs
		}

		env.Log("errors.log", env.ConstLogPrefixDebug, "found optionOptionsValue: "+utils.InterfaceToString(optionOptionsItemHashMap));
	}

	env.Log("errors.log", env.ConstLogPrefixDebug, "foundOptionsCount: "+utils.InterfaceToString(foundOptionsCount));
	env.Log("errors.log", env.ConstLogPrefixDebug, "selectedProductIDs: "+utils.InterfaceToString(selectedProductIDs));
	if foundOptionsCount > 0 {
		if len(selectedProductIDs) < 1 {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2ae359d4-36b8-4182-bb12-c302938b28ad", "no product specified for selected options")
		} else if len(selectedProductIDs) > 1 {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6e356239-5d7e-42a8-8392-b97c80e56fda", "more than one product specified for selected options")
		} else {
			simpleProduct, err := product.LoadProductByID(selectedProductIDs[0])
			if err != nil {
				return env.ErrorDispatch(err)
			}

			// required attributes of simple product
			it.Enabled = simpleProduct.GetEnabled()
			it.Sku = simpleProduct.GetSku()
			it.Name = simpleProduct.GetName()
			it.Price = simpleProduct.GetPrice()
			it.Weight = simpleProduct.GetWeight()

			// not required attributes of simple product
			if simpleProduct.GetShortDescription() != "" {
				it.ShortDescription = simpleProduct.GetShortDescription()
			}

			if simpleProduct.GetDescription() != "" {
				it.Description = simpleProduct.GetDescription()
			}

			if simpleProduct.GetDefaultImage() != "" {
				it.DefaultImage = simpleProduct.GetDefaultImage()
			}

			// store configurable id
			it.Options["configurable_id"] = it.GetID()

			// required ID attribute
			it.SetID(simpleProduct.GetID())

			return nil
		}
	}

	// loop over item applied option in right order
	for i := 0; i < len(optionsApplyOrder); i++ {
		itemOptionNameKey := optionsApplyOrder[i]
		itemOptionName := itemOptionNameKey[strings.Index(itemOptionNameKey, " ")+1:]
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

	return nil
}

// ----------------------------------------------------------------------------------
// InterfaceObject implementation (package "github.com/ottemo/foundation/app/models")
// ----------------------------------------------------------------------------------

// Get returns an object attribute value or nil
func (it *DefaultProduct) Get(attribute string) interface{} {

	if _, present := it.externalAttributes.ListExternalAttributes()[attribute]; present {
		return it.externalAttributes.Get(attribute)
	}

	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "enable", "enabled":
		return it.Enabled
	case "sku":
		return it.Sku
	case "name":
		return it.Name
	case "short_description":
		return it.ShortDescription
	case "description":
		return it.Description
	case "default_image", "defaultimage":
		return it.DefaultImage
	case "price":
		return it.Price
	case "weight":
		return it.Weight
	case "options":
		return it.GetOptions()
	case "related_pids":
		return it.GetRelatedProductIds()
	}

	return it.customAttributes.Get(attribute)
}

// Set will apply the given attribute value to the product or return an error
func (it *DefaultProduct) Set(attribute string, value interface{}) error {
	env.Log("errors.log", env.ConstLogPrefixDebug, "SET: "+attribute+"="+utils.InterfaceToString(value));
	lowerCaseAttribute := strings.ToLower(attribute)

	if _, present := it.externalAttributes.ListExternalAttributes()[lowerCaseAttribute]; present {
		if err := it.externalAttributes.Set(lowerCaseAttribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
		return nil
	}

	switch lowerCaseAttribute {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)
	case "enable", "enabled":
		it.Enabled = utils.InterfaceToBool(value)
	case "sku":
		it.Sku = utils.InterfaceToString(value)
	case "name":
		it.Name = utils.InterfaceToString(value)
	case "short_description":
		it.ShortDescription = utils.InterfaceToString(value)
	case "description":
		it.Description = utils.InterfaceToString(value)
	case "default_image", "defaultimage":
		it.DefaultImage = utils.InterfaceToString(value)
	case "price":
		it.Price = utils.InterfaceToFloat64(value)
	case "weight":
		it.Weight = utils.InterfaceToFloat64(value)
	case "options":
		it.Options = utils.InterfaceToMap(value)
	case "related_pids":
		it.RelatedProductIds = make([]string, 0)

		switch typedValue := value.(type) {
		case []product.InterfaceProduct:
			for _, productItem := range typedValue {
				it.RelatedProductIds = append(it.RelatedProductIds, productItem.GetID())
			}

		case []interface{}:

			for _, listItem := range typedValue {
				var relatedProductIDs []string

				currentProductID := it.GetID()
				if relatedProductID, ok := listItem.(string); ok &&
					relatedProductID != "" &&
					currentProductID != relatedProductID {

					relatedProductIDs = append(relatedProductIDs, relatedProductID)
				}

				// checking related products existence
				dbCollection, err := db.GetCollection(ConstCollectionNameProduct)
				if err != nil {
					return env.ErrorDispatch(err)
				}
				err = dbCollection.AddFilter("_id", "in", relatedProductIDs)
				if err != nil {
					return env.ErrorDispatch(err)
				}
				err = dbCollection.SetResultColumns("_id")
				if err != nil {
					return env.ErrorDispatch(err)
				}
				records, err := dbCollection.Load()
				if err != nil {
					return env.ErrorDispatch(err)
				}

				// adding only exist products to model
				for _, record := range records {
					productID := utils.InterfaceToString(record["_id"])
					it.RelatedProductIds = append(it.RelatedProductIds, productID)
				}
			}

		default:
			if value != nil {
				return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3c402ecc-7c7d-49ab-879e-16af5f4661ed", "unsupported 'related_pids' attribute value")
			}
		}

	default:
		err := it.customAttributes.Set(attribute, value)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// FromHashMap will populate object attributes from map[string]interface{}
func (it *DefaultProduct) FromHashMap(input map[string]interface{}) error {
	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.LogError(err)
		}
	}
	return nil
}

// ToHashMap returns a map[string]interface{}
func (it *DefaultProduct) ToHashMap() map[string]interface{} {
	result := it.customAttributes.ToHashMap()

	result["_id"] = it.id

	result["enabled"] = it.Enabled

	result["sku"] = it.Sku
	result["name"] = it.Name

	result["short_description"] = it.ShortDescription
	result["description"] = it.Description

	result["default_image"] = it.DefaultImage

	result["price"] = it.Price
	result["weight"] = it.Weight

	result["options"] = it.GetOptions()

	result["related_pids"] = it.Get("related_pids")

	for key, value := range it.externalAttributes.ToHashMap() {
		result[key] = value
	}

	return result
}

// GetAttributesInfo returns the requested object attributes
func (it *DefaultProduct) GetAttributesInfo() []models.StructAttributeInfo {
	result := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "_id",
			Type:       db.ConstTypeID,
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "enabled",
			Type:       db.ConstTypeBoolean,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Enabled",
			Group:      "General",
			Editors:    "boolean",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "sku",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "SKU",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
			Validators: "sku",
		},
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "name",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "short_description",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Short Description",
			Group:      "General",
			Editors:    "multiline_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "description",
			Type:       db.ConstTypeText,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Description",
			Group:      "General",
			Editors:    "multiline_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "default_image",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "DefaultImage",
			Group:      "General",
			Editors:    "image_selector",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "price",
			Type:       db.ConstTypeMoney,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Price",
			Group:      "General",
			Editors:    "price",
			Options:    "",
			Default:    "",
			Validators: "price",
		},
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "weight",
			Type:       db.ConstTypeDecimal,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Weight",
			Group:      "General",
			Editors:    "numeric",
			Options:    "",
			Default:    "",
			Validators: "numeric positive",
		},
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "options",
			Type:       db.ConstTypeJSON,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Options",
			Group:      "Options",
			Editors:    "product_options",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "related_pids",
			Type:       db.TypeArrayOf(db.ConstTypeInteger),
			IsRequired: false,
			IsStatic:   true,
			Label:      "Related Products",
			Group:      "General",
			Editors:    "product_selector",
			Options:    "",
			Default:    "",
		},
	}

	customAttributesInfo := it.customAttributes.GetAttributesInfo()
	for _, customAttribute := range customAttributesInfo {
		result = append(result, customAttribute)
	}

	externalAttributesInfo := it.externalAttributes.GetAttributesInfo()
	for _, externalAttribute := range externalAttributesInfo {
		result = append(result, externalAttribute)
	}

	return result
}

// ------------------------------------------------------------------------------------
// InterfaceStorable implementation (package "github.com/ottemo/foundation/app/models")
// ------------------------------------------------------------------------------------

// GetID returns current product id
func (it *DefaultProduct) GetID() string {
	return it.id
}

// SetID sets current product id
func (it *DefaultProduct) SetID(id string) error {
	it.id = id

	return it.externalAttributes.SetID(id)
}

// Load loads product information from DB
func (it *DefaultProduct) Load(id string) error {

	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := collection.LoadByID(id)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a671dee4-b95b-11e5-a86b-28cfe917b6c7", "Unable to find product by id; "+id)
	}

	err = it.FromHashMap(dbRecord)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.externalAttributes.Load(id)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Delete removes current product from DB
func (it *DefaultProduct) Delete() error {
	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.externalAttributes.Delete()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Save stores current product to DB
func (it *DefaultProduct) Save() error {
	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if it.GetName() == "" || it.GetSku() == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ac7cd02e-0722-4ac8-bbe0-ffa74d091a94", "sku and name should be specified")
	}

	valuesToStore := it.ToHashMap()

	for x := range it.externalAttributes.ListExternalAttributes() {
		delete(valuesToStore, x)
	}

	newID, err := collection.Save(valuesToStore)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// set new ID before saving external attributes, because external attributes requires it
	err = it.SetID(newID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.externalAttributes.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// ------------------------------------------------------------------------------------
// InterfaceListable implementation (package "github.com/ottemo/foundation/app/models")
// ------------------------------------------------------------------------------------

// GetCollection returns collection of current instance type
func (it *DefaultProduct) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(product.ConstModelNameProductCollection)
	if result, ok := model.(product.InterfaceProductCollection); ok {
		return result
	}

	return nil
}

// ---------------------------------------------------------------------------------
// InterfaceMedia implementation (package "github.com/ottemo/foundation/app/models")
// ---------------------------------------------------------------------------------

// AddMedia adds new media assigned to product
func (it *DefaultProduct) AddMedia(mediaType string, mediaName string, content []byte) error {
	productID := it.GetID()
	if productID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "85650715-3acf-4e47-a365-c6e8911d9118", "product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return mediaStorage.Save(it.GetModelName(), productID, mediaType, mediaName, content)
}

// RemoveMedia removes media assigned to product
func (it *DefaultProduct) RemoveMedia(mediaType string, mediaName string) error {
	productID := it.GetID()
	if productID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "87bb383a-cf35-48e0-9d50-ad517ed2e8f9", "product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return mediaStorage.Remove(it.GetModelName(), productID, mediaType, mediaName)
}

// ListMedia lists media assigned to product
func (it *DefaultProduct) ListMedia(mediaType string) ([]string, error) {
	var result []string

	productID := it.GetID()
	if productID == "" {
		return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "45b1ebde-3dd0-4c6c-9960-fddd89f4907f", "product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	return mediaStorage.ListMedia(it.GetModelName(), productID, mediaType)
}

// GetMedia returns content of media assigned to product
func (it *DefaultProduct) GetMedia(mediaType string, mediaName string) ([]byte, error) {
	productID := it.GetID()
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5f5d3c33-de82-4580-a6e7-f5c45e9281e5", "product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return mediaStorage.Load(it.GetModelName(), productID, mediaType, mediaName)
}

// GetMediaPath returns relative location of media assigned to product in media storage
func (it *DefaultProduct) GetMediaPath(mediaType string) (string, error) {
	productID := it.GetID()
	if productID == "" {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0055f93a-5d10-41db-8d93-ea2bb4bee216", "product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return mediaStorage.GetMediaPath(it.GetModelName(), productID, mediaType)
}

// --------------------------------------------------------------------------------------------
// InterfaceCustomAttributes implementation (package "github.com/ottemo/foundation/app/models")
// --------------------------------------------------------------------------------------------

// GetInstance returns current instance delegate attached to
func (it *DefaultProduct) GetInstance() interface{} {
	return it
}

// EditAttribute modifies custom attribute for collection
func (it *DefaultProduct) EditAttribute(attributeName string, attributeValues models.StructAttributeInfo) error {
	return it.customAttributes.EditAttribute(attributeName, attributeValues)
}

// RemoveAttribute removes custom attribute from collection
func (it *DefaultProduct) RemoveAttribute(attributeName string) error {
	return it.customAttributes.RemoveAttribute(attributeName)
}

// AddNewAttribute extends collection with new custom attribute
func (it *DefaultProduct) AddNewAttribute(newAttribute models.StructAttributeInfo) error {
	return it.customAttributes.AddNewAttribute(newAttribute)
}

// GetCustomAttributeCollectionName returns collection name you can use to fill ModelCustomAttributes struct
func (it *DefaultProduct) GetCustomAttributeCollectionName() string {
	return it.customAttributes.GetCustomAttributeCollectionName()
}

// ----------------------------------------------------------------------------------------------
// InterfaceExternalAttributes implementation (package "github.com/ottemo/foundation/app/models")
// ----------------------------------------------------------------------------------------------

// GetInstance() method was implemented before for InterfaceCustomAttributes

// GetExtendedInstance returns current instance delegate attached to
func (it *DefaultProduct) GetExtendedInstance() interface{} {
	return it
}

// AddExternalAttributes registers new delegate for a given attribute
func (it *DefaultProduct) AddExternalAttributes(delegate models.InterfaceAttributesDelegate) error {
	return it.externalAttributes.AddExternalAttributes(delegate)
}

// RemoveExternalAttributes registers new delegate for a given attribute
func (it *DefaultProduct) RemoveExternalAttributes(delegate models.InterfaceAttributesDelegate) error {
	return it.externalAttributes.RemoveExternalAttributes(delegate)
}

// ListExternalAttributes registers new delegate for a given attribute
func (it *DefaultProduct) ListExternalAttributes() map[string]models.InterfaceAttributesDelegate {
	return it.externalAttributes.ListExternalAttributes()
}
