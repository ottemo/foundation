package stripe

const (
	ConstPaymentCode = "stripe"
	ConstPaymentName = "Stripe"

	ConstConfigPathGroup   = "payment.stripe"
	ConstConfigPathEnabled = "payment.stripe.enabled"
	ConstConfigPathName    = "payment.stripe.name"
	ConstConfigPathAPIKey  = "payment.stripe.apiKey"

	ConstErrorModule = "payment/stripe"
)

type Payment struct{}
