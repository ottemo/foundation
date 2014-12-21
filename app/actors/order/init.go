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
	// api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		collection, err := dbEngine.GetCollection(ConstCollectionNameOrder)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		collection.AddColumn("increment_id", "varchar(50)", true)
		collection.AddColumn("status", "varchar(50)", true)

		collection.AddColumn("visitor_id", "id", true)
		collection.AddColumn("cart_id", "id", true)

		collection.AddColumn("billing_address", "text", true)
		collection.AddColumn("shipping_address", "text", true)

		collection.AddColumn("customer_email", "varchar(100)", true)
		collection.AddColumn("customer_name", "varchar(100)", false)

		collection.AddColumn("payment_method", "varchar(100)", false)
		collection.AddColumn("shipping_method", "varchar(100)", false)

		collection.AddColumn("subtotal", "decimal(10,2)", false)
		collection.AddColumn("discount", "decimal(10,2)", false)
		collection.AddColumn("tax_amount", "decimal(10,2)", false)
		collection.AddColumn("shipping_amount", "decimal(10,2)", false)
		collection.AddColumn("grand_total", "decimal(10,2)", false)

		collection.AddColumn("created_at", "datetime", false)
		collection.AddColumn("updated_at", "datetime", false)

		collection.AddColumn("description", "text", false)
		collection.AddColumn("payment_info", "text", false)

		collection, err = dbEngine.GetCollection(ConstCollectionNameOrderItems)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		collection.AddColumn("idx", "int", false)

		collection.AddColumn("order_id", "id", true)
		collection.AddColumn("product_id", "id", true)

		collection.AddColumn("qty", "int", false)

		collection.AddColumn("name", "varchar(150)", false)
		collection.AddColumn("sku", "varchar(100)", false)
		collection.AddColumn("short_description", "varchar(255)", false)

		collection.AddColumn("options", "text", false)

		collection.AddColumn("price", "decimal(10,2)", false)
		collection.AddColumn("weight", "decimal(10,2)", false)

	} else {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "9d0358f103ab44d1a0802c62fed5fd81", "Can't get database engine")
	}

	return nil
}
