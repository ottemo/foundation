package composer

import (
	"github.com/ottemo/foundation/env"
)

func (it *DefaultComposer) RegisterUnit(unit InterfaceComposeUnit) error {
	unitName := unit.GetName()

	if _, present := it.units[unitName]; !present {
		it.units[unitName] = unit
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "57d471d0-1fe0-40bf-999d-96ef803f62fa", "unit already registered")
	}

	return nil
}

func (it *DefaultComposer) UnRegisterUnit(unit InterfaceComposeUnit) error {
	unitName := unit.GetName()

	if _, present := it.units[unitName]; present {
		delete(it.units, unitName)
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3537c93c-4f22-466a-8c76-da47373a26ba", "unit not exists")
	}

	return nil
}

func (it *DefaultComposer) ListUnits() []InterfaceComposeUnit {

}

func (it *DefaultComposer) Process(unit InterfaceComposeUnit, in map[string]interface{}) (interface{}, error) {
	return nil, nil
}