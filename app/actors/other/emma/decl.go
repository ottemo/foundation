package emma

import "github.com/ottemo/foundation/env"

// Package constants for Emma module
const (
	ConstErrorModule = "emma"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstConfigPathEmma                = "general.emma"
	ConstConfigPathEmmaEnabled         = "general.emma.enabled"
	ConstConfigPathEmmaPublicAPIKey    = "general.emma.public_api_key"
	ConstConfigPathEmmaPrivateAPIKey   = "general.emma.private_api_key"
	ConstConfigPathEmmaAccountID       = "general.emma.account_id"
	ConstConfigPathEmmaSKU             = "general.emma.trigger_sku"

	ConstEmmaApiUrl = "https://api.e2ma.net/"
)

