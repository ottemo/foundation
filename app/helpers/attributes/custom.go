package attributes

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
)

// CustomAttributes initializes helper instance before usage
func CustomAttributes(model string, collection string) (*ModelCustomAttributes, error) {
	result := new(ModelCustomAttributes)

	result.model = model
	result.collection = collection
	result.values = make(map[string]interface{})

	modelCustomAttributesMutex.Lock()
	defer modelCustomAttributesMutex.Unlock()

	_, present := modelCustomAttributes[model]

	info, present := modelCustomAttributes[model]
	if !present {
		info = make(map[string]models.StructAttributeInfo)

		// retrieving information from DB
		//-------------------------------
		customAttributesCollection, err := db.GetCollection(ConstCollectionNameCustomAttributes)
		if err != nil {
			return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "460f57a5-3c39-4db2-ae41-bce6bad58857", "Can't get collection 'custom_attributes': "+err.Error())
		}

		if err := customAttributesCollection.AddFilter("model", "=", result.model); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c79a7389-4105-40ab-9d7a-22fcb5dc2402", "unable to add filter': "+err.Error())
		}
		records, err := customAttributesCollection.Load()
		if err != nil {
			env.LogError(err)
			return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "91e0f7e5-7234-4a33-b94b-bec7437200a5", "Can't load custom attributes information for '"+result.model+"'")
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

			info[attribute.Attribute] = attribute
		}

		modelCustomAttributes[result.model] = info
	}

	return result, nil
}

// --------------------------------------------------------------------------------------------
// InterfaceCustomAttributes implementation (package "github.com/ottemo/foundation/app/models")
// --------------------------------------------------------------------------------------------

// GetInstance returns current instance delegate attached to
func (it *ModelCustomAttributes) GetInstance() interface{} {
	return it.instance
}

// EditAttribute modifies custom attribute for collection
func (it *ModelCustomAttributes) EditAttribute(attributeName string, attributeValues models.StructAttributeInfo) error {
	modelCustomAttributesMutex.Lock()
	defer modelCustomAttributesMutex.Unlock()

	customAttributes, present := modelCustomAttributes[it.model]
	if !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ee8037e8-3728-44b8-a548-d074fa2afda3", "There is no attributes for model '"+it.model+"'")
	}

	customAttribute, present := customAttributes[attributeName]
	if !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5b3d8da7-fa80-4510-8d5c-0a32fef48f22", "There is no attribute '"+attributeName+"' for model '"+it.model+"'")
	}

	customAttributesCollection, err := db.GetCollection(ConstCollectionNameCustomAttributes)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "15fcfa45-4b48-417f-93b0-9d1ac63ffa12", "Can't get collection '"+ConstCollectionNameCustomAttributes+"': "+err.Error())
	}

	if err := customAttributesCollection.AddFilter("model", "=", customAttribute.Model); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dacfa3ff-597a-4560-acc5-d778e090e28e", "unable to add filter': "+err.Error())
	}
	if err := customAttributesCollection.AddFilter("attribute", "=", attributeName); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9aea6045-1040-4e6d-82c4-6a3ddb222ee8", "unable to add filter': "+err.Error())
	}
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

		modelCustomAttributes[it.model][attributeName] = customAttribute
	}

	return nil
}

// RemoveAttribute removes custom attribute from collection
func (it *ModelCustomAttributes) RemoveAttribute(attributeName string) error {

	modelCustomAttributesMutex.Lock()
	defer modelCustomAttributesMutex.Unlock()

	customAttributes, present := modelCustomAttributes[it.model]
	if !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f204f63d-0792-4e97-87a7-d612bc13d3b1", "There is no attributes for model '"+it.model+"'")
	}

	customAttribute, present := customAttributes[attributeName]
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

	delete(modelCustomAttributes, it.model)

	if err := customAttributesCollection.AddFilter("model", "=", customAttribute.Model); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "751f8e9e-9c2b-4e01-bed5-580fe46c1b3e", "unable to add filter': "+err.Error())
	}
	if err := customAttributesCollection.AddFilter("attribute", "=", attributeName); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3e771c53-c0ab-4695-8d30-77541c4b38d5", "unable to add filter': "+err.Error())
	}

	if _, err = customAttributesCollection.Delete(); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "da771c6b-402f-4816-a07d-9602f584b45d", "Can't remove attribute '"+attributeName+"' information from 'custom_attributes' collection '"+customAttribute.Collection+"': "+err.Error())
	}

	return nil
}

