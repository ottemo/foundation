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
	checkout.RegisterPriceAdjustment(instance)

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

	collection.AddColumn("code", db.ConstTypeVarchar, true)
	collection.AddColumn("name", db.ConstTypeVarchar, false)
	collection.AddColumn("amount", db.ConstTypeDecimal, false)
	collection.AddColumn("percent", db.ConstTypeDecimal, false)
	collection.AddColumn("times", db.ConstTypeInteger, false)
	collection.AddColumn("since", db.ConstTypeDatetime, false)
	collection.AddColumn("until", db.ConstTypeDatetime, false)
	collection.AddColumn("limits", db.ConstTypeJSON, false)
	collection.AddColumn("target", db.ConstTypeVarchar, false)

	return nil
}

// initListeners register event listeners
func initListeners() error {

	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)

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
	dbOrderCollection.AddFilter("visitor_id", "!=", nil)
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
