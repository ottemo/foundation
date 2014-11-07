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

// returns I_VisitorAddress model filled with values from DB or blank structure if no id found in DB
func (it *DefaultVisitor) passwdEncode(passwd string) string {
	salt := ":"
	if len(passwd) > 2 {
		salt += passwd[0:1]
	}

	hasher := md5.New()
	hasher.Write([]byte(passwd + salt))

	return hex.EncodeToString(hasher.Sum(nil))
}

// GetEmail returns visitor e-mail, which also used as login
func (it *DefaultVisitor) GetEmail() string {
	return it.Email
}

// GetFacebookID returns visitor facebook id
func (it *DefaultVisitor) GetFacebookID() string {
	return it.FacebookId
}

// returns visitor google id
func (it *DefaultVisitor) GetGoogleId() string {
	return it.GoogleId
}

// returns visitor full name
func (it *DefaultVisitor) GetFullName() string {
	return it.FirstName + " " + it.LastName
}

// returns visitor first name
func (it *DefaultVisitor) GetFirstName() string {
	return it.FirstName
}

// returns visitor last name
func (it *DefaultVisitor) GetLastName() string {
	return it.LastName
}

// returns visitor birthday
func (it *DefaultVisitor) GetBirthday() time.Time {
	return it.Birthday
}

// returns visitor created at date
func (it *DefaultVisitor) GetCreatedAt() time.Time {
	return it.CreatedAt
}

// returns shipping address for visitor
func (it *DefaultVisitor) GetShippingAddress() visitor.I_VisitorAddress {
	return it.ShippingAddress
}

// updates shipping address for visitor
func (it *DefaultVisitor) SetShippingAddress(address visitor.I_VisitorAddress) error {
	it.ShippingAddress = address
	return nil
}

// returns billing address for visitor
func (it *DefaultVisitor) GetBillingAddress() visitor.I_VisitorAddress {
	return it.BillingAddress
}

// updates billing address for visitor
func (it *DefaultVisitor) SetBillingAddress(address visitor.I_VisitorAddress) error {
	it.BillingAddress = address
	return nil
}

// returns true if visitor is admin
func (it *DefaultVisitor) IsAdmin() bool {
	return it.Admin
}

// returns true if visitor e-mail was validated
func (it *DefaultVisitor) IsValidated() bool {
	return it.ValidateKey == ""
}

// marks visitor e-mail as not validated, sends to visitor e-mail with new validation key
func (it *DefaultVisitor) Invalidate() error {

	if it.GetEmail() == "" {
		return env.ErrorNew("email is not specified")
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

	linkHref := app.GetFoundationUrl("visitor/validate/" + it.ValidateKey)

	err = app.SendMail(it.GetEmail(), "e-mail validation", "please follow the link to validate your e-mail: "+linkHref)

	return env.ErrorDispatch(err)
}

// validates visitors e-mails for given key
//   - if key was expired, user will receive new one validation code
func (it *DefaultVisitor) Validate(key string) error {

	// looking for visitors with given validation key in DB and collecting ids
	visitorIds := make([]string, 0)

	collection, err := db.GetCollection(COLLECTION_NAME_VISITOR)
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

	for _, record := range records {
		if visitorId, present := record["_id"]; present {
			if visitorId, ok := visitorId.(string); ok {
				visitorIds = append(visitorIds, visitorId)
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

	validationExpired := (timeNow - timeWas) > EMAIL_VALIDATE_EXPIRE

	// processing visitors for given validation key
	for _, visitorId := range visitorIds {

		visitorModel, err := visitor.LoadVisitorById(visitorId)
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

			return env.ErrorNew("validation period expired, new validation URL was sent")
		}
	}

	return nil
}

// updates password for visitor
func (it *DefaultVisitor) SetPassword(passwd string) error {
	if len(passwd) > 0 {
		if utils.IsMD5(passwd) {
			it.Password = passwd
		} else {
			it.Password = it.passwdEncode(passwd)
		}
	} else {
		return env.ErrorNew("password can't be blank")
	}

	return nil
}

// validates password for visitor
func (it *DefaultVisitor) CheckPassword(passwd string) bool {
	return it.passwdEncode(passwd) == it.Password
}

// generates new password for user
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

	linkHref := app.GetStorefrontUrl("login")
	err = app.SendMail(it.GetEmail(), "forgot password event", "Forgot password was requested for your account "+
		it.GetEmail()+"\n\n"+
		"New password: "+newPassword+"\n\n"+
		"Please change your password on next login "+linkHref)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// loads visitor information from DB based on google account id
func (it *DefaultVisitor) LoadByGoogleId(googleId string) error {

	collection, err := db.GetCollection(COLLECTION_NAME_VISITOR)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddFilter("google_id", "=", googleId)
	rows, err := collection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return env.ErrorNew("visitor not found")
	}

	if len(rows) > 1 {
		return env.ErrorNew("duplicated google account id")
	}

	err = it.FromHashMap(rows[0])
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// loads visitor information from DB based on facebook account id
func (it *DefaultVisitor) LoadByFacebookId(facebookId string) error {

	collection, err := db.GetCollection(COLLECTION_NAME_VISITOR)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddFilter("facebook_id", "=", facebookId)
	rows, err := collection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return env.ErrorNew("visitor not found")
	}

	if len(rows) > 1 {
		return env.ErrorNew("duplicated facebook account id")
	}

	err = it.FromHashMap(rows[0])
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// loads visitor information from DB based on email which must be unique
func (it *DefaultVisitor) LoadByEmail(email string) error {

	collection, err := db.GetCollection(COLLECTION_NAME_VISITOR)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddFilter("email", "=", email)
	rows, err := collection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return env.ErrorNew("visitor not found")
	}

	if len(rows) > 1 {
		return env.ErrorNew("duplicated email")
	}

	err = it.FromHashMap(rows[0])
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
