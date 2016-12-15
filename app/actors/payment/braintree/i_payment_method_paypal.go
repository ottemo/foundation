package braintree

import (
	//"github.com/lionelbarrow/braintree-go"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/checkout"
	//"github.com/ottemo/foundation/api"
	//"github.com/ottemo/foundation/app"
	"github.com/lionelbarrow/braintree-go"
	"fmt"
	//"github.com/stripe/stripe-go/customer"
	//"github.com/stripe/stripe-go/card"
	//"github.com/stripe/stripe-go/charge"
	//"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
)

// GetCode returns payment method code for use in business logic
func (it *BraintreePaypalPaymentMethod) GetCode() string {
	return ConstPaypalPaymentCode
}

// GetInternalName returns the human readable name of the payment method
func (it *BraintreePaypalPaymentMethod) GetInternalName() string {
	return ConstPaypalPaymentInternalName
}

// GetName returns the user customized name of the payment method
func (it *BraintreePaypalPaymentMethod) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPaypalName))
}

// GetType returns type of payment method according to "github.com/ottemo/foundation/app/models/checkout"
func (it *BraintreePaypalPaymentMethod) GetType() string {
	return checkout.ConstPaymentTypeScript
}

// IsAllowed checks for payment method applicability
func (it *BraintreePaypalPaymentMethod) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathBraintreePaypalEnabled))
}

// IsTokenable returns possibility to save token for this payment method
func (it *BraintreePaypalPaymentMethod) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return (false)
}

