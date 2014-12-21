package rts

import (
	"encoding/json"
	"fmt"
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
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e9ee22d7f62d4379b48eec9a59e388c8", "Invalid URL in referrer")
	}
	result := groups[2]

	for index := 0; index < len(excludeURLs); index++ {
		if strings.Contains(excludeURLs[index], result) {
			return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "841fa275e0fb4d29868f2bca20d5fe4e", "Invalid URL in referrer")
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
			year := datetime.Year()
			month := datetime.Month()
			day := datetime.Day()
			result := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
			if time.Now().Format("2006-01-02") == datetime.Format("2006-01-02") {
				return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ead00ed31d1e45e7b3305dcd73c88764", "Sales history has last data")
			}

			return result, nil
		}
	}

	orderCollectionModel, err := order.GetOrderCollectionModel()
	dateFrom := time.Now()
	if err != nil {
		year := dateFrom.Year()
		month := dateFrom.Month()
		day := dateFrom.Day()
		result := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		return result, nil
	}
	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.SetResultColumns("created_at")
	dbOrderCollection.AddSort("created_at", false)
	dbRecord, err := dbOrderCollection.Load()
	if len(dbRecord) > 0 {
		dateFrom = utils.InterfaceToTime(dbRecord[0]["created_at"])
		year := dateFrom.Year()
		month := dateFrom.Month()
		day := dateFrom.Day()
		result := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
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

	year := date.Year()
	month := date.Month()
	day := date.Day()
	todayFrom := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	todayTo := time.Date(year, month, day, 23, 59, 59, 0, time.Local)
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
			salesHistoryRow := make(map[string]interface{})
			salesHistoryRow["product_id"] = productID
			salesHistoryRow["created_at"] = date
			salesHistoryRow["count"] = count
			_, err = salesHistoryCollection.Save(salesHistoryRow)
			if err != nil {
				return env.ErrorDispatch(err)
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
	_range := "2014-01-01:"

	return _range
}

// GetTotalSales will create the totals for current sales data
func GetTotalSales(fromDate time.Time, toDate time.Time) error {

	orderCollectionModelT, err := order.GetOrderCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbCollection := orderCollectionModelT.GetDBCollection()

	year := fromDate.Year()
	month := fromDate.Month()
	day := fromDate.Day()
	todayFrom := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	todayTo := time.Date(year, month, day, 23, 59, 59, 0, time.Local)

	dbCollection.AddFilter("created_at", ">=", todayFrom)
	dbCollection.AddFilter("created_at", "<=", todayTo)

	// filters handle for today
	today, err := dbCollection.Count()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbCollection.ClearFilters()
	year = toDate.Year()
	month = toDate.Month()
	day = toDate.Day()
	yesterdayFrom := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	yesterdayTo := time.Date(year, month, day, 23, 59, 59, 0, time.Local)

	dbCollection.AddFilter("created_at", ">=", yesterdayFrom)
	dbCollection.AddFilter("created_at", "<=", yesterdayTo)

	// filters handle for yesterday
	yesterday, err := dbCollection.Count()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	ratio := float64(1)
	if 0 != yesterday {
		ratio = float64(today)/float64(yesterday) - float64(1)
	}

	sales.ratio = ratio
	sales.today = today
	sales.lastUpdate = time.Now().Unix()
	sales.yesterday = yesterday

	return nil
}

// GetSalesDetail will return the sale data for the given time period
func GetSalesDetail(fromDate time.Time, toDate time.Time, hash string) error {

	orderCollectionModelT, err := order.GetOrderCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbCollection := orderCollectionModelT.GetDBCollection()
	dbCollection.SetResultColumns("_id", "created_at")
	dbCollection.AddSort("created_at", false)

	year := fromDate.Year()
	month := fromDate.Month()
	day := fromDate.Day()
	dateFrom := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	dbCollection.AddFilter("created_at", ">=", dateFrom)

	year = toDate.Year()
	month = toDate.Month()
	day = toDate.Day()
	dateTo := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	dbCollection.AddFilter("created_at", "<=", dateTo)

	// filters handle for yesterday
	list, err := dbCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	delta := toDate.Sub(fromDate)
	salesDetail[hash] = &SalesDetailData{Data: make(map[string]int)}
	if delta.Hours() > 48 { // group by days
		// fills the data a zero
		for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {
			timestamp := int64(date.Unix())
			mapIndex := GetDayForTimestamp(timestamp, false)

			salesDetail[hash].Data[mapIndex] = 0
		}
		// counts items
		for _, order := range list {
			timestamp := int64(utils.InterfaceToTime(order["created_at"]).Unix())
			mapIndex := GetDayForTimestamp(timestamp, false)

			salesDetail[hash].Data[mapIndex]++
		}
	} else { // group by hours

		for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {
			timestamp := int64(date.Unix())
			year := time.Unix(int64(timestamp), 0).Year()
			month := time.Unix(int64(timestamp), 0).Month()
			day := time.Unix(int64(timestamp), 0).Day()
			for hour := 0; hour < 24; hour++ {
				timeGroup := time.Date(year, month, day, hour, 0, 0, 0, time.Local)
				if timeGroup.Unix() > time.Now().Unix() {
					break
				}
				mapIndex := fmt.Sprintf("%v", int32(timeGroup.Unix()))
				salesDetail[hash].Data[mapIndex] = 0
			}
		}
		for _, order := range list {
			timestamp := int64(utils.InterfaceToTime(order["created_at"]).Unix())
			mapIndex := GetDayForTimestamp(timestamp, true)

			salesDetail[hash].Data[mapIndex]++
		}
	}

	salesDetail[hash].lastUpdate = time.Now().Unix()

	return nil
}

// GetDayForTimestamp returns the day or hour for the given time
func GetDayForTimestamp(timestamp int64, byHour bool) string {
	hour := 0
	if byHour {
		hour = time.Unix(timestamp, 0).Hour()
	}

	year := time.Unix(timestamp, 0).Year()
	month := time.Unix(timestamp, 0).Month()
	day := time.Unix(timestamp, 0).Day()
	timeGroup := time.Date(year, month, day, hour, 0, 0, 0, time.Local)
	mapIndex := fmt.Sprintf("%v", int32(timeGroup.Unix()))

	return mapIndex
}

// GetTodayVisitorsData will return Visitor data for Today
func GetTodayVisitorsData() error {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	if visitorsInfoToday.Day == today.In(time.Local) {
		return nil
	}

	year = time.Now().AddDate(0, 0, -1).Year()
	month = time.Now().AddDate(0, 0, -1).Month()
	day = time.Now().AddDate(0, 0, -1).Day()
	yesterday := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if visitorsInfoToday.Day == yesterday {
		SaveVisitorData()
		visitorsInfoYesterday = visitorsInfoToday
	}
	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSVisitors)
	if err == nil {
		visitorInfoCollection.AddFilter("day", "=", today)
		dbRecord, _ := visitorInfoCollection.Load()

		if len(dbRecord) > 0 {
			visitorsInfoToday.ID = utils.InterfaceToString(dbRecord[0]["_id"])
			visitorsInfoToday.Day = utils.InterfaceToTime(dbRecord[0]["day"])
			visitorsInfoToday.Visitors = utils.InterfaceToInt(dbRecord[0]["visitors"])
			visitorsInfoToday.Cart = utils.InterfaceToInt(dbRecord[0]["cart"])
			visitorsInfoToday.Checkout = utils.InterfaceToInt(dbRecord[0]["checkout"])
			visitorsInfoToday.Sales = utils.InterfaceToInt(dbRecord[0]["sales"])
			visitorsInfoToday.Details = DecodeDetails(utils.InterfaceToString(dbRecord[0]["details"]))

			return nil
		}
	}

	visitorsInfoToday = new(dbVisitorsRow)
	visitorsInfoToday.ID = ""
	visitorsInfoToday.Day = today
	visitorsInfoToday.Details = make(map[string]*VisitorDetail)

	return nil
}

