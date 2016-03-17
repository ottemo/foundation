package attributes
import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// GenericProduct type implements:
// 	- InterfaceExternalAttributes
// 	- InterfaceObject
// 	- InterfaceStorable

// Init initializes per instance helper before usage
// {instance} is a reference to object which using helper
func (it *ModelExternalAttributes) Init(instance interface{}) (*ModelExternalAttributes, error) {
	// making new ModelExternalAttributes struct for a given instance
	result := &ModelExternalAttributes{instance: instance}

	// checking the instance model
	modelName := ""
	instanceAsModel, ok := instance.(models.InterfaceModel)
	if instanceAsModel != nil {
		modelName = instanceAsModel.GetModelName()
	}

	if !ok || instanceAsModel == nil || modelName == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fe42f2db-2d4b-444a-9891-dc4632ad6dff", "Invalid instance")
	}

	// getting external attributes for instance model, if not set - making empty list
	modelExternalAttributesMutex.Lock()
	attributesInfo, present := modelExternalAttributes[modelName];
	if !present {
		attributesInfo = new(map[string]ModelExternalAttributes)
		modelExternalAttributes[modelName] = attributesInfo
	}
	modelExternalAttributesMutex.Unlock()

	// updating instantized ModelExternalAttributes struct
	result.model  = modelName
	result.info   = attributesInfo

	return result, nil
}


// ----------------------------------------------------------------------------------------------
// InterfaceExternalAttributes implementation (package "github.com/ottemo/foundation/app/models")
// ----------------------------------------------------------------------------------------------


// GetCurrentInstance returns current instance delegate attached to
func (it *ModelExternalAttributes) GetCurrentInstance() interface{} {
	return it.instance
}

// AddExternalAttribute registers new delegate for a given attribute
func (it *ModelExternalAttributes) AddExternalAttribute(newAttribute models.StructAttributeInfo, delegate interface{}) error {
	modelName := it.model
	attributeName := newAttribute.Attribute

	modelCustomAttributesMutex.Lock()
	defer modelCustomAttributesMutex.Unlock()

	attributesInfo, present := modelCustomAttributes[modelName]
	if !present {
		modelCustomAttributes[modelName] = new(map[string]models.StructAttributeInfo)
	}

	_, present = attributesInfo[attributeName]
	if !present {
		modelCustomAttributes[modelName] = newAttribute
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c2175996-b5f1-40dc-9ce2-9df133c3a2c4", "Attribute already exist")
	}

	return nil
}

// RemoveExternalAttribute registers new delegate for a given attribute
func (it *ModelExternalAttributes) RemoveExternalAttribute(attributeName string) error {
	modelName := it.model

	modelCustomAttributesMutex.Lock()
	defer modelCustomAttributesMutex.Unlock()

	attributesInfo, present := modelCustomAttributes[modelName]
	if !present {
		modelCustomAttributes[modelName] = new(map[string]models.StructAttributeInfo)
	}

	_, present = attributesInfo[attributeName]
	if present {
		delete(attributesInfo, attributeName)
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c2175996-b5f1-40dc-9ce2-9df133c3a2c4", "Attribute not exist")
	}

	return nil
}

// ListExternalAttributes registers new delegate for a given attribute
func (it *ModelExternalAttributes) ListExternalAttributes() []string {

}


// ----------------------------------------------------------------------------------
// InterfaceObject implementation (package "github.com/ottemo/foundation/app/models")
// ----------------------------------------------------------------------------------


// Get returns object attribute value or nil
func (it *ModelExternalAttributes) Get(attribute string) interface{} {
	if delegate, present := it.delegates[attribute]; present {
		if delegate, ok := delegate.(interface{ Get(string) interface{} }); ok {
			return delegate.Get(attribute)
		}
	}

	return nil
}

// Set sets attribute value to object or returns error
func (it *ModelExternalAttributes) Set(attribute string, value interface{}) error {
	if delegate, present := it.delegates[attribute]; present {
		if delegate, ok := delegate.(interface{ Set(string, interface{}) error }); ok {
			return delegate.Set(attribute, value)
		}
	}

	return nil
}

// GetAttributesInfo represents object as map[string]interface{}
func (it *ModelExternalAttributes) GetAttributesInfo() []models.StructAttributeInfo {
	var result []models.StructAttributeInfo

	it.mutex.Lock()
	for _, x := range it.info {
		result = append(result, x)
	}
	it.mutex.Unlock()

	return result
}

// FromHashMap represents object as map[string]interface{}
func (it *ModelExternalAttributes) FromHashMap(input map[string]interface{}) error {
	if delegate, present := it.delegates[attribute]; present {
		if delegate, ok := delegate.(models.InterfaceObject); ok {
			return delegate.FromHashMap()
		}
	}

	return nil
}

// ToHashMap fills object attributes from map[string]interface{}
func (it *ModelExternalAttributes) ToHashMap() map[string]interface{} {
	return it.values
}
