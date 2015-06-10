package app

import (
	"time"

	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstVersionMajor = 0
	ConstVersionMinor = 9
	ConstSprintNumber = 8

	ConstConfigPathGroup      = "general"
	ConstConfigPathAppGroup   = "general.app"
	ConstConfigPathStoreGroup = "general.store"
	ConstConfigPathMailGroup  = "general.mail"

	ConstConfigPathStorefrontURL = "general.app.storefront_url"
	ConstConfigPathDashboardURL  = "general.app.dashboard_url"
	ConstConfigPathFoundationURL = "general.app.foundation_url"

	ConstConfigPathStoreName  = "general.store.name"
	ConstConfigPathStoreEmail = "general.store.email"

	ConstConfigPathStoreRootLogin    = "general.store.root_login"
	ConstConfigPathStoreRootPassword = "general.store.root_password"

	ConstConfigPathStoreCountry      = "general.store.country"
	ConstConfigPathStoreState        = "general.store.state"
	ConstConfigPathStoreCity         = "general.store.city"
	ConstConfigPathStoreAddressline1 = "general.store.addressline1"
	ConstConfigPathStoreAddressline2 = "general.store.addressline2"
	ConstConfigPathStoreZip          = "general.store.zip"

	ConstConfigPathMailFrom      = "general.mail.from"
	ConstConfigPathMailSignature = "general.mail.footer"
	ConstConfigPathMailServer    = "general.mail.server"
	ConstConfigPathMailPort      = "general.mail.port"
	ConstConfigPathMailUser      = "general.mail.user"
	ConstConfigPathMailPassword  = "general.mail.password"

	ConstErrorModule = "app"
	ConstErrorLevel  = env.ConstErrorLevelService

	ConstAllowGuest = true
)

// build related information supposed to be specified through -ldflags "-X key value"
//   - sample: go build -ldflags "-X github.com/ottemo/foundation/app.buildDate '`date`'"
var (
	buildTags   string
	buildDate   string
	buildNumber string
	buildBranch string

	startTime = time.Now().UTC().Truncate(time.Second)
)
