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
	ConstCollectionNameRTSReferrals    = "rts_referrals"

	ConstReferrerTypeDirect = 0
	ConstReferrerTypeSite   = 1
	ConstReferrerTypeSearch = 2

	ConstVisitorOnlineSeconds = 10

	ConstTimeDay = time.Hour * 24

	ConstErrorModule = "rts"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// Package global variables
var (
	updateSync sync.RWMutex

	referrers = make(map[string]int)         // collects and counts refers from external sites
	statistic = make(map[int64]*ActionsMade) // information about per hour site activity

	lastUpdate = time.Now()            // last update timer for day reset
	visitState = make(map[string]bool) // reflects session state: 1) not present - new visit, 2) false - addToCart not happened, 3) true - addToCart happened

	// OnlineSessions holds session based information about referrer type on first visit
	OnlineSessions = make(map[string]*OnlineReferrer)

	OnlineDirect      = 0
	OnlineSite        = 0
	OnlineSearch      = 0
	OnlineSessionsMax = 0
	OnlineDirectMax   = 0
	OnlineSiteMax     = 0
	OnlineSearchMax   = 0

	// knownSearchEngines is a list of search engines used to determinate referer type
	knownSearchEngines = []string{"www.daum.net", "www.google.com", "www.eniro.se", "www.naver.com", "www.yahoo.com",
		"www.msn.com", "www.bing.com", "www.aol.com", "www.aol.com", "www.lycos.com", "www.ask.com", "www.altavista.com",
		"search.netscape.com", "www.cnn.com", "www.about.com", "www.mamma.com", "www.alltheweb.com", "www.voila.fr",
		"search.virgilio.it", "www.bing.com", "www.baidu.com", "www.alice.com", "www.yandex.com", "www.najdi.org.mk",
		"www.aol.com", "www.mamma.com", "www.seznam.cz", "www.search.com", "www.wp.pl", "online.onetcenter.org",
		"www.szukacz.pl", "www.yam.com", "www.pchome.com", "www.kvasir.no", "sesam.no", "www.ozu.es", "www.terra.com",
		"www.mynet.com", "www.ekolay.net", "www.rambler.ru"}
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
type OnlineReferrer struct {
	referrerType int
	time         time.Time
}
