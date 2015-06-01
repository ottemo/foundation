// Package trustpilot implements trust pilot functions
package trustpilot

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstTestMode = false

	ConstErrorModule = "trustpilot"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstConfigPathTrustPilot                 = "general.trustpilot"
	ConstConfigPathTrustPilotEnabled          = "general.trustpilot.enabled"
	ConstConfigPathTrustPilotAPIKey           = "general.trustpilot.apiKey"
	ConstConfigPathTrustPilotAPISecret        = "general.trustpilot.apiSecret"
	ConstConfigPathTrustPilotBusinessUnitID   = "general.trustpilot.businessUnitID"
	ConstConfigPathTrustPilotUsername         = "general.trustpilot.username"
	ConstConfigPathTrustPilotPassword         = "general.trustpilot.password"
	ConstConfigPathTrustPilotAccessTokenURL   = "general.trustpilot.accessTokenURL"
	ConstConfigPathTrustPilotProductReviewURL = "general.trustpilot.productReviewURL"
)
