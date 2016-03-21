package order

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func (it DefaultOrder) SendShippingStatusUpdateEmail() error {
	subject := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathShippingEmailSubject))
	emailTemplate := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathShippingEmailTemplate))

	templateVariables := map[string]interface{}{
		"Site":  map[string]string{"Url": app.GetStorefrontURL("")},
		"Order": it.ToHashMap(),
	}

	body, err := utils.TextTemplate(emailTemplate, templateVariables)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	to := utils.InterfaceToString(it.Get("customer_email"))
	if to == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "370e99c1-727c-4ccf-a004-078d4ab343c7", "Couldn't figure out who to send an shipping status update email to. order_id: "+it.GetID())
	}

	err = app.SendMail(to, subject, body)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
