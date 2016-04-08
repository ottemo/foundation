package cart

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultCart)

	ifce := interface{}(instance)
	if _, ok := ifce.(models.InterfaceModel); !ok {
		panic("DefaultCart - InterfaceModel interface not implemented")
	}
	if _, ok := ifce.(models.InterfaceStorable); !ok {
		panic("DefaultCart - InterfaceStorable interface not implemented")
	}
	if _, ok := ifce.(cart.InterfaceCart); !ok {
		panic("DefaultCart - InterfaceCategory interface not implemented")
	}

	models.RegisterModel("Cart", instance)
	db.RegisterOnDatabaseStart(instance.setupDB)

	api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)

	app.OnAppStart(setupEventListeners)
	app.OnAppStart(cleanupGuestCarts)
	//TODO:
	// app.OnAppStart(scheduleAbandonCartEmails)
}

// setupEventListeners registers model related event listeners within system
func setupEventListeners() error {
	// on session close cart model should be also deleted
	sessionCloseListener := func(eventName string, data map[string]interface{}) bool {
		if data != nil {
			if sessionObject, present := data["session"]; present {
				if sessionInstance, ok := sessionObject.(api.InterfaceSession); ok {
					if cartID := sessionInstance.Get(cart.ConstSessionKeyCurrentCart); cartID != nil {

						cartModel, err := cart.GetCartModelAndSetID(utils.InterfaceToString(cartID))
						if err != nil {
							env.LogError(err)
						}

						err = cartModel.Delete()
						if err != nil {
							env.LogError(err)
						}
					}
				}
			}
		}
		return true
	}
	env.EventRegisterListener("session.close", sessionCloseListener)
	return nil
}

// cleanupGuestCarts cleanups guest carts
func cleanupGuestCarts() error {
	cartCollection, err := db.GetCollection(ConstCartCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	cartItemsCollection, err := db.GetCollection(ConstCartItemsCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	cartCollection.AddFilter("visitor_id", "=", nil)
	err = cartCollection.SetResultColumns("_id", "session_id")
	if err != nil {
		return env.ErrorDispatch(err)
	}

	records, err := cartCollection.Load()
	if err != nil {
		env.LogError(err)
	}
	for _, record := range records {
		sessionID := utils.InterfaceToString(record["session_id"])
		if sessionInstance, err := api.GetSessionByID(sessionID, false); err != nil || sessionInstance == nil {
			cartID := utils.InterfaceToString(record["_id"])
			err = cartCollection.DeleteByID(cartID)
			if err != nil {
				env.LogError(err)
			}

			cartItemsCollection.ClearFilters()
			cartItemsCollection.AddFilter("cart_id", "=", cartID)
			_, err = cartItemsCollection.Delete()
			if err != nil {
				env.LogError(err)
			}
		}
	}

	return nil
}

// setupDB prepares system database for package usage
func (it *DefaultCart) setupDB() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		collection, err := dbEngine.GetCollection(ConstCartCollectionName)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		collection.AddColumn("visitor_id", db.ConstTypeID, true)
		collection.AddColumn("session_id", db.ConstTypeID, true)
		collection.AddColumn("updated_at", db.ConstTypeDatetime, true)
		collection.AddColumn("active", db.ConstTypeBoolean, true)
		collection.AddColumn("info", db.ConstTypeJSON, false)
		collection.AddColumn("custom_info", db.ConstTypeJSON, false)

		collection, err = dbEngine.GetCollection(ConstCartItemsCollectionName)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		collection.AddColumn("idx", db.ConstTypeInteger, false)
		collection.AddColumn("cart_id", db.ConstTypeID, true)
		collection.AddColumn("product_id", db.ConstTypeID, true)
		collection.AddColumn("qty", db.ConstTypeInteger, false)
		collection.AddColumn("options", db.ConstTypeJSON, false)

	} else {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "33076d0b-5c65-41dd-aa84-e4b68e1efa5b", "Can't get database engine")
	}

	return nil
}

func scheduleAbandonCartEmails() error {
	if scheduler := env.GetScheduler(); scheduler != nil {
		scheduler.RegisterTask("abandonCartEmail", abandonCartTask)
		//TODO: REVIEW CRON SCHEDULE
		scheduler.ScheduleRepeat("0 9 * * *", "abandonCartEmail", nil)
	}

	return nil
}

func abandonCartTask(params map[string]interface{}) error {

	// get carts
	// - updated in that timeframe
	// - that we haven't emailed
	// - that have an email address

	// send email to each cart
	// - flag the cart as emailed?

	return nil
}
