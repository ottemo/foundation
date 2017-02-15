package visitor

import (
	"crypto/rand"
	"strings"
	"time"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetEmail returns the Visitor e-mail which also used as a login ID
func (it *DefaultVisitor) GetEmail() string {
	return it.Email
}

// GetFacebookID returns the Visitor's Facebook ID
func (it *DefaultVisitor) GetFacebookID() string {
	return it.FacebookID
}

// GetGoogleID returns the Visitor's Google ID
func (it *DefaultVisitor) GetGoogleID() string {
	return it.GoogleID
}

// GetFullName returns visitor full name
func (it *DefaultVisitor) GetFullName() string {
	return it.FirstName + " " + it.LastName
}

// GetFirstName returns the Visitor's first name
func (it *DefaultVisitor) GetFirstName() string {
	return it.FirstName
}

// GetLastName returns the Visitor's last name
func (it *DefaultVisitor) GetLastName() string {
	return it.LastName
}

// GetCreatedAt returns the Visitor creation date
func (it *DefaultVisitor) GetCreatedAt() time.Time {
	return it.CreatedAt
}

// GetShippingAddress returns the shipping address for the Visitor
func (it *DefaultVisitor) GetShippingAddress() visitor.InterfaceVisitorAddress {
	return it.ShippingAddress
}

// SetShippingAddress updates the shipping address for the Visitor
func (it *DefaultVisitor) SetShippingAddress(address visitor.InterfaceVisitorAddress) error {
	it.ShippingAddress = address
	return nil
}

// GetBillingAddress returns the billing address for the Visitor
func (it *DefaultVisitor) GetBillingAddress() visitor.InterfaceVisitorAddress {
	return it.BillingAddress
}

// SetBillingAddress updates the billing address for the Visitor
func (it *DefaultVisitor) SetBillingAddress(address visitor.InterfaceVisitorAddress) error {
	it.BillingAddress = address
	return nil
}

// GetToken returns the default card for the Visitor
func (it *DefaultVisitor) GetToken() visitor.InterfaceVisitorCard {
	return it.Token
}

// SetToken updates the default card for the Visitor
func (it *DefaultVisitor) SetToken(token visitor.InterfaceVisitorCard) error {
	it.Token = token
	return nil
}

// IsAdmin returns true if the visitor is an Admin (have admin rights)
func (it *DefaultVisitor) IsAdmin() bool {
	return it.Admin
}

// IsGuest returns true if instance represents guest visitor
func (it *DefaultVisitor) IsGuest() bool {
	return it.GetGoogleID() == "" && it.GetFacebookID() == "" && it.GetEmail() == ""
}

// IsVerified returns true if the Visitor's e-mail has been verified
func (it *DefaultVisitor) IsVerified() bool {
	return it.VerificationKey == ""
}

// Invalidate marks a visitor e-mail address as not verified, then sends an e-mail to the Visitor with a new verification key
func (it *DefaultVisitor) Invalidate() error {

	if it.GetEmail() == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9310827d-a15a-4abb-b47d-a9c1520c37ba", "The email address field cannot be blank.")
	}

	data, err := time.Now().MarshalBinary()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.VerificationKey = utils.CryptToURLString(data)
	err = it.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	linkHref := app.GetStorefrontURL("login?validate=" + it.VerificationKey)

	err = app.SendMail(it.GetEmail(), "e-mail verification", "Please follow the link to verify your e-mail address: <a href=\""+linkHref+"\">"+linkHref+"</a>")

	return env.ErrorDispatch(err)
}

