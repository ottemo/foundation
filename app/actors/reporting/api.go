package reporting

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/actors/discount/giftcard"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	// Admin only endpoint
	service := api.GetRestService()

	service.GET("reporting/product-performance", api.IsAdminHandler(listProductPerformance))
	service.GET("reporting/customer-activity", api.IsAdminHandler(listCustomerActivity))
	service.GET("reporting/payment-method", api.IsAdminHandler(listPaymentMethod))
	service.GET("reporting/shipping-method", api.IsAdminHandler(listShippingMethod))
	service.GET("reporting/location-country", api.IsAdminHandler(listLocationCountry))
	service.GET("reporting/location-us", api.IsAdminHandler(listLocationUS))
	service.GET("reporting/gift-cards", api.IsAdminHandler(listGiftCards))

	return nil
}

// listProductPerformance Handler that returns product performance information by date range
func listProductPerformance(context api.InterfaceApplicationContext) (interface{}, error) {

	// Expecting dates in UTC, and adjusted for your timezone
	// `2006-01-02 15:04`
	startDate := utils.InterfaceToTime(context.GetRequestArgument("start_date"))
	endDate := utils.InterfaceToTime(context.GetRequestArgument("end_date"))
	err := ValidateStartAndEndDate(startDate, endDate)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	foundOrders := order.GetOrdersCreatedBetween(startDate, endDate)
	foundOrderIds := getOrderIds(foundOrders)
	foundOrderItems := order.GetItemsForOrders(foundOrderIds)
	aggregatedResults := aggregateOrderItems(foundOrderItems)

	totalSales := getTotalSales(foundOrderItems)

	response := map[string]interface{}{
		"total_orders":    len(foundOrders),
		"total_items":     len(foundOrderItems),
		"total_sales":     totalSales,
		"aggregate_items": aggregatedResults,
	}

	return response, nil
}

// getOrderIds Create a list of order ids
func getOrderIds(foundOrders []models.StructListItem) []string {
	var orderIds []string
	for _, foundOrder := range foundOrders {
		orderIds = append(orderIds, foundOrder.ID)
	}
	return orderIds
}

// aggregateOrderItems Takes a list of order ids and aggregates their price / qty by their sku
func aggregateOrderItems(oitems []map[string]interface{}) []ProductPerfItem {
	keyedResults := make(map[string]ProductPerfItem)

	// Aggregate by sku
	for _, oitem := range oitems {
		sku := utils.InterfaceToString(oitem["sku"])
		item, ok := keyedResults[sku]

		// First time, set the static details
		if !ok {
			item.Name = utils.InterfaceToString(oitem["name"])
			item.Sku = sku
		}

		item.GrossSales += utils.InterfaceToFloat64(oitem["price"])
		item.UnitsSold += utils.InterfaceToInt(oitem["qty"])

		keyedResults[sku] = item
	}

	// map to slice
	var results ProductPerf
	for _, item := range keyedResults {
		// @TODO: Round money is bad
		item.GrossSales = utils.RoundPrice(item.GrossSales)
		results = append(results, item)
	}

	sort.Sort(results)

	return results
}

func getTotalSales(oitems []map[string]interface{}) float64 {
	var totalSales float64
	for _, oitem := range oitems {
		totalSales += utils.InterfaceToFloat64(oitem["price"])
	}

	return totalSales
}

