package rts

import (
	"fmt"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI configures package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.POST("rts/visit", APIRegisterVisit)

	service.GET("rts/visits", APIGetVisits)
	service.GET("rts/visits/detail/:from/:to", APIGetVisitsDetails)
	// service.GET("rts/visits/realtime", APIGetVisitsRealtime)

	service.GET("rts/sales", APIGetSales)
	service.GET("rts/sales/detail/:from/:to", APIGetSalesDetails)

	service.GET("rts/conversion", APIGetConversion)
	service.GET("rts/bestsellers", APIGetBestsellers)
	// service.GET("rts/referrers", APIGetReferrers)

	return nil
}

// APIRegisterVisit registers request for a statistics
func APIRegisterVisit(context api.InterfaceApplicationContext) (interface{}, error) {
	// Variables in post; path=/shop/cleaning-products, referrer=http://google.com
	// In headers; Referrer=http://karigran.com/shop/cleaning-products
	eventData := map[string]interface{}{"session": context.GetSession(), "context": context}
	env.Event("api.rts.visit", eventData)

	requestData, _ := api.GetRequestContentAsMap(context)
	path := utils.InterfaceToString(requestData["path"])
	var checkoutPath = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathCheckoutPath))
	if checkoutPath != "" && path == checkoutPath {
		// record rts event for  checkout
		currentCheckout, err := checkout.GetCurrentCheckout(context, false)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		eventData = map[string]interface{}{"session": context.GetSession(), "checkout": currentCheckout}
		env.Event("api.checkout.visit", eventData)
	}

	return nil, nil
}

// APIGetVisits returns site visit information for a specified local day
func APIGetVisits(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})

	timeZone, err := app.GetSessionTimeZone(context.GetSession())
	if err != nil || timeZone == "" {
		timeZone = utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))
	}

	// get a hours pasted for local day and count for them and for previous day
	todayTo := time.Now().Truncate(time.Hour)
	todayFrom, _ := utils.MakeUTCOffsetTime(todayTo, utils.InterfaceToString(timeZone))
	if utils.IsZeroTime(todayFrom) {
		todayFrom = todayTo
	}

	todayHoursPast := time.Duration(todayFrom.Hour()) * time.Hour
	todayFrom = todayTo.Add(-todayHoursPast)
	yesterdayFrom := todayFrom.AddDate(0, 0, -1)
	weekFrom := yesterdayFrom.AddDate(0, 0, -5)

	// get data for visits
	todayStats, err := GetRangeStats(todayFrom, todayTo)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	todayVisits := todayStats.Visit
	todayTotalVisits := todayStats.TotalVisits

	// excluding last our for yesterday range statistic
	yesterdayTo := todayFrom.Add(-time.Nanosecond)
	yesterdayStats, err := GetRangeStats(yesterdayFrom, yesterdayTo)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	yesterdayVisits := yesterdayStats.Visit
	yesterdayTotalVisits := yesterdayStats.TotalVisits

	// excluding last our for week range statistic
	weekTo := yesterdayFrom.Add(-time.Nanosecond)
	weekStats, err := GetRangeStats(weekFrom, weekTo)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	weekVisits := yesterdayVisits + todayVisits + weekStats.Visit
	weekTotalVisits := yesterdayVisits + todayVisits + weekStats.TotalVisits

	result["total"] = map[string]int{
		"today":       todayTotalVisits,
		"yesterday":   yesterdayTotalVisits,
		"week":        weekTotalVisits,
		"monthToDate": monthStatistic.TotalVisits,
	}
	result["unique"] = map[string]int{
		"today":       todayVisits,
		"yesterday":   yesterdayVisits,
		"week":        weekVisits,
		"monthToDate": monthStatistic.Visit,
	}

	return result, nil
}

