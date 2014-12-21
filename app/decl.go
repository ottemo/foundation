// Package app represents application object.
//
// That package contains routines to register callbacks on application start/end,
// API functions for administrator login, system configuration values, etc.
//
// In general this package contains the stuff addressed to system application.
package app

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
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
