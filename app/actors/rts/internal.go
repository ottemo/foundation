package rts

import (
	"regexp"
	"strings"
	"time"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetReferrer returns a string when provided a URL
func GetReferrer(url string) (string, error) {
	excludeURLs := []string{app.GetFoundationURL(""), app.GetDashboardURL("")}

	r := regexp.MustCompile(`^(http|https):\/\/(.+)\/.*$`)
	groups := r.FindStringSubmatch(url)
	if len(groups) == 0 {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e9ee22d7-f62d-4379-b48e-ec9a59e388c8", "Invalid URL in referrer")
	}
	result := groups[2]

	for index := 0; index < len(excludeURLs); index++ {
		if strings.Contains(excludeURLs[index], result) {
			return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "841fa275-e0fb-4d29-868f-2bca20d5fe4e", "Invalid URL in referrer")
		}
	}

	return result, nil
}

// IncreaseOnline is a method to increase the provided counter by 1
func IncreaseOnline(typeCounter int) {
	switch typeCounter {
	case ConstReferrerTypeDirect:
		OnlineDirect++
		if OnlineDirect > OnlineDirectMax {
			OnlineDirectMax = OnlineDirect
		}
		break
	case ConstReferrerTypeSearch:
		OnlineSearch++
		if OnlineSearch > OnlineSearchMax {
			OnlineSearchMax = OnlineSearch
		}
		break
	case ConstReferrerTypeSite:
		OnlineSite++
		if OnlineSite > OnlineSiteMax {
			OnlineSiteMax = OnlineSite
		}
		break
	}
}

// DecreaseOnline is a method to decrease the provided counter by 1
func DecreaseOnline(typeCounter int) {
	switch typeCounter {
	case ConstReferrerTypeDirect:
		if OnlineDirect != 0 {
			OnlineDirect--
		}
		break
	case ConstReferrerTypeSearch:
		if OnlineSearch != 0 {
			OnlineSearch--
		}

		break
	case ConstReferrerTypeSite:
		if OnlineSite != 0 {
			OnlineSite--
		}
		break
	}
}

// GetDateFrom returns the a time.Time of last record of sales history
func GetDateFrom() (time.Time, error) {
	result := time.Now()

	salesHistoryCollection, err := db.GetCollection(ConstCollectionNameRTSSalesHistory)
	if err == nil {
		salesHistoryCollection.SetResultColumns("created_at")
		salesHistoryCollection.AddSort("created_at", true)
		salesHistoryCollection.SetLimit(0, 1)
		dbRecord, err := salesHistoryCollection.Load()
		if err != nil {
			env.LogError(err)
		}

		if len(dbRecord) > 0 {
			return utils.InterfaceToTime(dbRecord[0]["created_at"]), nil
		}
	}

	orderCollectionModel, err := order.GetOrderCollectionModel()

	if err != nil {
		return result, env.ErrorDispatch(err)
	}
	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.SetResultColumns("created_at")
	dbOrderCollection.AddSort("created_at", false)
	dbOrderCollection.SetLimit(0, 1)
	dbRecord, err := dbOrderCollection.Load()
	if err != nil {
		env.LogError(err)
	}

	if len(dbRecord) > 0 {
		return utils.InterfaceToTime(dbRecord[0]["created_at"]), nil
	}

	return result, nil
}

