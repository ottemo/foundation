package reporting

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"sort"
	"time"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()
	service.GET("reporting/product-performance", listProductPerformance)

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

	foundOrders := getOrders(startDate, endDate)
	foundOrderIds := getOrderIds(foundOrders)
	foundOrderItems := getItemsForOrders(foundOrderIds)
	aggregatedResults := aggregateOrderItems(foundOrderItems)

	response := map[string]interface{}{
		"order_count":     len(foundOrders),
		"item_count":      len(foundOrderItems),
		"aggregate_items": aggregatedResults,
	}

	return response, nil
}

// getOrders Get the orders `created_at` a certain date range
func getOrders(startDate time.Time, endDate time.Time) []models.StructListItem {
	oModel, _ := order.GetOrderCollectionModel()
	oModel.GetDBCollection().AddFilter("created_at", ">=", startDate)
	oModel.GetDBCollection().AddFilter("created_at", "<", endDate)
	oModel.ListAddExtraAttribute("created_at")
	foundOrders, _ := oModel.List()

	return foundOrders
}

// getOrderIds Create a list of order ids
func getOrderIds(foundOrders []models.StructListItem) []string {
	var orderIds []string
	for _, foundOrder := range foundOrders {
		orderIds = append(orderIds, foundOrder.ID)
	}
	return orderIds
}

// getItemsForOrders Get the relavent order items given a slice of orders
func getItemsForOrders(orderIds []string) []map[string]interface{} {
	oiModel, _ := order.GetOrderItemCollectionModel()
	oiDB := oiModel.GetDBCollection()
	oiDB.AddFilter("order_id", "in", orderIds)
	oiResults, _ := oiDB.Load()

	return oiResults
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
		results = append(results, item)
	}

	sort.Sort(ByUnitsSold(results))

	return results
}
