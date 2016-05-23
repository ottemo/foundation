package flatweight

import (
	"fmt"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func setupConfig() error {
	config := env.GetConfig()

	// Group Title
	config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGroup,
		Type:        env.ConstConfigTypeGroup,
		Label:       "Flat Weight",
		Description: "static amount stipping method",
	}, nil)

	// Enabled
	config.RegisterItem(env.StructConfigItem{
		Path:   ConstConfigPathEnabled,
		Value:  false,
		Type:   env.ConstConfigTypeBoolean,
		Editor: "boolean",
		Label:  "Enabled",
	}, nil)

	// Rates
	// demo json
	// [{"title": "Standard Shipping","code": "std_1","price": 1.99,"weight_from": 0.0,"weight_to": 5.0}]
	config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathRates,
		Value:       `[]`,
		Type:        env.ConstConfigTypeText,
		Editor:      "multiline_text",
		Label:       "Rates",
		Description: `Configuration format: [{"title": "Standard Shipping",  "code": "std_1", "price": 1.99,  "weight_from": 0.0, "weight_to": 5.0}]`,
	}, validateAndApplyRates)

	return nil
}

// validateAndApplyRates validate rates and convert to Rates type
func validateAndApplyRates(rawRates interface{}) (interface{}, error) {
	fmt.Println("validating raw rates", rawRates)

	// Allow empty
	rawRatesString := utils.InterfaceToString(rawRates)
	isEmptyString := rawRatesString == ""
	isEmptyArray := rawRatesString == "[]"
	isEmptyObj := rawRatesString == "{}"
	if isEmptyString || isEmptyArray || isEmptyObj {
		fmt.Println("found this to be empty", rawRates)
		rawRates = ""
		rates = make(Rates, 0)
		return rawRates, nil
	}

	parsedRates, err := utils.DecodeJSONToArray(rawRates)
	if err != nil {
		return nil, err
	}

	for _, rawRate := range parsedRates {
		parsedRate := utils.InterfaceToMap(rawRate)

		// Make sure we have our keys
		if !utils.KeysInMapAndNotBlank(parsedRate, "title", "code", "price", "weight_from", "weight_to") {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "todo", "Missing keys in config object; title, code, price, weight_from, weight_to")
		}

		// Assemble new rate
		rate := Rate{
			Title:      utils.InterfaceToString(parsedRate["title"]),
			Code:       utils.InterfaceToString(parsedRate["code"]),
			Price:      utils.InterfaceToFloat64(parsedRate["price"]),
			WeightFrom: utils.InterfaceToFloat64(parsedRate["weight_from"]),
			WeightTo:   utils.InterfaceToFloat64(parsedRate["weight_to"]),
		}

		// Update global variable
		rates = append(rates, rate)
	}

	return rawRates, nil
}

func configIsEnabled() bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEnabled))
}

func configRates() interface{} {
	return env.ConfigGetValue(ConstConfigPathRates)
}
