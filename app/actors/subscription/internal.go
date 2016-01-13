package subscription

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
)

// sendConfirmationEmail used to send confirmation and submit emails about subscription change status or to proceed checkout
func sendConfirmationEmail(subscriptionRecord map[string]interface{}, storefrontConfirmationLink, emailTemplate, emailSubject string) error {

	visitorMap := make(map[string]interface{})
	templateMap := make(map[string]interface{})

	customInfo := map[string]interface{}{
		"link": storefrontConfirmationLink,
	}

	templateMap["Info"] = customInfo

	if value, present := subscriptionRecord["order_id"]; present {
		orderID := utils.InterfaceToString(value)

		orderModel, err := order.LoadOrderByID(orderID)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		orderMap := orderModel.ToHashMap()

		var orderItems []map[string]interface{}

		for _, item := range orderModel.GetItems() {
			options := make(map[string]interface{})

			for _, optionKeys := range item.GetOptions() {
				optionMap := utils.InterfaceToMap(optionKeys)
				options[utils.InterfaceToString(optionMap["label"])] = optionMap["value"]
			}
			orderItems = append(orderItems, map[string]interface{}{
				"name":    item.GetName(),
				"options": options,
				"sku":     item.GetSku(),
				"qty":     item.GetQty(),
				"price":   item.GetPrice()})
		}

		orderMap["items"] = orderItems

		templateMap["Order"] = orderMap

		visitorMap = map[string]interface{}{
			"name":  orderModel.Get("customer_name"),
			"email": orderModel.Get("customer_email"),
		}

		templateMap["Visitor"] = visitorMap

	} else {
		visitorID := utils.InterfaceToString(subscriptionRecord["visitor_id"])

		visitorModel, err := visitor.LoadVisitorByID(visitorID)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		visitorMap = map[string]interface{}{
			"name":  visitorModel.GetFullName(),
			"email": visitorModel.GetEmail(),
		}

		templateMap["Visitor"] = visitorMap
	}

	confirmationEmail, err := utils.TextTemplate(emailTemplate, templateMap)

	err = app.SendMail(utils.InterfaceToString(visitorMap["email"]), emailSubject, confirmationEmail)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// sendConfirmationEmail used to send confirmation and submit emails about subscription change status or to proceed checkout
// current day - 30.07 -- + 30 = 29.08 if 29.08 > (15.08) - > new next date
// if 01.09 !before 01.09 --> new date
func nextAllowedCreationDate() time.Time {
	currentDayWithOffset := time.Now().Truncate(ConstTimeDay).AddDate(0, 0, ConstCreationDaysDelay)
	if !currentDayWithOffset.Before(nextCreationDate) {
		nextCreationDate = currentDayWithOffset
		nextDay := nextCreationDate.Day()

		switch {
		case nextDay > 15:
			nextCreationDate = nextCreationDate.AddDate(0, 1, 1-nextDay)
			break
		case nextDay > 1:
			nextCreationDate = nextCreationDate.AddDate(0, 0, 15-nextDay)
			break
		}
	}

	return nextCreationDate
}

// isSubscriptionDateValid used for validation of new subscription date
// TODO: put logic to handle requirements for it
func validateSubscriptionDate(date time.Time) error {

	if date.Before(time.Now()) {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "d3754eb7-3679-4917-a0d9-ed33cb050081", "Subscription Date should be later then today.")
	}
	//
	//	if date.Day() != 15 && date.Day() != 1 {
	//		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "29c73d2f-0c85-4906-95b7-4812542e33a1", "schedule for either the 1st of the month or the 15th of the month")
	//	}

	return nil
}

// isSubscriptionPeriodValid used for validation of subscription period value
// TODO: update this with additional requirements and block map by mutex if it's allowed to change from config
func validateSubscriptionPeriod(days int) error {

	for _, allowedValue := range allowedSubscriptionPeriods {
		if days == allowedValue {
			return nil
		}
	}

	if days < 0 {
		return nil
	}

	return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "29c73d2f-0c85-4906-95b7-4812542e33a1", "Allowed period are: "+utils.InterfaceToString(allowedSubscriptionPeriods))
}

// getPeriodValue used to obtain valid period value from option value
func getPeriodValue(option string) int {

	if value, present := optionValues[option]; present {
		return value
	}

	if validateSubscriptionPeriod(utils.InterfaceToInt(option)) == nil {
		return utils.InterfaceToInt(option)
	}

	return 30
}
