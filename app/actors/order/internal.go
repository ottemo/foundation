package order

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// SendShippingStatusUpdateEmail will send an email to alert customers their order has been packed and shipped
func (it DefaultOrder) SendShippingStatusUpdateEmail() error {
	subject := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathShippingEmailSubject))
	emailTemplate := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathShippingEmailTemplate))
	timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))

	// Assemble template variables
	orderMap := it.ToHashMap()

	// convert date of order creation to store time zone
	if date, present := orderMap["created_at"]; present {
		convertedDate, _ := utils.MakeTZTime(utils.InterfaceToTime(date), timeZone)
		if !utils.IsZeroTime(convertedDate) {
			orderMap["created_at"] = convertedDate
		}
	}

	templateVariables := map[string]interface{}{
		"Site":  map[string]string{"Url": app.GetStorefrontURL("")},
		"Order": orderMap,
	}

	body, err := utils.TextTemplate(emailTemplate, templateVariables)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	to := utils.InterfaceToString(it.Get("customer_email"))
	if to == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "370e99c1-727c-4ccf-a004-078d4ab343c7", "Couldn't figure out who to send a shipping status update email to. order_id: "+it.GetID())
	}

	err = app.SendMail(to, subject, body)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