// Validate takes a visitors verification key and checks it against the database, a new verification email is sent if the key cannot be validated
func (it *DefaultVisitor) Validate(key string) error {

	// looking for visitors with given verification key in DB and collecting ids
	var visitorIDs []string

	collection, err := db.GetCollection(ConstCollectionNameVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.AddFilter("validate", "=", key)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	records, err := collection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(records) == 0 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "597c38a7-fae4-4eab-9c8e-380ecc626dd2", "Unable to validate the provided Verification Key, please request a new one.")
	}

	for _, record := range records {
		if visitorID, present := record["_id"]; present {
			if visitorID, ok := visitorID.(string); ok {
				visitorIDs = append(visitorIDs, visitorID)
			}
		}

	}

	// checking verification key expiration
	data, err := utils.DecryptURLString(key)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	stamp := time.Now()
	timeNow := stamp.Unix()
	if err := stamp.UnmarshalBinary([]byte(data)); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4e4f872c-1838-4a65-bf38-fb592f2c401c", err.Error())
	}
	timeWas := stamp.Unix()

	verificationExpired := (timeNow - timeWas) > ConstEmailVerifyExpire

	// processing visitors for given verification key
	for _, visitorID := range visitorIDs {

		visitorModel, err := visitor.LoadVisitorByID(visitorID)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		if !verificationExpired {
			visitorModel := visitorModel.(*DefaultVisitor)
			visitorModel.VerificationKey = ""
			if err := visitorModel.Save(); err != nil {
				return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ef91c830-bbbd-47fd-9801-7a86539ac771", err.Error())
			}
		} else {
			err = visitorModel.Invalidate()
			if err != nil {
				return env.ErrorDispatch(err)
			}

			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1ae869fa-0fa2-4ec0-b092-a2c18b963f2d", "The provided Verification Key had expired, a new verification link has been sent your email address.")
		}
	}

	return nil
}

// SetPassword updates the password for the current Visitor
func (it *DefaultVisitor) SetPassword(passwd string) error {
	if len(passwd) > 0 {

		tmp := strings.Split(passwd, ":")
		if len(tmp) == 2 {
			if utils.IsMD5(tmp[0]) {
				it.Password = passwd
			} else {
				it.Password = utils.PasswordEncode(passwd, "")
			}
		} else if utils.IsMD5(passwd) {
			it.Password = passwd
		} else {
			it.Password = utils.PasswordEncode(passwd, "")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c24bb166-0ffb-4abc-a8d5-ddacd859da72", "The password field cannot be blank.")
	}

	return nil
}

// CheckPassword validates password for the current Visitor
func (it *DefaultVisitor) CheckPassword(passwd string) bool {
	return utils.PasswordCheck(it.Password, passwd)
}

// GenerateNewPassword generates new password for the current Visitor
func (it *DefaultVisitor) GenerateNewPassword() error {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	const n = 10

	var bytes = make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "96fb318e-98e9-491c-9f03-7e7b35a2425f", err.Error())
	}
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}

	newPassword := string(bytes)
	err := it.SetPassword(newPassword)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	linkHref := app.GetStorefrontURL("login")
	err = app.SendMail(it.GetEmail(), "Password Recovery", "A new password was requested for your account: "+it.GetEmail()+"<br><br>"+
		"New password: "+newPassword+"<br><br>"+
		"Please remember to change your password upon next login "+linkHref)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// ResetPassword generates new password for the current Visitor
