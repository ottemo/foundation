package usps

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	if config := env.GetConfig(); config != nil {
		err := config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathGroup,
			Value:       nil,
			Type:        env.ConstConfigTypeGroup,
			Editor:      "",
			Options:     nil,
			Label:       "USPS",
			Description: "USPS - The United States Postal Service",
			Image:       "https://www.usps.com/ContentTemplates/assets/images/global/usps_logo.gif",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathEnabled,
			Value:       false,
			Type:        env.ConstConfigTypeBoolean,
			Editor:      "boolean",
			Options:     nil,
			Label:       "Enabled",
			Description: "enables/disables shipping method for storefront",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathUser,
			Value:       nil,
			Type:        env.ConstConfigTypeText,
			Editor:      "line_editor",
			Options:     nil,
			Label:       "USERID",
			Description: "Account ID to use in USPS API",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathOriginZip,
			Value:       nil,
			Type:        env.ConstConfigTypeText,
			Editor:      "line_editor",
			Options:     nil,
			Label:       "Origin zip",
			Description: "shipping origin zip code - needed for price calculation",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathContainer,
			Value:       "VARIABLE",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "select",
			Options:     "{'VARIABLE': 'Variable', 'RECTANGULAR': 'Rectangular', 'NONRECTANGULAR': 'Not rectangular'}",
			Label:       "Container",
			Description: "Container type will be sent to USPS service",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathSize,
			Value:       "REGULAR",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "select",
			Options:     "{'REGULAR': 'Regular', 'LARGE': 'Large'}",
			Label:       "Size",
			Description: "Package size will be sent to USPS service. Regular - dimensions are 12’’ or less; Large - dimensions must be specified in manually",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDefaultDimensions,
			Value:       "1.0 x 1.0 x 1.0 x 1.0",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "dimensions",
			Options:     nil,
			Label:       "Default dimensions",
			Description: "Will be used if product dimension are not specified (width x length x height x girth - in inches)",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDefaultWeight,
			Value:       0.1,
			Type:        env.ConstConfigTypeDecimal,
			Editor:      "decimal",
			Options:     nil,
			Label:       "Default weight",
			Description: "Will be used if product do not have this value (in pounds)",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathAllowedMethods,
			Value:       "",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "multi_select",
			Options:     ConstShippingMethods,
			Label:       "Allowed methods",
			Description: "To customer will be proposed only allowed methods",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDebugLog,
			Value:       false,
			Type:        env.ConstConfigTypeBoolean,
			Editor:      "boolean",
			Options:     nil,
			Label:       "Debug log",
			Description: "enables/disables shipping method debug log",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return env.ErrorDispatch(err)
		}
	} else {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "91775c4b-521a-49f3-90bf-c377f266b62e", "Unable to obtain configuration for USPS")
		return env.ErrorDispatch(err)
	}

	return nil
}
