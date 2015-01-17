package attributes

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Init initializes helper before usage
func (it *CustomAttributes) Init(model string, collection string) (*CustomAttributes, error) {
	it.model = model
	it.collection = collection
	it.values = make(map[string]interface{})

	globalCustomAttributesMutex.Lock()

	_, present := globalCustomAttributes[model]

	if present {
		it.attributes = globalCustomAttributes[model]
	} else {

		it.attributes = make(map[string]models.StructAttributeInfo)

		// retrieving information from DB
		//-------------------------------
		customAttributesCollection, err := db.GetCollection(ConstCollectionNameCustomAttributes)
		if err != nil {
			return it, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "460f57a5-3c39-4db2-ae41-bce6bad58857", "Can't get collection 'custom_attributes': "+err.Error())
		}

		customAttributesCollection.AddFilter("model", "=", it.model)
		records, err := customAttributesCollection.Load()
		if err != nil {
			env.ErrorDispatch(err)
			return it, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "91e0f7e5-7234-4a33-b94b-bec7437200a5", "Can't load custom attributes information for '"+it.model+"'")
		}

		// filling attribute info structure
		//---------------------------------
		for _, record := range records {
			attribute := models.StructAttributeInfo{}

			for key, value := range record {
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
				case "ispublic", "public":
					attribute.IsPublic = utils.InterfaceToBool(value)
				}
			}

			it.attributes[attribute.Attribute] = attribute
		}

		globalCustomAttributes[it.model] = it.attributes
	}

	globalCustomAttributesMutex.Unlock()

	return it, nil
}

// EditAttribute modifies custom attribute for collection
func (it *CustomAttributes) EditAttribute(attributeName string, attributeValues models.StructAttributeInfo) error {
	customAttribute, present := it.attributes[attributeName]
	if !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d4ba1021-eb4d-4f03-aafd-6a4e33efb5ed", "There is no attribute '"+attributeName+"' for model '"+it.model+"'")
	}

	customAttributesCollection, err := db.GetCollection(ConstCollectionNameCustomAttributes)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3b8b1e23-c2ad-45c5-9252-215084a8cd81", "Can't get collection '"+ConstCollectionNameCustomAttributes+"': "+err.Error())
	}

	customAttributesCollection.AddFilter("model", "=", customAttribute.Model)
	customAttributesCollection.AddFilter("attribute", "=", attributeName)
	records, err := customAttributesCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, record := range records {
		customAttribute.IsRequired = attributeValues.IsRequired
		record["required"] = attributeValues.IsRequired

		customAttribute.Label = attributeValues.Label
		record["label"] = attributeValues.Label

		if customAttribute.Group != "" {
			customAttribute.Group = attributeValues.Group
			record["group"] = attributeValues.Group
		}
		if customAttribute.Editors != "" {
			customAttribute.Editors = attributeValues.Editors
			record["editors"] = attributeValues.Editors
		}

		customAttribute.Options = attributeValues.Options
		record["options"] = attributeValues.Options

		customAttribute.Default = attributeValues.Default
		record["default"] = attributeValues.Default

		customAttribute.Validators = attributeValues.Validators
		record["validators"] = attributeValues.Validators

		customAttribute.IsLayered = attributeValues.IsLayered
		record["layered"] = attributeValues.IsLayered

		customAttribute.IsPublic = attributeValues.IsPublic
		record["public"] = attributeValues.IsPublic

		_, err := customAttributesCollection.Save(record)
		if err != nil {
			return err
		}

		it.attributes[attributeName] = customAttribute
	}

	return nil
}

// RemoveAttribute removes custom attribute from collection
func (it *CustomAttributes) RemoveAttribute(attributeName string) error {

	customAttribute, present := it.attributes[attributeName]
	if !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d4ba1021-eb4d-4f03-aafd-6a4e33efb5ed", "There is no attribute '"+attributeName+"' for model '"+it.model+"'")
	}

	customAttributesCollection, err := db.GetCollection(ConstCollectionNameCustomAttributes)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3b8b1e23-c2ad-45c5-9252-215084a8cd81", "Can't get collection '"+ConstCollectionNameCustomAttributes+"': "+err.Error())
	}

	modelCollection, err := db.GetCollection(customAttribute.Collection)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "861e9bf1-44ac-418b-8249-85613451fc9c", "Can't get attribute '"+customAttribute.Attribute+"' collection '"+customAttribute.Collection+"': "+err.Error())
	}

	err = modelCollection.RemoveColumn(attributeName)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "901ce41c-6802-4ecd-b654-93d96d34b361", "Can't remove attribute '"+attributeName+"' from collection '"+customAttribute.Collection+"': "+err.Error())
	}

	globalCustomAttributesMutex.Lock()
	delete(globalCustomAttributes, it.model)
	globalCustomAttributesMutex.Unlock()

	customAttributesCollection.AddFilter("model", "=", customAttribute.Model)
	customAttributesCollection.AddFilter("attribute", "=", attributeName)
	_, err = customAttributesCollection.Delete()
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "da771c6b-402f-4816-a07d-9602f584b45d", "Can't remove attribute '"+attributeName+"' information from 'custom_attributes' collection '"+customAttribute.Collection+"': "+err.Error())
	}

	return nil
}

