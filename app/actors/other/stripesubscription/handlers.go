package stripesubscription

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/stripesubscription"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"github.com/stripe/stripe-go"
)

//eventCancelHandler changes subscription status and sends email to a customer when subscription was canceled
func eventCancelHandler(evt *stripe.Event) error {
	stripeCustomerID := utils.InterfaceToString(evt.Data.Obj["customer"])
	stripeSubscriptionCollection, err := getSubscriptionsByStripeCustomerID(stripeCustomerID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, currentSubscription := range stripeSubscriptionCollection.ListSubscriptions() {
		status := currentSubscription.Get("status")
		if status != "canceled" && status != "unpaid" {
			currentSubscription.Set("status", evt.Data.Obj["status"])
			currentSubscription.Save()

			email := currentSubscription.GetCustomerEmail()
			emailTemplate := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmailCancelTemplate))
			emailSubject := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmailCancelSubject))

			if emailTemplate == "" {
				emailTemplate = `Dear {{.Visitor.name}},
Your subscription was canceled`
			}
			if emailSubject == "" {
				emailSubject = "Subscription Cancelation"
			}
			templateMap := map[string]interface{}{
				"Visitor": map[string]interface{}{"name": currentSubscription.Get("customer_name")},
				"Site":    map[string]interface{}{"url": app.GetStorefrontURL("")},
			}
			emailToVisitor, err := utils.TextTemplate(emailTemplate, templateMap)
			if err != nil {
				return env.ErrorDispatch(err)
			}

			if err = app.SendMail(email, emailSubject, emailToVisitor); err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}

	return nil
}

// eventUpdateHandler updates subscription
func eventUpdateHandler(evt *stripe.Event) error {
	stripeCustomerID := utils.InterfaceToString(evt.Data.Obj["customer"])
	stripeSubscriptionCollection, err := getSubscriptionsByStripeCustomerID(stripeCustomerID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, currentSubscription := range stripeSubscriptionCollection.ListSubscriptions() {
		// If subscription current period has been changed
		// we want to send a renewing notify email to a customer
		currPeriodEnd := currentSubscription.GetPeriodEnd()
		newPeriodEnd := utils.InterfaceToTime(evt.Data.Obj["current_period_end"])
		if newPeriodEnd.After(currPeriodEnd) {
			currentSubscription.Set("renew_notified", false)
		}

		currentSubscription.Set("period_end", newPeriodEnd)
		currentSubscription.Set("status", evt.Data.Obj["status"])
		currentSubscription.Save()
	}

	return nil
}

// eventPaymentSucceededHandler saves last payment information
func eventPaymentHandler(evt *stripe.Event) error {
	stripeCustomerID := utils.InterfaceToString(evt.Data.Obj["customer"])
	stripeSubscriptionCollection, err := getSubscriptionsByStripeCustomerID(stripeCustomerID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, currentSubscription := range stripeSubscriptionCollection.ListSubscriptions() {
		currentSubscription.Set("last_payment_info", evt.Data.Obj)
		currentSubscription.Save()
	}

	return nil
}

// getSubscriptionByStripeCustomerID returns subscriptions with specified stripe_customer_id
func getSubscriptionsByStripeCustomerID(stripeCustomerID string) (stripesubscription.InterfaceStripeSubscriptionCollection, error) {
	if stripeCustomerID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1d6956b2-67cd-40ea-ae13-c9505293369e", "Stripe customer ID is empty")
	}
	stripeSubscriptionCollection, err := stripesubscription.GetStripeSubscriptionCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stripeSubscriptionCollection.ListFilterAdd("stripe_customer_id", "=", stripeCustomerID)
	return stripeSubscriptionCollection, nil
}
