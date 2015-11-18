package friendmail

import (
	"bytes"
	"encoding/base64"
	"time"

	"github.com/dchest/captcha"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("friend/email", api.ConstRESTOperationCreate, APIFriendEmail)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("friend/captcha", api.ConstRESTOperationGet, APIFriendCaptcha)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIFriendCaptcha will generate a captcha for use in a form
func APIFriendCaptcha(context api.InterfaceApplicationContext) (interface{}, error) {

	var captchaDigits []byte
	var captchaValue string
	var result interface{}

	captchaValuesMutex.Lock()
	if len(captchaValues) < ConstMaxCaptchaItems {
		captchaDigits = captcha.RandomDigits(captcha.DefaultLen)
		for i := range captchaDigits {
			captchaValue += string(captchaDigits[i] + '0')
		}
	} else {
		for key := range captchaValues {
			captchaValue = key
			break
		}
		captchaDigits = make([]byte, len(captchaValue))
		for idx, digit := range []byte(captchaValue) {
			captchaDigits[idx] = digit - '0'
		}
	}
	captchaValues[captchaValue] = time.Now()
	captchaValuesMutex.Unlock()

	// generate a captcha image
	image := captcha.NewImage(captchaValue, captchaDigits, captcha.StdWidth, captcha.StdHeight)
	if image == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "7224c8fe-6079-4bb3-9dc3-ad0847db8e29", "Unable to generate a captcha image.")
	}

	if context.GetRequestArguments()["json"] != "" {
		buffer := new(bytes.Buffer)
		buffer.WriteString("data:image/png;base64,")
		image.WriteTo(base64.NewEncoder(base64.StdEncoding, buffer))

		result = map[string]interface{}{
			"captcha": buffer.String(),
		}
	} else {
		context.SetResponseContentType("image/png")
		image.WriteTo(context.GetResponseWriter())
	}

	return result, nil
}

// APIFriendEmail sends an email to a friend
func APIFriendEmail(context api.InterfaceApplicationContext) (interface{}, error) {

	var friendEmail string

	// checking request values
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(requestData, "captcha", "friend_email") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3c3d0918-b951-43b7-943c-54e8571d0c32", "'captcha' and/or 'friend_email' fields are not specified or blank")
	}

	friendEmail = utils.InterfaceToString(requestData["friend_email"])
	if !utils.ValidEmailAddress(friendEmail) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5821734e-4c84-449b-9f75-fd1154623c42", "The email address, "+friendEmail+", is not in valid format.")
	}

	// checking captcha
	captchaValue := utils.InterfaceToString(requestData["captcha"])

	captchaValuesMutex.Lock()
	if _, present := captchaValues[captchaValue]; !present {
		captchaValuesMutex.Unlock()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "8bd3ad79-e464-4355-8a13-27ff55980fbb", "The entered character sequence does not match the captcha value.")
	}
	captchaStore.Get(captchaValue, true)
	delete(captchaValues, captchaValue)
	captchaValuesMutex.Unlock()

	// sending an e-mail
	emailSubject := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathFriendMailEmailSubject))

	emailTemplate := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathFriendMailEmailTemplate))
	emailTemplate, err = utils.TextTemplate(emailTemplate, requestData)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = app.SendMail(friendEmail, emailSubject, emailTemplate)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// storing data to database
	saveData := map[string]interface{}{
		"date":  time.Now(),
		"email": friendEmail,
		"data":  requestData,
	}

	dbCollection, err := db.GetCollection(ConstCollectionNameFriendMail)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	_, err = dbCollection.Save(saveData)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}