// APIGetVisitsDetails returns detailed site visit information for a specified period
//   - period start and end dates should be specified in "from" and "to" attributes in YYYY-MM-DD format
func APIGetVisitsDetails(context api.InterfaceApplicationContext) (interface{}, error) {

	// getting initial values
	result := make(map[string]int)
	var arrayResult [][]int

	timeZone, err := app.GetSessionTimeZone(context.GetSession())
	if err != nil || timeZone == "" {
		timeZone = utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))
	}
	dateFrom := utils.InterfaceToTime(context.GetRequestArgument("from"))
	dateTo := utils.InterfaceToTime(context.GetRequestArgument("to"))

	// checking if user specified correct from and to dates
	if dateFrom.IsZero() {
		dateFrom = time.Now().Truncate(ConstTimeDay)
	}

	if dateTo.IsZero() {
		dateTo = time.Now().Truncate(ConstTimeDay)
	}

	if dateFrom == dateTo {
		dateTo = dateTo.Add(ConstTimeDay)
	}

	// time zone recognize routines save time difference to show in graph by local time
	hoursOffset := time.Hour * 0

	if timeZone != "" {
		dateFrom, hoursOffset = utils.MakeUTCTime(dateFrom, timeZone)
		dateTo, _ = utils.MakeUTCTime(dateTo, timeZone)
	}

	// determining required scope
	delta := dateTo.Sub(dateFrom)

	timeScope := time.Hour
	if delta.Hours() > 48 {
		timeScope = timeScope * 24
	}
	dateFrom = dateFrom.Truncate(time.Hour)
	dateTo = dateTo.Truncate(time.Hour)

	// making database request
	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSVisitors)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorInfoCollection.AddFilter("day", ">=", dateFrom)
	visitorInfoCollection.AddFilter("day", "<=", dateTo)

	dbRecords, err := visitorInfoCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// filling requested period
	timeIterator := dateFrom
	for timeIterator.Before(dateTo) {
		arrayResult = append(arrayResult, []int{utils.InterfaceToInt(timeIterator.Add(hoursOffset).Unix()), 0})
		result[fmt.Sprint(timeIterator.Add(hoursOffset).Unix())] = 0
		timeIterator = timeIterator.Add(timeScope)
	}

	// grouping database records
	for _, item := range dbRecords {
		timestamp := fmt.Sprint(utils.InterfaceToTime(item["day"]).Truncate(timeScope).Add(hoursOffset).Unix())
		visits := utils.InterfaceToInt(item["visitors"])

		if value, present := result[timestamp]; present {
			result[timestamp] = value + visits
		} else {
			env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "80666c27-e67a-420d-9625-004122523451", timestamp+" - not present in result"))
		}
	}

	for _, item := range arrayResult {
		item[1] = result[utils.InterfaceToString(item[0])]
	}

	return arrayResult, nil
}

// APIGetConversion returns site conversion information
func APIGetConversion(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})

	timeZone, err := app.GetSessionTimeZone(context.GetSession())
	if err != nil || timeZone == "" {
		timeZone = utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))
	}

	// get hours for current local day and count only for them
	todayTo := time.Now().Truncate(time.Hour).Add(time.Hour)
	todayFrom, _ := utils.MakeUTCOffsetTime(todayTo, timeZone)
	if utils.IsZeroTime(todayFrom) {
		todayFrom = todayTo
	}

	todayHoursPast := time.Duration(todayFrom.Hour()) * time.Hour
	todayFrom = todayTo.Add(-todayHoursPast)

	visitorCount := 0
	addToCartByVisitor := 0
	visitCheckoutByVisitor := 0
	setPaymentByVisitor := 0
	saleTransactionsByVisitor := 0

	// Go through period and summarize visits
	for todayFrom.Before(todayTo) {

		todayFromStamp := todayFrom.Unix()
		if _, present := statistic[todayFromStamp]; present && statistic[todayFromStamp] != nil {
			visitorCount += statistic[todayFromStamp].Visit
			addToCartByVisitor += statistic[todayFromStamp].Cart
			visitCheckoutByVisitor += statistic[todayFromStamp].VisitCheckout
			setPaymentByVisitor += statistic[todayFromStamp].SetPayment
			saleTransactionsByVisitor += statistic[todayFromStamp].Sales
		}

		todayFrom = todayFrom.Add(time.Hour)
	}

	result["totalVisitors"] = visitorCount
	result["addedToCart"] = addToCartByVisitor
	result["visitCheckout"] = visitCheckoutByVisitor
	result["setPayment"] = setPaymentByVisitor
	result["purchased"] = saleTransactionsByVisitor

	return result, nil
}

