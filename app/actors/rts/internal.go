package rts

import (
	"time"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetDateFrom returns the a time.Time of last record of sales history
func GetDateFrom() (time.Time, error) {
	result := time.Now()

	salesHistoryCollection, err := db.GetCollection(ConstCollectionNameRTSSalesHistory)
	if err == nil {
		if err := salesHistoryCollection.SetResultColumns("created_at"); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ed738eda-ddcb-4beb-b45c-98cf6c7676f1", err.Error())
		}
		if err := salesHistoryCollection.AddSort("created_at", true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "aeac1ed0-a4ac-474e-8b06-853ce0518c4c", err.Error())
		}
		if err := salesHistoryCollection.SetLimit(0, 1); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8e419217-c040-4eac-8a97-9cf84863d490", err.Error())
		}
		dbRecord, err := salesHistoryCollection.Load()
		if err != nil {
			_ = env.ErrorDispatch(err)
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
	if err := dbOrderCollection.SetResultColumns("created_at"); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1cc773ac-9e20-43ca-bd9b-75597b935d6b", err.Error())
	}
	if err := dbOrderCollection.AddSort("created_at", false); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "54966243-6ffb-4e56-9622-d74aeb6a494e", err.Error())
	}
	if err := dbOrderCollection.SetLimit(0, 1); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d5b6cce3-f1e4-41f5-a391-9165fadd0d87", err.Error())
	}
	dbRecord, err := dbOrderCollection.Load()
	if err != nil {
		_ = env.ErrorDispatch(err)
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
	if err := dbOrderCollection.SetResultColumns("_id", "created_at"); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6dc6ad78-156c-43a2-ac5e-b78af7445e63", err.Error())
	}
	if err := dbOrderCollection.AddFilter("created_at", ">", dateFrom); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d6e71ebd-4569-4379-89e7-b1e5590ae472", err.Error())
	}

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

		if err := dbOrderItemCollection.ClearFilters(); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f3d10114-dd6c-45e2-b385-0e92913a8fcd", err.Error())
		}
		if err := dbOrderItemCollection.AddFilter("order_id", "=", order["_id"]); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fc15ac2d-7d3f-466b-8450-4fad508e714e", err.Error())
		}
		if err := dbOrderItemCollection.SetResultColumns("product_id", "qty"); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2d4b396c-b975-46ba-ac7d-5b227b8ae4ba", err.Error())
		}
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

			if err := salesHistoryCollection.ClearFilters(); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "47d706c8-dee9-4777-92d9-e325a55815bd", err.Error())
			}
			if err := salesHistoryCollection.AddFilter("created_at", "=", orderTime); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4ae240f2-1ca7-4f3b-b6cf-813693861221", err.Error())
			}
			if err := salesHistoryCollection.AddFilter("product_id", "=", productID); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "74551cab-8227-48e0-ba4a-96741d136550", err.Error())
			}

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
	var result ActionsMade

	// making minimal offset to include dateTo timestamp,
	// dateFrom timestamp includes by default in time.Before() function
	dateTo = dateTo.Add(time.Nanosecond)

	// Go through period and summarise counters
	for dateFrom.Before(dateTo) {
		timestamp := dateFrom.Unix()
		if statisticValue, present := statistic[timestamp]; present {
			result.Visit += statisticValue.Visit
			result.Sales += statisticValue.Sales
			result.Cart += statisticValue.Cart
			result.TotalVisits += statisticValue.TotalVisits
			result.SalesAmount += statisticValue.SalesAmount
		}

		dateFrom = dateFrom.Add(time.Hour)
	}
	return result, nil
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

	if err := visitorInfoCollection.AddFilter("day", "<=", dateTo); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d6a6bead-3367-451f-b37b-8a833dfaa9c3", err.Error())
	}
	if err := visitorInfoCollection.AddFilter("day", ">=", dateFrom); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ce11e113-2184-4c9d-a6e5-c09f503d4620", err.Error())
	}

	dbRecords, err := visitorInfoCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	timeIterator := dateFrom.Unix()

	// add info from db record if not null to variable
	for _, item := range dbRecords {
		timeIterator = utils.InterfaceToTime(item["day"]).Unix()
		if _, present := statistic[timeIterator]; !present {
			updateSync.Lock()
			statistic[timeIterator] = &ActionsMade{}
			updateSync.Unlock()
		}
		// add info to hour
		statistic[timeIterator].TotalVisits += utils.InterfaceToInt(item["total_visits"])
		statistic[timeIterator].SalesAmount += utils.InterfaceToFloat64(item["sales_amount"])
		statistic[timeIterator].Visit += utils.InterfaceToInt(item["visitors"])
		statistic[timeIterator].Sales += utils.InterfaceToInt(item["sales"])
		statistic[timeIterator].VisitCheckout += utils.InterfaceToInt(item["visit_checkout"])
		statistic[timeIterator].SetPayment += utils.InterfaceToInt(item["set_payment"])
		statistic[timeIterator].Cart += utils.InterfaceToInt(item["cart"])
	}

	dateTo = time.Now()
	// beginning of current month
	dateFrom = time.Date(dateTo.Year(), dateTo.Month(), 0, 0, 0, 0, 0, time.UTC)

	if err := visitorInfoCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8e0ced5a-827d-486a-bb3f-9af6e15569e1", err.Error())
	}
	if err := visitorInfoCollection.AddFilter("day", "<", dateTo); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "716a7108-eec0-4f51-91c1-6bada62f20f9", err.Error())
	}
	if err := visitorInfoCollection.AddFilter("day", ">=", dateFrom); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7b655d07-b069-4ce7-870f-5640eb8f7562", err.Error())
	}

	dbRecords, err = visitorInfoCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, item := range dbRecords {
		monthStatistic.TotalVisits += utils.InterfaceToInt(item["total_visits"])
		monthStatistic.SalesAmount += utils.InterfaceToFloat64(item["sales_amount"])
		monthStatistic.Visit += utils.InterfaceToInt(item["visitors"])
		monthStatistic.Sales += utils.InterfaceToInt(item["sales"])
		monthStatistic.VisitCheckout += utils.InterfaceToInt(item["visit_checkout"])
		monthStatistic.SetPayment += utils.InterfaceToInt(item["set_payment"])
		monthStatistic.Cart += utils.InterfaceToInt(item["cart"])

	}

	return nil
}

