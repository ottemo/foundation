package coupon

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine
func init() {
	instance := new(Coupon)
	var _ checkout.InterfacePriceAdjustment = instance
	if err := checkout.RegisterPriceAdjustment(instance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2f05d681-eb95-43de-a1ed-7ceed18ff378", err.Error())
	}

	db.RegisterOnDatabaseStart(setupDB)
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)

	app.OnAppStart(initListeners)
}

// setupDB prepares system database for package usage
func setupDB() error {

	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddColumn("code", db.ConstTypeVarchar, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3bd7188b-f77b-4d46-b661-87f56fbf6e81", err.Error())
	}
	if err := collection.AddColumn("name", db.ConstTypeVarchar, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a4b7571d-e1b3-400c-89a1-79fcf8255da0", err.Error())
	}
	if err := collection.AddColumn("amount", db.ConstTypeDecimal, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e3f76159-c830-46a0-9d81-850748fa545a", err.Error())
	}
	if err := collection.AddColumn("percent", db.ConstTypeDecimal, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e69fa750-03b9-400c-8643-494bb8c6a1c7", err.Error())
	}
	if err := collection.AddColumn("times", db.ConstTypeInteger, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bc0f9f1f-b086-4c6e-9aef-0a88458436a1", err.Error())
	}
	if err := collection.AddColumn("since", db.ConstTypeDatetime, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "36107122-8e67-4645-9071-e57785e5fe26", err.Error())
	}
	if err := collection.AddColumn("until", db.ConstTypeDatetime, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "962103b3-44af-4175-91da-5bc5268cb667", err.Error())
	}
	if err := collection.AddColumn("limits", db.ConstTypeJSON, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7e3dedc2-4b36-42ec-9810-66c538d23685", err.Error())
	}
	if err := collection.AddColumn("target", db.ConstTypeVarchar, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "84d2f4e4-2e12-4a44-8ca4-79d1000860f6", err.Error())
	}

	return nil
}

// initListeners register event listeners
func initListeners() error {

	return nil
}

// initUsedCoupons adding from orders currently available coupon codes to usedCoupons variable with visitors ID's
func initUsedCoupons() error {
	usedCoupons = make(map[string][]string)

	// loading information about applied discounts
	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		return err
	}

	existingDiscounts, err := collection.Load()

	for _, discount := range existingDiscounts {
		if discountCode, present := discount["code"]; present {
			usedCoupons[utils.InterfaceToString(discountCode)] = make([]string, 0)
		}
	}

	// get orders that created after begin date
	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbOrderCollection := orderCollectionModel.GetDBCollection()
	if err := dbOrderCollection.AddFilter("visitor_id", "!=", nil); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c5addde9-220b-40ae-9e9d-72e7d5608595", err.Error())
	}
	//	filtering for array can't be applied
	//	dbOrderCollection.AddFilter("discounts", "!=", nil)

	orders, err := dbOrderCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// go through orders and append visitors to existing codes
	for _, order := range orders {
		visitorID := utils.InterfaceToString(order["visitor_id"])
		discounts := utils.InterfaceToArray(order["discounts"])

		if len(discounts) > 0 && visitorID != "" {
			for _, discount := range discounts {
				discount := utils.InterfaceToMap(discount)
				discountCode := utils.InterfaceToString(discount["Code"])

				if _, present := usedCoupons[discountCode]; present && !utils.IsInListStr(visitorID, usedCoupons[discountCode]) {
					usedCoupons[discountCode] = append(usedCoupons[discountCode], visitorID)
				}
			}
		}
	}

	return nil
}
