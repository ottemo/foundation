package app

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "f635e96c3cd74ae2a5074349021e9f13", "can't obtain config")
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGroup,
		Value:       nil,
		Type:        env.ConstConfigItemGroupType,
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
		Type:        env.ConstConfigItemGroupType,
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
		Type:        "varchar(255)",
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
		Type:        "varchar(255)",
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
		Type:        "varchar(255)",
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
		Type:        env.ConstConfigItemGroupType,
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
		Type:        "varchar(255)",
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
		Type:        "varchar(255)",
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
		Type:        "varchar(255)",
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
		Type:        "varchar(255)",
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
		Type:        "string",
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
		Type:        "string",
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
		Type:        "string",
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
		Type:        "string",
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
		Type:        "string",
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
		Type:        "string",
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
		Path:        ConstConfigPathMailGroup,
		Value:       nil,
		Type:        env.ConstConfigItemGroupType,
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
		Type:        "varchar(255)",
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
		Type:        "int",
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
		Type:        "varchar(100)",
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
		Type:        "varchar(100)",
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
		Type:        "varchar(100)",
		Editor:      "text",
		Options:     nil,
		Label:       "From",
		Description: "sending mail from field",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailSignature,
		Value:       "Sincerely, Ottemo",
		Type:        "text",
		Editor:      "multiline_text",
		Options:     nil,
		Label:       "Signature",
		Description: "sending mail signature",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
