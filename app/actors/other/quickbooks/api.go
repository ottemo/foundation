package quickbooks

import (
	"encoding/csv"
	"regexp"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI configures the package related API endpoints
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("export/quickbooks", api.ConstRESTOperationCreate, APIExportOrders)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIExportOrders returns a list of orders in Quickbooks IIF format
// - returns all orders in IIF format with no parameters
// - returns orders specified in url parameters
func APIExportOrders(context api.InterfaceApplicationContext) (interface{}, error) {
	var itemCSVRecords [][]string
	var orders []string

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(requestData, "orders") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f2602f73-7cae-4525-8405-9e470681c20e", "Specifiy a minimum of one order in the orders parameter.")
	}

	orderItemsCollectionModel, err := order.GetOrderItemCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbOrderItemsCollection := orderItemsCollectionModel.GetDBCollection()

	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.AddSort("created_at", false)

	orders = utils.InterfaceToStringArray(utils.InterfaceToArray(requestData["orders"]))
	if orders != nil && len(orders) > 0 && !utils.IsInListStr("all", orders) {
		dbOrderCollection.AddFilter("_id", "in", orders)
	}

	ordersRecords, err := dbOrderCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if len(ordersRecords) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "28eac91b-39ec-4034-b664-4004e940a6d1", "No orders were found.")
	}

	for _, columnsHeaders := range orderFields {
		itemCSVRecords = append(itemCSVRecords, columnsHeaders)
	}

	// WebGility importer bombs out on some characters even if they are properly
	// escaped in a csv: . , \ / ( )
	regexCleaner, _ := regexp.Compile(`[.,\/\\()]`)

	for _, orderRecord := range ordersRecords {

		dbOrderItemsCollection.ClearFilters()
		dbOrderItemsCollection.AddFilter("order_id", "=", orderRecord["_id"])

		orderItemsRecords, err := dbOrderItemsCollection.Load()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		orderRecord := utils.InterfaceToMap(orderRecord)
		shippingAddress := utils.InterfaceToMap(orderRecord["shipping_address"])
		billingAddress := utils.InterfaceToMap(orderRecord["billing_address"])

		for orderItemIndex, orderItem := range orderItemsRecords {

			for _, inputValues := range dataSeted {
				var rowData []string

				for _, value := range inputValues {
					cellValue := ""

					switch typedValue := value.(type) {
					case string:
						cellValue = typedValue
						switch {
						case strings.Index(cellValue, "$") == 0:
							cellValue = utils.InterfaceToString(orderRecord[strings.Replace(cellValue, "$", "", 1)])
							break

						case strings.HasPrefix(cellValue, "item."):
							cellValue = utils.InterfaceToString(orderItem[strings.Replace(cellValue, "item.", "", 1)])
							break

						case strings.HasPrefix(cellValue, "shipping."):
							addressKey := strings.Replace(cellValue, "shipping.", "", 1)
							cellValue = utils.InterfaceToString(shippingAddress[addressKey])
							cellValue = regexCleaner.ReplaceAllString(cellValue, "")
							break

						case strings.HasPrefix(cellValue, "billing."):
							addressKey := strings.Replace(cellValue, "billing.", "", 1)
							cellValue = utils.InterfaceToString(billingAddress[addressKey])
							cellValue = regexCleaner.ReplaceAllString(cellValue, "")
							break
						}
						break

					case func(record map[string]interface{}) string:
						cellValue = typedValue(orderRecord)
						break

					case func(int, map[string]interface{}) string:
						cellValue = typedValue(orderItemIndex, orderItem)
						break

					}
					// collect cellValue to row
					rowData = append(rowData, cellValue)
				}
				// collect row to table
				itemCSVRecords = append(itemCSVRecords, rowData)
			}
		}
	}

	// preparing csv writer
	csvWriter := csv.NewWriter(context.GetResponseWriter())

	exportFilename := "orders_export_" + time.Now().Format(time.RFC3339) + ".csv"

	context.SetResponseContentType("text/csv")
	context.SetResponseSetting("Content-disposition", "attachment;filename="+exportFilename)

	for _, csvRecord := range itemCSVRecords {
		csvWriter.Write(csvRecord)
	}
	csvWriter.Flush()

	return "", nil
}