// GetYesterdayVisitorsData will build a collection of data representing yesterdays Visitor statistics
func GetYesterdayVisitorsData() error {
	year := time.Now().AddDate(0, 0, -1).Year()
	month := time.Now().AddDate(0, 0, -1).Month()
	day := time.Now().AddDate(0, 0, -1).Day()
	yesterday := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	if visitorsInfoYesterday.Day == yesterday {
		return nil
	}

	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSVisitors)
	if err == nil {
		visitorInfoCollection.AddFilter("day", "=", yesterday)
		dbRecord, _ := visitorInfoCollection.Load()

		if len(dbRecord) > 0 {
			visitorsInfoYesterday.ID = utils.InterfaceToString(dbRecord[0]["_id"])
			visitorsInfoYesterday.Day = utils.InterfaceToTime(dbRecord[0]["day"])
			visitorsInfoYesterday.Visitors = utils.InterfaceToInt(dbRecord[0]["visitors"])
			visitorsInfoYesterday.Cart = utils.InterfaceToInt(dbRecord[0]["cart"])
			visitorsInfoYesterday.Checkout = utils.InterfaceToInt(dbRecord[0]["checkout"])
			visitorsInfoYesterday.Sales = utils.InterfaceToInt(dbRecord[0]["sales"])
			visitorsInfoYesterday.Details = DecodeDetails(utils.InterfaceToString(dbRecord[0]["details"]))

			return nil
		}
	}

	visitorsInfoYesterday = new(dbVisitorsRow)
	visitorsInfoYesterday.ID = ""
	visitorsInfoYesterday.Day = yesterday
	visitorsInfoYesterday.Details = make(map[string]*VisitorDetail)

	return nil
}

// SaveVisitorData will persist the Visitor data to the database
func SaveVisitorData() error {
	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSVisitors)
	if err == nil {
		visitorInfoRow := make(map[string]interface{})
		if "" != visitorsInfoToday.ID {
			visitorInfoRow["_id"] = visitorsInfoToday.ID
		}

		visitorInfoRow["day"] = visitorsInfoToday.Day
		visitorInfoRow["visitors"] = visitorsInfoToday.Visitors
		visitorInfoRow["cart"] = visitorsInfoToday.Cart
		visitorInfoRow["checkout"] = visitorsInfoToday.Checkout
		visitorInfoRow["sales"] = visitorsInfoToday.Sales
		visitorInfoRow["details"] = EncodeDetails(visitorsInfoToday.Details)

		_, err = visitorInfoCollection.Save(visitorInfoRow)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// EncodeDetails returns the Visitor data in a string when provided a VisitorDetail map[string]*
func EncodeDetails(details map[string]*VisitorDetail) string {
	jsonString, _ := json.Marshal(details)

	return string(jsonString)
}

// DecodeDetails returns the Visitor data in a VisitorDetail map[string]* when provieded an encoded string
func DecodeDetails(detailsString string) map[string]*VisitorDetail {
	var details map[string]*VisitorDetail
	_ = json.Unmarshal([]byte(detailsString), &details)

	return details
}
