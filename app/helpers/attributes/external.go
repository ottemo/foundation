package attributes

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// perDelegateAttributes returns new map where the attributes are grouped by delegate
// (as one delegate could serve couple attributes we should instantiate and notify on that delegate only once
// that grouping by delegate helps to solve this task)
func groupByDelegate(input map[string]models.InterfaceAttributesDelegate) map[models.InterfaceAttributesDelegate][]string {
	perDelegateGroup := make(map[models.InterfaceAttributesDelegate][]string)

	for attribute, delegate := range input {
		if _, present := perDelegateGroup[delegate]; present {
			perDelegateGroup[delegate] = append(perDelegateGroup[delegate], attribute)
		} else {
			perDelegateGroup[delegate] = []string{attribute}
		}
	}

	return perDelegateGroup
}

// ExternalAttributes type implements:
// 	- InterfaceExternalAttributes
// 	- InterfaceObject
// 	- InterfaceStorable

// ExternalAttributes initializes helper for model instance - should be called in model.New()
// 	- "instance" is a reference to object which using helper
func ExternalAttributes(instance interface{}) (*ModelExternalAttributes, error) {
	newInstance := &ModelExternalAttributes{instance: instance}

	// getting model name from given instance
	modelName := ""
	instanceAsModel, ok := instance.(models.InterfaceModel)
	if !ok || instanceAsModel == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fe42f2db-2d4b-444a-9891-dc4632ad6dff", "Invalid instance")
	}
	modelName = instanceAsModel.GetModelName()

	if modelName == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fe42f2db-2d4b-444a-9891-dc4632ad6dff", "Invalid instance")
	}
	newInstance.model = modelName
	newInstance.delegates = make(map[string]models.InterfaceAttributesDelegate)

	modelExternalAttributesMutex.Lock()
	perDelegateGroup := groupByDelegate(modelExternalAttributes[modelName])
	modelExternalAttributesMutex.Unlock()

	// instantiating delegates for instance
	for delegate, attributes := range perDelegateGroup {
		delegateInstance, err := delegate.New(instance)
		if err != nil {
			env.ErrorDispatch(err)
		}

		for _, attribute := range attributes {
			newInstance.delegates[attribute] = delegateInstance
		}
	}

	return newInstance, nil
}

// ----------------------------------------------------------------------------------------------
// InterfaceExternalAttributes implementation (package "github.com/ottemo/foundation/app/models")
// ----------------------------------------------------------------------------------------------

// GetInstance returns current instance delegate attached to
func (it *ModelExternalAttributes) GetInstance() interface{} {
	return it.instance
}

// AddExternalAttributes registers new delegate for a it's attributes - delegate.GetAttributesInfo()
func (it *ModelExternalAttributes) AddExternalAttributes(delegate models.InterfaceAttributesDelegate) error {
	modelName := it.model
	if modelName == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fe42f2db-2d4b-444a-9891-dc4632ad6dff", "Invalid instance")
	}

	if delegate == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "203a7c29-7792-404d-baa0-ec9f56e226d2", "Invalid delegate")
	}

	for _, attributeInfo := range delegate.GetAttributesInfo() {
		attributeName := attributeInfo.Attribute
		if attributeName == "" {
			continue
		}

		modelExternalAttributesMutex.Lock()
		if _, present := modelExternalAttributes[modelName]; !present {
			modelExternalAttributes[modelName] = make(map[string]models.InterfaceAttributesDelegate)
		}

		if _, present := modelExternalAttributes[modelName][attributeName]; present {
			modelExternalAttributes[modelName][attributeName] = delegate
		} else {
			modelExternalAttributes[modelName][attributeName] = delegate
		}
		modelExternalAttributesMutex.Unlock()
	}

	// updating current instance
	newInstance, err := ExternalAttributes(it.instance)
	if err != nil {
		return err
	}
	it = newInstance

	return nil
}

// RemoveExternalAttributes removes the delegate for it's attributes - delegate.GetAttributesInfo()
func (it *ModelExternalAttributes) RemoveExternalAttributes(delegate models.InterfaceAttributesDelegate) error {
	modelName := it.model
	if modelName == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fe42f2db-2d4b-444a-9891-dc4632ad6dff", "Invalid instance")
	}

	if delegate == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "203a7c29-7792-404d-baa0-ec9f56e226d2", "Invalid delegate")
	}

	for _, attributeInfo := range delegate.GetAttributesInfo() {
		attributeName := attributeInfo.Attribute
		if attributeName == "" {
			continue
		}

		modelExternalAttributesMutex.Lock()
		oldDelegate, present := modelExternalAttributes[modelName][attributeName]
		if present && oldDelegate == delegate {
			delete(modelExternalAttributes[modelName], attributeName)
		}
		modelExternalAttributesMutex.Unlock()
	}

	// updating current instance
	newInstance, err := ExternalAttributes(it.instance)
	if err != nil {
		return err
	}
	it = newInstance
	return nil
}

