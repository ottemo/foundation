package composer

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/composer"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultComposer)
	instance.units = make(map[string]InterfaceComposeUnit)

	composer = instance

	api.RegisterOnRestServiceStart(setupAPI)
	initBaseUnits()
}


func initBaseUnits() {

	action := func(in map[string]interface{}) (map[string]interface{}, error) {
		result := true
		if utils.InterfaceToString(in[ConstInKey]) == utils.InterfaceToString(in[InKey("cmp")]) {
			result = false
		}
		return map[string]interface{} {ConstOutKey: result}
	}

	composer.RegisterUnit( BasicUnit{
		Name: "same",
		Type: map[string]string{
			ConstInKey: ConstTypeAny, InKey("cmp"): ConstTypeAny, ConstOutKey: "bool"},
		Label: map[string]string{InKey("cmp"): "="},
		Description: map[string]string{ConstUnitDescriptionKey: "Checks if value same to other value"},
		Action: action,
	})


	action = func(in map[string]interface{}) (map[string]interface{}, error) {
		result := true
		if utils.InterfaceToString(in[ConstInKey]) == utils.InterfaceToString(in[InKey("cmp")]) {
			result = false
		}
		return map[string]interface{} {ConstOutKey: result}
	}

	composer.RegisterUnit( BasicUnit{
		Name: "gt",
		Type: map[string]string{
			ConstInKey: ConstTypeAny, InKey("cmp"): ConstTypeAny, ConstOutKey: "bool"},
		Label: map[string]string{InKey("cmp"): "="},
		Description: map[string]string{ConstUnitDescriptionKey: "Checks if value same to other value"},
		Action: action,
	})

}