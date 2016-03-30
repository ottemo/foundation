package reporting

import (
	"fmt"
	"github.com/ottemo/foundation/api"
	"time"
	// "github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
	// "github.com/ottemo/foundation/app/models/subscription"
	// "github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("reporting/product-performance", listProductPerformance)
	return nil
}

func listProductPerformance(context api.InterfaceApplicationContext) (interface{}, error) {

	// Expecting dates in UTC, and adjusted for your timezone
	startDate := utils.InterfaceToTime(context.GetRequestArgument("start_date"))
	endDate := utils.InterfaceToTime(context.GetRequestArgument("end_date"))
	if startDate.IsZero() || endDate.IsZero() {
		context.SetResponseStatusBadRequest()
		msg := "start_date or end_date missing from response, or not formatted in YYYY-MM-DD"
		return nil, env.ErrorNew("reporting", 6, "3ed77c0d-2c54-4401-9feb-6e1d04b8baef", msg)
	}

	// Debugging //TODO: RM
	fmt.Println(endDate, startDate)

	foundOrders := getOrderIds(startDate, endDate)
	foundOrderItems := getItemsForOrders(foundOrders)
	aggregatedResults := aggregateOrderItems(foundOrderItems)

	return aggregatedResults, nil
}

func getOrderIds(startDate time.Time, endDate time.Time) []models.StructListItem {
	oModel, _ := order.GetOrderCollectionModel()
	//TODO:
	// oModel.GetDBCollection().AddFilter("updated_at", ">=", startDate)
	// oModel.GetDBCollection().AddFilter("updated_at", "<", endDate)
	oModel.ListLimit(0, 10)
	oModel.ListAddExtraAttribute("created_at")
	foundOrders, _ := oModel.List()

	return foundOrders
}

func getItemsForOrders(foundOrders []models.StructListItem) []map[string]interface{} {
	// get list of order ids
	var orderIds []string
	for _, foundOrder := range foundOrders {
		orderIds = append(orderIds, foundOrder.ID)
	}

	// load the order items
	oiModel, _ := order.GetOrderItemCollectionModel()
	oiDB := oiModel.GetDBCollection()
	oiDB.AddFilter("order_id", "in", orderIds)
	oiResults, _ := oiDB.Load()

	return oiResults
}

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

	// Strip the keys off of this map
	var results []AggrOrderItems
	for _, item := range keyedResults {
		results = append(results, item)
	}

	return results
}

type AggrOrderItems struct {
	Name       string  `json:"name"`
	Sku        string  `json:"sku"`
	GrossSales float64 `json:"gross_sales"`
	UnitsSold  int     `json:"units_sold"`
}
