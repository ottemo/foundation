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
	if _, present := modelExternalAttributes[modelName]; !present {
		modelExternalAttributes[modelName] = make(map[string]models.StructAttributeInfo)
	}
	modelExternalAttributesMutex.Unlock()

	// updating instantized ModelExternalAttributes struct
	result.model  = modelName

	return result, nil
}

// ----------------------------------------------------------------------------------------------
// InterfaceExternalAttributes implementation (package "github.com/ottemo/foundation/app/models")
// ----------------------------------------------------------------------------------------------


// GetCurrentInstance returns current instance delegate attached to
func (it *ModelExternalAttributes) GetInstance() interface{} {
	return it.instance
}

// AddExternalAttribute registers new delegate for a given attribute
func (it *ModelExternalAttributes) AddExternalAttribute(newAttribute models.StructAttributeInfo, delegate interface{}) error {
	modelName := it.model
	attributeName := newAttribute.Attribute

	modelExternalAttributesMutex.Lock()
	defer modelExternalAttributesMutex.Unlock()

	attributesInfo, present := modelCustomAttributes[modelName]
	if !present {
		modelCustomAttributes[modelName] = make(map[string]models.StructAttributeInfo)
	}

	_, present = attributesInfo[attributeName]
	if !present {
		modelCustomAttributes[modelName][attributeName] = newAttribute
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c2175996-b5f1-40dc-9ce2-9df133c3a2c4", "Attribute already exist")
	}

	return nil
}

// RemoveExternalAttribute registers new delegate for a given attribute
func (it *ModelExternalAttributes) RemoveExternalAttribute(attributeName string) error {
	modelName := it.model

	modelExternalAttributesMutex.Lock()
	defer modelExternalAttributesMutex.Unlock()

	attributesInfo, present := modelCustomAttributes[modelName]
	if !present {
		modelCustomAttributes[modelName] = make(map[string]models.StructAttributeInfo)
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
	var result []string
	modelName := it.model

	modelExternalAttributesMutex.Lock()
	defer modelExternalAttributesMutex.Unlock()

	attributesInfo, present := modelCustomAttributes[modelName]
	if !present {
		modelCustomAttributes[modelName] = make(map[string]models.StructAttributeInfo)
	}

	for name := range attributesInfo {
		result = append(result, name)
	}

	return result
}


// ----------------------------------------------------------------------------------
// InterfaceObject implementation (package "github.com/ottemo/foundation/app/models")
// ----------------------------------------------------------------------------------


// Get returns object attribute value or nil
func (it *ModelExternalAttributes) Get(attribute string) interface{} {
	modelExternalDelegatesMutex.Lock()
	defer modelExternalDelegatesMutex.Unlock()

	if delegates, present := modelExternalDelegates[it.model]; present {
		if delegate, present := delegates[attribute]; present {
			if delegate, ok := delegate.(interface{ Get(string) interface{} }); ok {
				return delegate.Get(attribute)
			}
		}
	}

	return nil
}

// Set sets attribute value to object or returns error
func (it *ModelExternalAttributes) Set(attribute string, value interface{}) error {
	modelExternalDelegatesMutex.Lock()
	defer modelExternalDelegatesMutex.Unlock()

	if delegates, present := modelExternalDelegates[it.model]; present {
		if delegate, present := delegates[attribute]; present {
			if delegate, ok := delegate.(interface{ Set(string, interface{}) error }); ok {
				return delegate.Set(attribute, value)
			}
		}
	}

	return nil
}

// GetAttributesInfo represents object as map[string]interface{}
func (it *ModelExternalAttributes) GetAttributesInfo() []models.StructAttributeInfo {
	var result []models.StructAttributeInfo

	modelExternalAttributesMutex.Lock()
	defer modelExternalAttributesMutex.Unlock()

	if attributesInfo, present := modelExternalAttributes[it.model]; !present {
		for _, info := range attributesInfo {
			result = append(result, info)
		}
	}

	return result
}

// FromHashMap represents object as map[string]interface{}
func (it *ModelExternalAttributes) FromHashMap(input map[string]interface{}) error {
	modelExternalDelegatesMutex.Lock()
	delegates, present := modelExternalDelegates[it.model]
	modelExternalDelegatesMutex.Unlock()

	if present {
		for attribute, delegate := range delegates {
			if value, present := input[attribute]; present {
				if delegate, ok := delegate.(interface{ Set(string, interface{}) error }); ok {
					err := delegate.Set(attribute, value)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// ToHashMap fills object attributes from map[string]interface{}
func (it *ModelExternalAttributes) ToHashMap() map[string]interface{} {
	modelExternalDelegatesMutex.Lock()
	delegates, present := modelExternalDelegates[it.model]
	modelExternalDelegatesMutex.Unlock()

	result := make(map[string]interface{})
	if present {
		for attribute, delegate := range delegates {
			if delegate, ok := delegate.(interface{ Get(string) interface{} }); ok {
				result[attribute] = delegate.Get(attribute)
			}
		}
	}
	return result
}

// ------------------------------------------------------------------------------------
// InterfaceStorable implementation (package "github.com/ottemo/foundation/app/models")
// ------------------------------------------------------------------------------------


// GetID delegates call back to instance (stub method)
func (it *ModelExternalAttributes) GetID() string {
	if instance, ok := it.instance.(interface{ GetID() string }); ok {
		return instance.GetID()
	}
	return ""
}

// SetID callbacks all external attribute delegates
func (it *ModelExternalAttributes) SetID(id string) error {
	modelExternalDelegatesMutex.Lock()
	delegates, present := modelExternalDelegates[it.model]
	modelExternalDelegatesMutex.Unlock()

	if present {
		for _, delegate := range delegates {
			if delegate, ok := delegate.(interface{ SetID(newID string) error }); ok {
				if err := delegate.SetID(id); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Load callbacks all external attribute delegates
func (it *ModelExternalAttributes) Load(id string) error {
	modelExternalDelegatesMutex.Lock()
	delegates, present := modelExternalDelegates[it.model]
	modelExternalDelegatesMutex.Unlock()

	if present {
		for _, delegate := range delegates {
			if delegate, ok := delegate.(interface{ Load(loadID string) error }); ok {
				if err := delegate.Load(id); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Delete callbacks all external attribute delegates
func (it *ModelExternalAttributes) Delete() error {
	modelExternalDelegatesMutex.Lock()
	delegates, present := modelExternalDelegates[it.model]
	modelExternalDelegatesMutex.Unlock()

	if present {
		for _, delegate := range delegates {
			if delegate, ok := delegate.(interface{ Delete() error }); ok {
				if err := delegate.Delete(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Save callbacks all external attribute delegates
func (it *ModelExternalAttributes) Save() error {
	modelExternalDelegatesMutex.Lock()
	delegates, present := modelExternalDelegates[it.model]
	modelExternalDelegatesMutex.Unlock()

	if present {
		for _, delegate := range delegates {
			if delegate, ok := delegate.(interface{ Save() error }); ok {
				if err := delegate.Save(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

