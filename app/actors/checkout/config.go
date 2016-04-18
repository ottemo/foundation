package checkout

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "701e85e4-b63c-48f4-a990-673ba0ed6a2a", "can't obtain config")
	}

	// Checkout
	//---------
	err := config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Checkout",
		Description: "checkout related options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	config.RegisterItem(env.StructConfigItem{
		Path: checkout.ConstConfigPathConfirmationEmail,
		Value: `Dear {{.Visitor.last_name}} {{.Visitor.first_name}}
<br />
<br />
Thank for your order.
<br />
<h3>Order #{{.Order.increment_id}}: </h3><br />
Order summary<br />
Subtotal: ${{.Order.subtotal}}<br />
Tax: ${{.Order.tax_amount}}<br />
Shipping: ${{.Order.shipping_amount}}<br />
Total: ${{.Order.grand_total}}<br />`,
		Type:        env.ConstConfigTypeHTML,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Order confirmation e-mail: ",
		Description: "contents of email will be sent to customer on success checkout",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Payment
	//--------
	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathPaymentGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Payment",
		Description: "payment methods related group",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathPaymentOriginGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Payment Origin",
		Description: "payments methods origin information",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathPaymentOriginCountry,
		Value:       "US",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     models.ConstCountriesList,
		Label:       "Country",
		Description: "payment methods origin country",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathPaymentOriginState,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     models.ConstStatesList,
		Label:       "State",
		Description: "payment methods origin state",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathPaymentOriginCity,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "City",
		Description: "payment methods origin city",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathPaymentOriginAddressline1,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 1",
		Description: "payment methods origin address line 1",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathPaymentOriginAddressline2,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 2",
		Description: "payment methods origin address line 2",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathPaymentOriginZip,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "zip",
		Description: "payment methods origin zip code",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Shipping
	//---------
	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathShippingGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Shipping",
		Description: "shipping methods related group",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathShippingOriginGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Shipping Origin",
		Description: "shipping methods origin information",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathShippingOriginCountry,
		Value:       "US",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     models.ConstCountriesList,
		Label:       "Country",
		Description: "shipping methods origin country",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathShippingOriginState,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     models.ConstStatesList,
		Label:       "State",
		Description: "shipping methods origin state",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathShippingOriginCity,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "City",
		Description: "shipping methods origin city",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathShippingOriginAddressline1,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 1",
		Description: "shipping methods origin address line 1",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathShippingOriginAddressline2,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 2",
		Description: "shipping methods origin address line 2",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathShippingOriginZip,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "zip",
		Description: "shipping methods origin zip code",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	config.RegisterItem(env.StructConfigItem{
		Path:        checkout.ConstConfigPathOversell,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Limit Oversell",
		Description: "Do not allow product to oversell, (i.e. do not sell product when qty < 0)",
		Image:       "",
	}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
