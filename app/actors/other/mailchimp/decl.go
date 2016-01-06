package mailchimp

import "github.com/ottemo/foundation/env"

// Package constants for Mailchimp module
const (
	ConstErrorModule = "mailchimp"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstMailchimp        = "general.mailchimp"
	ConstMailchimpEnabled = "general.mailchimp.enabled"
	ConstMailchimpAPIKey  = "general.mailchimp.api_key"
	ConstMailchimpBaseURL = "general.mailchimp.base_url"
)
