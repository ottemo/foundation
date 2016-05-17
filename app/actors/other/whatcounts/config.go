package whatcounts

import (
	"github.com/ottemo/foundation/env"
)

func setupConfig() error {

	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "6b78d38a-35c5-4aa2-aec1-eaa16830ff61", "Error configuring Whatcounts module")
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathWhatcounts,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "WhatCounts",
		Description: "WhatCounts Settings",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathWhatcountsEnabled,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "WhatCounts Enabled",
		Description: "Enable WhatCounts integration(defaults to false)",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathWhatcountsRealm,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "WhatCounts Realm",
		Description: "Enter your WhatCounts Realm",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathWhatcountsAPIKey,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "WhatCounts API Key",
		Description: "Enter your WhatCounts API Key",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathWhatcountsBaseURL,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "WhatCounts Base URL",
		Description: "Defines the base url for this account",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path: ConstConfigPathWhatcountsEmailTemplate,
		Value: `Warning  ....
		<br />
		<br />
		The following email address could not be added to Whatcounts:
		{{.email_address}}`,
		Type:        env.ConstConfigTypeHTML,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Support Email Template",
		Description: "Template for sending support emails",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathWhatcountsSupportAddress,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Support Email Address",
		Description: "Email address to send errors encountered when adding to lists",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathWhatcountsSubjectLine,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Support Email Subject",
		Description: "Subject Line for emails describing whatcounts list addition failures",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathWhatcountsList,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "WhatCounts List ID",
		Description: "Enter your WhatCounts List ID",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathWhatcountsSKU,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Trigger SKU (comma seperated list of SKUs)",
		Description: "Enter the SKU you want to use as a trigger",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathWhatcountsNoConfirm,
		Value:       true,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Override Confirmation - 0 for Yes, 1 for No ",
		Description: "0 - send confirmation email, 1 - do not send confirmation email",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathWhatcountsForceSub,
		Value:       true,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Force Subscribe to List",
		Description: "0 - do not force add to list, 1 - force add to list",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
