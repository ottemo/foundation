package mailchimp

import "github.com/ottemo/foundation/env"

// Package constants for Mailchimp module
const (
	ConstErrorModule = "grouping"
	ConstErrorLevel  = env.ConstErrorLevelActor

	CollectionNameMailchimp = "mailchimp_data"

	MailchimpConfigPath    = "mailchimp"
	MailchimpEnabledConfig = "mailchimp.enabled"
	MailchimpAPIKeyConfig  = "mailchimp.api_key"
	MailchimpBaseURLConfig = "mailchimp.base_url"
)
