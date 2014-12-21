package authorizenet

import (
	"errors"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/api/session"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("authorizenet", "POST", "receipt", restReceipt)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("authorizenet", "POST", "relay", restRelay)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function to process Authorize.Net receipt result
func restReceipt(params *api.StructAPIHandlerParams) (interface{}, error) {

	postData := params.RequestContent.(map[string]interface{})

	status := postData["x_response_code"]

	session, err := session.GetSessionByID(utils.InterfaceToString(postData["x_session"]))
	if err != nil {
		return nil, errors.New("Wrong session ID")
	}
	params.Session = session

	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	checkoutOrder := currentCheckout.GetOrder()

	switch status {
	case ConstTransactionApproved:
		{
			currentCart := currentCheckout.GetCart()
			if currentCart == nil {
				return nil, errors.New("Cart is not specified")
			}
			if checkoutOrder != nil {
				checkoutOrder.NewIncrementID()

				checkoutOrder.Set("status", "pending")
				checkoutOrder.Set("payment_info", postData)

				err = currentCheckout.CheckoutSuccess(checkoutOrder, params.Session)
				if err != nil {
					return nil, err
				}

				// Send confirmation email
				err = currentCheckout.SendOrderConfirmationMail()
				if err != nil {
					return nil, err
				}

				env.Log("authorizenet.log", env.ConstLogPrefixInfo, "TRANSACTION APPROVED: "+
					"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
					"OrderID - "+checkoutOrder.GetID()+", "+
					"Card  - "+utils.InterfaceToString(postData["x_card_type"])+" "+utils.InterfaceToString(postData["x_account_number"])+", "+
					"Total - "+utils.InterfaceToString(postData["x_amount"])+", "+
					"Transaction ID - "+utils.InterfaceToString(postData["x_trans_id"]))

				return api.StructRestRedirect{Location: app.GetStorefrontURL("account/order/" + checkoutOrder.GetID()), DoRedirect: true}, nil
			}
		}
	case ConstTransactionDeclined:
	case ConstTransactionWaitingReview:
	default:
		{
			if checkoutOrder != nil {
				env.Log("authorizenet.log", env.ConstLogPrefixError, "TRANSACTION NOT APPROVED: "+
					"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
					"OrderID - "+checkoutOrder.GetID()+", "+
					"Card  - "+utils.InterfaceToString(postData["x_card_type"])+" "+utils.InterfaceToString(postData["x_account_number"])+", "+
					"Total - "+utils.InterfaceToString(postData["x_amount"])+", "+
					"Transaction ID - "+utils.InterfaceToString(postData["x_trans_id"]))
			}

			return []byte(`<html>
					 <head>
						 <noscript>
						 	<meta http-equiv='refresh' content='1;url=` + app.GetStorefrontURL("checkout") + `'>
						 </noscript>
					 </head>
					 <body>
					 	<h1>Something went wrong</h1>
					 	<p>` + utils.InterfaceToString(postData["x_response_reason_text"]) + `</p>

						<p><a href="` + app.GetStorefrontURL("checkout") + `">Back to store</a></p>

					 </body>
				</html>`), nil
		}
	}
	if checkoutOrder != nil {
		env.Log("authorizenet.log", env.ConstLogPrefixError, "TRANSACTION NOT APPROVED: (can't process authorize.net response) "+
			"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
			"OrderID - "+checkoutOrder.GetID()+", "+
			"Card  - "+utils.InterfaceToString(postData["x_card_type"])+" "+utils.InterfaceToString(postData["x_account_number"])+", "+
			"Total - "+utils.InterfaceToString(postData["x_amount"])+", "+
			"Transaction ID - "+utils.InterfaceToString(postData["x_trans_id"]))
	}
	return nil, errors.New("can't process authorize.net response")
}

// WEB REST API function to process Authorize.Net relay result
func restRelay(params *api.StructAPIHandlerParams) (interface{}, error) {

	postData := params.RequestContent.(map[string]interface{})

	status := postData["x_response_code"]

	session, err := session.GetSessionByID(utils.InterfaceToString(postData["x_session"]))
	if err != nil {
		return nil, errors.New("Wrong session ID")
	}
	params.Session = session

	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	checkoutOrder := currentCheckout.GetOrder()

	switch status {
	case ConstTransactionApproved:
		{
			currentCart := currentCheckout.GetCart()
			if currentCart == nil {
				return nil, errors.New("Cart is not specified")
			}
			if checkoutOrder != nil {
				checkoutOrder.NewIncrementID()

				checkoutOrder.Set("status", "pending")
				checkoutOrder.Set("payment_info", postData)

				err = currentCheckout.CheckoutSuccess(checkoutOrder, params.Session)
				if err != nil {
					return nil, err
				}

				// Send confirmation email
				err = currentCheckout.SendOrderConfirmationMail()
				if err != nil {
					return nil, err
				}

				params.ResponseWriter.Header().Set("Content-Type", "text/plain")

				env.Log("authorizenet.log", env.ConstLogPrefixInfo, "TRANSACTION APPROVED: "+
					"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
					"OrderID - "+checkoutOrder.GetID()+", "+
					"Card  - "+utils.InterfaceToString(postData["x_card_type"])+" "+utils.InterfaceToString(postData["x_account_number"])+", "+
					"Total - "+utils.InterfaceToString(postData["x_amount"])+", "+
					"Transaction ID - "+utils.InterfaceToString(postData["x_trans_id"]))

				return []byte(`<html>
					 <head>
						 <noscript>
						 	<meta http-equiv='refresh' content='1;url=` + app.GetStorefrontURL("account/order/"+checkoutOrder.GetID()) + `'>
						 </noscript>
					 </head>
					 <body>
					 	<h1>Thanks for your purchase.</h1>
					 	<p>Your transaction ID: <b>` + utils.InterfaceToString(postData["x_trans_id"]) + `</b></p>
					 	<p>You will  redirect to the store after <span id="sec"></span> sec.	<a href="` + app.GetStorefrontURL("account/order/"+checkoutOrder.GetID()) + `">Back to store</a></p>
					 </body>
					 <script type='text/javascript' charset='utf-8'>
					 	(function(){
							var seconds = 10;
							document.getElementById("sec").innerHTML = seconds;
							setInterval(function(){
								seconds -= 1;
								document.getElementById("sec").innerHTML = seconds;
								if(0 === seconds){
									window.location='` + app.GetStorefrontURL("account/order/"+checkoutOrder.GetID()) + `';
								}
							}, 1000);
					 	})();
					 </script>
				</html>`), nil
			}
		}
	case ConstTransactionDeclined:
	case ConstTransactionWaitingReview:
	default:
		{
			if checkoutOrder != nil {
				env.Log("authorizenet.log", env.ConstLogPrefixError, "TRANSACTION NOT APPROVED: "+
					"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
					"OrderID - "+checkoutOrder.GetID()+", "+
					"Card  - "+utils.InterfaceToString(postData["x_card_type"])+" "+utils.InterfaceToString(postData["x_account_number"])+", "+
					"Total - "+utils.InterfaceToString(postData["x_amount"])+", "+
					"Transaction ID - "+utils.InterfaceToString(postData["x_trans_id"]))
			}
			return []byte(`<html>
					 <head>
						 <noscript>
						 	<meta http-equiv='refresh' content='1;url=` + app.GetStorefrontURL("checkout") + `'>
						 </noscript>
					 </head>
					 <body>
					 	<h1>Something went wrong</h1>
					 	<p>` + utils.InterfaceToString(postData["x_response_reason_text"]) + `</p>

						<p><a href="` + app.GetStorefrontURL("checkout") + `">Back to store</a></p>

					 </body>
				</html>`), nil
		}
	}
	if checkoutOrder != nil {
		env.Log("authorizenet.log", env.ConstLogPrefixError, "TRANSACTION NOT APPROVED: (can't process authorize.net response) "+
			"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
			"OrderID - "+checkoutOrder.GetID()+", "+
			"Card  - "+utils.InterfaceToString(postData["x_card_type"])+" "+utils.InterfaceToString(postData["x_account_number"])+", "+
			"Total - "+utils.InterfaceToString(postData["x_amount"])+", "+
			"Transaction ID - "+utils.InterfaceToString(postData["x_trans_id"]))
	}

	return nil, errors.New("can't process authorize.net response")
}
