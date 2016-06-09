package app

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetVersion returns current version number
func GetVersion() string {
	return fmt.Sprintf("v%d.%d.%d", ConstVersionMajor, ConstVersionMinor, ConstSprintNumber)
}

// GetVerboseVersion returns verbose information about application build
func GetVerboseVersion() string {
	result := GetVersion()
	if buildBranch != "" {
		result += "-" + buildBranch
	}
	if buildNumber != "" {
		result += "-b" + buildNumber
	}
	if buildTags != "" {
		result += " (" + buildTags + ")"
	}
	if buildDate != "" {
		result += " [" + buildDate + "]"
	}
	return result
}

// GetDashboardURL returns url related to dashboard server
func GetDashboardURL(path string) string {
	baseURL := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathDashboardURL))
	return strings.TrimRight(baseURL, "#/") + "/#/" + path
}

// GetStorefrontURL returns url related to storefront server
func GetStorefrontURL(path string) string {
	baseURL := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathStorefrontURL))
	return strings.TrimRight(baseURL, "/") + "/" + path
}

// GetFoundationURL returns url related to foundation server
func GetFoundationURL(path string) string {
	baseURL := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathFoundationURL))
	return strings.TrimRight(baseURL, "/") + "/" + path
}

// SendMail sends mail via smtp server specified in config
func SendMail(to string, subject string, body string) error {

	userName := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailUser))
	password := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailPassword))

	mailServer := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailServer))
	mailPort := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailPort))
	if mailPort != "" && mailPort != "0" {
		mailPort = ":" + mailPort
	} else {
		return nil
	}

	context := map[string]interface{}{
		"From":      utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailFrom)),
		"To":        to,
		"Subject":   subject,
		"Body":      body,
		"Signature": utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailSignature)),
	}

	emailTemplateBody := `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
Content-Type: text/html

<p>{{.Body}}</p>

<p>{{.Signature}}</p>`

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

	var auth smtp.Auth
	if userName != "" {
		auth = smtp.PlainAuth("", userName, password, mailServer)
	}

	err = smtp.SendMail(mailServer+mailPort, auth, userName, []string{to}, doc.Bytes())
	return env.ErrorDispatch(err)
}

// SendMailEx sends mail via smtp server specified in config
// - use nil context and/or headers for default values
func SendMailEx(headers map[string]string, body string, context map[string]interface{}) error {

	userName := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailUser))
	password := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailPassword))

	mailServer := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailServer))
	mailPort := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailPort))
	if mailPort != "" && mailPort != "0" {
		mailPort = ":" + mailPort
	} else {
		return nil
	}

	if context == nil {
		context = make(map[string]interface{})
	}

	emailTo, present := headers["To"]
	if !present {
		return env.ErrorDispatch(env.ErrorNew(ConstErrorModule, env.ConstErrorLevelHelper, "d76f17be-407c-4e3d-acb7-74dc1d77f329", "send To email not specified in headers"))
	}

	if _, present := context["From"]; !present {
		context["From"] = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailFrom))
	}

	if _, present := context["Signature"]; !present {
		context["Signature"] = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathMailSignature))
	}

	if _, present := context["Subject"]; !present {
		context["Subject"] = "Ottemo email"
	}

	emailTemplateBody := `From: {{.From}}
To: ` + emailTo + `
Subject: {{.Subject}}`

	for key, value := range headers {
		if _, present := context[key]; !present {
			context[key] = value
		}

		if !utils.IsAmongStr(key, "From", "To", "Subject") {
			newHeaderLine := `
` + key + ": " + value

			emailTemplateBody = emailTemplateBody + newHeaderLine
		}
	}

	emailTemplateBody = emailTemplateBody + `
Content-Type: text/html

<p>` + body + `</p>

<p>{{.Signature}}</p>`

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

	var auth smtp.Auth
	if userName != "" {
		auth = smtp.PlainAuth("", userName, password, mailServer)
	}

	err = smtp.SendMail(mailServer+mailPort, auth, userName, []string{emailTo}, doc.Bytes())
	return env.ErrorDispatch(err)
}

// GetSessionTimeZone - return time zone of session
func GetSessionTimeZone(session api.InterfaceSession) (string, error) {

	if session != nil {
		return utils.InterfaceToString(session.Get(api.ConstSessionKeyTimeZone)), nil
	}

	return "", env.ErrorNew(ConstErrorModule, env.ConstErrorLevelHelper, "d29c25ff-ffd0-4fd0-a02b-6f22a1bf9969", "session is not specified")
}

// SetSessionTimeZone - validate time zone and set it to session
func SetSessionTimeZone(session api.InterfaceSession, zone string) (interface{}, error) {

	if session != nil && zone != "" {
		zoneName, zoneOffset := utils.ParseTimeZone(utils.InterfaceToString(zone))

		if zoneUTCOffset, present := utils.TimeZones[zoneName]; present {
			session.Set(api.ConstSessionKeyTimeZone, zone)
			zoneOffset += zoneUTCOffset

			return zoneName + zoneOffset.String(), nil
		}

		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelHelper, "962f865e-be19-44c6-8139-aca808f047d3", " specified time zone can't be parsed")
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelHelper, "b981731d-907b-4271-8087-095cff6edf3b", "session or zone is not specified")
}
