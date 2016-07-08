package stripe

// Stripe package constants
const (
	ConstPaymentCode = "stripe"
	ConstPaymentName = "Stripe"

	ConstConfigPathGroup   = "payment.stripe"
	ConstConfigPathEnabled = "payment.stripe.enabled"
	ConstConfigPathName    = "payment.stripe.name"
	ConstConfigPathAPIKey  = "payment.stripe.apiKey"

	ConstErrorModule = "payment/stripe"
)

// Payment is the struct to hold the payment information for a visitor's order
type Payment struct{}
