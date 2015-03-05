// Package rts implements Real Time Statistics calculations module
package rts

import (
	"github.com/ottemo/foundation/env"
)

import "time"

// Package global constants
const (
	ConstCollectionNameRTSSalesHistory = "rts_sales_history"
	ConstCollectionNameRTSSales        = "rts_sales"
	ConstCollectionNameRTSVisitors     = "rts_visitors"

	ConstReferrerTypeDirect = 0
	ConstReferrerTypeSite   = 1
	ConstReferrerTypeSearch = 2

	ConstVisitorAddToCart     = 1
	ConstVisitorCheckout      = 2
	ConstVisitorSales         = 3

	ConstVisitorOnlineSeconds = 10

	ConstErrorModule = "rts"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// Package global variables
var (
	referrers             = make(map[string]int) // collects and counts refers from external sites
	statistic      = make(map[int64]*ActionsMade) 	// information about per hour site activity

	visitState      = make(map[string]bool)			//checks a buying status of visitor by it sessionID

	// OnlineSessions holds session based information about referer type on first visit
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
	Visit 	int       	// count site visits
	Cart     int       // count times products was added to cart
	Sales    int       // count of orders visitors made
}

// TopSellers holds information about best sellers
type TopSellers struct {
	Data       map[string]*SellerInfo // product id based map holds best sellers
	lastUpdate int64                  // timestamp used to update struct once in a hour
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