// AddNewAttribute extends collection with new custom attribute
func (it *ModelCustomAttributes) AddNewAttribute(newAttribute models.StructAttributeInfo) error {

	modelCustomAttributesMutex.Lock()
	defer modelCustomAttributesMutex.Unlock()

	customAttributes, present := modelCustomAttributes[it.model]
	if !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a24f7729-58ef-42b1-865b-7e42542c108e", "There is no attributes for model '"+it.model+"'")
	}

	if _, present := customAttributes[newAttribute.Attribute]; present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "24aa5125-d8b3-4e55-b321-79eef3eeccb8", "There is already atribute '"+newAttribute.Attribute+"' for model '"+it.model+"'")
	}

	// getting collection where custom attribute information stores
	customAttributesCollection, err := db.GetCollection(ConstCollectionNameCustomAttributes)
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

	// Type verification
	parsedType := utils.DataTypeParse(newAttribute.Type)
	if !parsedType.IsKnown {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "26f15efc-8490-41f0-82f7-b677ea427af7", "unknown attribute type '"+newAttribute.Type+"'")
	}

	// Assemble the cleaned type
	newAttribute.Type = parsedType.Name
	if parsedType.IsArray {
		newAttribute.Type = "[]" + newAttribute.Type
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

	newCustomAttributeID, err := customAttributesCollection.Save(record)

	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ad98e5b0-7672-4d02-9744-e1beecb88922", "Can't insert attribute '"+newAttribute.Attribute+"' in collection '"+newAttribute.Collection+"': "+err.Error())
	}

	// inserting new attribute to supposed location
	err = modelCollection.AddColumn(newAttribute.Attribute, newAttribute.Type, false)
	if err != nil {
		if err := customAttributesCollection.DeleteByID(newCustomAttributeID); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7730b177-ae12-4211-b4b7-d0c8a5e8584b", "Unable to delete new attribute: "+err.Error())
		}

		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0c11b43b-ec29-4e08-b614-7a8ec8345c9b", "Can't insert attribute '"+newAttribute.Attribute+"' in collection '"+newAttribute.Collection+"': "+err.Error())
	}

	customAttributes[newAttribute.Attribute] = newAttribute

	// Release lock to update records with default value.
	// Here mutex is still locked.
	modelCustomAttributesMutex.Unlock()

	// Populate default value.
	err = populateDefaultValue(newAttribute)

	// Regain lock. Unlock is still deferred.
	modelCustomAttributesMutex.Lock()

	if err != nil {
		if err := customAttributesCollection.DeleteByID(newCustomAttributeID); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7dbfd42b-4af8-4768-b25a-97f478a06a43", "Unable to delete new attribute: "+err.Error())
		}

		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fb47e17b-152f-4afb-9254-74ff1e0aa320", "Unable to populate new attribute's default value: "+err.Error())
	}

	return env.ErrorDispatch(err)
}

// GetCustomAttributeCollectionName returns collection name you can use to fill ModelCustomAttributes struct
func (it *ModelCustomAttributes) GetCustomAttributeCollectionName() string {
	return it.collection
}

// ----------------------------------------------------------------------------------
// InterfaceObject implementation (package "github.com/ottemo/foundation/app/models")
// ----------------------------------------------------------------------------------

// Get returns object attribute value or nil
func (it *ModelCustomAttributes) Get(attribute string) interface{} {
	return it.values[attribute]
}

