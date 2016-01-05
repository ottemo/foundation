package mailchimp

import "github.com/ottemo/foundation/env"

// Package constants for Mailchimp module
const (
	ConstErrorModule = "mailchimp"
	ConstErrorLevel  = env.ConstErrorLevelActor

	//CollectionNameMailchimp = "mailchimp_data"
	ConstConfigPathMailchimp        = "general.mailchimp"
	ConstConfigPathMailchimpEnabled = "general.mailchimp.enabled"
	ConstConfigPathMailchimpAPIKey  = "general.mailchimp.api_key"
	ConstConfigPathMailchimpBaseURL = "general.mailchimp.base_url"
)
