package models

import (
	"github.com/ottemo/foundation/env"
)

// RegisterModel registers new model to system
func RegisterModel(ModelName string, Model InterfaceModel) error {
	if _, present := declaredModels[ModelName]; present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0300eb6b08b8497eafd0eda0ee358596", "The model with name '"+ModelName+"' has already been registered")
	}
	declaredModels[ModelName] = Model

	return nil
}

// UnRegisterModel removes registered model from system
func UnRegisterModel(ModelName string) error {
	if _, present := declaredModels[ModelName]; present {
		delete(declaredModels, ModelName)
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3d651c0ae4a4443fa6e8de3a95d89b5c", "Unable to find model to delete with name '"+ModelName+"'")
	}
	return nil
}

// GetModel returns registered in system model
func GetModel(ModelName string) (InterfaceModel, error) {
	if model, present := declaredModels[ModelName]; present {
		return model.New()
	}
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5d49fd0d1fed47dc8e722346f1e778c3", "Unable to find model with name '"+ModelName+"'")
}

// GetDeclaredModels returns all currently registered in system models
func GetDeclaredModels() map[string]InterfaceModel {
	return declaredModels
}
