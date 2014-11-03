package app

import (
	"bytes"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// returns url related to dashboard server
func GetDashboardUrl(path string) string {
	baseUrl := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_DASHBOARD_URL))
	return strings.TrimRight(baseUrl, "#/") + "/#/" + path
}

// returns url related to storefront server
func GetStorefrontUrl(path string) string {
	baseUrl := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_STOREFRONT_URL))
	return strings.TrimRight(baseUrl, "#/") + "/#/" + path
}

// returns url related to foundation server
func GetFoundationUrl(path string) string {
	baseUrl := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_FOUNDATION_URL))
	return strings.TrimRight(baseUrl, "/") + "/" + path
}

// sends mail via smtp server specified in config
func SendMail(to string, subject string, body string) error {

	userName := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_MAIL_USER))
	password := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_MAIL_PASSWORD))

	mailServer := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_MAIL_SERVER))
	mailPort := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_MAIL_PORT))
	if mailPort != "" {
		mailPort = ":" + mailPort
	} else {
		return nil
	}

	context := map[string]interface{}{
		"From":      utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_MAIL_FROM)),
		"To":        to,
		"Subject":   subject,
		"Body":      body,
		"Signature": utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_MAIL_SIGNATURE)),
	}

	emailTemplateBody := `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
Content-Type: text/html

{{.Body}}

{{.Signature}}`

	emailTemplate := template.New("emailTemplate")
	emailTemplate, err := emailTemplate.Parse(emailTemplateBody)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var doc bytes.Buffer
	err = emailTemplate.Execute(&doc, context)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var auth smtp.Auth = nil
	if userName != "" {
		auth = smtp.PlainAuth("", userName, password, mailServer)
	}

	err = smtp.SendMail(mailServer+mailPort, auth, userName, []string{to}, doc.Bytes())
	return env.ErrorDispatch(err)
}
