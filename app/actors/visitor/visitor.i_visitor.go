package visitor

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"

	"encoding/base64"
	"time"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/utils"
)

// returns InterfaceVisitorAddress model filled with values from DB or blank structure if no id found in DB
func (it *DefaultVisitor) passwdEncode(passwd string) string {
	salt := ":"
	if len(passwd) > 2 {
		salt += passwd[0:1]
	}

	hasher := md5.New()
	hasher.Write([]byte(passwd + salt))

	return hex.EncodeToString(hasher.Sum(nil))
}

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

// GetBirthday returns the Visitor's birthday
func (it *DefaultVisitor) GetBirthday() time.Time {
	return it.Birthday
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

// IsAdmin returns true if the visitor is an Admin (have admin rights)
func (it *DefaultVisitor) IsAdmin() bool {
	return it.Admin
}

// IsGuest returns true if instance represents guest visitor
func (it *DefaultVisitor) IsGuest() bool {
	return it.GetGoogleID() == "" && it.GetFacebookID() == "" && it.GetEmail() == ""
}

// IsValidated returns true if the Visitor's e-mail has been verified
func (it *DefaultVisitor) IsValidated() bool {
	return it.ValidateKey == ""
}

// Invalidate marks a visitor e-mail address as not validated, then sends an e-mail to the Visitor with a new validation key
func (it *DefaultVisitor) Invalidate() error {

	if it.GetEmail() == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bef673e979c142bcade0e870b3da0e2f", "email was not specified")
	}

	data, err := time.Now().MarshalBinary()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.ValidateKey = hex.EncodeToString([]byte(base64.StdEncoding.EncodeToString(data)))
	err = it.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	linkHref := app.GetStorefrontURL("login?validate=" + it.ValidateKey)

	err = app.SendMail(it.GetEmail(), "e-mail validation", "Please follow the link to validate your e-mail: <a href=\""+linkHref+"\">"+linkHref+"</a>")

	return env.ErrorDispatch(err)
}

// Validate takes a visitors validation key and checks it against the database, a new validation email is sent if the key cannot be validated
func (it *DefaultVisitor) Validate(key string) error {

	// looking for visitors with given validation key in DB and collecting ids
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
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "597c38a7fae44eab9c8e380ecc626dd2", "wrong validation key")
	}

	for _, record := range records {
		if visitorID, present := record["_id"]; present {
			if visitorID, ok := visitorID.(string); ok {
				visitorIDs = append(visitorIDs, visitorID)
			}
		}

	}

	// checking validation key expiration
	step1, err := hex.DecodeString(key)
	data, err := base64.StdEncoding.DecodeString(string(step1))
	if err != nil {
		return env.ErrorDispatch(err)
	}

	stamp := time.Now()
	timeNow := stamp.Unix()
	stamp.UnmarshalBinary(data)
	timeWas := stamp.Unix()

	validationExpired := (timeNow - timeWas) > ConstEmailValidateExpire

	// processing visitors for given validation key
	for _, visitorID := range visitorIDs {

		visitorModel, err := visitor.LoadVisitorByID(visitorID)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		if !validationExpired {
			visitorModel := visitorModel.(*DefaultVisitor)
			visitorModel.ValidateKey = ""
			visitorModel.Save()
		} else {
			err = visitorModel.Invalidate()
			if err != nil {
				return env.ErrorDispatch(err)
			}

			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1ae869fa0fa24ec0b092a2c18b963f2d", "validation key expired, new validation link was sent to visitor e-mail")
		}
	}

	return nil
}

// SetPassword updates the password for the current Visitor
func (it *DefaultVisitor) SetPassword(passwd string) error {
	if len(passwd) > 0 {
		if utils.IsMD5(passwd) {
			it.Password = passwd
		} else {
			it.Password = it.passwdEncode(passwd)
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c24bb1660ffb4abca8d5ddacd859da72", "password can't be blank")
	}

	return nil
}

// CheckPassword validates password for the current Visitor
func (it *DefaultVisitor) CheckPassword(passwd string) bool {
	return it.passwdEncode(passwd) == it.Password
}

// GenerateNewPassword generates new password for the current Visitor
func (it *DefaultVisitor) GenerateNewPassword() error {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	const n = 10

	var bytes = make([]byte, n)
	rand.Read(bytes)
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
	err = app.SendMail(it.GetEmail(), "forgot password event", "Forgot password was requested for your account "+
		it.GetEmail()+"\n\n"+
		"New password: "+newPassword+"\n\n"+
		"Please change your password on next login "+linkHref)
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

	collection.AddFilter("google_id", "=", googleID)
	rows, err := collection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4ffde5a66e8444cfacb6fb9714b82bcc", "visitor not found")
	}

	if len(rows) > 1 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "693e7c5afdcf47319e3941d6f6c849ae", "duplicated google account id")
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

	collection.AddFilter("facebook_id", "=", facebookID)
	rows, err := collection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c33d114e435a44fe80f1456c57a692b9", "visitor not found")
	}

	if len(rows) > 1 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b3b941c0fa6b47faac6010f27e3bd69c", "duplicated facebook account id")
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

	collection.AddFilter("email", "=", email)
	rows, err := collection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0a7063fe949549918a80dcfcfc6f5b92", "visitor not found")
	}

	if len(rows) > 1 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9c7abb4649d440eaa33a9c6790cdb0d8", "duplicated email")
	}

	err = it.FromHashMap(rows[0])
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
