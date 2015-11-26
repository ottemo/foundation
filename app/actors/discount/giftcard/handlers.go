package giftcard

import (
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"strings"
	"time"
)

// orderProceedHandler is fired during order creation to check the cart for
// gift cards, if a card is present add it to the table gift_card.  Next step is
// inspect the checkout object for applied discounts and record usage amounts in db
func orderProceedHandler(event string, eventData map[string]interface{}) bool {

	orderProceed, ok := eventData["order"].(order.InterfaceOrder)
	if !ok {
		env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4bb5d8a8-15bf-42d8-bd1d-1f9e715779e6", "Unable to find an order when firing event, order.proceed."))
		return false
	}

	giftCardCollection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		env.LogError(err)
		return false
	}

	orderID := orderProceed.GetID()

	giftCardCollection.AddFilter("orders_used", "LIKE", orderID)

	orderGiftCardApplying, err := giftCardCollection.Load()
	if err != nil {
		env.LogError(err)
		return false
	}

	// check is discounts are applied to this order if they make change of used gift card's
	orderAppliedDiscounts := orderProceed.GetDiscounts()

	// check used gift card's to update amount or if this procedure was already done
	if len(orderGiftCardApplying) == 0 && len(orderAppliedDiscounts) > 0 {

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
					env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "987929ab-8d20-4413-a0aa-bb4baae02aeb", "Discount code, "+orderAppliedDiscount.Code+" has been over credited."))
					giftCard["amount"] = 0
					giftCard["status"] = ConstGiftCardStatusOverCredited
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

// orderRollbackHandler inspects the order for presence of gift cards in the apply state
// - refill used amount on gift card and change status to 'refilled'
func orderRollbackHandler(event string, eventData map[string]interface{}) bool {

	rollbackOrder, ok := eventData["order"].(order.InterfaceOrder)
	if !ok {
		env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6d674d4d-be5e-42d0-a3d7-b9731dbcc207", "Unable to find an order when firing event, order.rollback."))
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

	return true
}

// checkoutSuccessHandler create gift cards object from placed order
func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {

	orderProceed, ok := eventData["order"].(order.InterfaceOrder)
	if !ok {
		env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4bb5d8a8-15bf-42d8-bd1d-1f9e715779e6", "Unable to find an order when firing event, order.success."))
		return false
	}

	giftCardCollection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		env.LogError(err)
		return false
	}

	// collect necessary info to variables
	// get a customer and his mail to set him as addressee
	visitorID := orderProceed.Get("visitor_id")
	orderID := orderProceed.GetID()

	cartProducts := orderProceed.GetItems()
	giftCardSkuElement := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathGiftCardSKU))

	giftCardsToSendImmediately := make([]string, 0)

	// check cart for gift card's and save in table if they present
	for _, cartItem := range cartProducts {
		giftCardSku := cartItem.GetSku()

		if strings.Contains(giftCardSku, giftCardSkuElement) {

			recipientEmail := utils.InterfaceToString(orderProceed.Get("customer_email"))
			recipientName := orderProceed.Get("customer_name")
			giftCardAmount := float64(0)
			deliveryDate := time.Now()
			currentTime := time.Now()
			customMessage := ""

			// split item SKU with config gift card SKU value and sign "-" take last element
			giftCardSplitedSKU := strings.Split(giftCardSku, giftCardSkuElement)
			giftCardSplitedSKU = strings.Split(giftCardSplitedSKU[len(giftCardSplitedSKU)-1], "-")

			// reed recipient options
			productOptions := cartItem.GetOptions()
			if recipientEmailOption := utils.GetFirstMapValue(productOptions, "Recipient Email", "Email", "recipient_mailbox"); recipientEmailOption != nil {

				recipientEmailOption := utils.InterfaceToMap(recipientEmailOption)
				emailValue, present := recipientEmailOption["value"]

				if present {
					email := utils.InterfaceToString(emailValue)
					if utils.ValidEmailAddress(email) && email != "" {
						recipientEmail = utils.InterfaceToString(emailValue)
					}
				}
			}

			if recipientNameOption := utils.GetFirstMapValue(productOptions, "Recipient", "Recipient Name", "Name", "name"); recipientNameOption != nil {

				recipientNameOption := utils.InterfaceToMap(recipientNameOption)
				nameValue, present := recipientNameOption["value"]

				if present && utils.InterfaceToString(nameValue) != "" {
					recipientName = utils.InterfaceToString(nameValue)
				}
			}

			if customMessageOption := utils.GetFirstMapValue(productOptions, "Message", "Gift Message", "Note", "message"); customMessageOption != nil {
				customMessageOption := utils.InterfaceToMap(customMessageOption)
				messageValue, present := customMessageOption["value"]
				if present {
					customMessage = utils.InterfaceToString(messageValue)
				}
			}

			if deliveryDateOption := utils.GetFirstMapValue(productOptions, "Date", "Delivery Date", "send_date", "Send Date", "date"); deliveryDateOption != nil {
				deliveryDateOption := utils.InterfaceToMap(deliveryDateOption)
				dateValue, present := deliveryDateOption["value"]
				if present && !utils.IsZeroTime(utils.InterfaceToTime(dateValue)) {
					deliveryDate = utils.InterfaceToTime(dateValue)
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
				giftCard["sku"] = giftCardSku

				giftCard["amount"] = giftCardAmount

				giftCard["order_id"] = orderID
				giftCard["visitor_id"] = visitorID

				giftCard["status"] = ConstGiftCardStatusNew
				giftCard["orders_used"] = make(map[string]float64)

				giftCard["name"] = recipientName
				giftCard["message"] = customMessage

				giftCard["recipient_mailbox"] = recipientEmail
				giftCard["delivery_date"] = deliveryDate

				giftCardID, err := giftCardCollection.Save(giftCard)
				if err != nil {
					env.LogError(err)
					return false
				}
				if deliveryDate.Truncate(time.Hour).Before(currentTime) {
					giftCardsToSendImmediately = append(giftCardsToSendImmediately, giftCardID)
				}
			}
		}
	}

	// run SendTask task to send immediately if delivery_date is today's date
	if len(giftCardsToSendImmediately) > 0 {
		params := map[string]interface{}{
			"giftCards":          giftCardsToSendImmediately,
			"ignoreDeliveryDate": true,
		}

		go SendTask(params)
	}

	return true
}
