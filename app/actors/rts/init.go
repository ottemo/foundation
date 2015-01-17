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

	seconds := time.Millisecond * ConstVisitorOnlineSeconds * 1000
	ticker := time.NewTicker(seconds)
	go func() {
		for _ = range ticker.C {
			for sessionID := range OnlineSessions {
				delta := time.Now().Sub(OnlineSessions[sessionID].time)
				if float64(ConstVisitorOnlineSeconds) < delta.Seconds() {
					DecreaseOnline(OnlineSessions[sessionID].referrerType)
					delete(OnlineSessions, sessionID)
				}
			}
		}
	}()
}

// DB preparations for current model implementation
func initListners() error {

	env.EventRegisterListener("api.rts.visit", referrerHandler)
	env.EventRegisterListener("api.rts.visit", visitsHandler)
	env.EventRegisterListener("api.cart.addToCart", addToCartHandler)
	env.EventRegisterListener("api.checkout.setPayment", reachedCheckoutHandler)
	env.EventRegisterListener("checkout.success", purchasedHandler)
	env.EventRegisterListener("checkout.success", salesHandler)
	env.EventRegisterListener("api.rts.visit", registerVisitorAsOnlineHandler)
	env.EventRegisterListener("api.request", registerVisitorAsOnlineHandler)
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

	collection, err = db.GetCollection(ConstCollectionNameRTSSales)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("product_id", db.ConstTypeID, true)
	collection.AddColumn("count", db.ConstTypeInteger, false)
	collection.AddColumn("range", db.TypeWPrecision(db.ConstTypeVarchar, 21), false)

	collection, err = db.GetCollection(ConstCollectionNameRTSVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("day", db.ConstTypeDatetime, false)
	collection.AddColumn("visitors", db.ConstTypeInteger, false)
	collection.AddColumn("cart", db.ConstTypeInteger, false)
	collection.AddColumn("checkout", db.ConstTypeInteger, false)
	collection.AddColumn("sales", db.ConstTypeInteger, false)
	collection.AddColumn("details", db.ConstTypeText, false)

	return nil
}