// ListExternalAttributes returns delegate per attribute mapping
func (it *ModelExternalAttributes) ListExternalAttributes() map[string]models.InterfaceAttributesDelegate {
	result := make(map[string]models.InterfaceAttributesDelegate)

	modelName := it.model
	if modelName == "" {
		env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fe42f2db-2d4b-444a-9891-dc4632ad6dff", "Invalid instance")
		return result
	}

	modelExternalAttributesMutex.Lock()
	for attribute, delegate := range modelExternalAttributes[modelName] {
		result[attribute] = delegate
	}
	modelExternalAttributesMutex.Unlock()

	return result
}

// ----------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models")
// ----------------------------------------------------------------------------------

// GetModelName stub func for Model interface - returns model name
func (it *ModelExternalAttributes) GetModelName() string {
	return it.model
}

// GetImplementationName stub func for Model interface - doing callback to model instance function if possible
func (it *ModelExternalAttributes) GetImplementationName() string {
	if instanceAsModel, ok := it.instance.(models.InterfaceModel); ok {
		return instanceAsModel.GetImplementationName()
	}
	return ""
}

// ----------------------------------------------------------------------------------
// InterfaceObject implementation (package "github.com/ottemo/foundation/app/models")
// ----------------------------------------------------------------------------------

// Get returns object attribute value or nil
func (it *ModelExternalAttributes) Get(attribute string) interface{} {
	if delegate, present := it.delegates[attribute]; present {
		return delegate.Get(attribute)
	}
	return nil
}

// Set sets attribute value to object or returns error
func (it *ModelExternalAttributes) Set(attribute string, value interface{}) error {
	if delegate, present := it.delegates[attribute]; present {
		return delegate.Set(attribute, value)
	}
	return nil
}

// GetAttributesInfo represents object as map[string]interface{}
func (it *ModelExternalAttributes) GetAttributesInfo() []models.StructAttributeInfo {
	var result []models.StructAttributeInfo

	for attribute, delegate := range it.delegates {
		for _, info := range delegate.GetAttributesInfo() {
			if info.Attribute == attribute {
				result = append(result, info)
				break
			}
		}
	}

	return result
}

// FromHashMap updates delegated attributes from given map
func (it *ModelExternalAttributes) FromHashMap(input map[string]interface{}) error {
	for attribute, delegate := range it.delegates {
		if value, present := input[attribute]; present {
			err := delegate.Set(attribute, value)
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}
	return nil
}

// ToHashMap returns delegated attributes in map
func (it *ModelExternalAttributes) ToHashMap() map[string]interface{} {
	result := make(map[string]interface{})
	for attribute, delegate := range it.delegates {
		if delegate, ok := delegate.(interface {
			Get(string) interface{}
		}); ok {
			result[attribute] = delegate.Get(attribute)
		}
	}
	return result
}

// ------------------------------------------------------------------------------------
// InterfaceStorable implementation (package "github.com/ottemo/foundation/app/models")
// ------------------------------------------------------------------------------------

// GetID stub function - callback to instance getID()
func (it *ModelExternalAttributes) GetID() string {
	if instance, ok := it.instance.(interface {
		GetID() string
	}); ok {
		return instance.GetID()
	}
	return ""
}

// SetID proxies method to external attribute delegates
func (it *ModelExternalAttributes) SetID(id string) error {
	for delegate := range groupByDelegate(it.delegates) {
		if delegate, ok := delegate.(interface {
			SetID(newID string) error
		}); ok {
			if err := delegate.SetID(id); err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}
	return nil
}

// Load proxies method to external attribute delegates
func (it *ModelExternalAttributes) Load(id string) error {
	for delegate := range groupByDelegate(it.delegates) {
		if delegate, ok := delegate.(interface {
			Load(loadID string) error
		}); ok {
			if err := delegate.Load(id); err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}
	return nil
}

// Delete proxies method to external attribute delegates
func (it *ModelExternalAttributes) Delete() error {
	for delegate := range groupByDelegate(it.delegates) {
		if delegate, ok := delegate.(interface {
			Delete() error
		}); ok {
			if err := delegate.Delete(); err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}
	return nil
}

// Save proxies method to external attribute delegates
func (it *ModelExternalAttributes) Save() error {
	for delegate := range groupByDelegate(it.delegates) {
		if delegate, ok := delegate.(interface {
			Save() error
		}); ok {
			if err := delegate.Save(); err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}
	return nil
}