// SaveStatisticsData save a statistic data row for last hour to database
func SaveStatisticsData() error {
	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	currentHour := time.Now().Truncate(time.Hour)

	// find last saved record time to start saving from it
	if err := visitorInfoCollection.AddFilter("day", "=", currentHour); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e80bdc72-a8d2-4984-bbf4-b400dd68264d", err.Error())
	}
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
		env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9712c601-662e-4744-b9fb-991a959cff32", "key "+currentHour.String()+" not present in memory statistic value"))
	}

	return nil
}

// CheckHourUpdateForStatistic if it's a new hour action we need renew all session as a new in this hour
// and remove old record from statistic
func CheckHourUpdateForStatistic() {
	currentTime := time.Now()
	currentHour := currentTime.Truncate(time.Hour).Unix()
	durationWeek := time.Hour * 168

	lastHour := time.Now().Add(-durationWeek).Truncate(time.Hour).Unix()

	timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))
	currentServerTime, _ := utils.MakeTZTime(currentTime, timeZone)
	lastServerTime, _ := utils.MakeTZTime(lastUpdate, timeZone)

	if currentServerTime.Month() > lastServerTime.Month() {
		monthStatistic.Visit = 0
		monthStatistic.Cart = 0
		monthStatistic.Sales = 0
		monthStatistic.TotalVisits = 0
		monthStatistic.SalesAmount = 0
		monthStatistic.VisitCheckout = 0
		monthStatistic.SetPayment = 0
	}

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

		updateSync.Lock()
		statistic[currentHour] = &ActionsMade{}
		updateSync.Unlock()
	}

	updateSync.Lock()
	for timeIn := range statistic {
		if timeIn < lastHour {
			delete(statistic, timeIn)
		}
	}
	updateSync.Unlock()

	lastUpdate = time.Now()
}