func listCustomerActivity(context api.InterfaceApplicationContext) (interface{}, error) {
	perfStart := time.Now()

	// Limit results count, not the query
	limit := utils.InterfaceToInt(context.GetRequestArgument("limit"))
	if limit == 0 {
		limit = 50
	}

	sortArg := utils.InterfaceToString(context.GetRequestArgument("sort"))

	// Fetch orders
	oModel, _ := order.GetOrderCollectionModel()
	if err := oModel.ListAddExtraAttribute("created_at"); err != nil {
		env.LogError(errors.New("35e5e271-a605-44f3-9bdc-0972cd9e14c8 : " + err.Error()))
	}
	if err := oModel.ListAddExtraAttribute("customer_name"); err != nil {
		env.LogError(errors.New("30426172-a9a1-4d8c-978e-944acce18dde : " + err.Error()))
	}
	if err := oModel.ListAddExtraAttribute("customer_email"); err != nil {
		env.LogError(errors.New("451cf083-dd39-4805-9732-115afd8bd347 : " + err.Error()))
	}
	if err := oModel.ListAddExtraAttribute("grand_total"); err != nil {
		env.LogError(errors.New("d6c9dd6d-40ed-4820-9197-fadea124e46c : " + err.Error()))
	}

	err := ApplyDateRangeFilter(context, oModel.GetDBCollection())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// This is the lite response StructListItem
	foundOrders, _ := oModel.List()
	aggregatedResults := aggregateCustomerActivity(foundOrders)
	resultCount := len(aggregatedResults)

	// Sorting
	switch sortArg {
	case "total_orders":
		sort.Sort(CustomerActivityByOrders(aggregatedResults))
	case "total_sales":
		fallthrough
	default:
		sort.Sort(CustomerActivityBySales(aggregatedResults))
	}

	// Apply the limit
	if resultCount > limit {
		aggregatedResults = aggregatedResults[:limit]
	}

	response := map[string]interface{}{
		"aggregate_items": aggregatedResults,
		"meta": map[string]interface{}{
			"limit": limit,
			"count": resultCount,
		},
		"perf_ms": time.Now().Sub(perfStart).Seconds() * 1e3, // in milliseconds
	}
	return response, nil
}

func aggregateCustomerActivity(foundOrders []models.StructListItem) []CustomerActivityItem {
	keyedResults := make(map[string]CustomerActivityItem)

	for _, o := range foundOrders {
		email := utils.InterfaceToString(o.Extra["customer_email"])
		createdAt := utils.InterfaceToTime(o.Extra["created_at"])
		item, ok := keyedResults[email]

		// First time, set some static props
		if !ok {
			item.Email = email
			item.EarliestPurchase = createdAt
			item.LatestPurchase = createdAt
		} else {
			if createdAt.Before(item.EarliestPurchase) {
				item.EarliestPurchase = createdAt
			}
			if createdAt.After(item.LatestPurchase) {
				item.LatestPurchase = createdAt
			}
		}

		// Name might be empty, early records had a bug
		if item.Name == "" {
			item.Name = utils.InterfaceToString(o.Extra["customer_name"])
		}

		// Aggregated props
		item.TotalSales += utils.InterfaceToFloat64(o.Extra["grand_total"])
		item.TotalOrders++

		// Save
		keyedResults[email] = item
	}

	// map to slice
	var results []CustomerActivityItem
	for _, i := range keyedResults {
		// Add in averaging stat now that aggregation is complete
		i.AverageSales = i.TotalSales / float64(i.TotalOrders)

		// Round money
		i.TotalSales = utils.RoundPrice(i.TotalSales)
		i.AverageSales = utils.RoundPrice(i.AverageSales)

		results = append(results, i)
	}

	return results
}

func listPaymentMethod(context api.InterfaceApplicationContext) (interface{}, error) {
	perfStart := time.Now()

	// Fetch orders
	oModel, _ := order.GetOrderCollectionModel()
	if err := oModel.ListAddExtraAttribute("created_at"); err != nil {
		env.LogError(errors.New("c63e01a1-2e2d-41a0-a8f9-09b55d25fadc : " + err.Error()))
	}
	if err := oModel.ListAddExtraAttribute("payment_method"); err != nil {
		env.LogError(errors.New("4a53c0a8-67d3-4a1f-97c1-9c2d1f6e625c : " + err.Error()))
	}
	if err := oModel.ListAddExtraAttribute("grand_total"); err != nil {
		env.LogError(errors.New("f7524ade-0983-45a6-9f31-ed11de180680 : " + err.Error()))
	}

	err := ApplyDateRangeFilter(context, oModel.GetDBCollection())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// This is the lite response StructListItem
	foundOrders, _ := oModel.List()
	aggregatedResults := aggregatePaymentMethod(foundOrders)

	// Sorting
	sort.Sort(StatsBySales(aggregatedResults))

	// Calculate extra data points
	var totalSales float64
	for _, m := range aggregatedResults {
		totalSales += m.TotalSales
	}
	totalSales = utils.RoundPrice(totalSales)

	response := map[string]interface{}{
		"aggregate_items": aggregatedResults,
		"total_sales":     totalSales,
		"perf_ms":         time.Now().Sub(perfStart).Seconds() * 1e3, // in milliseconds
	}
	return response, nil
}

