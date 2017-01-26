// Package rts implements Real Time Statistics calculations module
package rts

import (
	"github.com/ottemo/foundation/env"
)

import (
	"sync"
	"time"
)

// Package global constants
const (
	ConstCollectionNameRTSSalesHistory = "rts_sales_history"
	ConstCollectionNameRTSVisitors     = "rts_visitors"
	// ConstCollectionNameRTSReferrals    = "rts_referrals"

	ConstReferrerTypeDirect = 0
	ConstReferrerTypeSite   = 1
	ConstReferrerTypeSearch = 2

	ConstVisitorOnlineSeconds = 10

	ConstTimeDay = time.Hour * 24

	ConstErrorModule = "rts"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstConfigPathCheckoutPath = "general.app.checkout_path"
)

// Package global variables
var (
	updateSync sync.RWMutex

	// referrers      = make(map[string]int)         // collects and counts refers from external sites
	statistic      = make(map[int64]*ActionsMade) // information about per hour site activity
	monthStatistic = new(ActionsMade)             // information total month activity

	lastUpdate = time.Now() // last update timer for day reset

	// reflects session state:
	// 1) not present - new visit,
	// 2) false - addToCart not happened,
	// 3) true - addToCart happened
	visitState = make(map[string]bool)
)

// ActionsMade contains info of visits, cart create and sales made for a hour
type ActionsMade struct {
	Visit         int     // count site visits
	Cart          int     // count times products was added to cart
	Sales         int     // count of orders visitors made
	TotalVisits   int     // total visits count
	SalesAmount   float64 // count sales
	VisitCheckout int     // count visitors reached checkout
	SetPayment    int     // count payment methods used
}

// SellerInfo represents particular product in TopSellers struct
type SellerInfo struct {
	Name  string // product name
	Image string // product image
	Count int    // times bought
}

// OnlineReferrer holds information about visit referer type and visit time
// type OnlineReferrer struct {
// 	referrerType int
// 	time         time.Time
// }
