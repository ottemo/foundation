package visitor

import (
	"crypto/md5"
	"encoding/hex"

	"encoding/base64"
	"time"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"

	"github.com/ottemo/foundation/app/utils"
	"github.com/ottemo/foundation/app/utils/sendmail"
	"github.com/ottemo/foundation/db"

	"errors"
)

//---------------------------------
// IMPLEMENTATION SPECIFIC METHODS
//---------------------------------

// returns I_VisitorAddress model filled with values from DB or blank structure if no id found in DB
func (it *DefaultVisitor) getVisitorAddressById(addressId string) visitor.I_VisitorAddress {
	address_model, err := models.GetModel("VisitorAddress")
	if err != nil {
		return nil
	}

	if address_model, ok := address_model.(visitor.I_VisitorAddress); ok {
		if addressId != "" {
			address_model.Load(addressId)
		}

		return address_model
	}

	return nil
}

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

//--------------------------
// INTERFACE IMPLEMENTATION
//--------------------------

func (it *DefaultVisitor) GetEmail() string      { return it.Email }
func (it *DefaultVisitor) GetFacebookId() string { return it.FacebookId }
func (it *DefaultVisitor) GetGoogleId() string   { return it.GoogleId }

func (it *DefaultVisitor) GetFullName() string  { return it.FirstName + " " + it.LastName }
func (it *DefaultVisitor) GetFirstName() string { return it.FirstName }
func (it *DefaultVisitor) GetLastName() string  { return it.LastName }

func (it *DefaultVisitor) GetShippingAddress() visitor.I_VisitorAddress {
	return it.ShippingAddress
}

func (it *DefaultVisitor) SetShippingAddress(address visitor.I_VisitorAddress) error {
	it.ShippingAddress = address
	return nil
}

func (it *DefaultVisitor) GetBillingAddress() visitor.I_VisitorAddress {
	return it.BillingAddress
}

func (it *DefaultVisitor) SetBillingAddress(address visitor.I_VisitorAddress) error {
	it.BillingAddress = address
	return nil
}

// returns true if visitor e-mail was validated
func (it *DefaultVisitor) IsValidated() bool {
	return it.ValidateKey == ""
}

// marks visitor e-mail as not validated
//	- sends to visitor e-mail new validation key
func (it *DefaultVisitor) Invalidate() error {

	if it.GetEmail() == "" {
		return errors.New("email is not specified")
	}

	data, err := time.Now().MarshalBinary()
	if err != nil {
		return err
	}

	it.ValidateKey = hex.EncodeToString([]byte(base64.StdEncoding.EncodeToString(data)))
	err = it.Save()
	if err != nil {
		return err
	}

	// TODO: probably not a best solution to have it there
	linkHref := utils.GetSiteBackUrl() + "visitor/validate/" + it.ValidateKey
	//link := "<a href=\"" + linkHref + "\"/>" + linkHref + "</a>"

	sendmail.SendMail(it.GetEmail(), "e-mail validation", "please follow the link to validate your e-mail: "+linkHref)

	return nil
}

// validates visitors e-mails for given key
//   - if key was expired, user will receive new one validation code
func (it *DefaultVisitor) Validate(key string) error {

	// looking for visitors with given validation key in DB and collecting ids
	visitorIds := make([]string, 0)
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {
			collection.AddFilter("validate", "=", key)

			records, err := collection.Load()
			if err != nil {
				return err
			}

			for _, record := range records {
				if visitorId, present := record["_id"]; present {
					if visitorId, ok := visitorId.(string); ok {
						visitorIds = append(visitorIds, visitorId)
					}
				}

			}
		}
	}

	// checking validation key expiration
	step1, err := hex.DecodeString(key)
	data, err := base64.StdEncoding.DecodeString(string(step1))
	if err != nil {
		return err
	}

	stamp := time.Now()
	timeNow := stamp.Unix()
	stamp.UnmarshalBinary(data)
	timeWas := stamp.Unix()

	validationExpired := (timeNow - timeWas) > EMAIL_VALIDATE_EXPIRE

	// processing visitors for given validation key
	for _, visitorId := range visitorIds {
		model, _ := it.New()
		visitorModel := model.(*DefaultVisitor)

		err := visitorModel.Load(visitorId)
		if err != nil {
			return err
		}

		if !validationExpired {
			visitorModel.ValidateKey = ""
			visitorModel.Save()
		} else {
			visitorModel.Invalidate()
			return errors.New("validation period expired, new validation URL was sent")
		}
	}

	return nil
}

func (it *DefaultVisitor) SetPassword(passwd string) error {
	if len(passwd) > 0 {
		if utils.IsMD5(passwd) {
			it.Password = passwd
		} else {
			it.Password = it.passwdEncode(passwd)
		}
	} else {
		return errors.New("password can't be blank")
	}

	return nil
}

func (it *DefaultVisitor) CheckPassword(passwd string) bool {
	return it.passwdEncode(passwd) == it.Password
}

// loads visitor information from DB based on google account id
func (it *DefaultVisitor) LoadByGoogleId(googleId string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {

			collection.AddFilter("google_id", "=", googleId)
			rows, err := collection.Load()
			if err != nil {
				return err
			}

			if len(rows) == 0 {
				return errors.New("visitor was not found")
			}

			if len(rows) > 1 {
				return errors.New("duplicated google account id")
			} else {
				err = it.FromHashMap(rows[0])
				if err != nil {
					return err
				}

				return nil
			}
		}
	}

	return errors.New("Something went wrong")
}

// loads visitor information from DB based on facebook account id
func (it *DefaultVisitor) LoadByFacebookId(facebookId string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {

			collection.AddFilter("facebook_id", "=", facebookId)
			rows, err := collection.Load()
			if err != nil {
				return err
			}

			if len(rows) == 0 {
				return errors.New("visitor was not found")
			}

			if len(rows) > 1 {
				return errors.New("duplicated facebook account id")
			} else {
				err = it.FromHashMap(rows[0])
				if err != nil {
					return err
				}

				return nil
			}
		}
	}

	return errors.New("Something went wrong")
}

// loads visitor information from DB based on email which must be unique
func (it *DefaultVisitor) LoadByEmail(email string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {

			collection.AddFilter("email", "=", email)
			rows, err := collection.Load()
			if err != nil {
				return err
			}

			if len(rows) == 0 {
				return errors.New("visitor was not found")
			}

			if len(rows) > 1 {
				return errors.New("duplicated email")
			} else {
				err = it.FromHashMap(rows[0])
				if err != nil {
					return err
				}

				return nil
			}
		}
	}

	return errors.New("Something went wrong")
}