func aggregatePaymentMethod(foundOrders []models.StructListItem) []StatItem {

	paymentMethodNames := map[string]string{}

	for _, m := range checkout.GetRegisteredPaymentMethods() {
		code := m.GetCode()
		paymentMethodNames[code] = m.GetInternalName()
	}

	aggregateKey := "payment_method"

	return aggregateGeneral(foundOrders, paymentMethodNames, aggregateKey)
}

func aggregateShippingMethod(foundOrders []models.StructListItem) []StatItem {

	keyNameMap := map[string]string{}
	for _, method := range checkout.GetRegisteredShippingMethods() {
		methodCode := strings.ToLower(method.GetCode())
		methodName := method.GetName()
		rates := method.GetAllRates()

		for _, rate := range rates {
			key := methodCode + "/" + rate.Code
			keyNameMap[key] = methodName + " - " + rate.Name
		}
	}

	aggregateKey := "shipping_method"
	return aggregateGeneral(foundOrders, keyNameMap, aggregateKey)
}

func aggregateGeneral(foundOrders []models.StructListItem, keyNameMap map[string]string, aggregateKey string) []StatItem {

	keyedResults := make(map[string]StatItem)

	for _, o := range foundOrders {
		key := utils.InterfaceToString(o.Extra[aggregateKey])
		item, ok := keyedResults[key]

		// First time, set some static props
		if !ok {
			item.Key = key
			item.Name = keyNameMap[key]
		}

		// Aggregated props
		item.TotalSales += utils.InterfaceToFloat64(o.Extra["grand_total"])
		item.TotalOrders++

		// Save
		keyedResults[key] = item
	}

	// map to slice
	var results []StatItem
	for _, i := range keyedResults {
		// Add in averaging stat now that aggregation is complete
		i.AverageSales = i.TotalSales / float64(i.TotalOrders)

		// Round money
		i.TotalSales = utils.RoundPrice(i.TotalSales)
		i.AverageSales = utils.RoundPrice(i.AverageSales)

		results = append(results, i)
	}

	return results
}

func listShippingMethod(context api.InterfaceApplicationContext) (interface{}, error) {
	perfStart := time.Now()

	// Fetch orders
	oModel, _ := order.GetOrderCollectionModel()
	if err := oModel.ListAddExtraAttribute("created_at"); err != nil {
		env.LogError(errors.New("4ff4ac4c-8048-41e6-89b6-ab65c03ffbde : " + err.Error()))
	}
	if err := oModel.ListAddExtraAttribute("shipping_method"); err != nil {
		env.LogError(errors.New("8104d737-81a3-4059-8a65-1c5a513930a8 : " + err.Error()))
	}
	if err := oModel.ListAddExtraAttribute("grand_total"); err != nil {
		env.LogError(errors.New("d98106a4-63f9-4f7d-b854-74c6af6a0a33 : " + err.Error()))
	}

	err := ApplyDateRangeFilter(context, oModel.GetDBCollection())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// This is the lite response StructListItem
	foundOrders, _ := oModel.List()
	aggregatedResults := aggregateShippingMethod(foundOrders)

	// Sorting
	sort.Sort(StatsBySales(aggregatedResults))

	// Calculate extra data points
	var totalSales float64
	for _, m := range aggregatedResults {
		totalSales += m.TotalSales
	}
	totalSales = utils.RoundPrice(totalSales)

	response := map[string]interface{}{
		"aggregate_items": aggregatedResults,
		"total_sales":     totalSales,
		"perf_ms":         time.Now().Sub(perfStart).Seconds() * 1e3, // in milliseconds
	}
	return response, nil
}