// TODO Authorize makes payment method authorize operation
func (it *BraintreePaypalPaymentMethod) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	fmt.Println("BraintreePaypalPaymentMethod) Authorize")
	//fmt.Println("orderInstance: ", orderInstance)
	//fmt.Println("orderInstance: ", utils.InterfaceToString(orderInstance))
	//fmt.Println("paymentInfo: ", paymentInfo)
	//fmt.Println("paymentInfo: ", utils.InterfaceToString(paymentInfo))

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
		braintree.Sandbox,                  // "sandbox"
		"ddxtcwf5n3hvtz3g",                 // MerchantId
		"cfj6fzzrkc898mm6",                 // pubKey
		"24d8738ee7bc4331bbc3bac79f2a54c2", // privKey
	)
	_ = bt

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

	//action := paymentInfo[checkout.ConstPaymentActionTypeKey]
	//isCreateToken := utils.InterfaceToString(action) == checkout.ConstPaymentActionTypeCreateToken
	//fmt.Println("\n--- isCreateToken: ", isCreateToken, "\n\n", utils.InterfaceToString(isCreateToken))
	//if isCreateToken {
	//	// NOTE: `orderInstance = nil` when creating a token
	//
	//	// 1. Get our customer token
	//	extra := utils.InterfaceToMap(paymentInfo["extra"])
	//	visitorID := utils.InterfaceToString(extra["visitor_id"])
	//	stripeCID := getBraintreeCustomerToken(visitorID)
	//	fmt.Println("\n--- stripeCID: ", stripeCID, "\n\n", utils.InterfaceToString(stripeCID))
	//	if stripeCID == "" {
	//
	//		// 2. We don't have a stripe id on file, make a new customer
	//		//var customerPtr *braintree.Customer
	//		customerPtr, err := bt.Customer().Create(&braintree.Customer{
	//			Email: utils.InterfaceToString(extra["email"]),
	//		})
	//		//if err != nil {
	//		//	fmt.Println("Customer creation error")
	//		//	return nil, env.ErrorDispatch(err)
	//		//}
	//		//c, err := customer.New(&stripe.CustomerParams{
	//		//	Email: utils.InterfaceToString(extra["email"]),
	//		//	// TODO: coupons?
	//		//})
	//		if err != nil {
	//			fmt.Println("NO customerPtr")
	//			return nil, env.ErrorDispatch(err)
	//		}
	//		fmt.Println("\n--- customerPtr: ", customerPtr, "\n\n", utils.InterfaceToString(customerPtr))
	//
	//		stripeCID = customerPtr.Id
	//	}
	//
	//	// 3. Create a card
	//	ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
	//	ccInfo["billing_name"] = extra["billing_name"]
	//	//cp, err := getCardParams(ccInfo, stripeCID)
	//
	//	ccCVC := utils.InterfaceToString(ccInfo["cvc"])
	//	if ccCVC == "" {
	//		err := env.ErrorNew(ConstErrorModule, 1, "15edae76-1d3e-4e7a-a474-75ffb61d26cb", "CVC field was left empty")
	//		return nil, env.ErrorDispatch(err)
	//	}
	//
	//	cp := &braintree.CreditCard{
	//			CustomerId: 	stripeCID,
	//			Number:         utils.InterfaceToString(ccInfo["number"]),
	//			ExpirationYear:utils.InterfaceToString(ccInfo["expire_year"]),
	//			ExpirationMonth: utils.InterfaceToString(ccInfo["expire_month"]),
	//			CVV:            ccCVC,
	//			Options: &braintree.CreditCardOptions{
	//				VerifyCard: true,
	//			},
	//		}
	//	fmt.Println("\n--- cp: ", cp, "\n\n", utils.InterfaceToString(cp))
	//	//cp := &stripe.CardParams{
	//	//	Number: utils.InterfaceToString(ccInfo["number"]),
	//	//	Month:  utils.InterfaceToString(ccInfo["expire_month"]),
	//	//	Year:   utils.InterfaceToString(ccInfo["expire_year"]),
	//	//	CVC:    ccCVC, // Optional, highly recommended
	//	//
	//	//	// might not be passed in
	//	//	Customer: stripeCID,
	//	//	Name:     utils.InterfaceToString(ccInfo["billing_name"]), // Optional
	//	//
	//	//	// Address fields can be passed here as well to aid in fraud prevention
	//	//}
	//	//
	//	//if err != nil {
	//	//	return nil, env.ErrorDispatch(err)
	//	//}
	//
	//	//ca, err := card.New(cp)
	//	ca, err := bt.CreditCard().Create(cp)
	//	// env.LogEvent(env.LogFields{"api_response": ca, "err": err}, "card")
	//	if err != nil {
	//		fmt.Println("\nERROR bt.CreditCard().Create(cp)\n", err, "\n")
	//		return nil, env.ErrorDispatch(err)
	//	}
	//
	//	fmt.Println("\n--- ca: ", ca, "\n\n", utils.InterfaceToString(ca))
	//
	//	// This response looks like our normal authorize response
	//	// but this map is translated into other keys to store a token
	//	result := map[string]interface{}{
	//		"transactionID":      ca.Token,                        // token_id
	//		"creditCardLastFour": ca.Last4,                  // number
	//		"creditCardType":     ca.CardType, // type
	//		"creditCardExp":      formatCardExp(*ca),           // expiration_date
	//		"customerID":         ca.CustomerId,                    // customer_id
	//	}
	//
	//	fmt.Println("\n--- result: ", result, "\n\n", utils.InterfaceToString(result))
	//
	//	return result, nil
	//}
	//
	//// Charging: https://stripe.com/docs/api/go#create_charge
	////var ch *stripe.Charge
	//var tr *braintree.Transaction
	//ccInfo := paymentInfo["cc"]
	////ccInfoMap := utils.InterfaceToMap(ccInfo)
	//
	//// Token Charge
	//// - we have a Customer, and a Card
	//// - create a Charge with the Card as the Source
	//// - must reference Customer
	//// - email is stored on the Customer
	//if creditCard, ok := ccInfo.(visitor.InterfaceVisitorCard); ok && creditCard != nil {
	//	fmt.Println("\n--- creditCard, ok: ", utils.InterfaceToString(creditCard), "\n")
	//	var err error
	//	cardID := creditCard.GetToken()
	//	stripeCID := creditCard.GetCustomerID()
	//
	//	if cardID == "" || stripeCID == "" {
	//		fmt.Println("cardID == '' || stripeCID == ''")
	//		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "02128bc6-83d6-4c12-ae90-900a94adb3ad", "looks like we want to charge a token, but we don't have the fields we need")
	//		return nil, env.ErrorDispatch(err)
	//	}
	//
	//	//chParams := stripe.ChargeParams{
	//	//	Currency: "usd",
	//	//	Amount:   uint64(orderInstance.GetGrandTotal() * 100), // Amount is in cents
	//	//	Customer: stripeCID,                                   // Mandatory
	//	//}
	//	//ccCVC := utils.InterfaceToString(ccInfoMap["cvc"])
	//	//if ccCVC == "" {
	//	//	fmt.Println("ccCVC == ''")
	//	//	err := env.ErrorNew(ConstErrorModule, 1, "15edae76-1d3e-4e7a-a474-75ffb61d26cb", "CVC field was left empty")
	//	//	return nil, env.ErrorDispatch(err)
	//	//}
	//
	//	cc, err := bt.CreditCard().Find(cardID)
	//	if err != nil {
	//		fmt.Println("\n--- Can not find cc.")
	//		return nil, env.ErrorDispatch(err)
	//	}
	//	fmt.Println("\n--- found creditCard, ok: ", utils.InterfaceToString(cc), "\n")
	//
	//	//cc := &braintree.CreditCard{
	//	//	//Number:         utils.InterfaceToString(ccInfoMap["number"]),
	//	//	//ExpirationYear:utils.InterfaceToString(ccInfoMap["expire_year"]),
	//	//	//ExpirationMonth: utils.InterfaceToString(ccInfoMap["expire_month"]),
	//	//	//CVV:            ccCVC,
	//	//	//ExpirationYear:  "25",
	//	//	//CustomerId:stripeCID,
	//	//	Token:cardID,
	//	//	Options: &braintree.CreditCardOptions{
	//	//		VerifyCard: true,
	//	//		//FailOnDuplicatePaymentMethod: true,
	//	//	},
	//	//}
	//
	//	tx := &braintree.Transaction{
	//		Type: "sale",
	//		//Amount: uint64(orderInstance.GetGrandTotal() * 100),
	//		Amount: braintree.NewDecimal(int64(orderInstance.GetGrandTotal() * 100), 2),
	//		//CustomerID: utils.InterfaceToString(requestData["x_customer_id"]),
	//		//PaymentMethodNonce: utils.InterfaceToString(requestData["nonce"]),
	//		CustomerID:stripeCID,
	//		//CreditCard: cc,
	//		PaymentMethodToken:cardID,
	//
	//		Options: &braintree.TransactionOptions{
	//			SubmitForSettlement: true,
	//			StoreInVault: true,
	//		},
	//	}
	//	//chParams.SetSource(cardID)
	//	fmt.Println("\n--- tx: ", tx, "\n\n", utils.InterfaceToString(tx))
	//
	//	//ch, err = charge.New(&chParams)
	//	tr, err = bt.Transaction().Create(tx)
	//	if err != nil {
	//		return nil, env.ErrorDispatch(err)
	//	}
	//
	//	fmt.Println("\n--- tr: ", tr, "\n\n", utils.InterfaceToString(tr))
	//
	//} else {
	//	//fmt.Println("Regular Charge STOP")
	//	//return nil, env.ErrorDispatch(*new(error))
	//	//// Regular Charge
	//	//// - don't create a customer, or store a token
	//	//// - email is stored on the charge's meta hashmap
	//	var err error
	//	//chargeParams := stripe.ChargeParams{
	//	//	Currency: "usd",
	//	//	Amount:   uint64(orderInstance.GetGrandTotal() * 100), // Amount is in cents
	//	//}
	//
	//	//// Must attach either `customer` or `source` to charge
	//	//// source can be either a `token` or `cardParams`
	//	ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
	//	//
	//	if ba := orderInstance.GetBillingAddress(); ba != nil {
	//		ccInfo["billing_name"] = ba.GetFirstName() + " " + ba.GetLastName()
	//	}
	//	//
	//	ccCVC := utils.InterfaceToString(ccInfo["cvc"])
	//	if ccCVC == "" {
	//		err := env.ErrorNew(ConstErrorModule, 1, "15edae76-1d3e-4e7a-a474-75ffb61d26cb", "CVC field was left empty")
	//		return nil, env.ErrorDispatch(err)
	//	}
	//
	//	cp := &braintree.CreditCard{
	//		//CustomerId: 	stripeCID,
	//		Number:         utils.InterfaceToString(ccInfo["number"]),
	//		ExpirationYear:utils.InterfaceToString(ccInfo["expire_year"]),
	//		ExpirationMonth: utils.InterfaceToString(ccInfo["expire_month"]),
	//		CVV:            ccCVC,
	//		//Options: &braintree.CreditCardOptions{
	//		//	VerifyCard: true,
	//		//},
	//	}
	//	fmt.Println("\n--- cp: ", cp, "\n\n", utils.InterfaceToString(cp))
	//	//chargeParams.SetSource(cp)
	//
	//	tx := &braintree.Transaction{
	//		Type: "sale",
	//		//Amount: uint64(orderInstance.GetGrandTotal() * 100),
	//		Amount: braintree.NewDecimal(int64(orderInstance.GetGrandTotal() * 100), 2),
	//		//CustomerID: utils.InterfaceToString(requestData["x_customer_id"]),
	//		//PaymentMethodNonce: utils.InterfaceToString(requestData["nonce"]),
	//		//CustomerID:stripeCID,
	//		//Customer:&braintree.Customer{
	//		//	Email:	utils.InterfaceToString(orderInstance.Get("customer_email")),
	//		//},
	//		CreditCard: cp,
	//		//PaymentMethodToken:cardID,
	//
	//		Options: &braintree.TransactionOptions{
	//			SubmitForSettlement: true,
	//			//StoreInVault: true,
	//		},
	//	}
	//
	//	//chargeParams.AddMeta("email", utils.InterfaceToString(orderInstance.Get("customer_email")))
	//	//
	//	//
	//	//ch, err = charge.New(&chargeParams)
	//	//if err != nil {
	//	//	return nil, env.ErrorDispatch(err)
	//	//}
	//	fmt.Println("\n--- tx: ", tx, "\n\n", utils.InterfaceToString(tx))
	//
	//	//ch, err = charge.New(&chParams)
	//	tr, err = bt.Transaction().Create(tx)
	//	if err != nil {
	//		return nil, env.ErrorDispatch(err)
	//	}
	//
	//	fmt.Println("\n--- tr: ", tr, "\n\n", utils.InterfaceToString(tr))
	//}
	//
	//// Assemble the response
	//fmt.Println("\n--- Assemble the response\n")
	////fmt.Println("\n--- --- tr.CreditCard.Token\n", tr.CreditCard.Token)
	////fmt.Println("\n--- --- tr.CreditCard.Last4\n", tr.CreditCard.Last4)
	////fmt.Println("\n--- --- *tr.CreditCard\n", *tr.CreditCard)
	////fmt.Println("\n--- --- tr.CreditCard.CardType\n", tr.CreditCard.CardType)
	////fmt.Println("\n--- --- tr.Customer.Id\n", tr.Customer.Id)
	//orderPaymentInfo := map[string]interface{}{
	////	"transactionID":     tr.CreditCard.Token,
	////	"creditCardLastFour": tr.CreditCard.Last4,
	////	"creditCardExp":     formatCardExp(*tr.CreditCard),
	////	"creditCardType":    tr.CreditCard.CardType,
	////	"customerID":        tr.Customer.Id,
	////
	//}
	////fmt.Println("\n--- orderPaymentInfo: ", orderPaymentInfo, "\n\n", utils.InterfaceToString(orderPaymentInfo))
	//
	//return orderPaymentInfo, nil

	//
	//return nil, env.ErrorDispatch(*new(error))
	//
	token, err := bt.ClientToken().Generate();
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	//
	////gateway := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathFoundationURL)) + "checkout/submit-fake"
	//formValues := map[string]string{
	//	"x_amount": fmt.Sprintf("%.2f", orderInstance.GetGrandTotal()),
	//	"x_session":          utils.InterfaceToString(paymentInfo["sessionID"]),
	//	//"x_customer_id": (*customerPtr).Id,
	//}
	//
	htmlText := `
	<script>

	function submitBraintreeData(event) {
		console.log("submitBraintreeData");
		event.preventDefault();
		//return;

		braintree.client.create({
		  authorization: '`+token+`'
		}, function (clientErr, clientInstance) {
		  // Create PayPal component
		  console.log("Create PayPal component: ", clientInstance);

		  braintree.paypal.create({
		    client: clientInstance
		  }, function (err, paypalInstance) {
		      console.log("paypalInstance: ", paypalInstance);

		      paypalInstance.tokenize({
			flow: 'checkout', // Required
			//amount: 10.00, // Required
			//currency: 'USD', // Required
			//locale: 'en_US',
			//enableShippingAddress: true,
			//shippingAddressEditable: false,
			//shippingAddressOverride: {
			//  recipientName: 'Scruff McGruff',
			//  line1: '1234 Main St.',
			//  line2: 'Unit 1',
			//  city: 'Chicago',
			//  countryCode: 'US',
			//  postalCode: '60652',
			//  state: 'IL',
			//  phone: '123.456.7890'
			//}
		      }, function (tokenizeErr, tokenizationPayload) {
		      	    console.log("tokenizationPayload: ", tokenizationPayload);

			    if (tokenizeErr) {
			      // Handle tokenization errors or premature flow closure

			      switch (tokenizeErr.code) {
				case 'PAYPAL_POPUP_CLOSED':
				  console.error('Customer closed PayPal popup.');
				  break;
				case 'PAYPAL_ACCOUNT_TOKENIZATION_FAILED':
				  console.error('PayPal tokenization failed. See details:', tokenizeErr.details);
				  break;
				case 'PAYPAL_FLOW_FAILED':
				  console.error('Unable to initialize PayPal flow. Are your options correct?', tokenizeErr.details);
				  break;
				default:
				  console.error('Error!', tokenizeErr);
			      }
			    } else {
				// Tokenization complete
				// Send tokenizationPayload.nonce to server
				console.log("Tokenization complete");
				console.log("Send tokenizationPayload.nonce to server");
			    }
		      });
		  });
		});
	}

	function createBraintreeForm(){
		//var paypalButton = document.querySelector('.paypal-button');
		//paypalButton.addEventListener('click', submitBraintreeData);

		//var form = document.createElement('form');
		//
		////form.method = "POST";
		//form.action = "#";
		//form.id = "braintreeForm";
		//document.body.appendChild(form);
		////$('#braintreeForm').onsubmit = submitBraintreeData;
		////$('#braintreeForm').submit(submitBraintreeData);
		////$('#braintreeForm').submit(submitBraintreeData);
		////$('#braintreeForm').submit();

		var form = document.createElement('form');
		form.id = "braintreeForm";
		node = document.createElement("input");
		node.type = "submit";
		form.appendChild(node);
		//form.style.display = "none";
		document.body.appendChild(form);
		form.onsubmit = submitBraintreeData;
		// emulate user interaction to submit form - or braintree will generate MERCHANT error
		form.querySelector('input[type="submit"]').click()
	}

	$.getScript('https://js.braintreegateway.com/web/3.6.2/js/client.min.js')
		.done(function(script, status) {
			console.log("Success 1: ", status);

			$.getScript('https://js.braintreegateway.com/web/3.6.2/js/paypal.min.js')
				.done(function(script, status) {
					console.log("Success 2: ", status);

					createBraintreeForm();
				})
				.fail(function(jqxhr, settings, exception) {
					console.log("Fail 2: ", exception);
				});
		})
		.fail(function(jqxhr, settings, exception) {
			console.log("Fail 1: ", exception);
		});

	</script>
	`
	//htmlText := `$.getScript('https://js.braintreegateway.com/web/3.6.0/js/client.min.js', function(response, status){
	//	console.log('Status: ', status);
	//
	//	window.braintree.client.create({
	//		authorization: '`+token+`'
	//		},function (createErr, clientInstance) {
	//			if (createErr) { throw new Error(createErr); }`
	//
	//			var data = {
	//				creditCard: {
	//					number: '$CC_NUM',
	//					expirationDate: '$CC_MONTH/$CC_YEAR',
	//		        		options: {
        	//					validate: true
      	//					}
	//				}
      	//			};
	//
	//			clientInstance.request({
    	//				endpoint: 'payment_methods/credit_cards',
    	//				method: 'post',
    	//				data: data
  	//			}, function (requestErr, response) {
	//				console.debug("requestErr: ", requestErr);
	//				if (requestErr) {
	//					throw new Error(requestErr);
	//				}
	//
	//				console.log('Got nonce:', response.creditCards[0].nonce);
	//				console.log('Response:', response);
	//
	//				var form = document.createElement('form');
	//
	//				form.method = "POST";
    	//				form.action = "`+utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathFoundationURL)) + "/braintree/submit"+`";
    	//				form.id = "braintreeForm";
	//
	//				var elementNonce = document.createElement("input");
	//				elementNonce.name = "nonce";
	//				elementNonce.value = response.creditCards[0].nonce;
	//				form.appendChild(elementNonce);
	//
	//				var elementCardType = document.createElement("input");
	//				elementCardType.name = "cardType";
	//				elementCardType.value = response.creditCards[0].details.cardType;
	//				form.appendChild(elementCardType);
	//
	//				var elementLastTwo = document.createElement("input");
	//				elementLastTwo.name = "lastTwo";
	//				elementLastTwo.value = response.creditCards[0].details.lastTwo;
	//				form.appendChild(elementLastTwo);
	//				`
	//
	//for key, value := range formValues {
	//	var elementName = `element`+key;
	//	htmlText += `
	//	var `+elementName+` = document.createElement('input');
	//	`+elementName+`.name = '`+key+`';
	//	`+elementName+`.value = '`+value+`';
	//	form.appendChild(`+elementName+`);
	//	`
	//}
	//
	//htmlText += `
	//
	//				document.body.appendChild(form);
	//
	//				$('#braintreeForm').submit();
    	//				//form.submit();
	//htmlText += `		});
	//		}
	//	);
	//
	//	console.log('Client');
	//});
	//</script>
	//`

	env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "FORM: "+htmlText+"\n")

	return api.StructRestRedirect{Result: htmlText, Location: utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathFoundationURL)) + "/checkout/submit"}, nil
}

// Capture makes payment method capture operation
// - at time of implementation this method is not used anywhere
func (it *BraintreePaypalPaymentMethod) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e40b3354-5522-4027-b56d-ea9e736f637f", "Not implemented")
}

// Refund will return funds on the given order
// - at time of implementation this method is not used anywhere
func (it *BraintreePaypalPaymentMethod) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f28dab69-77a1-438c-a14f-70c700a70c3b", "Not implemented")
}

// Void will mark the order and capture as void
// - at time of implementation this method is not used anywhere
func (it *BraintreePaypalPaymentMethod) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "194a4323-4cc4-41a0-ae90-0b578dc4b73a", "Not implemented")
}

