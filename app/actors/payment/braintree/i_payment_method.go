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

	bt := braintree.New(
		braintree.Sandbox,
		"ddxtcwf5n3hvtz3g",
		"cfj6fzzrkc898mm6",
		"24d8738ee7bc4331bbc3bac79f2a54c2",
	)

	token, err := bt.ClientToken().Generate();
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	////gateway := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathFoundationURL)) + "checkout/submit-fake"
	//formValues := map[string]string{
	//	"CC": "4111 1111 1111 1111",
	//}

	htmlText := `
	<script>
	$.getScript('https://js.braintreegateway.com/web/3.6.0/js/client.min.js', function(response, status){
		console.log('Status: ', status);

		window.braintree.client.create({
			authorization: '`+token+`'
			},function (createErr, clientInstance) {
				var data = {
					creditCard: {
						number: '4111 1111 1111 1111',
						cvv: '111',
						expirationDate: '11/25'
					}
      				};

				clientInstance.request({
    					endpoint: 'payment_methods/credit_cards',
    					method: 'post',
    					data: data
  				}, function (requestErr, response) {
					if (requestErr) { throw new Error(requestErr); }

					console.log('Got nonce:', response.creditCards[0].nonce);
					console.log('Response:', response);

					var form = document.createElement('form');

					form.method = "POST";
    					form.action = "`+utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathFoundationURL)) + "/braintree/submit"+`";
    					form.id = "braintreeForm";

					var element1 = document.createElement("input");
					element1.name = "nonce";
					element1.value = response.creditCards[0].nonce;
					form.appendChild(element1);

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
