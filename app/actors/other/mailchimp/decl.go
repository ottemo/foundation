package mailchimp

import "github.com/ottemo/foundation/env"

// Package constants for Mailchimp module
const (
	ConstErrorModule = "mailchimp"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstConfigPathMailchimp               = "general.mailchimp"
	ConstConfigPathMailchimpEnabled        = "general.mailchimp.enabled"
	ConstConfigPathMailchimpAPIKey         = "general.mailchimp.api_key"
	ConstConfigPathMailchimpBaseURL        = "general.mailchimp.base_url"
	ConstConfigPathMailchimpSupportAddress = "general.mailchimp.support_addr"
	ConstConfigPathMailchimpEmailTemplate  = "general.mailchimp.template"
	ConstConfigPathMailchimpSubjectLine    = "general.mailchimp.subject_line"
	ConstConfigPathMailchimpList           = "general.mailchimp.subscribe_to_list"
	ConstConfigPathMailchimpSKU            = "general.mailchimp.trigger_sku"
)

// Registration is a struct to hold a single registation for a Mailchimp mailing list.
type Registration struct {
	EmailAddress string            `json:"email_address"`
	Status       string            `json:"status"`
	MergeFields  map[string]string `json:"merge_fields"`
}
