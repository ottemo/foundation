package rts

import (
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine before app start
func init() {
	api.RegisterOnRestServiceStart(setupAPI)
	db.RegisterOnDatabaseStart(setupDB)
	app.OnAppStart(initListners)
	app.OnAppStart(initSalesHistory)
	app.OnAppStart(initStatistic)
	app.OnAppStart(initReferrals)

	seconds := time.Millisecond * ConstVisitorOnlineSeconds * 1000
	ticker := time.NewTicker(seconds)
	go func() {
		for _ = range ticker.C {
			updateSync.Lock()
			for sessionID := range OnlineSessions {
				if OnlineSessions[sessionID] != nil {
					delta := time.Now().Sub(OnlineSessions[sessionID].time)
					if float64(ConstVisitorOnlineSeconds) < delta.Seconds() {
						DecreaseOnline(OnlineSessions[sessionID].referrerType)
						delete(OnlineSessions, sessionID)
					}
				}
			}
			updateSync.Unlock()
		}
	}()
}

// DB preparations for current model implementation
func initListners() error {
	env.EventRegisterListener("api.rts.visit", referrerHandler)
	env.EventRegisterListener("api.rts.visit", visitsHandler)
	env.EventRegisterListener("api.cart.addToCart", addToCartHandler)
	env.EventRegisterListener("api.checkout.visit", visitCheckoutHandler)
	env.EventRegisterListener("api.checkout.setPayment", setPaymentHandler)
	env.EventRegisterListener("checkout.success", purchasedHandler)
	env.EventRegisterListener("checkout.success", salesHandler)
	env.EventRegisterListener("api.request", visitorOnlineActionHandler)

	return nil
}

// DB preparations for current model implementation
func setupDB() error {

	collection, err := db.GetCollection(ConstCollectionNameRTSSalesHistory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("product_id", db.ConstTypeID, true)
	collection.AddColumn("created_at", db.ConstTypeDatetime, false)
	collection.AddColumn("count", db.ConstTypeInteger, false)

	collection, err = db.GetCollection(ConstCollectionNameRTSVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("day", db.ConstTypeDatetime, true)
	collection.AddColumn("visitors", db.ConstTypeInteger, false)
	collection.AddColumn("total_visits", db.ConstTypeInteger, false)
	collection.AddColumn("cart", db.ConstTypeInteger, false)
	collection.AddColumn("visit_checkout", db.ConstTypeInteger, false)
	collection.AddColumn("set_payment", db.ConstTypeInteger, false)
	collection.AddColumn("sales", db.ConstTypeInteger, false)
	collection.AddColumn("sales_amount", db.ConstTypeFloat, false)

	collection, err = db.GetCollection(ConstCollectionNameRTSReferrals)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("referral", db.TypeWPrecision(db.ConstTypeVarchar, 150), false)
	collection.AddColumn("count", db.ConstTypeInteger, false)

	return nil
}
