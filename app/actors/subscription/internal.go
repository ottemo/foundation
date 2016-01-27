package subscription

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
)

// sendNotificationEmail used to send emails in case when payment declined
func sendNotificationEmail(subscriptionInstance subscription.InterfaceSubscription) error {

	email := subscriptionInstance.GetCustomerEmail()
	emailTemplate := utils.InterfaceToString(env.ConfigGetValue(subscription.ConstConfigPathSubscriptionEmailTemplate))
	emailSubject := utils.InterfaceToString(env.ConfigGetValue(subscription.ConstConfigPathSubscriptionEmailSubject))

	if emailTemplate == "" {
		emailTemplate = `Dear {{.Visitor.name}},
		Yours subscription can't be processed couse you have insufficient funds on Credit Card
		please create new subscription using valid credit card.`
	}

	if emailSubject == "" {
		emailSubject = "Subscription"
	}

	templateMap := map[string]interface{}{
		"Visitor":      map[string]interface{}{"name": subscriptionInstance.GetCustomerName()},
		"Subscription": subscriptionInstance.ToHashMap(),
		"Site":         map[string]interface{}{"url": app.GetStorefrontURL("")},
	}
	emailToVisitor, err := utils.TextTemplate(emailTemplate, templateMap)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err = app.SendMail(email, emailSubject, emailToVisitor); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// nextAllowedCreationDate calculates next allowed day value,
// current day - 30.07 -- + 30 = 29.08 if 29.08 > (15.08) - > new next date
// if 01.09 !before 01.09 --> new date
//func nextAllowedCreationDate() time.Time {
//	currentDayWithOffset := time.Now().Truncate(ConstTimeDay).AddDate(0, 0, ConstCreationDaysDelay)
//	if !currentDayWithOffset.Before(nextCreationDate) {
//		nextCreationDate = currentDayWithOffset
//		nextDay := nextCreationDate.Day()
//
//		switch {
//		case nextDay > 15:
//			nextCreationDate = nextCreationDate.AddDate(0, 1, 1-nextDay)
//			break
//		case nextDay > 1:
//			nextCreationDate = nextCreationDate.AddDate(0, 0, 15-nextDay)
//			break
//		}
//	}
//
//	return nextCreationDate
//}

// isSubscriptionDateValid used for validation of new subscription date
// TODO: put logic to handle requirements for it
func validateSubscriptionDate(date time.Time) error {

	if date.Before(time.Now().Truncate(ConstTimeDay)) {
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
func validateSubscriptionPeriod(period int) error {

	for _, allowedValue := range subscription.GetSubscriptionOptionValues() {
		if period == allowedValue {
			return nil
		}
	}

	if period < 0 {
		return nil
	}

	return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "29c73d2f-0c85-4906-95b7-4812542e33a1", "Period value '"+utils.InterfaceToString(period)+"' is not allowed for subscription.")
}

// validateCheckoutToSubscribe check values related to checkout to be valid from checkout submit (cart not checked)
func validateCheckoutToSubscribe(currentCheckout checkout.InterfaceCheckout) error {

	if currentCheckout == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "a7a4b756-e7e3-4902-a07a-d103fd601420", "No checkout")
	}

	customerEmail := ""
	if customer := currentCheckout.GetVisitor(); customer != nil {
		customerEmail = customer.GetEmail()
	}

	if emailValue := currentCheckout.GetInfo("customer_email"); emailValue != nil {
		customerEmail = utils.InterfaceToString(emailValue)
	}

	if !utils.ValidEmailAddress(customerEmail) {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "7c78ba76-647f-4e16-abd3-5e2d5afeb8cf", "Customer email invalid")
	}

	if currentCheckout.GetShippingAddress() == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "e8b6c4cd-123a-4ec4-b413-55e66def1652", "No shipping address")
	}

	if currentCheckout.GetBillingAddress() == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "b5ecb475-cf90-4d56-99e9-2dcc5a772c54", "No billing address")
	}

	if currentCart := currentCheckout.GetCart(); currentCart == nil || currentCart.GetSubtotal() <= 0 {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "054459d2-0b6b-4526-b0a7-92e7dfce43b4", "Cart with items should be provided")
	}

	paymentMethod := currentCheckout.GetPaymentMethod()
	if paymentMethod == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "2c50d7c5-9f14-4ca2-8c79-1a14f068fe75", "Payment method not set")
	}

	if !paymentMethod.IsAllowed(currentCheckout) {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "8d2a66d2-5d47-43e3-bfaa-17ea0e796667", "Payment method not allowed")
	}

	if !paymentMethod.IsTokenable(currentCheckout) {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "dd38d1d7-3921-4da1-8a2a-420067341511", "Payment method not support subsciptions")
	}

	// checking shipping method an shipping rates
	if currentCheckout.GetShippingMethod() == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "b847fd19-81e6-44fa-946b-fc4c7c45a38b", "Shipping method not set")
	}

	if currentCheckout.GetShippingRate() == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "962291f9-3673-4ec6-b95a-a92a80c11eb4", "Shipping rate not set")
	}

	return nil
}

//retrieveCreditCard try to obtain credit card used in checkout or order
func retrieveCreditCard(currentCheckout checkout.InterfaceCheckout, currentOrder order.InterfaceOrder) visitor.InterfaceVisitorCard {

	if currentCheckout != nil {
		paymentMethod := currentCheckout.GetPaymentMethod()
		if creditCardValue := currentCheckout.GetInfo("cc"); paymentMethod != nil && creditCardValue != nil {
			if checkoutCreditCard, ok := creditCardValue.(visitor.InterfaceVisitorCard); ok && checkoutCreditCard != nil {
				if checkoutCreditCard.GetToken() != "" && checkoutCreditCard.GetPaymentMethodCode() == paymentMethod.GetCode() {
					checkoutCreditCard.SetID("")
					return checkoutCreditCard
				}
			}
		}
	}

	if currentOrder == nil {
		return nil
	}

	// trying to obtain payment card of visitor
	if paymentInfo := currentOrder.Get("payment_info"); paymentInfo != nil {
		cardInfoMap := utils.InterfaceToMap(paymentInfo)

		if creditCardID, present := cardInfoMap["creditCardID"]; present {
			orderCreditCard, err := visitor.LoadVisitorCardByID(utils.InterfaceToString(creditCardID))
			if err == nil {
				orderCreditCard.SetID("")
				return orderCreditCard
			}
		}
		orderCreditCard, err := visitor.GetVisitorCardModel()
		if err == nil {
			tokenRecord := map[string]interface{}{
				"payment":         currentOrder.GetPaymentMethod(),
				"type":            cardInfoMap["creditCardType"],
				"number":          cardInfoMap["creditCardNumbers"],
				"expiration_date": cardInfoMap["creditCardExp"],
				"token_id":        cardInfoMap["transactionID"],
			}

			if err = orderCreditCard.FromHashMap(tokenRecord); err == nil {
				return orderCreditCard
			}
		}
	}

	return nil
}