//APIGetSales returns information about sales in the recent period, taking into account time zone
func APIGetSales(context api.InterfaceApplicationContext) (interface{}, error) {

	result := make(map[string]interface{})
	timeZone, err := app.GetSessionTimeZone(context.GetSession())
	if err != nil || timeZone == "" {
		timeZone = utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))
	}

	// get a hours pasted for local day and count for them and for previous day
	todayTo := time.Now().Truncate(time.Hour)
	todayFrom, _ := utils.MakeUTCOffsetTime(todayTo, timeZone)
	if utils.IsZeroTime(todayFrom) {
		todayFrom = todayTo
	}

	todayHoursPast := time.Duration(todayFrom.Hour()) * time.Hour
	todayFrom = todayTo.Add(-todayHoursPast)
	yesterdayFrom := todayFrom.AddDate(0, 0, -1)
	weekFrom := yesterdayFrom.AddDate(0, 0, -5)

	// get data for sales
	todayStats, err := GetRangeStats(todayFrom, todayTo)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	todaySales := todayStats.Sales
	todaySalesAmount := todayStats.SalesAmount

	yesterdayTo := todayFrom.Add(-time.Nanosecond)
	yesterdayStats, err := GetRangeStats(yesterdayFrom, yesterdayTo)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	yesterdaySales := yesterdayStats.Sales
	yesterdaySalesAmount := yesterdayStats.SalesAmount

	weekTo := yesterdayFrom.Add(-time.Nanosecond)
	weekStats, err := GetRangeStats(weekFrom, weekTo)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	weekSales := todaySales + yesterdaySales + weekStats.Sales
	weekSalesAmount := todaySalesAmount + yesterdaySalesAmount + weekStats.SalesAmount

	result["sales"] = map[string]float64{
		"today":       todaySalesAmount,
		"yesterday":   yesterdaySalesAmount,
		"week":        weekSalesAmount,
		"monthToDate": monthStatistic.SalesAmount,
	}

	result["orders"] = map[string]int{
		"today":       todaySales,
		"yesterday":   yesterdaySales,
		"week":        weekSales,
		"monthToDate": monthStatistic.Sales,
	}

	return result, nil
}

// APIGetSalesDetails returns site sales information for a specified period
//   - period start and end dates should be specified in "from" and "to" attributes in DD-MM-YYY format
func APIGetSalesDetails(context api.InterfaceApplicationContext) (interface{}, error) {

	// getting initial values
	result := make(map[string]int)
	var arrayResult [][]int

	timeZone, err := app.GetSessionTimeZone(context.GetSession())
	if err != nil || timeZone == "" {
		timeZone = utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))
	}
	dateFrom := utils.InterfaceToTime(context.GetRequestArgument("from"))
	dateTo := utils.InterfaceToTime(context.GetRequestArgument("to"))

	// checking if user specified correct from and to dates
	if dateFrom.IsZero() {
		dateFrom = time.Now().Truncate(time.Hour)
	}

	if dateTo.IsZero() {
		dateTo = time.Now().Truncate(time.Hour)
	}

	if dateFrom == dateTo {
		dateTo = dateTo.Add(ConstTimeDay)
	}

	// time zone recognize routines save time difference to show in graph by local time
	hoursOffset := time.Hour * 0

	if timeZone != "" {
		dateFrom, hoursOffset = utils.MakeUTCTime(dateFrom, timeZone)
		dateTo, _ = utils.MakeUTCTime(dateTo, timeZone)
	}

	// determining required scope
	delta := dateTo.Sub(dateFrom)

	timeScope := time.Hour
	if delta.Hours() > 48 {
		timeScope = timeScope * 24
	}
	dateFrom = dateFrom.Truncate(time.Hour)
	dateTo = dateTo.Truncate(time.Hour)

	// set database request settings
	// making database request
	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSVisitors)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorInfoCollection.AddFilter("day", ">=", dateFrom)
	visitorInfoCollection.AddFilter("day", "<=", dateTo)

	dbRecords, err := visitorInfoCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// filling requested period
	for dateFrom.Before(dateTo) {
		arrayResult = append(arrayResult, []int{utils.InterfaceToInt(dateFrom.Add(hoursOffset).Unix()), 0})
		result[fmt.Sprint(dateFrom.Add(hoursOffset).Unix())] = 0
		dateFrom = dateFrom.Add(timeScope)
	}

	// grouping database records
	for _, item := range dbRecords {
		timestamp := fmt.Sprint(utils.InterfaceToTime(item["day"]).Truncate(timeScope).Add(hoursOffset).Unix())
		subtotal := utils.InterfaceToInt(item["sales_amount"])

		if value, present := result[timestamp]; present {
			result[timestamp] = value + subtotal
		} else {
			env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e8d31584-2ef2-4f00-9510-4913b4b1d6e6", timestamp+" - not present in result"))
		}
	}

	for _, item := range arrayResult {
		item[1] = result[utils.InterfaceToString(item[0])]
	}

	return arrayResult, nil
}

