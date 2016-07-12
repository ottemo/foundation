// Package trustpilot implements trust pilot functions
package trustpilot

import (
	"time"
)

// Package global constants
const (
	ConstEmailSubject = "Purchase feedback"

	ConstErrorModule = "trustpilot"
	ConstErrorLevel  = 1 // if i tell you to log, then do it

	ConstOrderCustomInfoLinkKey = "trustpilot_link"
	ConstOrderCustomInfoSentKey = "trustpilot_sent"

	ConstConfigPathTrustPilot               = "general.trustpilot"
	ConstConfigPathTrustPilotEnabled        = "general.trustpilot.enabled"
	ConstConfigPathTrustPilotAPIKey         = "general.trustpilot.apiKey"
	ConstConfigPathTrustPilotAPISecret      = "general.trustpilot.apiSecret"
	ConstConfigPathTrustPilotBusinessUnitID = "general.trustpilot.businessUnitID"
	ConstConfigPathTrustPilotUsername       = "general.trustpilot.username"
	ConstConfigPathTrustPilotPassword       = "general.trustpilot.password"
	ConstConfigPathTrustPilotEmailTemplate  = "general.trustpilot.emailTemplate"
	ConstConfigPathTrustPilotProductBrand   = "general.trustpilot.productBrand"

	ConstRatingSummaryURL = "https://api.trustpilot.com/v1/private/product-reviews/business-units/{businessUnitId}/summaries"
)

var (
	lastTimeSummariesUpdate time.Time
	summariesCache          interface{}
)
