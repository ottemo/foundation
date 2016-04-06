// Package trustpilot implements trust pilot functions
package trustpilot

// Package global constants
const (
	ConstProductBrand = "Kari Gran"
	ConstEmailSubject = "Purchase feedback"

	ConstErrorModule = "trustpilot"
	ConstErrorLevel  = 1 // if i tell you to log, then do it

	ConstOrderCustomInfoLinkKey = "trustpilot_link"
	ConstOrderCustomInfoSentKey = "trustpilot_sent"

	ConstConfigPathTrustPilot                 = "general.trustpilot"
	ConstConfigPathTrustPilotEnabled          = "general.trustpilot.enabled"
	ConstConfigPathTrustPilotAPIKey           = "general.trustpilot.apiKey"
	ConstConfigPathTrustPilotAPISecret        = "general.trustpilot.apiSecret"
	ConstConfigPathTrustPilotBusinessUnitID   = "general.trustpilot.businessUnitID"
	ConstConfigPathTrustPilotUsername         = "general.trustpilot.username"
	ConstConfigPathTrustPilotPassword         = "general.trustpilot.password"
	ConstConfigPathTrustPilotAccessTokenURL   = "general.trustpilot.accessTokenURL"
	ConstConfigPathTrustPilotProductReviewURL = "general.trustpilot.productReviewURL"
	ConstConfigPathTrustPilotServiceReviewURL = "general.trustpilot.serviceReviewURL"
	ConstConfigPathTrustPilotEmailTemplate    = "general.trustpilot.emailTemplate"
)
