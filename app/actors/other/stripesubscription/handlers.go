package stripesubscription

import (
	"github.com/ottemo/foundation/app/models/stripesubscription"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"github.com/stripe/stripe-go"
)

// eventCancelHandler processes 'customer.subscription.deleted' event from Stripe
// evt.Data.Obj describes Stripe subscription object https://stripe.com/docs/api#subscription_object
func eventCancelHandler(evt *stripe.Event) error {
	stripeSub := evt.Data.Obj
	stripeSubscriptionCollection, err := getSubscriptionsByStripeCustomerID(utils.InterfaceToString(stripeSub["customer"]))
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, currentSubscription := range stripeSubscriptionCollection.ListSubscriptions() {
		currentSubscription.Set("status", stripeSub["status"])
		if err = currentSubscription.Save(); err != nil {
			return err
		}
	}

	return nil
}

// eventUpdateHandler processes 'customer.subscription.updated' event from Stripe
// evt.Data.Obj describes Stripe subscription object https://stripe.com/docs/api#subscription_object
// updates subscription status, period, renew_notified flag
func eventUpdateHandler(evt *stripe.Event) error {
	stripeSub := evt.Data.Obj
	stripeCustomerID := utils.InterfaceToString(stripeSub["customer"])
	stripeSubscriptionCollection, err := getSubscriptionsByStripeCustomerID(stripeCustomerID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, currentSubscription := range stripeSubscriptionCollection.ListSubscriptions() {
		currentSubscription.Set("period_end", stripeSub["current_period_end"])
		currentSubscription.Set("status", stripeSub["status"])
		if err = currentSubscription.Save(); err != nil {
			return err
		}
	}

	return nil
}

// eventPaymentHandler processes 'invoice.payment_failed' and 'invoice.payment_succeeded' events from Stripe
// evt.Data.Obj describes Stripe invoice object https://stripe.com/docs/api#invoice_object
// saves all invoice data into last_payment_info
func eventPaymentHandler(evt *stripe.Event) error {
	stripeInvoice := evt.Data.Obj
	stripeCustomerID := utils.InterfaceToString(stripeInvoice["customer"])
	stripeSubscriptionCollection, err := getSubscriptionsByStripeCustomerID(stripeCustomerID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, currentSubscription := range stripeSubscriptionCollection.ListSubscriptions() {
		// TODO: what data we want to save here?
		currentSubscription.Set("last_payment_info", stripeInvoice)

		if err = currentSubscription.Save(); err != nil {
			return err
		}

		// TODO: send emails
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
