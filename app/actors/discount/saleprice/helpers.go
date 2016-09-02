package saleprice

import (
	"github.com/ottemo/foundation/app/models/discount/saleprice"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/db"
)

// Helper to produce new module level error
func newErrorHelper(msg, code string) error {
	return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, code, msg)
}

func logDebugHelper(msg string) {
	env.GetLogger().Log("errors.log", env.ConstLogPrefixDebug, msg)
}

func CreateSalePriceFromHashMapHelper(inputHashMap map[string]interface{}) (map[string]interface{}, error) {
	salePriceModel, err := saleprice.GetSalePriceModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = salePriceModel.FromHashMap(inputHashMap)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = salePriceModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return salePriceModel.ToHashMap(), nil
}

func UpdateSalePriceByHashMapHelper(id string, inputHashMap map[string]interface{}) (map[string]interface{}, error) {
	salePriceModel, err := saleprice.GetSalePriceModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = salePriceModel.Load(id)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	delete(inputHashMap, "id")
	err = salePriceModel.FromHashMap(inputHashMap)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = salePriceModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return salePriceModel.ToHashMap(), nil
}

func ReadSalePriceHashMapByIDHelper(id string) (map[string]interface{}, error) {
	salePriceModel, err := saleprice.GetSalePriceModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = salePriceModel.Load(id)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return salePriceModel.ToHashMap(), nil
}

func DeleteSalePriceByIDHelper(id string) error {
	salePriceModel, err := saleprice.GetSalePriceModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = salePriceModel.Load(id)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return salePriceModel.Delete()
}

func ReadSalePriceListHelper() ([]map[string]interface{}, error) {
	collection, err := db.GetCollection(ConstCollectionNameSalePrices)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return collection.Load()
}