// Set sets attribute value to object or returns error
func (it *ModelCustomAttributes) Set(attribute string, value interface{}) error {
	modelCustomAttributesMutex.Lock()
	defer modelCustomAttributesMutex.Unlock()

	customAttributes, present := modelCustomAttributes[it.model]
	if !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7972f82f-c8a5-43b6-a4da-a12b44cf7072", "There is no attributes for model '"+it.model+"'")
	}

	if _, present := customAttributes[attribute]; present {
		it.values[attribute] = value
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "154b03ed-0e75-416d-890b-8775fcd74063", "attribute '"+attribute+"' invalid")
	}

	return nil
}

// GetAttributesInfo represents object as map[string]interface{}
func (it *ModelCustomAttributes) GetAttributesInfo() []models.StructAttributeInfo {
	var result []models.StructAttributeInfo

	modelCustomAttributesMutex.Lock()
	defer modelCustomAttributesMutex.Unlock()

	if customAttributes, present := modelCustomAttributes[it.model]; present {
		for _, attribute := range customAttributes {
			result = append(result, attribute)
		}
	}

	return result
}

// FromHashMap represents object as map[string]interface{}
func (it *ModelCustomAttributes) FromHashMap(input map[string]interface{}) error {
	it.values = input
	return nil
}

// ToHashMap fills object attributes from map[string]interface{}
func (it *ModelCustomAttributes) ToHashMap() map[string]interface{} {
	return it.values
}

// ----------------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------------

// populateDefaultValue populate db records for preset attribute with default value
func populateDefaultValue(attributeInfo models.StructAttributeInfo) error {

	// check if model implements required interfaces
	attrModel, err := models.GetModel(attributeInfo.Model)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2896ab3c-983d-44c1-8a54-8a2413c65600", "Unable to get '"+attributeInfo.Model+"' model: "+err.Error())
	}

	if _, ok := attrModel.(models.InterfaceObject); !ok {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f1f38c27-d864-4025-8936-bfa755b3dcc4", "Model '"+attributeInfo.Model+"' does not implement Object: "+err.Error())
	}
	if _, ok := attrModel.(models.InterfaceStorable); !ok {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dd5d2ebc-91fe-4a2d-a401-51794cf07022", "Model '"+attributeInfo.Model+"' does not implement Storable: "+err.Error())
	}

	// load data from collection
	modelCollection, err := db.GetCollection(attributeInfo.Collection)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "962987f8-ccdc-473b-80de-617f9e61c235", "Unable to get '"+attributeInfo.Collection+"' collection: "+err.Error())
	}

	dbRecords, err := modelCollection.Load()
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8777df19-d2ef-496c-959f-f7c86ea48939", "Unable to load '"+attributeInfo.Collection+"' collection data: "+err.Error())
	}

	// update records
	for _, dbRecordData := range dbRecords {
		attrModel, err := models.GetModel(attributeInfo.Model)
		if err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "eefdb38a-0150-456d-b2a8-a10f1112d52b", "Unable to get '"+attributeInfo.Model+"' model: "+err.Error())
		}

		// type already checked
		modelObject := attrModel.(models.InterfaceObject)
		if err = modelObject.FromHashMap(dbRecordData); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "79119ca5-703d-4dd9-b142-247ea499c1c4", "Unable to populate '"+attributeInfo.Model+"' model: "+err.Error())
		}

		err = modelObject.Set(attributeInfo.Attribute, attributeInfo.Default)
		if err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c2c34a88-1409-4d01-b477-d1fe9449d664", "Unable to set value '"+attributeInfo.Default+"' for '"+attributeInfo.Attribute+"' in collection '"+attributeInfo.Collection+"': "+err.Error())
		}

		// type already checked
		err = modelObject.(models.InterfaceStorable).Save()
		if err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "51d6730d-0ff9-4a5a-bc83-8f0038d8c893", "Unable to save model '"+attributeInfo.Model+"': "+err.Error())
		}
	}

	return nil
}
