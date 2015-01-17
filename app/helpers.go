package app

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

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
	return strings.TrimRight(baseURL, "#/") + "/#/" + path
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
