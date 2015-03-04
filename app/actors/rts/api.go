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

	err = api.GetRestService().RegisterAPI("rts/visit", api.ConstRESTOperationCreate, APIRegisterVisit)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/referrers", api.ConstRESTOperationGet, APIGetReferrers)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/visits", api.ConstRESTOperationGet, APIGetVisits)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/visits/detail/:from/:to", api.ConstRESTOperationGet, APIGetVisitsDetails)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/conversion", api.ConstRESTOperationGet, APIGetConversion)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/sales", api.ConstRESTOperationGet, APIGetSales)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/sales/detail/:from/:to", api.ConstRESTOperationGet, APIGetSalesDetails)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/bestsellers", api.ConstRESTOperationGet, APIGetBestsellers)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/visits/realtime", api.ConstRESTOperationGet, APIGetVisitsRealtime)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIRegisterVisit registers request for a statistics
func APIRegisterVisit(context api.InterfaceApplicationContext) (interface{}, error) {
	if httpRequest, ok := context.GetRequest().(*http.Request); ok && httpRequest != nil {
		if httpResponseWriter, ok := context.GetResponseWriter().(http.ResponseWriter); ok && httpResponseWriter != nil {
			xReferrer := utils.InterfaceToString(httpRequest.Header.Get("X-Referer"))

			http.SetCookie(httpResponseWriter, &http.Cookie{Name: "X_Referrer", Value: xReferrer, Path: "/"})

			eventData := map[string]interface{}{"session": context.GetSession(), "context": context}
			env.Event("api.rts.visit", eventData)

			return nil, nil
		}
	}
	return nil, nil
}

// APIGetReferrers returns list of unique referrers were registered
func APIGetReferrers(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]int)

	for url, count := range referrers {
		result[url] = count
	}

	return result, nil
}

// APIGetVisits returns site visit information for current day
func APIGetVisits(context api.InterfaceApplicationContext) (interface{}, error) {
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

// APIGetVisitsDetails returns detailed site visit information for a specified period
//   - period start and end dates should be specified in "from" and "to" attributes in DD-MM-YYY format
func APIGetVisitsDetails(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]int)

	requestFromDate := context.GetRequestArgument("from")
	if requestFromDate == "" {
		requestFromDate = time.Now().Format("2006-01-02")
	}
	fromDate, _ := time.Parse("2006-01-02", requestFromDate)

	requestToDate := context.GetRequestArgument("to")
	if requestToDate == "" {
		requestToDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	toDate, _ := time.Parse("2006-01-02", requestToDate)

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
				currentTime := time.Unix(int64(timestamp), 0)
				for hour := 0; hour < 24; hour++ {
					timeGroup := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), hour, 0, 0, 0, currentTime.Location())
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

// APIGetConversion returns site conversion information
func APIGetConversion(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})

	result["totalVisitors"] = visitorsInfoToday.Visitors
	result["addedToCart"] = visitorsInfoToday.Cart
	result["reachedCheckout"] = visitorsInfoToday.Checkout
	result["purchased"] = visitorsInfoToday.Sales

	return result, nil
}

// APIGetSales returns information on site sales for today
func APIGetSales(context api.InterfaceApplicationContext) (interface{}, error) {
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

// APIGetSalesDetails returns site sales information for a specified period
//   - period start and end dates should be specified in "from" and "to" attributes in DD-MM-YYY format
func APIGetSalesDetails(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]int)
	currentTime := time.Now()

	requestFromDate := context.GetRequestArgument("from")
	if requestFromDate == "" {
		requestFromDate = currentTime.Format("2006-01-02")
	}
	fromDate, _ := time.Parse("2006-01-02", requestFromDate)

	requestToDate := context.GetRequestArgument("to")
	if requestToDate == "" {
		requestToDate = currentTime.AddDate(0, 0, -1).Format("2006-01-02")
	}
	toDate, _ := time.Parse("2006-01-02", requestToDate)

	hashCode := md5.New()
	io.WriteString(hashCode, requestFromDate+"/"+requestToDate)
	periodHash := fmt.Sprintf("%x", hashCode.Sum(nil))

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

// APIGetBestsellers returns information on site bestsellers
func APIGetBestsellers(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]*SellerInfo)

	salesCollection, err := db.GetCollection(ConstCollectionNameRTSSales)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	salesCollection.AddFilter("count", ">", 0)
	salesCollection.AddSort("count", true)
	salesCollection.SetLimit(0, 5)
	collectionRecords, _ := salesCollection.Load()

	for _, item := range collectionRecords {
		productID := utils.InterfaceToString(item["product_id"])
		result[productID] = &SellerInfo{}

		if _, ok := topSellers.Data[productID]; !ok {

			productInstance, err := product.LoadProductByID(productID)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			mediaPath, err := productInstance.GetMediaPath("image")
			if err != nil {
				return result, env.ErrorDispatch(err)
			}

			if productInstance.GetDefaultImage() != "" {
				result[productID].Image = mediaPath + productInstance.GetDefaultImage()
			}

			result[productID].Name = productInstance.GetName()
		}

		result[productID].Count = utils.InterfaceToInt(item["count"])
	}

	return result, nil
}

// APIGetVisitsRealtime returns real-time information on current visits
func APIGetVisitsRealtime(context api.InterfaceApplicationContext) (interface{}, error) {
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
