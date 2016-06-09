package order

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	orderInstance := new(DefaultOrder)
	var _ order.InterfaceOrder = orderInstance
	var _ order.InterfaceOrderItem = new(DefaultOrderItem)
	models.RegisterModel(order.ConstModelNameOrder, orderInstance)

	orderCollectionInstance := new(DefaultOrderCollection)
	var _ order.InterfaceOrderCollection = orderCollectionInstance
	models.RegisterModel(order.ConstModelNameOrderCollection, orderCollectionInstance)

	orderItemCollectionInstance := new(DefaultOrderItemCollection)
	var _ order.InterfaceOrderItemCollection = orderItemCollectionInstance
	models.RegisterModel(order.ConstModelNameOrderItemCollection, orderItemCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	env.RegisterOnConfigStart(setupConfig)

	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		collection, err := dbEngine.GetCollection(ConstCollectionNameOrder)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		collection.AddColumn("increment_id", db.TypeWPrecision(db.ConstTypeVarchar, 50), true)
		collection.AddColumn("status", db.TypeWPrecision(db.ConstTypeVarchar, 50), true)

		collection.AddColumn("visitor_id", db.ConstTypeID, true)
		collection.AddColumn("session_id", db.ConstTypeVarchar, true)
		collection.AddColumn("cart_id", db.ConstTypeID, true)

		collection.AddColumn("billing_address", db.ConstTypeJSON, true)
		collection.AddColumn("shipping_address", db.ConstTypeJSON, true)

		collection.AddColumn("customer_email", db.TypeWPrecision(db.ConstTypeVarchar, 100), true)
		collection.AddColumn("customer_name", db.TypeWPrecision(db.ConstTypeVarchar, 100), false)

		collection.AddColumn("payment_method", db.TypeWPrecision(db.ConstTypeVarchar, 100), false)
		collection.AddColumn("shipping_method", db.TypeWPrecision(db.ConstTypeVarchar, 100), false)

		collection.AddColumn("subtotal", db.ConstTypeMoney, false)
		collection.AddColumn("discount", db.ConstTypeMoney, false)
		collection.AddColumn("tax_amount", db.ConstTypeMoney, false)
		collection.AddColumn("shipping_amount", db.ConstTypeMoney, false)
		collection.AddColumn("grand_total", db.ConstTypeMoney, false)

		collection.AddColumn("discounts", db.TypeArrayOf(db.ConstTypeJSON), false)
		collection.AddColumn("taxes", db.TypeArrayOf(db.ConstTypeJSON), false)

		collection.AddColumn("created_at", db.ConstTypeDatetime, false)
		collection.AddColumn("updated_at", db.ConstTypeDatetime, false)

		collection.AddColumn("description", db.ConstTypeText, false)
		collection.AddColumn("payment_info", db.ConstTypeJSON, false)

		collection.AddColumn("custom_info", db.ConstTypeJSON, false)
		collection.AddColumn("shipping_info", db.ConstTypeJSON, false)

		collection.AddColumn("notes", db.TypeArrayOf(db.ConstTypeVarchar), false)

		collection, err = dbEngine.GetCollection(ConstCollectionNameOrderItems)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		collection.AddColumn("idx", db.ConstTypeInteger, false)

		collection.AddColumn("order_id", db.ConstTypeID, true)
		collection.AddColumn("product_id", db.ConstTypeID, true)

		collection.AddColumn("qty", db.ConstTypeInteger, false)

		collection.AddColumn("name", db.TypeWPrecision(db.ConstTypeVarchar, 150), false)
		collection.AddColumn("sku", db.TypeWPrecision(db.ConstTypeVarchar, 100), false)
		collection.AddColumn("short_description", db.ConstTypeVarchar, false)

		collection.AddColumn("options", db.ConstTypeJSON, false)

		collection.AddColumn("price", db.ConstTypeMoney, false)
		collection.AddColumn("weight", db.ConstTypeDecimal, false)

	} else {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "9d0358f1-03ab-44d1-a080-2c62fed5fd81", "Can't get database engine")
	}

	return nil
}
