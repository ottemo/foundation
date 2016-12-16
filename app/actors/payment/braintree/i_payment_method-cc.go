package braintree

import (
	"strings"

	"github.com/lionelbarrow/braintree-go"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
)

// GetCode returns payment method code for use in business logic
func (it *braintreeCCMethod) GetCode() string {
	return constCCMethodCode
}

// GetInternalName returns the human readable name of the payment method
func (it *braintreeCCMethod) GetInternalName() string {
	return constCCMethodInternalName
}

// GetName returns the user customized name of the payment method
func (it *braintreeCCMethod) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(constCCMethodConfigPathName))
}

// GetType returns type of payment method according to "github.com/ottemo/foundation/app/models/checkout"
func (it *braintreeCCMethod) GetType() string {
	return checkout.ConstPaymentTypeCreditCard
}

// IsAllowed checks for payment method applicability
func (it *braintreeCCMethod) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(constCCMethodConfigPathEnabled))
}

// IsTokenable returns possibility to save token for this payment method
func (it *braintreeCCMethod) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return (true)
}

// Authorize makes payment method authorize operations
//  - just create token if set in paymentInfo
//  - otherwise create transaction
func (it *braintreeCCMethod) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

	braintreeInstance := braintree.New(
		braintree.Environment(utils.InterfaceToString(env.ConfigGetValue(constGeneralConfigPathEnvironment))),
		utils.InterfaceToString(env.ConfigGetValue(constGeneralConfigPathMerchantID)),
		utils.InterfaceToString(env.ConfigGetValue(constGeneralConfigPathPublicKey)),
		utils.InterfaceToString(env.ConfigGetValue(constGeneralConfigPathPrivateKey)),
	)

	action := paymentInfo[checkout.ConstPaymentActionTypeKey]
	isCreateToken := utils.InterfaceToString(action) == checkout.ConstPaymentActionTypeCreateToken
	creditCardInfo := paymentInfo["cc"]
	creditCardMap := utils.InterfaceToMap(creditCardInfo)

	if isCreateToken {
		// NOTE: `orderInstance = nil` when creating a token

		// 1. Get our customer token
		extra := utils.InterfaceToMap(paymentInfo["extra"])
		visitorID := utils.InterfaceToString(extra["visitor_id"])
		customerID := getBraintreeCustomerToken(visitorID)

		if customerID == "" {
			// 2. We don't have a braintree client id on file, make a new customer

			var customerParamsPtr *braintree.Customer

			if visitorID == "" {
				var nameParts = strings.SplitN(utils.InterfaceToString(extra["billing_name"])+" ", " ", 2)
				var firstName = strings.TrimSpace(nameParts[0])
				var lastName = strings.TrimSpace(nameParts[1])

				customerParamsPtr = &braintree.Customer{
					FirstName: firstName,
					LastName:  lastName,
					Email:     utils.InterfaceToString(extra["email"]),
				}
			} else {
				visitorData, err := visitor.LoadVisitorByID(visitorID)
				if err != nil {
					return nil, env.ErrorDispatch(err)
				}

				customerParamsPtr = &braintree.Customer{
					FirstName: visitorData.GetFirstName(),
					LastName:  visitorData.GetLastName(),
					Email:     visitorData.GetEmail(),
				}
			}

			customerPtr, err := braintreeInstance.Customer().Create(customerParamsPtr)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			customerID = customerPtr.Id
		}

		// 3. Create a card
		creditCardCVC := utils.InterfaceToString(creditCardMap["cvc"])
		if creditCardCVC == "" {
			return nil, env.ErrorDispatch(env.ErrorNew(constErrorModule, constErrorLevel, "bd0a78bf-065a-462b-92c7-d5a1529797c4", "CVC field was left empty"))
		}

		creditCardParams := &braintree.CreditCard{
			CustomerId:      customerID,
			Number:          utils.InterfaceToString(creditCardMap["number"]),
			ExpirationYear:  utils.InterfaceToString(creditCardMap["expire_year"]),
			ExpirationMonth: utils.InterfaceToString(creditCardMap["expire_month"]),
			CVV:             creditCardCVC,
			Options: &braintree.CreditCardOptions{
				VerifyCard: true,
			},
		}

		createdCreditCard, err := braintreeInstance.CreditCard().Create(creditCardParams)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// This response looks like our normal authorize response
		// but this map is translated into other keys to store a token
		tokenCreationResult := map[string]interface{}{
			"transactionID":      createdCreditCard.Token,                      // token_id
			"creditCardLastFour": createdCreditCard.Last4,                      // number
			"creditCardType":     createdCreditCard.CardType,                   // type
			"creditCardExp":      formatCardExpirationDate(*createdCreditCard), // expiration_date
			"customerID":         createdCreditCard.CustomerId,                 // customer_id
		}

		return tokenCreationResult, nil
	}

	// Charging
	var transaction *braintree.Transaction

	// Token Charge
	// - we have a Customer, and a Card
	// - create a Transaction using Card token
	// - must reference Customer
	// - email is stored on the Customer
	if creditCard, ok := creditCardInfo.(visitor.InterfaceVisitorCard); ok && creditCard != nil {
		var err error
		cardToken := creditCard.GetToken()
		customerID := creditCard.GetCustomerID()

		if cardToken == "" || customerID == "" {
			err := env.ErrorNew(constErrorModule, constErrorLevel, "6b43e527-9bc7-48f7-8cdd-320ceb6d77e6", "looks like we want to charge a token, but we don't have the fields we need")
			return nil, env.ErrorDispatch(err)
		}

		if _, err := braintreeInstance.CreditCard().Find(cardToken); err != nil {
			return nil, env.ErrorDispatch(err)
		}

		transactionParams := &braintree.Transaction{
			Type:               "sale",
			Amount:             braintree.NewDecimal(int64(orderInstance.GetGrandTotal()*100), 2),
			CustomerID:         customerID,
			PaymentMethodToken: cardToken,
			OrderId:            orderInstance.GetID(),

			Options: &braintree.TransactionOptions{
				SubmitForSettlement: true,
				StoreInVault:        true,
			},
		}

		transactionParams.BillingAddress = newBraintreeAddress(orderInstance.GetBillingAddress())
		transactionParams.ShippingAddress = newBraintreeAddress(orderInstance.GetShippingAddress())

		transaction, err = braintreeInstance.Transaction().Create(transactionParams)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

	} else {
		// Regular Charge
		// - don't create a customer, or store a token
		var err error

		creditCardCVC := utils.InterfaceToString(creditCardMap["cvc"])
		if creditCardCVC == "" {
			return nil, env.ErrorDispatch(env.ErrorNew(constErrorModule, constErrorLevel, "7d4c3aca-8c51-4eec-aa7c-bd860944697d", "CVC field was left empty"))
		}

		creditCardParams := &braintree.CreditCard{
			Number:          utils.InterfaceToString(creditCardMap["number"]),
			ExpirationYear:  utils.InterfaceToString(creditCardMap["expire_year"]),
			ExpirationMonth: utils.InterfaceToString(creditCardMap["expire_month"]),
			CVV:             creditCardCVC,
		}

		transactionParams := &braintree.Transaction{
			Type:       "sale",
			Amount:     braintree.NewDecimal(int64(orderInstance.GetGrandTotal()*100), 2),
			CreditCard: creditCardParams,
			OrderId:    orderInstance.GetID(),
			Options: &braintree.TransactionOptions{
				SubmitForSettlement: true,
			},
		}

		transactionParams.BillingAddress = newBraintreeAddress(orderInstance.GetBillingAddress())
		transactionParams.ShippingAddress = newBraintreeAddress(orderInstance.GetShippingAddress())

		transaction, err = braintreeInstance.Transaction().Create(transactionParams)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	// Assemble the response
	paymentResult := map[string]interface{}{
		"transactionID":      transaction.CreditCard.Token,
		"creditCardLastFour": transaction.CreditCard.Last4,
		"creditCardExp":      formatCardExpirationDate(*transaction.CreditCard),
		"creditCardType":     transaction.CreditCard.CardType,
		"customerID":         transaction.Customer.Id,
	}

	return paymentResult, nil
}

// Capture makes payment method capture operation
// - at time of implementation this method is not used anywhere
func (it *braintreeCCMethod) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(constErrorModule, constErrorLevel, "772bc737-f025-4c81-a85a-c10efb67e1b3", " Capture method not implemented")
}

// Refund will return funds on the given order
// - at time of implementation this method is not used anywhere
func (it *braintreeCCMethod) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(constErrorModule, constErrorLevel, "26febf8b-7e26-44d4-bfb4-e9b29126fe5a", "Refund method not implemented")
}

// Void will mark the order and capture as void
// - at time of implementation this method is not used anywhere
func (it *braintreeCCMethod) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(constErrorModule, constErrorLevel, "561e0cc8-3bee-4ec4-bf80-585fa566abd4", "Void method not implemented")
}
