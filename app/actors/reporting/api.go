package reporting

import (
	"sort"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"errors"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	// Admin only endpoint
	service := api.GetRestService()
	service.GET("reporting/product-performance", api.IsAdmin(listProductPerformance))
	service.GET("reporting/customer-activity", api.IsAdmin(listCustomerActivity))
	service.GET("reporting/payment-method", api.IsAdmin(listPaymentMethod))
	service.GET("reporting/shipping-method", api.IsAdmin(listShippingMethod))
	service.GET("reporting/location-country", api.IsAdmin(listLocationCountry))
	service.GET("reporting/location-us", api.IsAdmin(listLocationUS))

	return nil
}

// listProductPerformance Handler that returns product performance information by date range
func listProductPerformance(context api.InterfaceApplicationContext) (interface{}, error) {

	// Expecting dates in UTC, and adjusted for your timezone
	// `2006-01-02 15:04`
	startDate := utils.InterfaceToTime(context.GetRequestArgument("start_date"))
	endDate := utils.InterfaceToTime(context.GetRequestArgument("end_date"))
	if startDate.IsZero() || endDate.IsZero() {
		context.SetResponseStatusBadRequest()
		msg := "start_date or end_date missing from response, or not formatted in YYYY-MM-DD"
		return nil, env.ErrorNew("reporting", 6, "3ed77c0d-2c54-4401-9feb-6e1d04b8baef", msg)
	}
	if startDate.After(endDate) || startDate.Equal(endDate) {
		context.SetResponseStatusBadRequest()
		msg := "the start_date must come before the end_date"
		return nil, env.ErrorNew("reporting", 6, "2eb9680c-d9a8-42ce-af63-fd6b0b742d0d", msg)
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

	// Expecting dates in UTC, and adjusted for your timezone `2006-01-02 15:04`
	startDate := utils.InterfaceToTime(context.GetRequestArgument("start_date"))
	endDate := utils.InterfaceToTime(context.GetRequestArgument("end_date"))
	hasDateRange := !startDate.IsZero() || !endDate.IsZero()

	// Date range validation
	if hasDateRange {
		if startDate.IsZero() || endDate.IsZero() {
			context.SetResponseStatusBadRequest()
			msg := "start_date or end_date missing from response, or not formatted in YYYY-MM-DD"
			return nil, env.ErrorNew("reporting", 6, "88b6fe1f-0e2f-4e63-b0a4-2e767c27dfd8", msg)
		}
		if startDate.After(endDate) || startDate.Equal(endDate) {
			context.SetResponseStatusBadRequest()
			msg := "the start_date must come before the end_date"
			return nil, env.ErrorNew("reporting", 6, "fb30b99c-0648-4219-8b9c-933edf9f7ed3", msg)
		}
	}

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
	if hasDateRange {
		if err := oModel.GetDBCollection().AddFilter("created_at", ">=", startDate); err != nil {
			env.LogError(errors.New("e1298ee1-9a10-47df-8404-908a4d9f4981 : " + err.Error()))
		}
		if err := oModel.GetDBCollection().AddFilter("created_at", "<", endDate); err != nil {
			env.LogError(errors.New("4afb1a1c-b8c9-42db-9e38-4286cc914190 : " + err.Error()))
		}
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

	// Expecting dates in UTC, and adjusted for your timezone `2006-01-02 15:04`
	startDate := utils.InterfaceToTime(context.GetRequestArgument("start_date"))
	endDate := utils.InterfaceToTime(context.GetRequestArgument("end_date"))
	hasDateRange := !startDate.IsZero() || !endDate.IsZero()

	// Date range validation
	if hasDateRange {
		if startDate.IsZero() || endDate.IsZero() {
			context.SetResponseStatusBadRequest()
			msg := "start_date or end_date missing from response, or not formatted in YYYY-MM-DD"
			return nil, env.ErrorNew("reporting", 6, "5731fd92-b2b1-44e7-8940-395575bca081", msg)
		}
		if startDate.After(endDate) || startDate.Equal(endDate) {
			context.SetResponseStatusBadRequest()
			msg := "the start_date must come before the end_date"
			return nil, env.ErrorNew("reporting", 6, "d771dc9d-de2a-4e84-8d5e-9bd050cc5d2d", msg)
		}
	}

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
	if hasDateRange {
		if err := oModel.GetDBCollection().AddFilter("created_at", ">=", startDate); err != nil {
			env.LogError(errors.New("b4206d27-7f73-47de-9e7a-4afe0600886d : " + err.Error()))
		}
		if err := oModel.GetDBCollection().AddFilter("created_at", "<", endDate); err != nil {
			env.LogError(errors.New("be1b13f1-db44-47b9-b6a5-a05253b100bf : " + err.Error()))
		}
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

	// Expecting dates in UTC, and adjusted for your timezone `2006-01-02 15:04`
	startDate := utils.InterfaceToTime(context.GetRequestArgument("start_date"))
	endDate := utils.InterfaceToTime(context.GetRequestArgument("end_date"))
	hasDateRange := !startDate.IsZero() || !endDate.IsZero()

	// Date range validation
	if hasDateRange {
		if startDate.IsZero() || endDate.IsZero() {
			context.SetResponseStatusBadRequest()
			msg := "start_date or end_date missing from response, or not formatted in YYYY-MM-DD"
			return nil, env.ErrorNew("reporting", 6, "48beae0f-f6fb-49f0-adff-1d0a2f8b5fff", msg)
		}
		if startDate.After(endDate) || startDate.Equal(endDate) {
			context.SetResponseStatusBadRequest()
			msg := "the start_date must come before the end_date"
			return nil, env.ErrorNew("reporting", 6, "eb050c5e-6ee8-4d09-9869-1813b252d3aa", msg)
		}
	}

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
	if hasDateRange {
		if err := oModel.GetDBCollection().AddFilter("created_at", ">=", startDate); err != nil {
			env.LogError(errors.New("2e38dab2-dedb-4e91-9538-b9d732c90408 : " + err.Error()))
		}
		if err := oModel.GetDBCollection().AddFilter("created_at", "<", endDate); err != nil {
			env.LogError(errors.New("eccd9296-fb5c-476b-ba05-efbbbb21a3b1 : " + err.Error()))
		}
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

	// Expecting dates in UTC, and adjusted for your timezone `2006-01-02 15:04`
	startDate := utils.InterfaceToTime(context.GetRequestArgument("start_date"))
	endDate := utils.InterfaceToTime(context.GetRequestArgument("end_date"))
	hasDateRange := !startDate.IsZero() || !endDate.IsZero()

	// Date range validation
	if hasDateRange {
		if startDate.IsZero() || endDate.IsZero() {
			context.SetResponseStatusBadRequest()
			msg := "start_date or end_date missing from response, or not formatted in YYYY-MM-DD"
			return nil, env.ErrorNew("reporting", 6, "8b9940cc-cc45-4ab8-af92-4f3cff0db5b1", msg)
		}
		if startDate.After(endDate) || startDate.Equal(endDate) {
			context.SetResponseStatusBadRequest()
			msg := "the start_date must come before the end_date"
			return nil, env.ErrorNew("reporting", 6, "dd9fad29-4321-4951-976f-9204890437bd", msg)
		}
	}

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
	if hasDateRange {
		if err := oModel.GetDBCollection().AddFilter("created_at", ">=", startDate); err != nil {
			env.LogError(errors.New("267ae749-0c75-4845-bc45-799978671084 : " + err.Error()))
		}
		if err := oModel.GetDBCollection().AddFilter("created_at", "<", endDate); err != nil {
			env.LogError(errors.New("c6c49618-49ab-4bf9-b123-6d2631ab3810 : " + err.Error()))
		}
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

	// Expecting dates in UTC, and adjusted for your timezone `2006-01-02 15:04`
	startDate := utils.InterfaceToTime(context.GetRequestArgument("start_date"))
	endDate := utils.InterfaceToTime(context.GetRequestArgument("end_date"))
	hasDateRange := !startDate.IsZero() || !endDate.IsZero()

	// Date range validation
	if hasDateRange {
		if startDate.IsZero() || endDate.IsZero() {
			context.SetResponseStatusBadRequest()
			msg := "start_date or end_date missing from response, or not formatted in YYYY-MM-DD"
			return nil, env.ErrorNew("reporting", 6, "2eb1c70a-3d37-46ab-91b1-9e6124685406", msg)
		}
		if startDate.After(endDate) || startDate.Equal(endDate) {
			context.SetResponseStatusBadRequest()
			msg := "the start_date must come before the end_date"
			return nil, env.ErrorNew("reporting", 6, "bc431db2-77ea-45da-9e61-3ee156ec62b6", msg)
		}
	}

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
	if hasDateRange {
		if err := oModel.GetDBCollection().AddFilter("created_at", ">=", startDate); err != nil {
			env.LogError(errors.New("29e271f7-97c4-4be1-a758-4c5360591814 : " + err.Error()))
		}
		if err := oModel.GetDBCollection().AddFilter("created_at", "<", endDate); err != nil {
			env.LogError(errors.New("8648b2b7-1450-4409-9324-17d7316d023c : " + err.Error()))
		}
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
