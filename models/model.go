package models

import (
	"errors"
)

type Model interface {
	GetModelName() string
	GetImplementationName() string
	New() (Model, error)
}

type Storable interface {
	GetId() string
	SetId(string) error

	Save() error
	Load(id string) error
	Delete(id string) error
}

type Object interface {
	Get(Attribute string) interface{}
	Set(Attribute string, Value interface{}) error

	GetAttributesInfo() []AttributeInfo
}

type Mapable interface {
	FromHashMap(HashMap map[string]interface{}) error
	ToHashMap() map[string]interface{}
}

type Attribute interface {
	AddNewAttribute(newAttribute AttributeInfo) error
	RemoveAttribute(attributeName string) error
}

type AttributeInfo struct {
	Model      string
	Collection string
	Attribute  string
	Type       string
	Label      string
	Group      string
	Editors    string
	Options    string
	Default    string
}

var declaredModels = map[string]Model{}

func RegisterModel(ModelName string, Model Model) error {
	if _, present := declaredModels[ModelName]; present {
		return errors.New("model with name '" + ModelName + "' already registered")
	} else {
		declaredModels[ModelName] = Model
	}
	return nil
}

func UnRegisterBlock(ModelName string) error {
	if _, present := declaredModels[ModelName]; present {
		delete(declaredModels, ModelName)
	} else {
		return errors.New("can't find module with name '" + ModelName + "'")
	}
	return nil
}

func GetModel(ModelName string) (Model, error) {
	if model, present := declaredModels[ModelName]; present {
		return model.New()
	} else {
		return nil, errors.New("can't find module with name '" + ModelName + "'")
	}
}
