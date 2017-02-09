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
	if err := models.RegisterModel(order.ConstModelNameOrder, orderInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3fe2c671-f3ef-4621-b783-0ce413abc961", err.Error())
	}

	orderCollectionInstance := new(DefaultOrderCollection)
	var _ order.InterfaceOrderCollection = orderCollectionInstance
	if err := models.RegisterModel(order.ConstModelNameOrderCollection, orderCollectionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "289cc0d5-1d7d-4c68-942d-9c671880841d", err.Error())
	}

	orderItemCollectionInstance := new(DefaultOrderItemCollection)
	var _ order.InterfaceOrderItemCollection = orderItemCollectionInstance
	if err := models.RegisterModel(order.ConstModelNameOrderItemCollection, orderItemCollectionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bb8fb60f-9aee-4a86-a760-8284b25555ae", err.Error())
	}

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

		if err := collection.AddColumn("increment_id", db.TypeWPrecision(db.ConstTypeVarchar, 50), true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "04fcdb69-5246-456d-8b94-ec350d14f1fd", err.Error())
		}
		if err := collection.AddColumn("status", db.TypeWPrecision(db.ConstTypeVarchar, 50), true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "aec7d04f-6aa8-446a-8b7b-dce537caaf75", err.Error())
		}

		if err := collection.AddColumn("visitor_id", db.ConstTypeID, true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "edb575f4-a6e8-49bd-852a-783207a1dc5e", err.Error())
		}
		if err := collection.AddColumn("session_id", db.ConstTypeVarchar, true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b3fd5b45-23d1-4349-9109-b2b53864b59d", err.Error())
		}
		if err := collection.AddColumn("cart_id", db.ConstTypeID, true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5879e9f5-b367-4e52-9f3c-a73e4249f4e6", err.Error())
		}

		if err := collection.AddColumn("billing_address", db.ConstTypeJSON, true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "cbd77fc8-e734-4739-86a1-e4c7a620a2de", err.Error())
		}
		if err := collection.AddColumn("shipping_address", db.ConstTypeJSON, true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c386b21e-8d5c-4d67-bcf3-11c408a39328", err.Error())
		}

		if err := collection.AddColumn("customer_email", db.TypeWPrecision(db.ConstTypeVarchar, 100), true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4aaaa0ae-67f8-4e43-8210-e09d9ef8bdf3", err.Error())
		}
		if err := collection.AddColumn("customer_name", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "873d34bb-f6e4-42c2-8ff8-1db812f3f495", err.Error())
		}

		if err := collection.AddColumn("payment_method", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "593c968d-5d35-4d17-a526-1c4608849e34", err.Error())
		}
		if err := collection.AddColumn("shipping_method", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "09f61aa2-2734-4823-b49d-56fc45b4c6ae", err.Error())
		}

		if err := collection.AddColumn("subtotal", db.ConstTypeMoney, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e65fd188-d67b-4ec8-9fe3-a98979f58d22", err.Error())
		}
		if err := collection.AddColumn("discount", db.ConstTypeMoney, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8bba3c26-7f86-4ff4-b577-07df6b7961fc", err.Error())
		}
		if err := collection.AddColumn("tax_amount", db.ConstTypeMoney, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e87a740e-7bbf-49a2-aace-e2b4aa4fc894", err.Error())
		}
		if err := collection.AddColumn("shipping_amount", db.ConstTypeMoney, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e9bb6bf8-c6cb-4f52-808b-11dc2359cca9", err.Error())
		}
		if err := collection.AddColumn("grand_total", db.ConstTypeMoney, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "628a5e5f-b8e4-4d1a-84a6-c4a1bda772e3", err.Error())
		}

		if err := collection.AddColumn("discounts", db.TypeArrayOf(db.ConstTypeJSON), false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "19937cab-5488-4eea-aed3-c64e19c481b8", err.Error())
		}
		if err := collection.AddColumn("taxes", db.TypeArrayOf(db.ConstTypeJSON), false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4e3ad595-7fd3-4ca7-b695-3fc05935a7ae", err.Error())
		}

		if err := collection.AddColumn("created_at", db.ConstTypeDatetime, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f726846b-fa74-4369-af5a-86a71737c001", err.Error())
		}
		if err := collection.AddColumn("updated_at", db.ConstTypeDatetime, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a1450985-36b9-408a-bdcc-f883606b0460", err.Error())
		}

		if err := collection.AddColumn("description", db.ConstTypeText, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4060eb83-28d8-434b-a867-cfe0bda38b62", err.Error())
		}
		if err := collection.AddColumn("payment_info", db.ConstTypeJSON, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "24d703a1-9caf-4051-883f-414b754f8a2a", err.Error())
		}

		if err := collection.AddColumn("custom_info", db.ConstTypeJSON, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "99bf3cad-3637-4b5b-8fbb-d172512c9732", err.Error())
		}
		if err := collection.AddColumn("shipping_info", db.ConstTypeJSON, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f89dc6e5-74a5-4ff3-8f5a-d8db50f49808", err.Error())
		}

		if err := collection.AddColumn("notes", db.TypeArrayOf(db.ConstTypeVarchar), false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f4484325-2b27-4551-bf8e-f58d6a0c6cd7", err.Error())
		}

		collection, err = dbEngine.GetCollection(ConstCollectionNameOrderItems)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		if err := collection.AddColumn("idx", db.ConstTypeInteger, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "940abf0a-95c0-44f0-98a2-d9838f4cf849", err.Error())
		}

		if err := collection.AddColumn("order_id", db.ConstTypeID, true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fe174102-bc28-434c-8a9b-21e8a17875b3", err.Error())
		}
		if err := collection.AddColumn("product_id", db.ConstTypeID, true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bf075f45-7170-41c7-820a-469a363f136a", err.Error())
		}

		if err := collection.AddColumn("qty", db.ConstTypeInteger, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ae48a68f-a07b-464b-8964-5a07bdcbc811", err.Error())
		}

		if err := collection.AddColumn("name", db.TypeWPrecision(db.ConstTypeVarchar, 150), false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f15de309-1c2a-4922-b5bf-37df07e0b0d6", err.Error())
		}
		if err := collection.AddColumn("sku", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "34e898a6-1655-40fb-b6de-579521c25c59", err.Error())
		}
		if err := collection.AddColumn("short_description", db.ConstTypeVarchar, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4a703ad8-f478-4a3f-a39b-0851db8e82f8", err.Error())
		}

		if err := collection.AddColumn("options", db.ConstTypeJSON, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f9ca8959-044f-480a-b89d-6b86646d6316", err.Error())
		}

		if err := collection.AddColumn("price", db.ConstTypeMoney, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "34658f6d-77bf-40a2-bafe-de50648d09b3", err.Error())
		}
		if err := collection.AddColumn("weight", db.ConstTypeDecimal, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f15e7054-37b3-46ee-a073-d4ff8581fa47", err.Error())
		}

	} else {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "9d0358f1-03ab-44d1-a080-2c62fed5fd81", "Can't get database engine")
	}

	return nil
}
