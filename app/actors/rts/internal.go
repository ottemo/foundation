package rts

import (
	"regexp"
	"strings"
	"time"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/product"
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

// GetProducts takes a []map[string]interface as input and returns a list of products
func GetProducts() ([]map[string]interface{}, error) {
	productCollectionModel, err := product.GetProductCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbProductCollection := productCollectionModel.GetDBCollection()
	dbProductCollection.SetResultColumns("_id")
	return dbProductCollection.Load()
}

// GetDateFrom returns the a time.Time object for the given sales record
func GetDateFrom() (time.Time, error) {
	result := time.Now()

	salesHistoryCollection, err := db.GetCollection(ConstCollectionNameRTSSalesHistory)
	if err == nil {
		salesHistoryCollection.SetResultColumns("created_at")
		salesHistoryCollection.AddSort("created_at", true)
		salesHistoryCollection.SetLimit(0, 1)
		dbRecord, _ := salesHistoryCollection.Load()

		if len(dbRecord) > 0 {
			datetime := utils.InterfaceToTime(dbRecord[0]["created_at"])
			result := time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 0, 0, 0, 0, datetime.Location())
			if time.Now().Format("2006-01-02") == datetime.Format("2006-01-02") {
				return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ead00ed3-1d1e-45e7-b330-5dcd73c88764", "Sales history has last data")
			}

			return result, nil
		}
	}

	orderCollectionModel, err := order.GetOrderCollectionModel()
	dateFrom := time.Now()
	if err != nil {
		result := time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), dateFrom.Hour(), 0, 0, 0, dateFrom.Location())
		return result, nil
	}
	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.SetResultColumns("created_at")
	dbOrderCollection.AddSort("created_at", false)
	dbRecord, err := dbOrderCollection.Load()
	if len(dbRecord) > 0 {
		dateFrom = utils.InterfaceToTime(dbRecord[0]["created_at"])
		result := time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), dateFrom.Hour(), 0, 0, 0, dateFrom.Location())
		return result, nil
	}

	return result, nil
}

// GetOrderItems returns a []map[string]interface for the given product
func GetOrderItems(date time.Time, productID string) ([]map[string]interface{}, error) {

	orderItemCollectionModel, err := order.GetOrderItemCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dbOrderItemCollection := orderItemCollectionModel.GetDBCollection()
	dbOrderItemCollection.SetResultColumns("qty")

	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.SetResultColumns("_id")

	todayFrom := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	todayTo := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
	dbOrderCollection.AddFilter("created_at", ">=", todayFrom)
	dbOrderCollection.AddFilter("created_at", "<=", todayTo)

	dbOrderItemCollection.AddFilter("order_id", "IN", dbOrderCollection)
	dbOrderItemCollection.AddFilter("product_id", "=", productID)

	return dbOrderItemCollection.Load()
}