func listLocationCountry(context api.InterfaceApplicationContext) (interface{}, error) {
	perfStart := time.Now()

	// Fetch orders
	oModel, _ := order.GetOrderCollectionModel()
	if err := oModel.ListAddExtraAttribute("created_at"); err != nil {
		env.LogError(errors.New("fbcc01df-049c-418b-8a15-f3acb595a26f : " + err.Error()))
	}
	if err := oModel.ListAddExtraAttribute("billing_address"); err != nil {
		env.LogError(errors.New("8d7c47d2-382e-4269-8633-fd4c7f5254fc : " + err.Error()))
	}
	if err := oModel.ListAddExtraAttribute("grand_total"); err != nil {
		env.LogError(errors.New("26dbf8cd-b884-4d2d-af4d-d0328ec4e24a : " + err.Error()))
	}

	err := ApplyDateRangeFilter(context, oModel.GetDBCollection())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// This is the lite response StructListItem
	foundOrders, _ := oModel.List()
	aggregatedResults := aggregateLocationCountry(foundOrders)

	// Sorting
	sort.Sort(StatsBySales(aggregatedResults))

	// Calculate extra data points
	var totalSales float64
	for _, m := range aggregatedResults {
		totalSales += m.TotalSales
	}
	totalSales = utils.RoundPrice(totalSales)

	response := map[string]interface{}{
		"aggregate_items": aggregatedResults,
		"total_sales":     totalSales,
		"perf_ms":         time.Now().Sub(perfStart).Seconds() * 1e3, // in milliseconds
	}
	return response, nil
}

func aggregateLocationCountry(foundOrders []models.StructListItem) []StatItem {
	return aggregateGeneralNested(foundOrders, "billing_address", "country")
}

func aggregateLocationUS(foundOrders []models.StructListItem) []StatItem {
	return aggregateGeneralNested(foundOrders, "billing_address", "state")
}

func aggregateGeneralNested(foundOrders []models.StructListItem, aggKeyContainer string, aggKey string) []StatItem {
	keyedResults := make(map[string]StatItem)

	for _, o := range foundOrders {
		container := utils.InterfaceToMap(o.Extra[aggKeyContainer])
		key := utils.InterfaceToString(container[aggKey])

		item, ok := keyedResults[key]
		if !ok {
			item.Name = key
		}

		// Aggregate props
		item.TotalSales += utils.InterfaceToFloat64(o.Extra["grand_total"])
		item.TotalOrders++

		// Save
		keyedResults[key] = item
	}

	// map to slice
	var results []StatItem
	for _, i := range keyedResults {
		// Add in averaging stat now that aggregation is complete
		i.AverageSales = i.TotalSales / float64(i.TotalOrders)

		// Round money
		i.TotalSales = utils.RoundPrice(i.TotalSales)
		i.AverageSales = utils.RoundPrice(i.AverageSales)

		results = append(results, i)
	}

	return results
}

