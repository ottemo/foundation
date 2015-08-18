package giftcard

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"strings"
	"time"
)

// orderProceedHandler check cart for gift cards, if present create them into table gift_card and
// check checkout for applied discounts and make change of gift card amount in database
func orderProceedHandler(event string, eventData map[string]interface{}) bool {

	orderProceed, ok := eventData["order"].(order.InterfaceOrder)
	if !ok {
		env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4bb5d8a8-15bf-42d8-bd1d-1f9e715779e6", "order can't be used"))
		return false
	}

	giftCardCollection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		env.LogError(err)
		return false
	}

	// collect necessary info to variables
	// get a customer and his mail to set him as addressee
	giftCardRecipientEmail := utils.InterfaceToString(orderProceed.Get("customer_email"))
	visitorID := orderProceed.Get("visitor_id")
	orderID := orderProceed.GetID()
	cartProducts := orderProceed.GetItems()
	giftCardsSKUElement := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathGiftCardSKU))

	// check is operations with this order was already done
	giftCardCollection.AddFilter("order_id", "=", orderID)

	orderGiftCardCreation, err := giftCardCollection.Load()
	if err != nil {
		env.LogError(err)
		return false
	}

	giftCardCollection.ClearFilters()
	giftCardCollection.AddFilter("orders_used", "LIKE", orderID)

	orderGiftCardApplying, err := giftCardCollection.Load()
	if err != nil {
		env.LogError(err)
		return false
	}

	if len(orderGiftCardCreation) > 0 || len(orderGiftCardApplying) > 0 {
		return true
	}

	// check cart for gift card's and save in table if they present
	for _, cartItem := range cartProducts {
		giftCardSKU := cartItem.GetSku()

		if strings.Contains(giftCardSKU, giftCardsSKUElement) {

			recipientName := orderProceed.Get("customer_name")
			giftCardAmount := float64(0)

			// split item SKU with config gift card SKU value and sign "-" take last element
			giftCardSplitedSKU := strings.Split(giftCardSKU, giftCardsSKUElement)
			giftCardSplitedSKU = strings.Split(giftCardSplitedSKU[len(giftCardSplitedSKU)-1], "-")

			// reed recipient options
			productOptions := cartItem.GetOptions()
			if recipientEmailOption := utils.GetFirstMapValue(productOptions, "Recipient email", "Email", "recipient_mailbox"); recipientEmailOption != nil {

				recipientEmailOption := utils.InterfaceToMap(recipientEmailOption)
				emailValue, present := recipientEmailOption["value"]

				if present {
					giftCardRecipientEmail = utils.InterfaceToString(emailValue)
				}
			}

			if recipientNameOption := utils.GetFirstMapValue(productOptions, "Recipient", "Recipient name", "Name", "name"); recipientNameOption != nil {

				recipientNameOption := utils.InterfaceToMap(recipientNameOption)
				nameValue, present := recipientNameOption["value"]

				if present && utils.InterfaceToString(nameValue) != "" {
					recipientName = utils.InterfaceToString(nameValue)
				}
			}

			// check is result value is number if not take amount as a price of item
			if giftCardAmount = utils.InterfaceToFloat64(giftCardSplitedSKU[len(giftCardSplitedSKU)-1]); giftCardAmount <= 0 {
				giftCardAmount = cartItem.GetPrice()
			}

			for i := 0; i < cartItem.GetQty(); i++ {

				// generate unique code by unix nano time
				giftCardUniqueCode := utils.InterfaceToString(time.Now().UnixNano())

				giftCard := make(map[string]interface{})

				giftCard["code"] = giftCardUniqueCode
				giftCard["name"] = recipientName
				giftCard["sku"] = giftCardSKU

				giftCard["amount"] = giftCardAmount

				giftCard["order_id"] = orderID
				giftCard["visitor_id"] = visitorID

				giftCard["status"] = ConstGiftCardStatusNew
				giftCard["orders_used"] = make(map[string]float64)
				giftCard["recipient_mailbox"] = giftCardRecipientEmail

				if _, err = giftCardCollection.Save(giftCard); err != nil {
					env.LogError(err)
					return false
				}
			}
		}
	}

	// check is discounts are applied to this order if they make change of used gift card's
	orderAppliedDiscounts := orderProceed.GetDiscounts()

	if len(orderAppliedDiscounts) > 0 {

		for _, orderAppliedDiscount := range orderAppliedDiscounts {

			if err := giftCardCollection.ClearFilters(); err != nil {
				env.LogError(err)
			}
			if err := giftCardCollection.AddFilter("code", "=", orderAppliedDiscount.Code); err != nil {
				env.LogError(err)
			}

			records, err := giftCardCollection.Load()
			if err != nil {
				env.LogError(err)
				return false
			}

			// change amount, status and orders_used information for gift card
			if len(records) > 0 {
				giftCard := records[0]

				// calculate the amount that will be on cart after apply and add order used record with orderID and amount
				giftCardAmountAfterApply := utils.InterfaceToFloat64(giftCard["amount"]) - orderAppliedDiscount.Amount

				ordersGiftCardUsedMap := utils.InterfaceToMap(giftCard["orders_used"])
				ordersGiftCardUsedMap[orderID] = orderAppliedDiscount.Amount

				giftCard["amount"] = giftCardAmountAfterApply
				giftCard["status"] = ConstGiftCardStatusApplied

				if giftCardAmountAfterApply < 0 {
					env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "987929ab-8d20-4413-a0aa-bb4baae02aeb", "discount "+orderAppliedDiscount.Code+" is over used"))
					giftCard["amount"] = 0
					giftCard["status"] = ConstGiftCardStatusOverUsed
				}

				if giftCardAmountAfterApply == 0 {
					giftCard["status"] = ConstGiftCardStatusUsed
				}

				giftCard["orders_used"] = ordersGiftCardUsedMap

				_, err := giftCardCollection.Save(giftCard)
				if err != nil {
					env.LogError(err)
					continue
				}
			}
		}
	}

	return true
}

