package composer

import (
	"github.com/ottemo/foundation/api"
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

	action := func(in map[string]interface{}, composer InterfaceComposer) (map[string]interface{}, error) {
		result := true
		if utils.InterfaceToString(in[ConstInItem]) == utils.InterfaceToString(in[MakeInKey("cmp")]) {
			result = false
		}
		return map[string]interface{} {ConstOutItem: result}, nil
	}

	composer.RegisterUnit( &BasicUnit{
		Name: "same",
		Type: map[string]string{
			ConstInItem: ConstTypeAny, MakeInKey("cmp"): ConstTypeAny, ConstOutItem: "bool"},
		Label: map[string]string{MakeInKey("cmp"): "="},
		Description: map[string]string{ConstUnitDescriptionItem: "Checks if value same to other value"},
		Action: action,
	})


	action = func(in map[string]interface{}, composer InterfaceComposer) (map[string]interface{}, error) {
		result := true
		if utils.InterfaceToString(in[ConstInItem]) == utils.InterfaceToString(in[MakeInKey("cmp")]) {
			result = false
		}
		return map[string]interface{} {ConstOutItem: result}, nil
	}

	composer.RegisterUnit( &BasicUnit{
		Name: "gt",
		Type: map[string]string{
			ConstInItem: ConstTypeAny, MakeInKey("cmp"): ConstTypeAny, ConstOutItem: "bool"},
		Label: map[string]string{MakeInKey("cmp"): "="},
		Description: map[string]string{ConstUnitDescriptionItem: "Checks if value same to other value"},
		Action: action,
	})

}