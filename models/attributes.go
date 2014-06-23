package models

import (
	"errors"

	"github.com/ottemo/foundation/database"
	"github.com/ottemo/foundation/models"
)

const (
	CustomAttributesCollection = "custom_attributes"
)

var globalCustomAttributes = map[string]map[string]models.AttributeInfo{}

type CustomAttributes struct {
	model      string
	attributes map[string]models.AttributeInfo

	values map[string]interface{}
}

func (it *CustomAttributes) Init(model string) (*CustomAttributes, error) {
	it.model = model
	it.values = make(map[string]interface{})

	_, present := globalCustomAttributes[model]

	if present {
		it.attributes = globalCustomAttributes[model]
	} else {

		it.attributes = make(map[string]models.AttributeInfo)

		dbEngine := database.GetDBEngine()
		if dbEngine == nil {
			return it, errors.New("There is no database engine")
		}

		caCollection, err := dbEngine.GetCollection(CustomAttributesCollection)
		if err != nil {
			return it, errors.New("Can't get collection 'custom_attributes': " + err.Error())
		}

		caCollection.AddFilter("model", "=", it.model)
		dbValues, err := caCollection.Load()
		if err != nil {
			return it, errors.New("Can't load custom attributes information for '" + it.model + "'")
		}

		for _, row := range dbValues {
			attribute := models.AttributeInfo{
				Model:      row["model"].(string),
				Collection: row["collection"].(string),
				Attribute:  row["attribute"].(string),
				Type:       row["type"].(string),
				Label:      row["label"].(string),
				Group:      row["group"].(string),
				Editors:    row["editors"].(string),
				Options:    row["options"].(string),
				Default:    row["default"].(string),
			}

			it.attributes[attribute.Attribute] = attribute
		}

		globalCustomAttributes[it.model] = it.attributes
	}

	return it, nil
}

func (it *CustomAttributes) RemoveAttribute(attributeName string) error {

	dbEngine := database.GetDBEngine()
	if dbEngine == nil {
		return errors.New("There is no database engine")
	}

	attribute, present := it.attributes[attributeName]
	if !present {
		return errors.New("There is no attribute '" + attributeName + "' for model '" + it.model + "'")
	}

	caCollection, err := dbEngine.GetCollection(CustomAttributesCollection)
	if err != nil {
		return errors.New("Can't get collection 'custom_attributes': " + err.Error())
	}

	attrCollection, err := dbEngine.GetCollection(attribute.Collection)
	if err != nil {
		return errors.New("Can't get attribute '" + attribute.Attribute + "' collection '" + attribute.Collection + "': " + err.Error())
	}

	err = attrCollection.RemoveColumn(attributeName)
	if err != nil {
		return errors.New("Can't remove attribute '" + attributeName + "' from collection '" + attribute.Collection + "': " + err.Error())
	}

	caCollection.AddFilter("model", "=", attribute.Collection)
	caCollection.AddFilter("attr", "=", attributeName)
	_, err = caCollection.Delete()
	if err != nil {
		return errors.New("Can't remove attribute '" + attributeName + "' information from 'custom_attributes' collection '" + attribute.Collection + "': " + err.Error())
	}

	return nil
}

func (it *CustomAttributes) AddNewAttribute(newAttribute models.T_AttributeInfo) error {

	if _, present := it.attributes[newAttribute.Attribute]; present {
		return errors.New("There is already atribute '" + newAttribute.Attribute + "' for model '" + it.model + "'")
	}

	dbEngine := database.GetDBEngine()
	if dbEngine == nil {
		return errors.New("There is no database engine")
	}

	// getting collection where custom attribute information stores
	caCollection, err := dbEngine.GetCollection(CustomAttributesCollection)
	if err != nil {
		return errors.New("Can't get collection 'custom_attributes': " + err.Error())
	}

	// getting collection where attribute supposed to be
	attrCollection, err := dbEngine.GetCollection(newAttribute.Collection)
	if err != nil {
		return errors.New("Can't get attribute '" + newAttribute.Attribute + "' collection '" + newAttribute.Collection + "': " + err.Error())
	}

	// inserting attribute information in custom_attributes collection
	hashMap := make(map[string]interface{})

	hashMap["model"] = newAttribute.Model
	hashMap["collection"] = newAttribute.Collection
	hashMap["attribute"] = newAttribute.Attribute
	hashMap["type"] = newAttribute.Type
	hashMap["label"] = newAttribute.Label
	hashMap["group"] = newAttribute.Group
	hashMap["editors"] = newAttribute.Editors
	hashMap["options"] = newAttribute.Options
	hashMap["default"] = newAttribute.Default

	newCustomAttributeId, err := caCollection.Save(hashMap)

	if err != nil {
		return errors.New("Can't insert attribute '" + newAttribute.Attribute + "' in collection '" + newAttribute.Collection + "': " + err.Error())
	}

	// inserting new attribute to supposed location
	err = attrCollection.AddColumn(newAttribute.Attribute, newAttribute.Type, false)
	if err != nil {
		caCollection.DeleteById(newCustomAttributeId)

		return errors.New("Can't insert attribute '" + newAttribute.Attribute + "' in collection '" + newAttribute.Collection + "': " + err.Error())
	}

	it.attributes[newAttribute.Attribute] = newAttribute

	return err
}

func (it *CustomAttributes) FromHashMap(input map[string]interface{}) error {
	it.values = input
	return nil
}

func (it *CustomAttributes) ToHashMap() map[string]interface{} {
	return it.values
}

func (it *CustomAttributes) Get(attribute string) interface{} {
	return it.values[attribute]
}

func (it *CustomAttributes) Set(attribute string, value interface{}) error {
	if _, present := it.attributes[attribute]; present {
		it.values[attribute] = value
	} else {
		return errors.New("attribute '" + attribute + "' invalid")
	}

	return nil
}

func (it *CustomAttributes) GetAttributesInfo() []models.AttributeInfo {
	info := make([]models.AttributeInfo, len(it.attributes))
	for _, attribute := range it.attributes {
		info = append(info, attribute)
	}
	return info
}

func init() {
	database.RegisterOnDatabaseStart(SetupModel)
}

func SetupModel() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("custom_attributes"); err == nil {
			collection.AddColumn("model", "text", true)
			collection.AddColumn("collection", "text", true)
			collection.AddColumn("attribute", "text", true)
			collection.AddColumn("type", "text", false)
			collection.AddColumn("label", "text", true)
			collection.AddColumn("group", "text", false)
			collection.AddColumn("editors", "text", false)
			collection.AddColumn("options", "text", false)
			collection.AddColumn("default", "text", false)

		} else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}
