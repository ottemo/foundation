package xdomain

import "github.com/ottemo/foundation/env"

// xdomain package constants
const (
	ConstErrorModule = "xdomain"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// xdomain package level globals
var (
	xdomainMasterURL = "http://*.staging.ottemo.io/" // default to all stores in staging
)
