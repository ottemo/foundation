package order

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
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
		itemOptions := item.GetOptions()
		for key, value := range itemOptions {
			if utils.IsAmongStr(key, "Date", "Delivery Date", "send_date", "Send Date", "date") {
				localizedDate, _ := utils.MakeTZTime(utils.InterfaceToTime(value), timeZone)
				if !utils.IsZeroTime(localizedDate) {
					itemOptions[key] = localizedDate
				}
			}
		}
		items = append(items, item.ToHashMap())
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

	// sending the email notification
	emailAddress := utils.InterfaceToString(visitor["email"])
	err = app.SendMail(emailAddress, "Order confirmation", confirmationEmail)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
