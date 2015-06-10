package giftcard

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("giftcard/:giftcode", api.ConstRESTOperationGet, APIGetGiftCard)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("giftcards", api.ConstRESTOperationGet, APIGetGiftCardsList)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("giftcard/:giftcode/apply", api.ConstRESTOperationGet, APIApplyGiftCard)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("giftcard/:giftcode/neglect", api.ConstRESTOperationGet, APINeglectGiftCard)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIGetGiftCard return gift card info buy it's code
func APIGetGiftCard(context api.InterfaceApplicationContext) (interface{}, error) {

	giftCardID := context.GetRequestArgument("giftcode")
	if giftCardID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "06792fd7-c838-4acc-9c6f-cb8fcff833dd", "gift card code was not specified")
	}

	collection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection.AddFilter("code", "=", giftCardID)
	rows, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dd7b2130-b5ed-4b26-b1fc-2d36c3bf147f", "gift card with such code not found")
	}

	return rows[0], nil
}

// APIGetGiftCardsList return list of gift cards
func APIGetGiftCardsList(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	rows, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return rows, nil
}

// APIApplyGiftCard applies gift card to current checkout
//   - Gift Card code should be specified in "giftcode" argument
func APIApplyGiftCard(context api.InterfaceApplicationContext) (interface{}, error) {

	giftCardCode := context.GetRequestArgument("giftcode")

	// getting applied gift codes array for current session
	appliedGiftCardCodes := utils.InterfaceToStringArray(context.GetSession().Get(ConstSessionKeyAppliedGiftCardCodes))

	// checking if gift codes was already applied
	if utils.IsInArray(giftCardCode, appliedGiftCardCodes) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1c310f79-0f79-493a-b761-ad4f24542559", "gift cart already applied")
	}

	// loading gift codes for specified code
	collection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	err = collection.AddFilter("code", "=", giftCardCode)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	records, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// checking and applying obtained gift card codes
	if len(records) == 1 && utils.InterfaceToString(records[0]["code"]) == giftCardCode {
		if utils.InterfaceToFloat64(records[0]["amount"]) <= 0 {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "ce349f59-51c7-43ec-a64c-80f7d4af6d3c", "gift cart amount is '0'")
		}

		// gift card codes is working - applying it
		appliedGiftCardCodes = append(appliedGiftCardCodes, giftCardCode)
		context.GetSession().Set(ConstSessionKeyAppliedGiftCardCodes, appliedGiftCardCodes)

	} else {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2b55d714-2cba-49f8-ad7d-fdc542bfc2a3", "gift cart code not found")
	}

	return "ok", nil
}

// APINeglectGiftCard neglects (un-apply) gift card code promotion to current checkout
//   - gift card code should be specified in "giftcode" argument
//   - use "*" as gift card code to neglect all gift card discounts
func APINeglectGiftCard(context api.InterfaceApplicationContext) (interface{}, error) {

	giftCardID := context.GetRequestArgument("giftcode")
	if giftCardID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e2bad33a-36e7-41d4-aea7-8fe1b97eb31c", "gift card code was not specified")
	}

	if giftCardID == "*" {
		context.GetSession().Set(ConstSessionKeyAppliedGiftCardCodes, make([]string, 0))
		return "ok", nil
	}

	appliedCoupons := utils.InterfaceToStringArray(context.GetSession().Get(ConstSessionKeyAppliedGiftCardCodes))
	if len(appliedCoupons) > 0 {
		var newAppliedCoupons []string
		for _, value := range appliedCoupons {
			if value != giftCardID {
				newAppliedCoupons = append(newAppliedCoupons, value)
			}
		}
		context.GetSession().Set(ConstSessionKeyAppliedGiftCardCodes, newAppliedCoupons)
	}

	return "ok", nil
}