func (it *DefaultVisitor) ResetPassword() error {

	if it.GetEmail() == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bef673e9-79c1-42bc-ade0-e870b3da0e2f", "The email address field cannot be blank.")
	}

	verificationCode := utils.InterfaceToString(time.Now().Unix()) + ":" + it.GetID()
	verificationCode = utils.CryptAsURLString(verificationCode)

	linkHref := app.GetStorefrontURL("reset-password") + "?key=" + verificationCode

	customerInfo := map[string]string{
		"name":               it.GetFullName(),
		"reset_password_url": linkHref,
		"reset_time_length":  "30",
	}
	siteInfo := map[string]string{
		"name": utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreName)),
		"url": app.GetStorefrontURL(""),
	}

	subject := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathLostPasswordEmailSubject))
	if subject == "" {
		subject = "Password Recovery"
	}

	emailTemplate := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathLostPasswordEmailTemplate))
	if emailTemplate == "" {
		emailTemplate := `<p>You have requested to have your password reset for your account at {{.Site.name}}</p>
		<p></p>
		<p>Please visit this url within the next {{.Customer.reset_time_length}} minutes to reset your password. </p>
		<p></p>
		<p><a href="{{.Customer.reset_password_url}}">{{.Customer.reset_password_url}}</a></p>
		<p></p>
		<p>If you received this email in error, you may safely ignore this request.</p>`

		if it.GetFullName() != "" {
			emailTemplate = `<p>Dear {{.Customer.name}},</p>
			<p></p>` + emailTemplate
		}
	}

	passwordRecoveryEmail, err := utils.TextTemplate(emailTemplate,
		map[string]interface{}{
			"Customer": customerInfo,
			"Site":     siteInfo,
		})

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = app.SendMail(it.GetEmail(), "Password Recovery", passwordRecoveryEmail)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// UpdateResetPassword takes a visitors verification key and checks it against the database, a new verification email is sent if the key cannot be validated
func (it *DefaultVisitor) UpdateResetPassword(key string, passwd string) error {

	// checking verification key expiration
	code, err := utils.DecryptURLString(key)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	verificationCode := strings.SplitN(string(code), ":", 2)

	// in this case code is invalid
	if len(verificationCode) != 2 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "875bbdf9-f746-4fbc-9950-d6cf98e681b7", "verification key mismatch")
	}

	timeWas := utils.InterfaceToTime(verificationCode[0]).Unix()
	timeNow := time.Now().Unix()
	timeDifference := (timeNow - timeWas)

	if timeDifference > ConstEmailPasswordResetExpire || timeDifference < 0 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d672f112-b74e-432f-b377-37ea2885a9fc", "Your password reset link has already expired, please submit a new request to reset your password")
	}

	err = it.Load(verificationCode[1])
	if err != nil || it.GetEmail() == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fe53edeb-85a6-4252-b705-ca4aeaabddf6", "verification key mismatch")
	}

	err = it.SetPassword(passwd)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// LoadByGoogleID loads the Visitor information from the database based on Google account ID
func (it *DefaultVisitor) LoadByGoogleID(googleID string) error {

	collection, err := db.GetCollection(ConstCollectionNameVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddFilter("google_id", "=", googleID); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0a0f84b2-518c-45ba-acd4-318716792f8f", err.Error())
	}
	rows, err := collection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4ffde5a6-6e84-44cf-acb6-fb9714b82bcc", "Unable to find an account associated with the provided Google ID.")
	}

	if len(rows) > 1 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "693e7c5a-fdcf-4731-9e39-41d6f6c849ae", "Found more than one account associated with the provided Google ID.")
	}

	err = it.FromHashMap(rows[0])
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// LoadByFacebookID loads the Visitor information from the database based on Facebook account ID
func (it *DefaultVisitor) LoadByFacebookID(facebookID string) error {

	collection, err := db.GetCollection(ConstCollectionNameVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddFilter("facebook_id", "=", facebookID); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0d4afa8b-0870-4f98-a00a-c55b1a92079e", err.Error())
	}
	rows, err := collection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c33d114e-435a-44fe-80f1-456c57a692b9", "Unable to find an account associated with provided Facebook ID.")
	}

	if len(rows) > 1 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b3b941c0-fa6b-47fa-ac60-10f27e3bd69c", "Found more than one account associated with the provided Facebook ID.")
	}

	err = it.FromHashMap(rows[0])
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// LoadByEmail loads the Visitor information from the database based on their email address, which must be unique
func (it *DefaultVisitor) LoadByEmail(email string) error {

	collection, err := db.GetCollection(ConstCollectionNameVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddFilter("email", "=", email); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3bee796f-ae90-42dd-bf4c-48f9d8731e2a", err.Error())
	}
	rows, err := collection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0a7063fe-9495-4991-8a80-dcfcfc6f5b92", "Unable to find an account associated with the provided email address, "+email+".")
	}

	if len(rows) > 1 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9c7abb46-49d4-40ea-a33a-9c6790cdb0d8", "Found more than one account associated with the provided email address.")
	}

	err = it.FromHashMap(rows[0])
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
