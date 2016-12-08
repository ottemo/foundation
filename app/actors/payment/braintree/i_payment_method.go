package braintree

import (
	//"github.com/lionelbarrow/braintree-go"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/lionelbarrow/braintree-go"
	"fmt"
)

// GetCode returns payment method code for use in business logic
func (it *BraintreePaymentMethod) GetCode() string {
	return ConstPaymentCode
}

// GetInternalName returns the human readable name of the payment method
func (it *BraintreePaymentMethod) GetInternalName() string {
	return ConstPaymentInternalName
}

// GetName returns the user customized name of the payment method
func (it *BraintreePaymentMethod) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathName))
}

// GetType returns type of payment method according to "github.com/ottemo/foundation/app/models/checkout"
func (it *BraintreePaymentMethod) GetType() string {
	return checkout.ConstPaymentTypePostCC // TODO decide this, other or new type
}

// IsAllowed checks for payment method applicability
func (it *BraintreePaymentMethod) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathBraintreeEnabled))
}

// IsTokenable returns possibility to save token for this payment method
func (it *BraintreePaymentMethod) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return (false)
}

// TODO Authorize makes payment method authorize operation
func (it *BraintreePaymentMethod) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

	fmt.Println("orderInstance: ", orderInstance)
	fmt.Println("orderInstance: ", utils.InterfaceToString(orderInstance))
	fmt.Println("paymentInfo: ", paymentInfo)
	fmt.Println("paymentInfo: ", utils.InterfaceToString(paymentInfo))

	//var transactionID string
	//var visitorCreditCard visitor.InterfaceVisitorCard

	// try to obtain visitor token info
	//if cc, present := paymentInfo["cc"]; present {
	//	if creditCard, ok := cc.(visitor.InterfaceVisitorCard); ok && creditCard != nil {
	//		transactionID = creditCard.GetToken()
	//		visitorCreditCard = creditCard
	//	}
	//}
	//
	//// id presense in credit card means it was saved so we can update token for it
	//if visitorCreditCard != nil && visitorCreditCard.GetID() != "" {
	//	//orderPaymentInfo["creditCardID"] = visitorCreditCard.GetID()
	//
	//	visitorCreditCard.Set("token_id", orderTransactionID)
	//	visitorCreditCard.Set("token_updated", time.Now())
	//	visitorCreditCard.Save()
	//}

	bt := braintree.New(
		braintree.Sandbox,
		"ddxtcwf5n3hvtz3g",
		"cfj6fzzrkc898mm6",
		"24d8738ee7bc4331bbc3bac79f2a54c2",
	)

	//query := new(braintree.SearchQuery)
	//f := query.AddTextField("email")
	//extra := utils.InterfaceToMap(paymentInfo["extra"])
	//f.Is = utils.InterfaceToString(extra["email"])
	//searchResult, err := bt.Customer().Search(query)
	//if err != nil {
	//	return nil, env.ErrorDispatch(err)
	//}
	//var customerPtr *braintree.Customer
	//if len(searchResult.Customers) > 0 {
	//	customerPtr = searchResult.Customers[0]
	//	fmt.Println("\n Found Customer: ", *customerPtr, "\n")
	//	fmt.Println("\n Found Customer: ", utils.InterfaceToString(*customerPtr), "\n")
	//} else {
	//	//return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "87b2f6a4-ee02-4142-9ff9-239e564f5b37", "could not search for a customer")
	//	customerPtr, err = bt.Customer().Create(&braintree.Customer{
	//		Email: utils.InterfaceToString(extra["email"]),
	//	})
	//	if err != nil {
	//		fmt.Println("Customer creation error")
	//		return nil, env.ErrorDispatch(err)
	//	}
	//	fmt.Println("\n New Customer: ", *customerPtr, "\n")
	//	fmt.Println("\n New Customer: ", utils.InterfaceToString(*customerPtr), "\n")
	//}

	//// Check if we are just supposed to create a Customer (aka a token)
	//action := paymentInfo[checkout.ConstPaymentActionTypeKey]
	//isCreateToken := utils.InterfaceToString(action) == checkout.ConstPaymentActionTypeCreateToken
	//if isCreateToken {
	//	extra := utils.InterfaceToMap(paymentInfo["extra"])
	//	visitorEmail := utils.InterfaceToString(extra["email"])
	//
	//	cust, err := bt.Customer().Find(&braintree.Customer{})
	//	if err != nil {
	//		fmt.Println("Customer creation error")
	//		return nil, env.ErrorDispatch(err)
	//	}
	//}
	//
	//cust, err := bt.Customer().Create(&braintree.Customer{})
	//if err != nil {
	//	fmt.Println("Customer creation error")
	//	return nil, env.ErrorDispatch(err)
	//}
	//
	////extra := utils.InterfaceToMap(paymentInfo["extra"])
	////visitorID := utils.InterfaceToString(extra["visitor_id"])
	//ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
	//card, err := bt.CreditCard().Create(&braintree.CreditCard{
	//	CustomerId: 	cust.Id,
	//	Number:         utils.InterfaceToString(ccInfo["number"]),
	//	ExpirationDate: utils.InterfaceToString(ccInfo["expire_month"])+"/"+utils.InterfaceToString(ccInfo["expire_year"]),
	//	CVV:            utils.InterfaceToString(ccInfo["cvc"]),
	//	Options: &braintree.CreditCardOptions{
	//		VerifyCard: true,
	//	},
	//})
	//if err != nil {
	//	fmt.Println("VERIFICATION ERROR for CARD: ", card)
	//	return nil, env.ErrorDispatch(err)
	//}

	//token, err := bt.ClientToken().GenerateWithCustomer((*customerPtr).Id);
	//if err != nil {
	//	return nil, env.ErrorDispatch(err)
	//}

	token, err := bt.ClientToken().Generate();
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	//gateway := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathFoundationURL)) + "checkout/submit-fake"
	formValues := map[string]string{
		"x_amount": fmt.Sprintf("%.2f", orderInstance.GetGrandTotal()),
		"x_session":          utils.InterfaceToString(paymentInfo["sessionID"]),
		//"x_customer_id": (*customerPtr).Id,
	}

	htmlText := `
	<script>
	$.getScript('https://js.braintreegateway.com/web/3.6.0/js/client.min.js', function(response, status){
		console.log('Status: ', status);

		window.braintree.client.create({
			authorization: '`+token+`'
			},function (createErr, clientInstance) {
				if (createErr) { throw new Error(createErr); }

				var data = {
					creditCard: {
						number: '$CC_NUM',
						expirationDate: '$CC_MONTH/$CC_YEAR',
			        		options: {
        						validate: true
      						}
					}
      				};

				clientInstance.request({
    					endpoint: 'payment_methods/credit_cards',
    					method: 'post',
    					data: data
  				}, function (requestErr, response) {
					console.debug("requestErr: ", requestErr);
					if (requestErr) {
						throw new Error(requestErr);
					}

					console.log('Got nonce:', response.creditCards[0].nonce);
					console.log('Response:', response);

					var form = document.createElement('form');

					form.method = "POST";
    					form.action = "`+utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathFoundationURL)) + "/braintree/submit"+`";
    					form.id = "braintreeForm";

					var elementNonce = document.createElement("input");
					elementNonce.name = "nonce";
					elementNonce.value = response.creditCards[0].nonce;
					form.appendChild(elementNonce);

					var elementCardType = document.createElement("input");
					elementCardType.name = "cardType";
					elementCardType.value = response.creditCards[0].details.cardType;
					form.appendChild(elementCardType);

					var elementLastTwo = document.createElement("input");
					elementLastTwo.name = "lastTwo";
					elementLastTwo.value = response.creditCards[0].details.lastTwo;
					form.appendChild(elementLastTwo);
					`

	for key, value := range formValues {
		var elementName = `element`+key;
		htmlText += `
		var `+elementName+` = document.createElement('input');
		`+elementName+`.name = '`+key+`';
		`+elementName+`.value = '`+value+`';
		form.appendChild(`+elementName+`);
		`
	}

	htmlText += `

					document.body.appendChild(form);

					$('#braintreeForm').submit();
    					//form.submit();
				});
			}
		);

		console.log('Client');
	});
	</script>
	`

	env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "FORM: "+htmlText+"\n")

	return api.StructRestRedirect{Result: htmlText, Location: utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathFoundationURL)) + "/checkout/submit"}, nil
}

// Capture makes payment method capture operation
// - at time of implementation this method is not used anywhere
func (it *BraintreePaymentMethod) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e40b3354-5522-4027-b56d-ea9e736f637f", "Not implemented")
}

// Refund will return funds on the given order
// - at time of implementation this method is not used anywhere
func (it *BraintreePaymentMethod) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f28dab69-77a1-438c-a14f-70c700a70c3b", "Not implemented")
}

// Void will mark the order and capture as void
// - at time of implementation this method is not used anywhere
func (it *BraintreePaymentMethod) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "194a4323-4cc4-41a0-ae90-0b578dc4b73a", "Not implemented")
}
