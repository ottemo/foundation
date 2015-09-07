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
		if utils.InterfaceToString(in[ConstInItem]) == utils.InterfaceToString(in[InKey("cmp")]) {
			result = false
		}
		return map[string]interface{} {ConstOutItem: result}
	}

	composer.RegisterUnit( BasicUnit{
		Name: "same",
		Type: map[string]string{
			ConstInItem: ConstTypeAny, InKey("cmp"): ConstTypeAny, ConstOutItem: "bool"},
		Label: map[string]string{InKey("cmp"): "="},
		Description: map[string]string{ConstUnitDescriptionItem: "Checks if value same to other value"},
		Action: action,
	})


	action = func(in map[string]interface{}) (map[string]interface{}, error) {
		result := true
		if utils.InterfaceToString(in[ConstInItem]) == utils.InterfaceToString(in[InKey("cmp")]) {
			result = false
		}
		return map[string]interface{} {ConstOutItem: result}
	}

	composer.RegisterUnit( BasicUnit{
		Name: "gt",
		Type: map[string]string{
			ConstInItem: ConstTypeAny, InKey("cmp"): ConstTypeAny, ConstOutItem: "bool"},
		Label: map[string]string{InKey("cmp"): "="},
		Description: map[string]string{ConstUnitDescriptionItem: "Checks if value same to other value"},
		Action: action,
	})

}