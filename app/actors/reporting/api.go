package reporting

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"sort"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	// Admin only endpoint
	service := api.GetRestService()
	service.GET("reporting/product-performance", listProductPerformance)

	return nil
}

// listProductPerformance Handler that returns product performance information by date range
func listProductPerformance(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

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
func aggregateOrderItems(oitems []map[string]interface{}) []AggrOrderItems {
	keyedResults := make(map[string]AggrOrderItems)

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
	var results []AggrOrderItems
	for _, item := range keyedResults {
		// @TODO: Round money is bad
		item.GrossSales = utils.RoundPrice(item.GrossSales)
		results = append(results, item)
	}

	sort.Sort(ByUnitsSold(results))

	return results
}

func getTotalSales(oitems []map[string]interface{}) float64 {
	var totalSales float64
	for _, oitem := range oitems {
		totalSales += utils.InterfaceToFloat64(oitem["price"])
	}

	return totalSales
}
