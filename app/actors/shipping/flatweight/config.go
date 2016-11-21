package flatweight

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "8a54aa93-8f6a-4e92-8398-d6e6d05ee2af", "can't obtain config")
		return env.ErrorDispatch(err)
	}

	// Group Title
	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGroup,
		Type:        env.ConstConfigTypeGroup,
		Label:       "Flat Weight",
		Description: "static amount stipping method",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Enabled
	err = config.RegisterItem(env.StructConfigItem{
		Path:   ConstConfigPathEnabled,
		Value:  false,
		Type:   env.ConstConfigTypeBoolean,
		Editor: "boolean",
		Label:  "Enabled",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Rates
	// demo json [{"title": "Standard Shipping","code": "std_1","price": 1.99,"weight_from": 0.0,"weight_to": 5.0}]
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathRates,
		Value:       `[]`,
		Type:        env.ConstConfigTypeText,
		Editor:      "multiline_text",
		Label:       "Rates",
		Description: `Configuration format: [{"title": "Standard Shipping",  "code": "std_1", "price": 1.99,  "weight_from": 0.0, "weight_to": 5.0}]`,
	}, validateAndApplyRates)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// validateAndApplyRates validate rates and convert to Rates type
func validateAndApplyRates(rawRates interface{}) (interface{}, error) {

	// Allow empty
	rawRatesString := utils.InterfaceToString(rawRates)
	isEmptyString := rawRatesString == ""
	isEmptyArray := rawRatesString == "[]"
	isEmptyObj := rawRatesString == "{}"
	if isEmptyString || isEmptyArray || isEmptyObj {
		// Reset our global variable
		rates = make(Rates, 0)
		rawRates = ""
		return rawRates, nil
	}

	parsedRates, err := utils.DecodeJSONToArray(rawRates)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Validate each new rate
	validRates := make(Rates, 0)
	for _, rawRate := range parsedRates {
		parsedRate := utils.InterfaceToMap(rawRate)

		// Make sure we have our keys
		if !utils.KeysInMapAndNotBlank(parsedRate, "title", "code", "price", "weight_from", "weight_to") {
			err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "todo", "Missing keys in config object; title, code, price, weight_from, weight_to")
			return nil, env.ErrorDispatch(err)
		}

		// Assemble new rate
		rate := Rate{
			Title:            utils.InterfaceToString(parsedRate["title"]),
			Code:             utils.InterfaceToString(parsedRate["code"]),
			Price:            utils.InterfaceToFloat64(parsedRate["price"]),
			WeightFrom:       utils.InterfaceToFloat64(parsedRate["weight_from"]),
			WeightTo:         utils.InterfaceToFloat64(parsedRate["weight_to"]),
			AllowedCountries: utils.InterfaceToString(parsedRate["allowed_countries"]),
			BannedCountries:  utils.InterfaceToString(parsedRate["banned_countries"]),
		}

		validRates = append(validRates, rate)
	}

	// We didn't hit any validation errors, update our global var
	rates = validRates

	return rawRates, nil
}

func configIsEnabled() bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEnabled))
}

func configRates() interface{} {
	return env.ConfigGetValue(ConstConfigPathRates)
}
