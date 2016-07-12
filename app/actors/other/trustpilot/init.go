package trustpilot

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
)

// init makes package self-initialization routine before app start
func init() {
	app.OnAppStart(onAppStart)
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)
}

// Function for every day checking for email sent to customers who order is already two week
func schedulerFunc(params map[string]interface{}) error {
	timeDay := time.Hour * 24
	twoWeeksAgo := time.Now().Truncate(timeDay).Add(-timeDay*14 + time.Nanosecond)

	ordersFrom := twoWeeksAgo
	ordersTo := ordersFrom.Add(timeDay)

	// Allows to use params values date_from and date_to to set records selecting range
	// date_to wouldn't be used until date_from is set and by default is calculating as date_from + day
	if fromDate, present := params["date_from"]; present {
		ordersFrom = utils.InterfaceToTime(fromDate)
		ordersTo = ordersFrom.Add(timeDay)

		if toDate, ok := params["date_to"]; ok {
			ordersTo = utils.InterfaceToTime(toDate)
		}
	}

	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.AddFilter("created_at", ">=", ordersFrom)
	dbOrderCollection.AddFilter("created_at", "<", ordersTo)

	validOrderStates := [2]string{order.ConstOrderStatusProcessed, order.ConstOrderStatusCompleted}
	dbOrderCollection.AddFilter("status", "in", validOrderStates)

	// Allows to use params value orders for specifying an array of orders by ID which would be processed
	if ordersID, present := params["orders"]; present {
		orders := utils.InterfaceToArray(ordersID)
		if len(orders) > 0 {
			if ordersFrom == twoWeeksAgo {
				dbOrderCollection.ClearFilters()
			}
			dbOrderCollection.AddFilter("_id", "in", orders)
		}
	}

	dbRecords, err := dbOrderCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	emailTemplate := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotEmailTemplate))
	for _, dbRecord := range dbRecords {
		currentOrder := utils.InterfaceToMap(dbRecord)
		customInfo := utils.InterfaceToMap(currentOrder["custom_info"])
		emailSent := utils.InterfaceToBool(customInfo[ConstOrderCustomInfoSentKey])

		if trustpilotLink, present := customInfo[ConstOrderCustomInfoLinkKey]; present && !emailSent {

			visitorMap := make(map[string]interface{})
			visitorEmail := utils.InterfaceToString(currentOrder["customer_email"])

			visitorMap["name"] = currentOrder["customer_name"]
			visitorMap["link"] = utils.InterfaceToString(trustpilotLink)
			visitorMap["email"] = visitorEmail
			siteMap := map[string]interface{}{
				"url": app.GetStorefrontURL(""),
			}

			emailToVisitor, err := utils.TextTemplate(emailTemplate,
				map[string]interface{}{
					"Visitor": visitorMap,
					"Order":   currentOrder,
					"Site":    siteMap})

			if err != nil {
				env.ErrorDispatch(err)
				continue
			}

			err = app.SendMail(visitorEmail, ConstEmailSubject, emailToVisitor)
			if err != nil {
				env.ErrorDispatch(err)
				continue
			}
			customInfo[ConstOrderCustomInfoSentKey] = true

			currentOrder["custom_info"] = customInfo
			_, err = dbOrderCollection.Save(currentOrder)
			if err != nil {
				env.ErrorDispatch(err)
				continue
			}
		}
	}

	return nil
}

// onAppStart makes module initialization on application startup
func onAppStart() error {

	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)

	if scheduler := env.GetScheduler(); scheduler != nil {
		scheduler.RegisterTask("trustPilotReview", schedulerFunc)
		scheduler.ScheduleRepeat("0 9 * * *", "trustPilotReview", nil)
	}

	return nil
}
