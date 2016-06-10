package app

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"sort"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/api/rest"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "f635e96c-3cd7-4ae2-a507-4349021e9f13", "can't obtain config")
		return env.ErrorDispatch(err)
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "General",
		Description: "application general options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathAppGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Application",
		Description: "application related options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStorefrontURL,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Storefront host URL",
		Description: "URL application will use to generate storefront resources links",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathDashboardURL,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Dashboard host URL",
		Description: "URL application will use to generate dashboard resources links",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathFoundationURL,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Foundation host URL",
		Description: "URL application will use to generate foundation resources links",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Store",
		Description: "web store related options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreName,
		Value:       "Ottemo store",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Name",
		Description: "name of your web store",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreEmail,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "E-mail",
		Description: "e-mail of your web store",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreRootLogin,
		Value:       "admin",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Root login",
		Description: "login to enter admin panel as root",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreRootPassword,
		Value:       "admin",
		Type:        env.ConstConfigTypeSecret,
		Editor:      "password",
		Options:     nil,
		Label:       "Root password",
		Description: "password to enter admin panel as root",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreCountry,
		Value:       "US",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     models.ConstCountriesList,
		Label:       "Country",
		Description: "store location country",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreState,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     models.ConstStatesList,
		Label:       "State",
		Description: "store location state",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreCity,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "City",
		Description: "store location city",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreAddressline1,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 1",
		Description: "store location address line 1",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreAddressline2,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 2",
		Description: "store location address line 2",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreZip,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "zip",
		Description: "store location zip code",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreTimeZone,
		Value:       "UTC",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     models.ConstTimeZonesList,
		Label:       "Time zone",
		Description: "store location time zone",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Mail",
		Description: "web store mailing options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailServer,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Host",
		Description: "web store mailing server",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailPort,
		Value:       587,
		Type:        env.ConstConfigTypeInteger,
		Editor:      "integer",
		Options:     nil,
		Label:       "Port",
		Description: "web store mailing server port",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailUser,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "User",
		Description: "mail server username",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailPassword,
		Value:       nil,
		Type:        env.ConstConfigTypeSecret,
		Editor:      "password",
		Options:     nil,
		Label:       "Password",
		Description: "mail server password",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailFrom,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "From",
		Description: "full name for from field",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailSignature,
		Value:       "Sincerely, Ottemo",
		Type:        env.ConstConfigTypeText,
		Editor:      "multiline_text",
		Options:     nil,
		Label:       "Signature",
		Description: "mail signature",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Email Verfication
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathVerfifyEmail,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Enable Email Verification for Registration",
		Description: "send email verification link",
		Image:       "",
	}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Hide levels
	emailValidator := func(newValue interface{}) (interface{}, error) {
		newEmail := utils.InterfaceToString(newValue)
		if newEmail == "" {
			err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "d5abe68b-5bde-4b14-a3a7-b89507c14597", "recipient e-mail can not be blank")
			return "support+contactus@ottemo.io", env.ErrorDispatch(err)
		}
		if !utils.ValidEmailAddress(newEmail) {
			err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "720c5b82-4aa9-405b-b52d-b94d1f31e49d", "recipient e-mail is not in a valid format")
			return "support+contactus@ottemo.io", env.ErrorDispatch(err)

		}

		return newEmail, nil
	}
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathContactUsRecipient,
		Value:       "support+contactus@ottemo.io",
		Type:        utils.DataTypeWPrecision(env.ConstConfigTypeVarchar, 255),
		Editor:      "text",
		Options:     nil,
		Label:       "Contact Us Recipient",
		Description: "email of contact form recipient",
		Image:       "",
	}, emailValidator)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// API settings
	err = config.RegisterItem(env.StructConfigItem{
		Path:        rest.ConstConfigPathAPI,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "API",
		Description: "API relarted options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        rest.ConstConfigPathAPILog,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Log",
		Description: "API logger relarted options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        rest.ConstConfigPathAPILogEnable,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Enable API Logger",
		Description: "enable/disable api logger",
		Image:       "",
	}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

	if err != nil {
		return env.ErrorDispatch(err)
	}

	APIURIs := map[string]string{}

	// sorting handlers before output
	apiRestService := api.GetRestService()
	if restService, ok := apiRestService.(*rest.DefaultRestService); ok {
		sort.Strings(restService.Handlers)
		for _, handlerPath := range restService.Handlers {
			path := string(handlerPath)
			APIURIs[path] = path
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        rest.ConstConfigPathAPILogExclude,
			Value:       "",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "multi_select",
			Options:     APIURIs,
			Label:       "Exclude the following URI from being logged",
			Description: "exclude selected URI from being logged",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}
	} else {
		return env.ErrorDispatch(env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "890ae42a-b890-4481-81ba-5165a78acd11", "can't read config of Excluded URIs"))
	}

	return nil
}
