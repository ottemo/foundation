package rts

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("rts", "GET", "visit", restRegisterVisit)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "referrers", restGetReferrers)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "visits", restGetVisits)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "visits/details/:from/:to", restGetVisitsDetails)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "conversions", restGetConversions)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "sales", restGetSales)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "sales/details/:from/:to", restGetSalesDetails)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "top_sellers", restGetTopSellers)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "visitors/realtime", restGetVisitsRealtime)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API used for registering a new storefront visit and obtain referer info
func restRegisterVisit(params *api.StructAPIHandlerParams) (interface{}, error) {
	xReferrer := utils.InterfaceToString(params.Request.Header.Get("X-Referer"))

	http.SetCookie(params.ResponseWriter, &http.Cookie{Name: "X_Referrer", Value: xReferrer, Path: "/"})

	eventData := map[string]interface{}{"session": params.Session, "apiParams": params}
	env.Event("api.rts.visit", eventData)

	return nil, nil
}

// WEB REST API used to obtain site referers list
func restGetReferrers(params *api.StructAPIHandlerParams) (interface{}, error) {
	result := make(map[string]int)

	for url := range referrers {
		result[utils.InterfaceToString(url)] = referrers[utils.InterfaceToString(url)]
	}

	return result, nil
}

// WEB REST API used to obtain site visit information for current day
func restGetVisits(params *api.StructAPIHandlerParams) (interface{}, error) {
	result := make(map[string]interface{})

	err := GetTodayVisitorsData()
	if err != nil {
		return result, nil
	}

	result["visitsToday"] = visitorsInfoToday.Visitors
	result["ratio"] = 1

	err = GetYesterdayVisitorsData()
	if err != nil {
		return result, nil
	}
	countYesterday := visitorsInfoYesterday.Visitors
	countToday := visitorsInfoToday.Visitors
	if countYesterday != 0 {
		ratio := float64(countToday)/float64(countYesterday) - float64(1)
		result["ratio"] = utils.Round(ratio, 0.5, 2)
	}

	return result, nil
}

// WEB REST API used to obtain detailed site visit information for a specified period
func restGetVisitsDetails(params *api.StructAPIHandlerParams) (interface{}, error) {
	result := make(map[string]int)

	fromDatetmp, present := params.RequestURLParams["from"]
	if !present {
		fromDatetmp = time.Now().Format("2006-01-02")
	}
	fromDate, _ := time.Parse("2006-01-02", fromDatetmp)

	toDatetmp, present := params.RequestURLParams["to"]
	if !present {
		toDatetmp = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	toDate, _ := time.Parse("2006-01-02", toDatetmp)

	delta := toDate.Sub(fromDate)
	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSVisitors)
	if err == nil {
		visitorInfoCollection.AddFilter("day", ">=", fromDate)
		visitorInfoCollection.AddFilter("day", "<", toDate)
		visitorInfoCollection.AddSort("day", false)
		dbRecord, _ := visitorInfoCollection.Load()
		dbResult := make(map[string]int)
		if delta.Hours() > 48 {
			if len(dbRecord) > 0 {
				for _, item := range dbRecord {
					timestamp := fmt.Sprintf("%v", int32(utils.InterfaceToTime(item["day"]).Unix()))
					dbResult[timestamp] = utils.InterfaceToInt(item["visitors"])
				}
			}
			// group by days
			for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {
				timestamp := fmt.Sprintf("%v", int32(date.Unix()))
				result[timestamp] = dbResult[timestamp]
			}

		} else {
			if len(dbRecord) > 0 {
				details := DecodeDetails(utils.InterfaceToString(dbRecord[0]["details"]))
				for _, item := range details {
					timestamp := fmt.Sprintf("%v", int32(utils.InterfaceToTime(item.Time).Unix()))
					dbResult[timestamp]++
				}
			}
			//	group by hours
			for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {

				timestamp := int32(date.Unix())
				year := time.Unix(int64(timestamp), 0).Year()
				month := time.Unix(int64(timestamp), 0).Month()
				day := time.Unix(int64(timestamp), 0).Day()
				for hour := 0; hour < 24; hour++ {
					timeGroup := time.Date(year, month, day, hour, 0, 0, 0, time.Local)
					if timeGroup.Unix() > time.Now().Unix() {
						break
					}
					timestamp := fmt.Sprintf("%v", int32(timeGroup.Unix()))
					result[timestamp] = dbResult[timestamp]
				}
			}
		}
	}

	return result, nil
}

// WEB REST API used to obtain site conversation information
func restGetConversions(params *api.StructAPIHandlerParams) (interface{}, error) {
	result := make(map[string]interface{})

	result["totalVisitors"] = visitorsInfoToday.Visitors
	result["addedToCart"] = visitorsInfoToday.Cart
	result["reachedCheckout"] = visitorsInfoToday.Checkout
	result["purchased"] = visitorsInfoToday.Sales

	return result, nil
}

