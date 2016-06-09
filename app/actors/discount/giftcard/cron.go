package giftcard

import (
	"time"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// SendTask will send email with purchased gift cards info
func SendTask(params map[string]interface{}) error {
	giftCardCollection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// send an email if we find gift cards in the order
	// use the email template from the configuration value
	giftCardTemplateEmail := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathGiftEmailTemplate))
	if giftCardTemplateEmail == "" {
		return env.ErrorDispatch(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "cccda57c-7f13-48be-a1fb-d29e545051ce", "No giftcard email template found."))
	}

	currentTime := time.Now()
	ignoreDeliveryDate := false

	// handle usage of params
	// key "giftCards" is required and should be an array of giftCard ids
	// key "ignoreDeliveryDate" is optional and can be true or false
	if params != nil {
		if giftCardsToSend, present := params["giftCards"]; present {
			giftCardsToSend := utils.InterfaceToArray(giftCardsToSend)
			if len(giftCardsToSend) >= 1 {
				if err := giftCardCollection.AddFilter("_id", "in", giftCardsToSend); err != nil {
					return env.ErrorDispatch(err)
				}
			}
		}

		if ignoreDateValue, present := params["ignoreDeliveryDate"]; present {
			ignoreDeliveryDate = utils.InterfaceToBool(ignoreDateValue)
		}
	}
	if !ignoreDeliveryDate {
		if err := giftCardCollection.AddFilter("delivery_date", "<=", currentTime); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	if err := giftCardCollection.AddFilter("status", "=", ConstGiftCardStatusNew); err != nil {
		return env.ErrorDispatch(err)
	}

	records, err := giftCardCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(records) < 1 {
		return nil
	}

	customInfo := map[string]interface{}{
		"Url": app.GetStorefrontURL(""),
	}

	giftCardEmailSubject := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathGiftEmailSubject))
	if len(giftCardEmailSubject) < 2 {
		giftCardEmailSubject = "Your giftcard has arrived"
	}

	for _, record := range records {

		giftCardRecipientEmail := utils.InterfaceToString(record["recipient_mailbox"])

		giftCardInfo := map[string]interface{}{
			"Amount":         utils.RoundPrice(utils.InterfaceToFloat64(record["amount"])),
			"Code":           utils.InterfaceToString(record["code"]),
			"RecipientName":  utils.InterfaceToString(record["name"]),
			"RecipientEmail": giftCardRecipientEmail,
			"Message":        utils.InterfaceToString(record["message"]),
		}

		buyerInfo := make(map[string]interface{})
		if orderID, present := record["order_id"]; present {
			currentOrder, err := order.LoadOrderByID(utils.InterfaceToString(orderID))
			if err == nil {
				buyerInfo["Name"] = currentOrder.Get("customer_name")
				buyerInfo["Email"] = currentOrder.Get("customer_email")
			}
		}

		recipientInfo := map[string]interface{}{
			"Name":  utils.InterfaceToString(record["name"]),
			"Email": giftCardRecipientEmail,
		}

		giftCardEmail, err := utils.TextTemplate(giftCardTemplateEmail,
			map[string]interface{}{
				"Recipient": recipientInfo,
				"Buyer":     buyerInfo,
				"GiftCard":  giftCardInfo,
				"Site":      customInfo,
			})

		if err != nil {
			env.ErrorDispatch(err)
			continue
		}

		err = app.SendMail(giftCardRecipientEmail, giftCardEmailSubject, giftCardEmail)
		if err != nil {
			env.ErrorDispatch(err)
			continue
		}

		record["status"] = ConstGiftCardStatusDelivered

		_, err = giftCardCollection.Save(record)
		if err != nil {
			env.ErrorDispatch(err)
		}
	}

	return nil
}