// AddNewAttribute extends collection with new custom attribute
func (it *CustomAttributes) AddNewAttribute(newAttribute models.StructAttributeInfo) error {

	if _, present := it.attributes[newAttribute.Attribute]; present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "24aa5125-d8b3-4e55-b321-79eef3eeccb8", "There is already atribute '"+newAttribute.Attribute+"' for model '"+it.model+"'")
	}

	// getting collection where custom attribute information stores
	customAttribuesCollection, err := db.GetCollection(ConstCollectionNameCustomAttributes)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5fcc7f5e-d469-4315-acf2-1ec04391a7f5", "Can't get collection '"+ConstCollectionNameCustomAttributes+"': "+err.Error())
	}

	// getting collection where attribute supposed to be
	modelCollection, err := db.GetCollection(newAttribute.Collection)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "da1337be-d710-4c1c-9e3e-04cca84cb82b", "Can't get attribute '"+newAttribute.Attribute+"' collection '"+newAttribute.Collection+"': "+err.Error())
	}

	// checking collection already existent columns
	if modelCollection.HasColumn(newAttribute.Attribute) {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0402b818-03e9-4c56-bee8-0c1471b8d2ba", "There is already atribute '"+newAttribute.Attribute+"' in collection '"+it.collection+"'")
	}

	// type validation
	if !utils.IsAmongStr(strings.Trim(newAttribute.Type, "[]"),
		db.ConstTypeID,
		db.ConstTypeBoolean,
		db.ConstTypeVarchar,
		db.ConstTypeText,
		db.ConstTypeInteger,
		db.ConstTypeDecimal,
		db.ConstTypeMoney,
		db.ConstTypeFloat,
		db.ConstTypeDatetime,
		db.ConstTypeJSON) {

		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "26f15efc-8490-41f0-82f7-b677ea427af7", "unknown attribute type '"+newAttribute.Type+"'")
	}

	// inserting attribute information in custom_attributes collection
	record := make(map[string]interface{})

	record["model"] = newAttribute.Model
	record["collection"] = newAttribute.Collection
	record["attribute"] = newAttribute.Attribute
	record["type"] = newAttribute.Type
	record["required"] = newAttribute.IsRequired
	record["label"] = newAttribute.Label
	record["group"] = newAttribute.Group
	record["editors"] = newAttribute.Editors
	record["options"] = newAttribute.Options
	record["default"] = newAttribute.Default
	record["validators"] = newAttribute.Validators
	record["layered"] = newAttribute.IsLayered
	record["public"] = newAttribute.IsPublic

	newCustomAttributeID, err := customAttribuesCollection.Save(record)

	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ad98e5b0-7672-4d02-9744-e1beecb88922", "Can't insert attribute '"+newAttribute.Attribute+"' in collection '"+newAttribute.Collection+"': "+err.Error())
	}

	// inserting new attribute to supposed location
	err = modelCollection.AddColumn(newAttribute.Attribute, newAttribute.Type, false)
	if err != nil {
		customAttribuesCollection.DeleteByID(newCustomAttributeID)

		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0c11b43b-ec29-4e08-b614-7a8ec8345c9b", "Can't insert attribute '"+newAttribute.Attribute+"' in collection '"+newAttribute.Collection+"': "+err.Error())
	}

	it.attributes[newAttribute.Attribute] = newAttribute

	return env.ErrorDispatch(err)
}

// GetCustomAttributeCollectionName returns collection name you can use to fill CustomAttributes struct
func (it *CustomAttributes) GetCustomAttributeCollectionName() string {
	return it.collection
}
