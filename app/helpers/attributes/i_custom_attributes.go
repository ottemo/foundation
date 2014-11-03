package attributes

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// initializes helper before usage
func (it *CustomAttributes) Init(model string) (*CustomAttributes, error) {
	it.model = model
	it.values = make(map[string]interface{})

	globalCustomAttributesMutex.Lock()

	_, present := globalCustomAttributes[model]

	if present {
		it.attributes = globalCustomAttributes[model]
	} else {

		it.attributes = make(map[string]models.T_AttributeInfo)

		// retrieving information from DB
		//-------------------------------
		customAttributesCollection, err := db.GetCollection(COLLECTION_NAME_CUSTOM_ATTRIBUTES)
		if err != nil {
			return it, env.ErrorNew("Can't get collection 'custom_attributes': " + err.Error())
		}

		customAttributesCollection.AddFilter("model", "=", it.model)
		dbValues, err := customAttributesCollection.Load()
		if err != nil {
			env.ErrorDispatch(err)
			return it, env.ErrorNew("Can't load custom attributes information for '" + it.model + "'")
		}

		// filling attribute info structure
		//---------------------------------
		for _, row := range dbValues {
			attribute := models.T_AttributeInfo{}

			for key, value := range row {
				switch key {
				case "model":
					attribute.Model = utils.InterfaceToString(value)
				case "collection":
					attribute.Collection = utils.InterfaceToString(value)
				case "attribute":
					attribute.Attribute = utils.InterfaceToString(value)
				case "type":
					attribute.Type = utils.InterfaceToString(value)
				case "label":
					attribute.Label = utils.InterfaceToString(value)
				case "group":
					attribute.Group = utils.InterfaceToString(value)
				case "editors":
					attribute.Editors = utils.InterfaceToString(value)
				case "options":
					attribute.Options = utils.InterfaceToString(value)
				case "default":
					attribute.Default = utils.InterfaceToString(value)
				case "validators":
					attribute.Validators = utils.InterfaceToString(value)

				case "isrequired", "required":
					attribute.IsRequired = utils.InterfaceToBool(value)
				case "islayered", "layered":
					attribute.IsLayered = utils.InterfaceToBool(value)
				}
			}

			it.attributes[attribute.Attribute] = attribute
		}

		globalCustomAttributes[it.model] = it.attributes
	}

	globalCustomAttributesMutex.Unlock()

	return it, nil
}

// removes custom attribute from collection
func (it *CustomAttributes) RemoveAttribute(attributeName string) error {

	customAttribute, present := it.attributes[attributeName]
	if !present {
		return env.ErrorNew("There is no attribute '" + attributeName + "' for model '" + it.model + "'")
	}

	customAttributesCollection, err := db.GetCollection(COLLECTION_NAME_CUSTOM_ATTRIBUTES)
	if err != nil {
		return env.ErrorNew("Can't get collection '" + COLLECTION_NAME_CUSTOM_ATTRIBUTES + "': " + err.Error())
	}

	modelCollection, err := db.GetCollection(customAttribute.Collection)
	if err != nil {
		return env.ErrorNew("Can't get attribute '" + customAttribute.Attribute + "' collection '" + customAttribute.Collection + "': " + err.Error())
	}

	err = modelCollection.RemoveColumn(attributeName)
	if err != nil {
		return env.ErrorNew("Can't remove attribute '" + attributeName + "' from collection '" + customAttribute.Collection + "': " + err.Error())
	}

	globalCustomAttributesMutex.Lock()
	delete(globalCustomAttributes, it.model)
	globalCustomAttributesMutex.Unlock()

	customAttributesCollection.AddFilter("model", "=", customAttribute.Model)
	customAttributesCollection.AddFilter("attribute", "=", attributeName)
	_, err = customAttributesCollection.Delete()
	if err != nil {
		return env.ErrorNew("Can't remove attribute '" + attributeName + "' information from 'custom_attributes' collection '" + customAttribute.Collection + "': " + err.Error())
	}

	return nil
}

// extends collection with new custom attribute
func (it *CustomAttributes) AddNewAttribute(newAttribute models.T_AttributeInfo) error {

	if _, present := it.attributes[newAttribute.Attribute]; present {
		return env.ErrorNew("There is already atribute '" + newAttribute.Attribute + "' for model '" + it.model + "'")
	}

	// getting collection where custom attribute information stores
	customAttribuesCollection, err := db.GetCollection(COLLECTION_NAME_CUSTOM_ATTRIBUTES)
	if err != nil {
		return env.ErrorNew("Can't get collection '" + COLLECTION_NAME_CUSTOM_ATTRIBUTES + "': " + err.Error())
	}

	// getting collection where attribute supposed to be
	modelCollection, err := db.GetCollection(newAttribute.Collection)
	if err != nil {
		return env.ErrorNew("Can't get attribute '" + newAttribute.Attribute + "' collection '" + newAttribute.Collection + "': " + err.Error())
	}

	// inserting attribute information in custom_attributes collection
	hashMap := make(map[string]interface{})

	hashMap["model"] = newAttribute.Model
	hashMap["collection"] = newAttribute.Collection
	hashMap["attribute"] = newAttribute.Attribute
	hashMap["type"] = newAttribute.Type
	hashMap["required"] = newAttribute.IsRequired
	hashMap["label"] = newAttribute.Label
	hashMap["group"] = newAttribute.Group
	hashMap["editors"] = newAttribute.Editors
	hashMap["options"] = newAttribute.Options
	hashMap["default"] = newAttribute.Default
	hashMap["validators"] = newAttribute.Validators
	hashMap["layered"] = newAttribute.IsLayered

	newCustomAttributeId, err := customAttribuesCollection.Save(hashMap)

	if err != nil {
		return env.ErrorNew("Can't insert attribute '" + newAttribute.Attribute + "' in collection '" + newAttribute.Collection + "': " + err.Error())
	}

	// inserting new attribute to supposed location
	err = modelCollection.AddColumn(newAttribute.Attribute, newAttribute.Type, false)
	if err != nil {
		customAttribuesCollection.DeleteById(newCustomAttributeId)

		return env.ErrorNew("Can't insert attribute '" + newAttribute.Attribute + "' in collection '" + newAttribute.Collection + "': " + err.Error())
	}

	it.attributes[newAttribute.Attribute] = newAttribute

	return env.ErrorDispatch(err)
}