// orderRollbackHandler check order for present gift cards in apply
// - refill used amount on gift card and change status to 'refilled'
// and in order:
// - set status of gift card to cancelled and amount to '0'
func orderRollbackHandler(event string, eventData map[string]interface{}) bool {

	rollbackOrder, ok := eventData["order"].(order.InterfaceOrder)
	if !ok {
		env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6d674d4d-be5e-42d0-a3d7-b9731dbcc207", "order can't be used"))
		return false
	}

	giftCardCollection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		env.LogError(err)
		return false
	}

	// get gift cards that was applied to this order and refill amount that was used in this order
	orderID := rollbackOrder.GetID()

	if err := giftCardCollection.AddFilter("orders_used", "LIKE", orderID); err != nil {
		env.LogError(err)
	}

	records, err := giftCardCollection.Load()
	if err != nil {
		env.LogError(err)
		return false
	}

	// check all records from gift_cards and restoring their balance
	for _, record := range records {

		ordersUsage := utils.InterfaceToMap(record["orders_used"])

		if refillAmount, present := ordersUsage[orderID]; present {

			newAmount := utils.InterfaceToFloat64(refillAmount) + utils.InterfaceToFloat64(record["amount"])

			// refill gift card amount, change status and orders_used information
			delete(ordersUsage, orderID)

			record["status"] = ConstGiftCardStatusRefilled
			record["orders_used"] = ordersUsage
			record["amount"] = newAmount

			_, err := giftCardCollection.Save(record)
			if err != nil {
				env.LogError(err)
				return false
			}
		}
	}

	// get gift cards that was buyed and change their amount to 0 and status to canceled
	if err := giftCardCollection.ClearFilters(); err != nil {
		env.LogError(err)
	}

	if err := giftCardCollection.AddFilter("order_id", "=", orderID); err != nil {
		env.LogError(err)
	}

	records, err = giftCardCollection.Load()
	if err != nil {
		env.LogError(err)
		return false
	}

	for _, record := range records {

		// if gift card is not used, we remove it else change status to cancelled and amount to "0"
		if len(utils.InterfaceToMap(record["orders_used"])) == 0 && record["status"] == ConstGiftCardStatusNew {
			giftCardCollection.DeleteByID(utils.InterfaceToString(record["_id"]))

		} else {
			env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e6d02f68-157f-421c-904d-7f33c8408d6f", "cancell gift card that was already apllied or delivered"))
			record["status"] = ConstGiftCardStatusCanceled
			record["amount"] = 0

			_, err := giftCardCollection.Save(record)
			if err != nil {
				env.LogError(err)
				return false
			}
		}
	}

	return true
}

// checkoutSuccessHandler send email to customer that purchased a gift cards with codes of them
func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {

	orderPlaced, ok := eventData["order"].(order.InterfaceOrder)
	if !ok {
		env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4bb5d8a8-15bf-42d8-bd1d-1f9e715779e6", "order can't be used"))
		return false
	}

	giftCardCollection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		env.LogError(err)
		return false
	}

	// set a customer mail as addressee
	buyerEmail := utils.InterfaceToString(orderPlaced.Get("customer_email"))
	orderID := orderPlaced.GetID()

	// send email if we have gift cards in order
	// mail template get from config value
	giftCardTemplateEmail := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathGiftEmailTemplate))

	if err := giftCardCollection.AddFilter("order_id", "=", orderID); err != nil {
		env.LogError(err)
		return false
	}

	records, err := giftCardCollection.Load()
	if err != nil {
		env.LogError(err)
		return false
	}

	if len(records) > 0 && giftCardTemplateEmail != "" {

		// use buyer info for giving name to gifted person
		buyerInfo := make(map[string]interface{})
		buyerInfo["email"] = buyerEmail
		buyerInfo["name"] = orderPlaced.Get("customer_name")

		customInfo := map[string]interface{}{
			"Url": app.GetStorefrontURL(""),
		}

		giftCardEmailSubject := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathGiftEmailSubject))
		if len(giftCardEmailSubject) < 2 {
			giftCardEmailSubject = "Your giftcard has arrived"
		}

		var giftCardRecipientEmail string

		for _, record := range records {
			giftCardRecipientEmail = utils.InterfaceToString(record["recipient_mailbox"])
			if !utils.ValidEmailAddress(giftCardRecipientEmail) {
				record["recipient_mailbox"] = buyerEmail
				giftCardRecipientEmail = buyerEmail
			}

			buyerInfo["name"] = utils.InterfaceToString(record["name"])

			giftCardInfo := map[string]interface{}{
				"Amount": utils.RoundPrice(utils.InterfaceToFloat64(record["amount"])),
				"Code":   utils.InterfaceToString(record["code"]),
			}

			giftCardsEmail, err := utils.TextTemplate(giftCardTemplateEmail,
				map[string]interface{}{
					"Visitor":  buyerInfo,
					"GiftCard": giftCardInfo,
					"Site":     customInfo,
				})

			if err != nil {
				env.LogError(err)
				continue
			}

			err = app.SendMail(giftCardRecipientEmail, giftCardEmailSubject, giftCardsEmail)
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
	}

	return true
}
