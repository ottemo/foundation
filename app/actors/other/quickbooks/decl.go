// Package quickbooks implements exporting function to iif format for Quickbooks
package quickbooks

import (
	"github.com/ottemo/foundation/app/models/order"
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
			"!TRNS", "TRNSID", "TRNSTYPE", "DATE", "ACCNT", "NAME", "CLASS", "AMOUNT", "DOCNUM", "MEMO", "CLEAR", "ADDR1", "ADDR2", "ADDR3", "ADDR4", "SADDR1", "SADDR2", "SADDR3", "SADDR4", "PAID",
		},

		{
			"!SPL", "SPLID", "TRNSTYPE", "DATE", "ACCNT", "NAME", "CLASS", "AMOUNT", "DOCNUM", "MEMO", "CLEAR", "QNTY", "PRICE", "INVITEM", "EXTRA",
		},

		{
			"!ENDTRNS",
		},
	}

	dataSeted = [][]interface{}{
		{
			// "!TRNS", "TRNSID", "TRNSTYPE",
			"TRNS", "$_id", "INVOICE",
			// Date
			func(record map[string]interface{}) string {
				return utils.InterfaceToTime(record["created_at"]).Format("01/02/06")
			},
			// "ACCNT", 		"NAME", "CLASS", "AMOUNT", "DOCNUM", "MEMO", "CLEAR",
			"Accounts Receivable", "$customer_name", "", "$grand_total", "$increment_id", "$customer_email", "N",
			// "ADDR1",
			func(record map[string]interface{}) string {
				address := utils.InterfaceToMap(record["billing_address"])
				return utils.InterfaceToString(address["first_name"]) + utils.InterfaceToString(address["last_name"])
			},
			// "ADDR2",
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["billing_address"])["address_line1"])
			},
			//	"ADDR3",
			func(record map[string]interface{}) string {
				address := utils.InterfaceToMap(record["billing_address"])
				return utils.InterfaceToString(address["city"]) + ", " + utils.InterfaceToString(address["state"]) +
					" " + utils.InterfaceToString(address["zip_code"])
			},
			// "ADDR4",
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["billing_address"])["country"])
			},
			// "SADDR1",
			func(record map[string]interface{}) string {
				address := utils.InterfaceToMap(record["shipping_address"])
				return utils.InterfaceToString(address["first_name"]) + utils.InterfaceToString(address["last_name"])
			},
			// "SADDR2",
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["shipping_address"])["address_line1"])
			},
			//	"SADDR3",
			func(record map[string]interface{}) string {
				address := utils.InterfaceToMap(record["shipping_address"])
				return utils.InterfaceToString(address["city"]) + ", " + utils.InterfaceToString(address["state"]) +
					" " + utils.InterfaceToString(address["zip_code"])
			},
			// "SADDR4",
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["shipping_address"])["country"])
			},

			// PAID Y or N depends from status
			func(record map[string]interface{}) string {
				status := utils.InterfaceToString(record["status"])
				if status == order.ConstOrderStatusCompleted || status == order.ConstOrderStatusProcessed {
					return "Y"
				}
				return "N"
			},
		},
		{
			// "!SPL", "SPLID", "TRNSTYPE"
			"SPL", "$_id", "INVOICE",
			// Date
			func(record map[string]interface{}) string {
				return utils.InterfaceToTime(record["created_at"]).Format("01/02/06")
			},
			// "ACCNT", "NAME", "CLASS"
			"Income", "", "",
			// Amount with negative value
			func(record map[string]interface{}) string {
				return "-" + utils.InterfaceToString(utils.InterfaceToFloat64(record["grand_total"])-utils.InterfaceToFloat64(record["tax_amount"]))
			},
			// "DOCNUM", "MEMO", "CLEAR", "QNTY", "PRICE", "INVITEM", "EXTRA"
			"$increment_id", "$description", "N", "1", "", "Product", "$description",
		},
		{
			// "!SPL", "SPLID", "TRNSTYPE"
			"SPL", "$_id", "INVOICE",
			// Date
			func(record map[string]interface{}) string {
				return utils.InterfaceToTime(record["created_at"]).Format("01/02/06")
			},
			// "ACCNT", "NAME", "CLASS"
			"Sales Tax Payable", "Sales Tax Vendor", "",
			// Amount with negative value
			func(record map[string]interface{}) string {
				return "-" + utils.InterfaceToString(record["tax_amount"])
			},
			// "DOCNUM",
			"$increment_id",
			// MEMO - description for taxes
			func(record map[string]interface{}) string {
				result := "Taxes: "
				for _, taxValue := range utils.InterfaceToArray(record["taxes"]) {
					result += utils.InterfaceToString(utils.InterfaceToMap(taxValue)["Code"]) + " "
				}
				return result
			},
			// "CLEAR", "QNTY", "PRICE", "INVITEM", "EXTRA"
			"N", "", "", "", "Tax",
		},
		{
			"ENDTRNS",
		},
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