// list aggregate sales by state for sales in the US
func listLocationUS(context api.InterfaceApplicationContext) (interface{}, error) {
	perfStart := time.Now()

	// Fetch orders
	oModel, _ := order.GetOrderCollectionModel()
	if err := oModel.ListAddExtraAttribute("created_at"); err != nil {
		env.LogError(errors.New("afee38b6-7744-45d1-a202-fa3e38cebb77 : " + err.Error()))
	}
	if err := oModel.ListAddExtraAttribute("billing_address"); err != nil {
		env.LogError(errors.New("d425730d-3298-44db-8312-686b6a31fa2b : " + err.Error()))
	}
	if err := oModel.ListFilterAdd("billing_address.country", "=", "US"); err != nil {
		env.LogError(errors.New("ca701edb-dad9-41ff-afa5-b8b71e9dbef8 : " + err.Error()))
	}
	if err := oModel.ListAddExtraAttribute("grand_total"); err != nil {
		env.LogError(errors.New("c2a7884f-30ec-46cf-9453-723fba5c4a72 : " + err.Error()))
	}

	err := ApplyDateRangeFilter(context, oModel.GetDBCollection())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// This is the lite response StructListItem
	foundOrders, _ := oModel.List()
	aggregatedResults := aggregateLocationUS(foundOrders)

	// Sorting
	sort.Sort(StatsBySales(aggregatedResults))

	// Calculate extra data points
	var totalSales float64
	for _, m := range aggregatedResults {
		totalSales += m.TotalSales
	}
	totalSales = utils.RoundPrice(totalSales)

	response := map[string]interface{}{
		"aggregate_items": aggregatedResults,
		"total_sales":     totalSales,
		"perf_ms":         time.Now().Sub(perfStart).Seconds() * 1e3, // in milliseconds
	}
	return response, nil
}

// listGiftCards returns information about gift cards
func listGiftCards(context api.InterfaceApplicationContext) (interface{}, error) {
	perfStart := time.Now()

	giftCardCollection, err := db.GetCollection(giftcard.ConstCollectionNameGiftCard)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	err = ApplyDateRangeFilter(context, giftCardCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collectionRecords, err := giftCardCollection.Load()
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	var giftCards []map[string]interface{}
	var total = 0.0
	var count = 0

	// get gift cards information
	for _, item := range collectionRecords {
		amount := utils.InterfaceToFloat64(item["amount"])
		giftCardItem := map[string]interface{}{
			"code":   utils.InterfaceToString(item["code"]),
			"name":   utils.InterfaceToString(item["name"]),
			"amount": amount,
			"date":   utils.InterfaceToTime(item["created_at"]),
		}
		total = total + amount
		count = count + 1
		giftCards = append(giftCards, giftCardItem)
	}

	results := map[string]interface{}{
		"gift-cards": giftCards,
		"total":      total,
		"count":      count,
		"perf_ms":    time.Now().Sub(perfStart).Seconds() * 1e3, // in milliseconds
	}

	return results, nil
}

func ApplyDateRangeFilter(context api.InterfaceApplicationContext, collection db.InterfaceDBCollection) error {
	// Expecting dates in UTC, and adjusted for your timezone `2006-01-02 15:04`
	startDate := utils.InterfaceToTime(context.GetRequestArgument("start_date"))
	endDate := utils.InterfaceToTime(context.GetRequestArgument("end_date"))
	hasDateRange := !startDate.IsZero() || !endDate.IsZero()

	// Date range validation
	if hasDateRange {
		err := ValidateStartAndEndDate(startDate, endDate)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		if err := collection.AddFilter("created_at", ">=", startDate); err != nil {
			env.LogError(errors.New("30c49961-ab76-4eaa-a7ba-24d8aa591536 : " + err.Error()))
		}
		if err := collection.AddFilter("created_at", "<", endDate); err != nil {
			env.LogError(errors.New("47403827-0a63-4067-bbb3-f00b85d6e775 : " + err.Error()))
		}
	}
	return nil
}

// ValidateStartAndEndDate - date range validation
func ValidateStartAndEndDate(startDate time.Time, endDate time.Time) error {
	if startDate.IsZero() || endDate.IsZero() {
		msg := "start_date or end_date missing from response, or not formatted in YYYY-MM-DD"
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "88b6fe1f-0e2f-4e63-b0a4-2e767c27dfd8", msg)
	}

	if startDate.After(endDate) || startDate.Equal(endDate) {
		msg := "the start_date must come before the end_date"
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fb30b99c-0648-4219-8b9c-933edf9f7ed3", msg)
	}

	return nil
}
