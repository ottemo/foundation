package order

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// SendShippingStatusUpdateEmail will send an email to alert customers their order has been packed and shipped
func (it DefaultOrder) SendShippingStatusUpdateEmail() error {
	subject := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathShippingEmailSubject))
	emailTemplate := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathShippingEmailTemplate))
	timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))

	// Assemble template variables
	orderMap := it.ToHashMap()

	// convert date of order creation to store time zone
	if date, present := orderMap["created_at"]; present {
		convertedDate, _ := utils.MakeTZTime(utils.InterfaceToTime(date), timeZone)
		if !utils.IsZeroTime(convertedDate) {
			orderMap["created_at"] = convertedDate
		}
	}

	templateVariables := map[string]interface{}{
		"Site":  map[string]string{"Url": app.GetStorefrontURL("")},
		"Order": orderMap,
	}

	body, err := utils.TextTemplate(emailTemplate, templateVariables)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	to := utils.InterfaceToString(it.Get("customer_email"))
	if to == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "370e99c1-727c-4ccf-a004-078d4ab343c7", "Couldn't figure out who to send a shipping status update email to. order_id: "+it.GetID())
	}

	err = app.SendMail(to, subject, body)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// SendOrderConfirmationEmail will send an order confirmation based on the detail of the current order
func (it DefaultOrder) SendOrderConfirmationEmail() error {

	// preparing template object "Info"
	customInfo := make(map[string]interface{})
	customInfo["base_storefront_url"] = utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStorefrontURL))

	// preparing template object "Visitor"
	visitor := make(map[string]interface{})
	visitor["first_name"] = it.Get("customer_name")
	visitor["email"] = it.Get("customer_email")

	// preparing template object "Order"
	order := it.ToHashMap()
	order["payment_method_title"] = it.GetPaymentMethod()
	order["shipping_method_title"] = it.GetShippingMethod()

	// the dates in order should be converted to clients locale
	// TODO: the dates to locale conversion should not happens there - it should be either part of order helper or utilities routine over resulting map
	timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))

	// "created_at" date conversion
	if date, present := order["created_at"]; present {
		convertedDate, _ := utils.MakeTZTime(utils.InterfaceToTime(date), timeZone)
		if !utils.IsZeroTime(convertedDate) {
			order["created_at"] = convertedDate
		}
	}

	// order items extraction
	var items []map[string]interface{}

	for _, item := range it.GetItems() {
		// the item options could also contain the date, which should be converted to local time
		itemMap := item.ToHashMap()
		options := item.GetOptionValues(true) // return "Size": "Small", "Color": "Blue"
		// TODO: this convertation should depend on 'type' of option ('date') or we can add additional method to time object that will return converted value (?)
		for option, value := range options {
			if strings.Index(strings.ToLower(option), "date") >= 0 {
				tempDate, _ := utils.MakeTZTime(utils.InterfaceToTime(value), timeZone)
				options[option] = tempDate
				// format the date if not zero
				if !utils.IsZeroTime(tempDate) {
					options[option] = tempDate.Format("Monday Jan 2 15:04")
				}
			}
		}
		// this will override default options
		itemMap["options"] = options
		items = append(items, itemMap)
	}
	order["items"] = items

	// processing email template
	template := utils.InterfaceToString(env.ConfigGetValue(checkout.ConstConfigPathConfirmationEmail))
	confirmationEmail, err := utils.TextTemplate(template, map[string]interface{}{
		"Order":   order,
		"Visitor": visitor,
		"Info":    customInfo,
	})
	if err != nil {
		return env.ErrorDispatch(err)
	}

	orderID := utils.InterfaceToString(it.Get("_id"))

	storeName := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreName))
	if storeName == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "81e24346-786e-4528-a3d6-85fc514917cc", "store name is not set in config")
	}

	subject := "Your " + storeName + " Order, #" + orderID

	// sending the email notification
	emailAddress := utils.InterfaceToString(visitor["email"])
	err = app.SendMail(emailAddress, subject, confirmationEmail)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