// WEB REST API used to get information on site sales for today
func restGetSales(params *api.StructAPIHandlerParams) (interface{}, error) {
	result := make(map[string]interface{})

	currDate, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	yesterdayDate, _ := time.Parse("2006-01-02", currDate.AddDate(0, 0, -1).Format("2006-01-02"))

	if sales.lastUpdate == 0 { // Init sales data
		GetTotalSales(currDate, yesterdayDate)
		result["today"] = sales.today
		result["ratio"] = sales.ratio
	} else {
		lastUpdate, _ := time.Parse("2006-01-02", time.Unix(int64(sales.lastUpdate), 0).Format("2006-01-02"))
		delta := currDate.Sub(lastUpdate)
		if delta > 1 { // Updates the sales data if they older 1 hour
			GetTotalSales(currDate, yesterdayDate)
			result["today"] = sales.today
			result["ratio"] = sales.ratio
		} else {
			// Returns the  existing data
			result["today"] = sales.today
			result["ratio"] = sales.ratio
		}
	}

	return result, nil
}

// WEB REST API used to get information on site sales for a specified period
func restGetSalesDetails(params *api.StructAPIHandlerParams) (interface{}, error) {
	result := make(map[string]int)

	// check request params
	//---------------------
	fromDatetmp, present := params.RequestURLParams["from"]
	if !present {
		fromDatetmp = time.Now().Format("2006-01-02")
	}
	fromDate, _ := time.Parse("2006-01-02", fromDatetmp)
	// check request params
	//---------------------
	toDatetmp, present := params.RequestURLParams["to"]
	if !present {
		toDatetmp = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	toDate, _ := time.Parse("2006-01-02", toDatetmp)

	h := md5.New()
	io.WriteString(h, fromDatetmp+"/"+toDatetmp)
	periodHash := fmt.Sprintf("%x", h.Sum(nil))

	if _, ok := salesDetail[periodHash]; !ok {
		salesDetail[periodHash] = &SalesDetailData{Data: make(map[string]int)}

		GetSalesDetail(fromDate, toDate, periodHash)

	} else {
		// check last updates
		if salesDetail[periodHash].lastUpdate == 0 {

			GetSalesDetail(fromDate, toDate, periodHash)

		} else {
			currDate, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
			lastUpdate, _ := time.Parse("2006-01-02", time.Unix(int64(salesDetail[periodHash].lastUpdate), 0).Format("2006-01-02"))
			delta := currDate.Sub(lastUpdate)
			if delta > 1 { // Updates the sales data if they older than 1 hour
				GetSalesDetail(fromDate, toDate, periodHash)
			}
		}
	}

	result = salesDetail[periodHash].Data

	return result, nil
}

// WEB REST API used to get information on site top sellers
func restGetTopSellers(params *api.StructAPIHandlerParams) (interface{}, error) {
	result := make(map[string]*SellerInfo)

	salesCollection, err := db.GetCollection(ConstCollectionNameRTSSales)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	salesCollection.AddFilter("count", ">", 0)
	salesCollection.AddSort("count", true)
	salesCollection.SetLimit(0, 5)
	items, _ := salesCollection.Load()

	productModel, err := product.GetProductModel()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, item := range items {
		productID := utils.InterfaceToString(item["product_id"])
		result[productID] = &SellerInfo{}
		if _, ok := topSellers.Data[productID]; !ok {

			product, err := product.LoadProductByID(productID)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			productModel.SetID(productID)
			mediaPath, err := productModel.GetMediaPath("image")
			if err != nil {
				return result, env.ErrorDispatch(err)
			}

			if product.GetDefaultImage() != "" {
				result[productID].Image = mediaPath + product.GetDefaultImage()
			}

			result[productID].Name = product.GetName()
		}

		result[productID].Count = utils.InterfaceToInt(item["count"])
	}

	return result, nil
}

// WEB REST API used to get information on site real time visitors
func restGetVisitsRealtime(params *api.StructAPIHandlerParams) (interface{}, error) {
	result := make(map[string]interface{})
	ratio := float64(0)

	result["Online"] = len(OnlineSessions)
	if OnlineSessionsMax == 0 || len(OnlineSessions) == 0 {
		ratio = float64(0)
	} else {
		ratio = float64(len(OnlineSessions)) / float64(OnlineSessionsMax)
	}
	result["OnlineRatio"] = utils.Round(ratio, 0.5, 2)

	result["Direct"] = OnlineDirect
	if OnlineDirectMax == 0 || OnlineDirect == 0 {
		ratio = float64(0)
	} else {
		ratio = float64(OnlineDirect) / float64(OnlineDirectMax)
	}
	result["DirectRatio"] = utils.Round(ratio, 0.5, 2)

	result["Search"] = OnlineSearch
	if OnlineSearchMax == 0 || OnlineSearch == 0 {
		ratio = float64(0)
	} else {
		ratio = float64(OnlineSearch) / float64(OnlineSearchMax)
	}
	result["SearchRatio"] = utils.Round(ratio, 0.5, 2)

	result["Site"] = OnlineSite
	if OnlineSiteMax == 0 || OnlineSite == 0 {
		ratio = float64(0)
	} else {
		ratio = float64(OnlineSite) / float64(OnlineSiteMax)
	}
	result["SiteRatio"] = utils.Round(ratio, 0.5, 2)

	return result, nil
}