// APIGetBestsellers returns information about bestsellers for some period
//     possible periods: "today", "yesterday", "week", "month"
func APIGetBestsellers(context api.InterfaceApplicationContext) (interface{}, error) {

	bestsellersRange := utils.InterfaceToString(context.GetRequestArgument("period"))

	timeZone, err := app.GetSessionTimeZone(context.GetSession())
	if err != nil || timeZone == "" {
		timeZone = utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))
	}

	// get a hours passed for local day and base from it
	todayTo := time.Now().Truncate(time.Hour).Add(time.Hour) // last hour of current day
	todayFrom, _ := utils.MakeUTCOffsetTime(todayTo, timeZone)
	if utils.IsZeroTime(todayFrom) {
		todayFrom = todayTo
	}

	todayHoursPast := time.Duration(todayFrom.Hour()) * time.Hour
	todayFrom = todayTo.Add(-todayHoursPast) // beginning of current day

	rangeFrom := todayFrom
	rangeTo := todayTo

	switch bestsellersRange {

	case "yesterday", "2":
		rangeTo = rangeFrom
		rangeFrom = rangeFrom.AddDate(0, 0, -1)
		break

	case "week", "7":
		rangeFrom = rangeFrom.AddDate(0, 0, -6)
		break

	case "month", "30":
		rangeFrom = rangeFrom.AddDate(0, 0, -30)
		break
	}

	salesHistoryCollection, err := db.GetCollection(ConstCollectionNameRTSSalesHistory)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	salesHistoryCollection.AddFilter("count", ">", 0)
	salesHistoryCollection.AddFilter("created_at", ">", rangeFrom)
	salesHistoryCollection.AddFilter("created_at", "<=", rangeTo)

	collectionRecords, err := salesHistoryCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// map  and arrays to hold sales data
	productSales := make(map[string]int)
	var productsToSort, bestSellers []map[string]interface{}

	// count the products sales by product id
	for _, item := range collectionRecords {
		pid := utils.InterfaceToString(item["product_id"])
		count := utils.InterfaceToInt(item["count"])
		productSales[pid] = count + productSales[pid]
	}

	// populate the bestseller data
	for id, count := range productSales {

		productInstance, err := product.LoadProductByID(id)
		if err != nil {
			continue
		}

		mediaPath, err := productInstance.GetMediaPath("image")
		if err != nil {
			continue
		}

		bestsellerItem := make(map[string]interface{})
		bestsellerItem["pid"] = id
		bestsellerItem["name"] = productInstance.GetName()
		bestsellerItem["count"] = count

		if productInstance.GetDefaultImage() != "" {
			bestsellerItem["image"] = mediaPath + productInstance.GetDefaultImage()
		}

		productsToSort = append(productsToSort, bestsellerItem)
	}

	// sort list of products by sales
	descending := true    // sort in descending order
	bestsellerLimit := 12 // limit on returned bestsellers
	productsSorted := utils.SortMapByKeys(productsToSort, descending, "count", "name")

	// pass back only bestsellerLimit or less
	if len(productsSorted) <= bestsellerLimit {
		bestSellers = productsSorted
	} else {
		bestSellers = productsSorted[:bestsellerLimit]
	}

	return bestSellers, nil
}

// // APIGetVisitsRealtime returns real-time information on current visits
// func APIGetVisitsRealtime(context api.InterfaceApplicationContext) (interface{}, error) {
// 	result := make(map[string]interface{})
// 	ratio := float64(0)

// 	onlineSessionCount := len(OnlineSessions)

// 	result["Online"] = onlineSessionCount
// 	if OnlineSessionsMax == 0 || onlineSessionCount == 0 {
// 		ratio = float64(0)
// 	} else {
// 		ratio = float64(onlineSessionCount) / float64(OnlineSessionsMax)
// 	}
// 	result["OnlineRatio"] = utils.Round(ratio, 0.5, 2)

// 	result["Direct"] = OnlineDirect
// 	if OnlineDirectMax == 0 || OnlineDirect == 0 {
// 		ratio = float64(0)
// 	} else {
// 		ratio = float64(OnlineDirect) / float64(OnlineDirectMax)
// 	}
// 	result["DirectRatio"] = utils.Round(ratio, 0.5, 2)

// 	result["Search"] = OnlineSearch
// 	if OnlineSearchMax == 0 || OnlineSearch == 0 {
// 		ratio = float64(0)
// 	} else {
// 		ratio = float64(OnlineSearch) / float64(OnlineSearchMax)
// 	}
// 	result["SearchRatio"] = utils.Round(ratio, 0.5, 2)

// 	result["Site"] = OnlineSite
// 	if OnlineSiteMax == 0 || OnlineSite == 0 {
// 		ratio = float64(0)
// 	} else {
// 		ratio = float64(OnlineSite) / float64(OnlineSiteMax)
// 	}
// 	result["SiteRatio"] = utils.Round(ratio, 0.5, 2)

// 	return result, nil
// }