func initSalesHistory() error {

	// GetDateFrom return data from where need to update our rts_sales_history
	dateFrom, err := GetDateFrom()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// get orders that created after begin date
	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}
	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.SetResultColumns("_id", "created_at")
	dbOrderCollection.AddFilter("created_at", ">", dateFrom)

	ordersForPeriod, err := dbOrderCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// get order items collection
	orderItemCollectionModel, err := order.GetOrderItemCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbOrderItemCollection := orderItemCollectionModel.GetDBCollection()

	// get sales history collection
	salesHistoryCollection, err := db.GetCollection(ConstCollectionNameRTSSalesHistory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	salesHistoryData := make(map[string]map[int64]int)

	// collect data from all orders into salesHistoryData
	// in format map[pid][time]qty
	for _, order := range ordersForPeriod {

		dbOrderItemCollection.ClearFilters()
		dbOrderItemCollection.AddFilter("order_id", "=", order["_id"])
		dbOrderItemCollection.SetResultColumns("product_id", "qty")
		orderItems, err := dbOrderItemCollection.Load()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		// collect records by time with rounding top on hour basics -- all orders which are saved to sales_history
		// would be rounded on one hour up order at time 17;16 -> 18;00
		currentDateUnix := utils.InterfaceToTime(order["created_at"]).Truncate(time.Hour).Add(time.Hour).Unix()

		for _, orderItem := range orderItems {
			currentProductID := utils.InterfaceToString(orderItem["product_id"])
			count := utils.InterfaceToInt(orderItem["qty"])

			// collect data to salesHistoryData
			if productInfo, present := salesHistoryData[currentProductID]; present {
				if oldCounter, ok := productInfo[currentDateUnix]; ok {
					salesHistoryData[currentProductID][currentDateUnix] = count + oldCounter
				} else {
					salesHistoryData[currentProductID][currentDateUnix] = count
				}
			} else {
				salesHistoryData[currentProductID] = map[int64]int{currentDateUnix: count}
			}
		}
	}

	// save records to database
	for productID, productStats := range salesHistoryData {
		for orderTime, count := range productStats {

			salesRow := make(map[string]interface{})

			salesHistoryCollection.ClearFilters()
			salesHistoryCollection.AddFilter("created_at", "=", orderTime)
			salesHistoryCollection.AddFilter("product_id", "=", productID)

			dbSaleRow, err := salesHistoryCollection.Load()
			if err != nil {
				return env.ErrorDispatch(err)
			}

			if len(dbSaleRow) > 0 {
				salesRow["_id"] = utils.InterfaceToString(dbSaleRow[0]["_id"])
				count = count + utils.InterfaceToInt(dbSaleRow[0]["count"])
			}

			salesRow["created_at"] = orderTime
			salesRow["product_id"] = productID
			salesRow["count"] = count
			_, err = salesHistoryCollection.Save(salesRow)
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}

	return nil
}

// GetRangeStats returns stats for range
func GetRangeStats(dateFrom, dateTo time.Time) (ActionsMade, error) {

	var stats ActionsMade

	// Go thru period and summarise a visits
	for dateFrom.Before(dateTo.Add(time.Nanosecond)) {
		if actions, present := statistic[dateFrom.Unix()]; present {
			stats.Visit = actions.Visit + stats.Visit
			stats.Sales = actions.Sales + stats.Sales
			stats.Cart = actions.Cart + stats.Cart
			stats.TotalVisits = actions.TotalVisits + stats.TotalVisits
			stats.SalesAmount = actions.SalesAmount + stats.SalesAmount
		}

		dateFrom = dateFrom.Add(time.Hour)
	}
	return stats, nil
}

// initStatistic get info from visitor database for 60 hours
func initStatistic() error {
	// convert to utc time and work with variables
	timeScope := time.Hour
	durationWeek := time.Hour * 168

	dateTo := time.Now().Truncate(timeScope)
	dateFrom := dateTo.Add(-durationWeek)

	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	visitorInfoCollection.AddFilter("day", "<=", dateTo)
	visitorInfoCollection.AddFilter("day", ">=", dateFrom)

	dbRecords, err := visitorInfoCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	timeIterator := dateFrom.Unix()

	// add info from db record if not null to variable
	for _, item := range dbRecords {
		timeIterator = utils.InterfaceToTime(item["day"]).Unix()
		if _, present := statistic[timeIterator]; !present {
			statistic[timeIterator] = new(ActionsMade)
		}
		// add info to hour
		statistic[timeIterator].TotalVisits = statistic[timeIterator].TotalVisits + utils.InterfaceToInt(item["total_visits"])
		statistic[timeIterator].SalesAmount = statistic[timeIterator].SalesAmount + utils.InterfaceToFloat64(item["sales_amount"])
		statistic[timeIterator].Visit = statistic[timeIterator].Visit + utils.InterfaceToInt(item["visitors"])
		statistic[timeIterator].Sales = statistic[timeIterator].Sales + utils.InterfaceToInt(item["sales"])
		statistic[timeIterator].VisitCheckout = statistic[timeIterator].VisitCheckout + utils.InterfaceToInt(item["visit_checkout"])
		statistic[timeIterator].SetPayment = statistic[timeIterator].SetPayment + utils.InterfaceToInt(item["set_payment"])
		statistic[timeIterator].Cart = statistic[timeIterator].Cart + utils.InterfaceToInt(item["cart"])
	}

	return nil
}

// SaveStatisticsData save a statistic data row gor last hour to database
func SaveStatisticsData() error {
	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	currentHour := time.Now().Truncate(time.Hour)

	// find last saved record time to start saving from it
	visitorInfoCollection.AddFilter("day", "=", currentHour)
	dbRecord, err := visitorInfoCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	visitorInfoRow := make(map[string]interface{})

	// write current records to database with rewrite of last
	if len(dbRecord) > 0 {
		visitorInfoRow = utils.InterfaceToMap(dbRecord[0])
	}

	if lastActions, present := statistic[currentHour.Unix()]; present {
		visitorInfoRow["day"] = currentHour
		visitorInfoRow["visitors"] = lastActions.Visit
		visitorInfoRow["cart"] = lastActions.Cart
		visitorInfoRow["sales"] = lastActions.Sales
		visitorInfoRow["visit_checkout"] = lastActions.VisitCheckout
		visitorInfoRow["set_payment"] = lastActions.SetPayment
		visitorInfoRow["sales_amount"] = lastActions.SalesAmount
		visitorInfoRow["total_visits"] = lastActions.TotalVisits

		// save data to database
		_, err = visitorInfoCollection.Save(visitorInfoRow)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	} else {
		env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9712c601-662e-4744-b9fb-991a959cff32", "last rts_visitors statistic save failed"))
	}

	return nil
}

// CheckHourUpdateForStatistic if it's a new hour action we need renew all session as a new in this hour
// and remove old record from statistic
func CheckHourUpdateForStatistic() {
	currentHour := time.Now().Truncate(time.Hour).Unix()
	durationWeek := time.Hour * 168

	lastHour := time.Now().Add(-durationWeek).Truncate(time.Hour).Unix()

	// if last our not present in statistic we need to update visitState
	// if it's a new day so we make clear a visitor state stats
	// and create clear record for this hour
	if _, present := statistic[currentHour]; !present {

		if lastUpdate.Truncate(time.Hour*24) != time.Now().Truncate(time.Hour*24) {
			visitState = make(map[string]bool)
		} else {
			cartCreatedPersons := make(map[string]bool)

			for sessionID, addToCartPresent := range visitState {
				if addToCartPresent {
					cartCreatedPersons[sessionID] = addToCartPresent
				}
			}
			visitState = cartCreatedPersons
		}
		statistic[currentHour] = new(ActionsMade)
	}

	for timeIn := range statistic {
		if timeIn < lastHour {
			delete(statistic, timeIn)
		}
	}

	lastUpdate = time.Now()
}

// saveNewReferrer make save a new referral to data base
func saveNewReferrer(referral string) error {
	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSReferrals)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// rewrite existing referral with new count
	visitorInfoCollection.AddFilter("referral", "=", referral)
	visitorInfoCollection.SetLimit(0, 1)
	dbRecord, err := visitorInfoCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	newRecord := make(map[string]interface{})

	if len(dbRecord) > 0 {
		newRecord["_id"] = dbRecord[0]["_id"]
	}
	newRecord["referral"] = referral
	newRecord["count"] = referrers[referral]

	// save data to database
	_, err = visitorInfoCollection.Save(newRecord)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// initReferrals get info from referrals database to variable
func initReferrals() error {

	rtsReferralsCollection, err := db.GetCollection(ConstCollectionNameRTSReferrals)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecords, err := rtsReferralsCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, record := range dbRecords {
		referrers[utils.InterfaceToString(record["referral"])] = utils.InterfaceToInt(record["count"])
	}

	return nil
}

// sortArrayOfMapByKey sort array from biggest to lowest value of map[key] element
func sortArrayOfMapByKey(data []map[string]interface{}, key string) []map[string]interface{} {

	var result []map[string]interface{}
	var indexOfMaxValueItem int
	var maxValue float64

	for len(data) > 0 {
		for index, item := range data {
			if utils.InterfaceToFloat64(item[key]) > maxValue {
				maxValue = utils.InterfaceToFloat64(item[key])
				indexOfMaxValueItem = index
			}
		}
		result = append(result, data[indexOfMaxValueItem])
		data = append(data[:indexOfMaxValueItem], data[indexOfMaxValueItem+1:]...)
		maxValue = 0
	}
	return result
}
