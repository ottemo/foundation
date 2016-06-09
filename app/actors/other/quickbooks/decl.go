// Package quickbooks implements exporting function to iif format for Quickbooks
package quickbooks

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Package global constants
const (
	ConstErrorModule = "quickbooks"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

var (
	orderFields = [][]string{

		{
			"Order No", "Order Date", "Order Time", "Email", "Payment Status", "Order Status", "Discount", "Sub Total", "Shipping Cost", "Tax", "Total", "Billing FirstName", "Billing LastName", "Billing Company", "Billing Address 1", "Billing Address 2", "Billing City", "Billing State", "Billing Country", "Billing Zip", "Billing Phone", "Shipping FirstName", "Shipping LastName", "Shipping Company", "Shipping Address 1", "Shipping Address 2", "Shipping City", "Shipping State", "Shipping Country", "Shipping Zip", "Shipping Phone", "Shipping Method", "Shipping Carrier", "Payment Method", "Transaction ID", "Credit Card No", "Card Holder Name", "Card Type", "SKU", "Product Name", "Item Price", "Item Quantity", "Item Weight", "Item Amount", "Item Order No",
		},

		//		{
		//			"!TRNS", "TRNSID", "TRNSTYPE", "DATE", "ACCNT", "NAME", "CLASS", "AMOUNT", "DOCNUM", "MEMO", "CLEAR", "ADDR1", "ADDR2", "ADDR3", "ADDR4", "SADDR1", "SADDR2", "SADDR3", "SADDR4", "PAID",
		//		},
		//
		//		{
		//			"!SPL", "SPLID", "TRNSTYPE", "DATE", "ACCNT", "NAME", "CLASS", "AMOUNT", "DOCNUM", "MEMO", "CLEAR", "QNTY", "PRICE", "INVITEM", "EXTRA",
		//		},
		//
		//		{
		//			"!ENDTRNS",
		//		},
	}

	dataSeted = [][]interface{}{

		{
			"$increment_id",
			func(record map[string]interface{}) string {
				timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))

				createdAt := utils.InterfaceToTime(record["created_at"])
				convertedDate, _ := utils.MakeUTCOffsetTime(createdAt, timeZone)
				return convertedDate.Format("01/02/2006")
			},
			func(record map[string]interface{}) string {
				timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))

				createdAt := utils.InterfaceToTime(record["created_at"])
				convertedDate, _ := utils.MakeUTCOffsetTime(createdAt, timeZone)
				return convertedDate.Format("15:04:05")
			},

			"$customer_email", "$status", "$status", "$discount", "$subtotal", "$shipping_amount", "$tax_amount", "$grand_total",
			"billing.first_name", "billing.last_name", "billing.company", "billing.address_line1", "billing.address_line2",
			"billing.city", "billing.state", "billing.country", "billing.zip_code", "billing.phone",
			"shipping.first_name", "shipping.last_name", "shipping.company", "shipping.address_line1", "shipping.address_line2",
			"shipping.city", "shipping.state", "shipping.country", "shipping.zip_code", "shipping.phone",

			// "Shipping Method" shipping_info.shipping_method_name
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["shipping_info"])["shipping_method_name"])
			},

			"$shipping_method",

			// "Payment Method", payment_info.payment_method_name
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["payment_info"])["payment_method_name"])
			},

			// "Transaction ID", payment_info.transactionID
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["payment_info"])["transactionID"])
			},

			// "Credit Card No", payment_info.creditCardNumbers
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["payment_info"])["creditCardNumbers"])
			},

			// "Card Holder Name"
			func(record map[string]interface{}) string {
				address := utils.InterfaceToMap(record["billing_address"])
				return utils.InterfaceToString(address["first_name"]) + " " + utils.InterfaceToString(address["last_name"])
			},

			// "Card Type",, payment_info.creditCardType
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["payment_info"])["creditCardType"])
			},
			"item.sku", "item.name", "item.price", "item.qty", "item.weight",
			//			"item.amount",
			func(orderItemIndex int, orderItem map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToFloat64(orderItem["price"]) * utils.InterfaceToFloat64(orderItem["qty"]))
			},
			"item.order_id",
		},

		//		{
		//			// "!TRNS", "TRNSID", "TRNSTYPE",
		//			"TRNS", "$_id", "INVOICE",
		//			// Date
		//			func(record map[string]interface{}) string {
		//				return utils.InterfaceToTime(record["created_at"]).Format("01/02/06")
		//			},
		//			// "ACCNT", 		"NAME", "CLASS", "AMOUNT", "DOCNUM", "MEMO", "CLEAR",
		//			"Accounts Receivable", "$customer_name", "", "$grand_total", "$increment_id", "$customer_email", "N",
		//			// "ADDR1",
		//			func(record map[string]interface{}) string {
		//				address := utils.InterfaceToMap(record["billing_address"])
		//				return utils.InterfaceToString(address["first_name"]) + utils.InterfaceToString(address["last_name"])
		//			},
		//			// "ADDR2",
		//			func(record map[string]interface{}) string {
		//				return utils.InterfaceToString(utils.InterfaceToMap(record["billing_address"])["address_line1"])
		//			},
		//			//	"ADDR3",
		//			func(record map[string]interface{}) string {
		//				address := utils.InterfaceToMap(record["billing_address"])
		//				return utils.InterfaceToString(address["city"]) + ", " + utils.InterfaceToString(address["state"]) +
		//					" " + utils.InterfaceToString(address["zip_code"])
		//			},
		//			// "ADDR4",
		//			func(record map[string]interface{}) string {
		//				return utils.InterfaceToString(utils.InterfaceToMap(record["billing_address"])["country"])
		//			},
		//			// "SADDR1",
		//			func(record map[string]interface{}) string {
		//				address := utils.InterfaceToMap(record["shipping_address"])
		//				return utils.InterfaceToString(address["first_name"]) + utils.InterfaceToString(address["last_name"])
		//			},
		//			// "SADDR2",
		//			func(record map[string]interface{}) string {
		//				return utils.InterfaceToString(utils.InterfaceToMap(record["shipping_address"])["address_line1"])
		//			},
		//			//	"SADDR3",
		//			func(record map[string]interface{}) string {
		//				address := utils.InterfaceToMap(record["shipping_address"])
		//				return utils.InterfaceToString(address["city"]) + ", " + utils.InterfaceToString(address["state"]) +
		//					" " + utils.InterfaceToString(address["zip_code"])
		//			},
		//			// "SADDR4",
		//			func(record map[string]interface{}) string {
		//				return utils.InterfaceToString(utils.InterfaceToMap(record["shipping_address"])["country"])
		//			},
		//
		//			// PAID Y or N depends from status
		//			func(record map[string]interface{}) string {
		//				status := utils.InterfaceToString(record["status"])
		//				if status == order.ConstOrderStatusCompleted || status == order.ConstOrderStatusProcessed {
		//					return "Y"
		//				}
		//				return "N"
		//			},
		//		},
		//		{
		//			// "!SPL", "SPLID", "TRNSTYPE"
		//			"SPL", "$_id", "INVOICE",
		//			// Date
		//			func(record map[string]interface{}) string {
		//				return utils.InterfaceToTime(record["created_at"]).Format("01/02/06")
		//			},
		//			// "ACCNT", "NAME", "CLASS"
		//			"Income", "", "",
		//			// Amount with negative value
		//			func(record map[string]interface{}) string {
		//				return "-" + utils.InterfaceToString(record["grand_total"])
		//			},
		//			// "DOCNUM", "MEMO", "CLEAR", "QNTY", "PRICE", "INVITEM", "EXTRA"
		//			"$increment_id", "$description", "N", "1", "$grand_total", "Product", "$description",
		//		},
		//		{
		//			"ENDTRNS",
		//		},
	}

	blocksHeaders = map[string][][]string{
		"customers": {{}},
		"invoices":  orderFields,
	}

	blocksRecords = map[string][][]interface{}{
		"customers": {{}},
		"invoices":  dataSeted,
	}
)
