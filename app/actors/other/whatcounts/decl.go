package whatcounts

import "github.com/ottemo/foundation/env"

// Package constants for Whatcounts module
const (
	ConstErrorModule = "whatcounts"
	ConstErrorLevel  = env.ConstErrorLevelStartStop

	ConstConfigPathWhatcounts               = "general.whatcounts"
	ConstConfigPathWhatcountsEnabled        = "general.whatcounts.enabled"
	ConstConfigPathWhatcountsRealm          = "general.whatcounts.realm"
	ConstConfigPathWhatcountsAPIKey         = "general.whatcounts.api_key"
	ConstConfigPathWhatcountsBaseURL        = "general.whatcounts.base_url"
	ConstConfigPathWhatcountsSupportAddress = "general.whatcounts.support_addr"
	ConstConfigPathWhatcountsEmailTemplate  = "general.whatcounts.template"
	ConstConfigPathWhatcountsSubjectLine    = "general.whatcounts.subject_line"
	ConstConfigPathWhatcountsList           = "general.whatcounts.subscribe_to_list"
	ConstConfigPathWhatcountsSKU            = "general.whatcounts.trigger_sku"
	ConstConfigPathWhatcountsNoConfirm      = "general.whatcounts.no_confirm"
	ConstConfigPathWhatcountsForceSub       = "general.whatcounts.force_subscribe"
)

// Registration is a struct to hold a single registation for a Whatcounts mailing list.
type Registration struct {
	EmailAddress string `json:"email_address"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}
