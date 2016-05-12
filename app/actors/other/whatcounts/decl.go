package whatcounts

import "github.com/ottemo/foundation/env"

// Package constants for Whatcounts module
const (
	ConstErrorModule = "whatcounts"
	ConstErrorLevel  = env.ConstErrorLevelStartStop

	ConstWhatcountsSubscribeStatus = "subscribed"

	ConstConfigPathWhatcounts               = "general.whatcounts"
	ConstConfigPathWhatcountsEnabled        = "general.whatcounts.enabled"
	ConstConfigPathWhatcountsAPIKey         = "general.whatcounts.api_key"
	ConstConfigPathWhatcountsBaseURL        = "general.whatcounts.base_url"
	ConstConfigPathWhatcountsSupportAddress = "general.whatcounts.support_addr"
	ConstConfigPathWhatcountsEmailTemplate  = "general.whatcounts.template"
	ConstConfigPathWhatcountsSubjectLine    = "general.whatcounts.subject_line"
	ConstConfigPathWhatcountsList           = "general.whatcounts.subscribe_to_list"
	ConstConfigPathWhatcountsSKU            = "general.whatcounts.trigger_sku"
)

// Registration is a struct to hold a single registation for a Whatcounts mailing list.
type Registration struct {
	EmailAddress string            `json:"email_address"`
	Status       string            `json:"status"`
	MergeFields  map[string]string `json:"merge_fields"`
}