// DeleteExistingRowHistory will remove the row entry for the given productID
func DeleteExistingRowHistory(date time.Time, productID string) error {
	salesHistoryCollection, err := db.GetCollection(ConstCollectionNameRTSSalesHistory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	salesHistoryCollection.ClearFilters()
	salesHistoryCollection.AddFilter("product_id", "=", productID)
	salesHistoryCollection.AddFilter("created_at", "=", date)
	dbSalesHist, _ := salesHistoryCollection.Load()
	if len(dbSalesHist) > 0 {
		err = salesHistoryCollection.DeleteByID(utils.InterfaceToString(dbSalesHist[0]["_id"]))
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

func initSalesHistory() error {
	sales := make(map[string]int)
	dateTo := time.Now()

	dbProductRecord, _ := GetProducts()
	dateFrom, err := GetDateFrom()
	if err != nil {
		return nil
	}

	salesHistoryCollection, err := db.GetCollection(ConstCollectionNameRTSSalesHistory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for date := dateFrom; int32(date.Unix()) < int32(dateTo.Unix()); date = date.AddDate(0, 0, 1) {
		for _, productItem := range dbProductRecord {
			productID := utils.InterfaceToString(productItem["_id"])
			DeleteExistingRowHistory(date, productID)
			count := 0

			items, _ := GetOrderItems(date, productID)
			for _, item := range items {
				count += utils.InterfaceToInt(item["qty"])
			}

			// Add history row
			if count > 0 {
				salesHistoryRow := make(map[string]interface{})
				salesHistoryRow["product_id"] = productID
				salesHistoryRow["created_at"] = date
				salesHistoryRow["count"] = count
				_, err = salesHistoryCollection.Save(salesHistoryRow)
				if err != nil {
					return env.ErrorDispatch(err)
				}
			}
			sales[productID] += count
		}
		SaveSalesData(sales)
		sales = make(map[string]int)
	}

	return nil
}

// SaveSalesData will persist the given map[string]int representing sales data
func SaveSalesData(data map[string]int) error {

	if len(data) == 0 {
		return nil
	}

	salesCollection, err := db.GetCollection(ConstCollectionNameRTSSales)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for productID, count := range data {
		// Add history row
		salesRow := make(map[string]interface{})

		salesCollection.ClearFilters()
		salesCollection.AddFilter("range", "=", GetSalesRange())
		salesCollection.AddFilter("product_id", "=", productID)
		dbSaleRow, _ := salesCollection.Load()
		if len(dbSaleRow) > 0 {
			salesRow["_id"] = utils.InterfaceToString(dbSaleRow[0]["_id"])
			oldCount := utils.InterfaceToInt(dbSaleRow[0]["count"])
			count += oldCount
		}

		salesRow["product_id"] = productID
		salesRow["count"] = count
		salesRow["range"] = GetSalesRange()
		_, err = salesCollection.Save(salesRow)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// GetSalesRange will return the date range for the sales data
func GetSalesRange() string {
	_range := "2015-01-01:"

	return _range
}

// GetRangeSales returns Get Sales for range
func GetRangeSales(dateFrom, dateTo time.Time) (int, error) {

	sales := 0

	// Go thrue period and summarise a visits
	for dateFrom.Before(dateTo) {

		if _, present := statistic[dateFrom.Unix()]; present {
			sales = sales + statistic[dateFrom.Unix()].Sales
		}

		dateFrom = dateFrom.Add(time.Hour)
	}

	return sales, nil
}

// GetRangeVisits get visits for determinated range
func GetRangeVisits(dateFrom, dateTo time.Time) (int, error) {

	visits := 0

	// Go thrue period and summarise a visits
	for dateFrom.Before(dateTo) {

		if _, present := statistic[dateFrom.Unix()]; present {
			visits = visits + statistic[dateFrom.Unix()].Visit
		}

		dateFrom = dateFrom.Add(time.Hour)
	}

	return visits, nil
}

// initStatistic get info from visitor database for 60 hours
func initStatistic() error {
	// convert to utc time and work with variables
	time.Local = time.UTC
	timeScope := time.Hour
	dateTo := time.Now().Add(time.Hour).Truncate(timeScope)
	dateFrom := dateTo.Add(-60 * time.Hour)

	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	visitorInfoCollection.SetResultColumns("day", "visitors", "cart", "sales")
	visitorInfoCollection.AddSort("day", false)

	for dateFrom.Before(dateTo) {

		timeIterator := dateFrom.Unix()
		// get database records for every hour
		visitorInfoCollection.ClearFilters()
		visitorInfoCollection.AddFilter("day", "=", timeIterator)
		dbRecords, err := visitorInfoCollection.Load()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		// add info from db record if not null to variable
		for _, item := range dbRecords {

			// create record for non existing hour
			if _, present := statistic[timeIterator]; !present {
				statistic[timeIterator] = new(ActionsMade)
			}

			// add info to hour
			statistic[timeIterator].Visit = statistic[timeIterator].Visit + utils.InterfaceToInt(item["visitors"])
			statistic[timeIterator].Sales = statistic[timeIterator].Sales + utils.InterfaceToInt(item["sales"])
			statistic[timeIterator].Cart = statistic[timeIterator].Cart + utils.InterfaceToInt(item["cart"])

		}

		dateFrom = dateFrom.Add(timeScope)
	}

	return nil
}

// SaveStatisticsData make save a statistic data row to database from last updated record in database to current hour
func SaveStatisticsData() error {
	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// find last saved record time to start saving from it
	visitorInfoCollection.SetResultColumns("day")
	visitorInfoCollection.AddSort("day", true)
	visitorInfoCollection.SetLimit(0, 1)
	dbRecord, err := visitorInfoCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}
	if len(dbRecord) == 0 {
		return nil
	}

	lastRecordTime := utils.InterfaceToTime(dbRecord[0]["day"])

	// delete last record from database to prevent duplicates
	if _, present := statistic[lastRecordTime.Unix()]; present {
		visitorInfoCollection.ClearFilters()
		visitorInfoCollection.ClearSort()
		visitorInfoCollection.SetLimit(0, 100)
		visitorInfoCollection.SetResultColumns("day", "_id")
		visitorInfoCollection.AddFilter("day", "=", lastRecordTime.Unix())
		dbLastRecord, err := visitorInfoCollection.Load()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		// deleting all records with the same hour
		for _, item := range dbLastRecord {
			err = visitorInfoCollection.DeleteByID(utils.InterfaceToString(item["_id"]))
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}

	}

	// begin process of writing to database
	visitorInfoCollection.ClearFilters()
	visitorInfoRow := make(map[string]interface{})

	// to be sure that all time is in hour format
	dateTo := time.Now().Add(time.Hour).Truncate(time.Hour)
	dateFrom := lastRecordTime.Truncate(time.Hour)

	// save data to database for every hour that data for present in statistic
	// beginning from last database record to current time
	for dateFrom.Before(dateTo) {
		if value, present := statistic[dateFrom.Unix()]; present {
			visitorInfoRow["day"] = dateFrom.Unix()
			visitorInfoRow["visitors"] = value.Visit
			visitorInfoRow["cart"] = value.Cart
			visitorInfoRow["sales"] = value.Sales

			// save data to database
			_, err = visitorInfoCollection.Save(visitorInfoRow)
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}

		dateFrom = dateFrom.Add(time.Hour)
	}

	return nil
}

// CheckHourUpdateForStatistic if it's a new hour action we need renew all session as a new in this hour
// and remove old record from statistic
func CheckHourUpdateForStatistic() {
	currentHour := time.Now().Truncate(time.Hour).Unix()
	lastHour := time.Now().Add(-time.Hour * 60).Truncate(time.Hour).Unix()

	if _, present := statistic[currentHour]; !present {

		cartCreatedPersons := make(map[string]bool)
		for sessionID, addToCartPresent := range visitState {
			if addToCartPresent {
				cartCreatedPersons[sessionID] = addToCartPresent
			}
		}

		visitState = cartCreatedPersons
		statistic[currentHour] = new(ActionsMade)
	}

	for timeIn := range statistic {
		if timeIn <= lastHour {
			delete(statistic, timeIn)
		}
	}

}
