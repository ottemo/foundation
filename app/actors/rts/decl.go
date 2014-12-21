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
	visitorsInfoToday     = new(dbVisitorsRow)   // information on current day site visits
	visitorsInfoYesterday = new(dbVisitorsRow)   // information on previous day site visits

	sales       = new(Sales)                        // amount of sales for today (created orders)
	salesDetail = make(map[string]*SalesDetailData) // sales details on day/hour basis (used for dashboard graphs)
	topSellers  = new(TopSellers)                   // top sellers information

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

// Visits - unknown purpose structure
type Visits struct {
	Data      map[string]map[string]int32
	Yesterday string
	Today     string
}

// VisitorDetail - unknown purpose structure
type VisitorDetail struct {
	Time     time.Time
	Checkout int
}

// dbVisitorsRow represents database record on site visit activity for a day
type dbVisitorsRow struct {
	ID       string    // database record _id
	Day      time.Time // day when event happened
	Visitors int       // count site visits
	Cart     int       // count times product was added to cart
	Checkout int       // count times checkout acheeved
	Sales    int       // count of orders visitor made
	Details  map[string]*VisitorDetail
}

// Sales represents statistics on created orders for a day
type Sales struct {
	lastUpdate int64   // timestamp used to track current day
	today      int     // today orders made
	yesterday  int     // yesterday orders made
	ratio      float64 // % sales for today in compare to yesterday: (today / yesterday)-1
}

// SalesDetailData holds hour/day detailed information about sales
type SalesDetailData struct {
	Data       map[string]int // count of sales for specified in key time
	lastUpdate int64          // timestamp used to update struct once in a hour
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
