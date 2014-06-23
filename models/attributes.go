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

func (ca *CustomAttributes) Init(model string) (*CustomAttributes, error) {
	ca.model = model
	ca.values = make(map[string]interface{})

	_, present := globalCustomAttributes[model]

	if present {
		ca.attributes = globalCustomAttributes[model]
	} else {

		ca.attributes = make(map[string]models.AttributeInfo)

		dbEngine := database.GetDBEngine()
		if dbEngine == nil {
			return it, errors.New("There is no database engine")
		}

		caCollection, err := dbEngine.GetCollection(CustomAttributesCollection)
		if err != nil {
			return it, errors.New("Can't get collection 'custom_attributes': " + err.Error())
		}

		caCollection.AddFilter("model", "=", ca.model)
		dbValues, err := caCollection.Load()
		if err != nil {
			return ca, errors.New("Can't load custom attributes information for '" + ca.model + "'")
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

			ca.attributes[attribute.Attribute] = attribute
		}

		globalCustomAttributes[ca.model] = ca.attributes
	}

	return ca, nil
}

func (ca *CustomAttributes) RemoveAttribute(attributeName string) error {

	dbEngine := database.GetDBEngine()
	if dbEngine == nil {
		return errors.New("There is no database engine")
	}

	attribute, present := ca.attributes[attributeName]
	if !present {
		return errors.New("There is no attribute '" + attributeName + "' for model '" + ca.model + "'")
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

func (ca *CustomAttributes) AddNewAttribute(newAttribute models.AttributeInfo) error {

	if _, present := it.attributes[newAttribute.Attribute]; present {
		return errors.New("There is already atribute '" + newAttribute.Attribute + "' for model '" + ca.model + "'")
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

	ca.attributes[newAttribute.Attribute] = newAttribute

	return err
}

func (ca *CustomAttributes) FromHashMap(input map[string]interface{}) error {
	ca.values = input
	return nil
}

func (ca *CustomAttributes) ToHashMap() map[string]interface{} {
	return ca.values
}

func (ca *CustomAttributes) Get(attribute string) interface{} {
	return ca.values[attribute]
}

func (ca *CustomAttributes) Set(attribute string, value interface{}) error {
	if _, present := ca.attributes[attribute]; present {
		ca.values[attribute] = value
	} else {
		return errors.New("attribute '" + attribute + "' invalid")
	}

	return nil
}

func (ca *CustomAttributes) GetAttributesInfo() []models.AttributeInfo {
	info := make([]models.AttributeInfo, len(ca.attributes))
	for _, attribute := range ca.attributes {
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
