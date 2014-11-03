package order

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/env"
)

// module entry point before app start
func init() {
	orderInstance := new(DefaultOrder)
	var _ order.I_Order = orderInstance
	var _ order.I_OrderItem = new(DefaultOrderItem)
	models.RegisterModel(order.MODEL_NAME_ORDER, orderInstance)

	orderCollectionInstance := new(DefaultOrderCollection)
	var _ order.I_OrderCollection = orderCollectionInstance
	models.RegisterModel(order.MODEL_NAME_ORDER_COLLECTION, orderCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	env.RegisterOnConfigStart(setupConfig)

	api.RegisterOnRestServiceStart(setupAPI)
	// api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func setupDB() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		collection, err := dbEngine.GetCollection(COLLECTION_NAME_ORDER)
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
		collection.AddColumn("payment_info", "json", false)

		collection, err = dbEngine.GetCollection(COLLECTION_NAME_ORDER_ITEMS)
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
		return env.ErrorNew("Can't get database engine")
	}

	return nil
}
