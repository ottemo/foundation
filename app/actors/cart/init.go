package cart

import (
	"fmt"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
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
	app.OnAppStart(scheduleAbandonCartEmails)
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
		scheduler.ScheduleRepeat("0 * * * *", "abandonCartEmail", nil)
	}

	return nil
}

func abandonCartTask(params map[string]interface{}) error {
	fmt.Println("Abandon carts scheduler firing")

	// Get config for time to send
	beforeDate, isEnabled := getConfigSendBefore()
	if !isEnabled {
		// context.SetResponseStatusNotFound()
		// return "endpoint not enabled", nil
		return nil
	}

	resultCarts := getAbandonedCarts(beforeDate)
	actionableCarts := getActionableCarts(resultCarts)

	fmt.Println("carts found in query:", len(resultCarts))
	fmt.Println("carts that are actionable:", len(actionableCarts))

	for _, aCart := range actionableCarts {
		err := sendAbandonEmail(aCart)
		if err != nil {
			continue
		}

		flagCartAsEmailed(aCart.Cart.ID)
	}

	// return actionableCarts, nil
	return nil
}

type AbandonCartEmailData struct {
	Visitor AbandonVisitor
	Cart    AbandonCart
}

type AbandonVisitor struct {
	Email     string
	FirstName string
	LastName  string
}

type AbandonCart struct {
	ID string
	// Items []AbandonCartItem
}

// type AbandonCartItem struct {
// 	Name  string
// 	SKU   string
// 	Price float64
// 	Image string
// }

func getConfigSendBefore() (time.Time, bool) {
	var isEnabled bool
	beforeConfig := utils.InterfaceToInt(env.ConfigGetValue(ConstConfigPathCartAbandonEmailSendTime))

	// Flag it as enabled
	if beforeConfig != 0 {
		isEnabled = true
	}

	// Build out the time to send before, we are expecting a config
	// that is a negative int
	beforeDuration := time.Duration(beforeConfig) * time.Hour
	beforeDate := time.Now().Add(beforeDuration)

	return beforeDate, isEnabled
}

// Get the abandoned carts
// - active
// - were updated in our time frame
// - have not been sent an abandon cart email
func getAbandonedCarts(beforeDate time.Time) []map[string]interface{} {
	dbEngine := db.GetDBEngine()
	cartCollection, _ := dbEngine.GetCollection(ConstCartCollectionName)
	cartCollection.AddFilter("active", "=", true)
	cartCollection.AddFilter("custom_info.is_abandon_email_sent", "!=", true)
	cartCollection.AddFilter("updated_at", "<", beforeDate)
	cartCollection.AddSort("updated_at", true)
	resultCarts, _ := cartCollection.Load()

	return resultCarts
}

func getActionableCarts(resultCarts []map[string]interface{}) []AbandonCartEmailData {
	allCartEmailData := []AbandonCartEmailData{}

	// Determine which carts have an email we can use
	for _, resultCart := range resultCarts {
		var email, firstName, lastName string
		cartID := utils.InterfaceToString(resultCart["_id"])
		sessionID := utils.InterfaceToString(resultCart["session_id"])
		visitorID := utils.InterfaceToString(resultCart["visitor_id"])

		// try to get by visitor_id
		if visitorID != "" {
			vModel, _ := visitor.LoadVisitorByID(visitorID)
			email = vModel.GetEmail()
			firstName = vModel.GetFirstName()
			lastName = vModel.GetLastName()

		} else if sessionID != "" {
			create := false
			sessionWrapper, _ := api.GetSessionService().Get(sessionID, create)
			sCheckout := utils.InterfaceToMap(sessionWrapper.Get(checkout.ConstSessionKeyCurrentCheckout))

			scInfo := utils.InterfaceToMap(sCheckout["Info"])
			email = utils.InterfaceToString(scInfo["customer_email"])
			//NOTE: We have customer_name here as well, which we could split
			//      or we could look to see if the address is filled out yet
		}

		// TODO: if we don't have an email then flag this cart as don't update?

		// no email address for us to contact, move along
		if email == "" {
			continue
		}

		// Assemble the details needed for further actions
		cartEmailData := AbandonCartEmailData{
			Visitor: AbandonVisitor{
				Email:     email,
				FirstName: firstName,
				LastName:  lastName,
			},
			Cart: AbandonCart{
				ID: cartID,
			},
		}

		// NOTE: In v1 we aren't including cart item details
		// Get the cart items for the carts we are about to email
		// cartItemsCollection, err := dbEngine.GetCollection(ConstCartItemsCollectionName)
		// cartItemsCollection.AddFilter("cart_id", "=", it.GetID())
		// cartItems, err := cartItemsCollection.Load()

		allCartEmailData = append(allCartEmailData, cartEmailData)
	}

	return allCartEmailData
}

func sendAbandonEmail(emailData AbandonCartEmailData) error {
	subject := "It looks like you forgot something in your cart"
	template := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathCartAbandonEmailTemplate))
	if template == "" {
		return env.ErrorDispatch(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1756ec63-7cd7-4764-a8ff-64b142fc3f9f", "Abandon cart emails want to send but the template is empty"))
	}

	templateData := utils.InterfaceToMap(emailData)
	templateData["Site"] = map[string]interface{}{
		"Url": app.GetStorefrontURL(""),
	}

	body, err := utils.TextTemplate(template, templateData)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = app.SendMail(emailData.Visitor.Email, subject, body)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

func flagCartAsEmailed(cartID string) {
	iCart, _ := cart.LoadCartByID(cartID)

	info := iCart.GetCustomInfo()
	info["is_abandon_email_sent"] = true
	info["abandon_email_sent_at"] = time.Now()

	iCart.SetCustomInfo(info)
	iCart.Save()
}
