// Package authorizenet is a Authorize.Net implementation of payment method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package authorizenet

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	// ConstTransactionApproved constants for transaction state
	ConstTransactionApproved      = "1"
	ConstTransactionDeclined      = "2"
	ConstTransactionError         = "3"
	ConstTransactionWaitingReview = "4"

	ConstLogStorage = "authorizenet.log"

	ConstPaymentCodeDPM = "authorizeNetDPM"
	ConstPaymentNameDPM = "Authorize.Net (Direct Post)"

	ConstDPMActionAuthorizeOnly       = "AUTH_ONLY"
	ConstDPMActionAuthorizeAndCapture = "AUTH_CAPTURE"

	ConstConfigPathDPMGroup = "payment.authorizeNetDPM"

	ConstConfigPathDPMEnabled = "payment.authorizeNetDPM.enabled"
	ConstConfigPathDPMAction  = "payment.authorizeNetDPM.action"
	ConstConfigPathDPMTitle   = "payment.authorizeNetDPM.title"

	ConstConfigPathDPMLogin   = "payment.authorizeNetDPM.login"
	ConstConfigPathDPMKey     = "payment.authorizeNetDPM.key"
	ConstConfigPathDPMGateway = "payment.authorizeNetDPM.gateway"

	ConstConfigPathDPMReceiptURL  = "payment.authorizeNetDPM.receiptURL"
	ConstConfigPathDPMDeclineURL  = "payment.authorizeNetDPM.declineURL"
	ConstConfigPathDPMReceiptHTML = "payment.authorizeNetDPM.receiptHTML"
	ConstConfigPathDPMDeclineHTML = "payment.authorizeNetDPM.declineHTML"

	ConstConfigPathDPMTest     = "payment.authorizeNetDPM.test"
	ConstConfigPathDPMDebug    = "payment.authorizeNetDPM.debug"
	ConstConfigPathDPMCheckout = "payment.authorizeNetDPM.checkout"

	ConstErrorModule = "payment/authorizenet"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstDefaultDeclineTemplate = `<html>
		 <head>
			<noscript>
				<meta http-equiv='refresh' content='1;url={{ .backURL }}'>
			</noscript>
		 </head>

		 <body>
			<h1>Something went wrong</h1>
			<p>Response text: {{ .response.x_response_reason_text }}</p>

			<p>You will be redirected beck to the store in <span id="countdown"></span> sec.
			<p><a href="{{ .backURL }}">Back to store</a></p>
		 </body>

		 <script type='text/javascript' charset='utf-8'>
				 var seconds = 13;
				 document.getElementById("countdown").innerHTML = seconds;
				 setInterval(function(){
					 seconds -= 1;
					 document.getElementById("countdown").innerHTML = seconds;
					 if (seconds == 0) {
						 window.location='{{ .backURL }}';
					 }
				 }, 1000);
		 </script>
	</html>`

	ConstDefaultReceiptTemplate = `<html>
		 <head>
			 <noscript>
				 <meta http-equiv='refresh' content='1;url={{ .backURL }}'>
			 </noscript>
		 </head>

		 <body>
			<h1>Thanks for your purchase.</h1>
			<p>Your transaction ID: <b>{{ .response.x_trans_id }}</b></p>
			<p>Order #{{ .order.increment_id }}</p>

			<p>You will be redirected back to the store in <span id="countdown"></span> sec.
			<a href="{{ .backURL }}">Back to store</a></p>
		 </body>

		 <script type='text/javascript' charset='utf-8'>
			 var seconds = 13;
			 document.getElementById("countdown").innerHTML = seconds;
			 setInterval(function(){
				 seconds -= 1;
				 document.getElementById("countdown").innerHTML = seconds;
				 if (seconds == 0) {
					 window.location='{{ .backURL }}';
				 }
			 }, 1000);
		 </script>
	</html>`
)

// DirectPostMethod is a implementer of InterfacePaymentMethod for a Authorize.Net Direct Post method
type DirectPostMethod struct{}
