package trustpilot

import (
	"github.com/ottemo/foundation/env"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "b2c1c442-36b9-4994-b5d1-7c948a7552bd", "can't obtain config")
	}

	// Trust pilot config elements
	//----------------------------

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilot,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Trust Pilot",
		Description: "Trust Pilot settings",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotEnabled,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Trust Pilot",
		Description: `Enabled Trust Pilot sent order data`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotTestMode,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Test mode",
		Description: `Enabled Trust Pilot sent order data in test mode (add _test@ to email)`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotBusinessUnitID,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Business Unit ID",
		Description: `Trustpilot Business Unit ID`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotUsername,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Trustpilot username",
		Description: `Trustpilot b2b login email`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotPassword,
		Value:       "",
		Type:        env.ConstConfigTypeSecret,
		Editor:      "password",
		Options:     "",
		Label:       "Trustpilot password",
		Description: `Trustpilot b2b login password`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotAPIKey,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "API Key",
		Description: `Trustpilot API Key`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotAPISecret,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "API Secret",
		Description: `Trustpilot API Secret`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotAccessTokenURL,
		Value:       "https://api.trustpilot.com/v1/oauth/oauth-business-users-for-applications/accesstoken",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Access token URL",
		Description: `Trustpilot URL for getting access token and appending it to product review request`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	//

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotServiceReviewURL,
		Value:       "https://invitations-api.trustpilot.com/v1/private/business-units/{businessUnitId}/invitation-links",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Service review URL",
		Description: `Trustpilot service review URL, {businessUnitId} - will be rewrited by "Business Unit ID" config value`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTrustPilotProductReviewURL,
		Value:       "https://api.trustpilot.com/v1/private/product-reviews/business-units/{businessUnitId}/invitation-links",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Product review URL",
		Description: `Trustpilot product review URL, {businessUnitId} - will be rewrited by "Business Unit ID" config value`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path: ConstConfigPathTrustPilotEmailTemplate,
		Value: `Dear {{.Visitor.name}},
		please provide your feedback on recently purchase
		<br />
		{{.Visitor.link}}`,
		Type:        env.ConstConfigTypeHTML,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Trustpilot data send e-mail: ",
		Description: "contents of email will be sent to cutomers two weeks after purchase",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
