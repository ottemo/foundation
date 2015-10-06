package giftcard

import (
	"time"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GiftCardSendTask send email with purchased gift cards info
func GiftCardSendTask(params map[string]interface{}) error {
	giftCardCollection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// send an email if we find gift cards in the order
	// use the email template from the configuration value
	giftCardTemplateEmail := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathGiftEmailTemplate))
	if giftCardTemplateEmail == "" {
		return env.ErrorDispatch(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "cccda57c-7f13-48be-a1fb-d29e545051ce", "No GiftCard Email Template found."))
	}

	currentTime := time.Now()

	if err := giftCardCollection.AddFilter("delivery_date", "<=", currentTime); err != nil {
		return env.ErrorDispatch(err)
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
			"Amount":        utils.RoundPrice(utils.InterfaceToFloat64(record["amount"])),
			"Code":          utils.InterfaceToString(record["code"]),
			"RecieverName":  utils.InterfaceToString(record["name"]),
			"RecieverEmail": giftCardRecipientEmail,
			"Message":       utils.InterfaceToString(record["message"]),
		}

		buyerInfo := make(map[string]interface{})
		if orderID, present := record["order_id"]; present {
			currentOrder, err := order.LoadOrderByID(utils.InterfaceToString(orderID))
			if err == nil {
				buyerInfo["Name"] = currentOrder.Get("customer_name")
				buyerInfo["Email"] = currentOrder.Get("customer_email")
			}
		}

		giftCardEmail, err := utils.TextTemplate(giftCardTemplateEmail,
			map[string]interface{}{
				"Visitor":  buyerInfo,
				"GiftCard": giftCardInfo,
				"Site":     customInfo,
			})

		if err != nil {
			env.LogError(err)
			continue
		}

		err = app.SendMail(giftCardRecipientEmail, giftCardEmailSubject, giftCardEmail)
		if err != nil {
			env.LogError(err)
			continue
		}

		record["status"] = ConstGiftCardStatusDelivered

		_, err = giftCardCollection.Save(record)
		if err != nil {
			env.LogError(err)
		}
	}

	return nil
}
