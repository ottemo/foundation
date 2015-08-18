package trustpilot

import (
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
}

// Function for every day checking for email sent to customers who order is already two week
func schedulerFunc(params map[string]interface{}) error {
	timeDay := time.Hour * 24
	currentTime := time.Now()
	twoWeeksAgo := currentTime.Add(-timeDay * 14)

	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.AddFilter("created_at", ">=", twoWeeksAgo)
	dbOrderCollection.AddFilter("created_at", "<", twoWeeksAgo.Add(timeDay))

	dbRecords, err := dbOrderCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	emailTemplate := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotEmailTemplate))
	for _, dbRecord := range dbRecords {
		currentOrder := utils.InterfaceToMap(dbRecord)
		customInfo := utils.InterfaceToMap(currentOrder["custom_info"])
		emailSent := utils.InterfaceToBool(customInfo[ConstOrderCustomInfoSentKey])
		orderStatus := utils.InterfaceToString(currentOrder["status"])

		if trustpilotLink, present := customInfo[ConstOrderCustomInfoLinkKey]; present && !emailSent && orderStatus != "new" {

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
				env.LogError(err)
				continue
			}

			err = app.SendMail(visitorEmail, ConstEmailSubject, emailToVisitor)
			if err != nil {
				env.LogError(err)
				continue
			}
			customInfo[ConstOrderCustomInfoSentKey] = true

			currentOrder["custom_info"] = customInfo
			_, err = dbOrderCollection.Save(currentOrder)
			if err != nil {
				env.LogError(err)
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
		scheduler.RegisterTask("checkOrdersToSent", schedulerFunc)
		scheduler.ScheduleRepeat("0 9 * * *", "checkOrdersToSent", nil)
	}

	return nil
}
