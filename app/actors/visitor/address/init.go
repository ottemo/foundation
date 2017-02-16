package address

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// init makes package self-initialization routine
func init() {
	visitorAddressInstance := new(DefaultVisitorAddress)
	var _ visitor.InterfaceVisitorAddress = visitorAddressInstance
	if err := models.RegisterModel(visitor.ConstModelNameVisitorAddress, visitorAddressInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "601618dd-9467-4f24-80ae-64d0d314c754", err.Error())
	}

	visitorAddressCollectionInstance := new(DefaultVisitorAddressCollection)
	var _ visitor.InterfaceVisitorAddressCollection = visitorAddressCollectionInstance
	if err := models.RegisterModel(visitor.ConstModelNameVisitorAddressCollection, visitorAddressCollectionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "86f56e92-90d1-4abf-8c8e-cdf9e3762274", err.Error())
	}

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddColumn("visitor_id", db.ConstTypeID, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2f9378d2-752b-4c4a-a547-6a3087e2e30b", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("first_name", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "925ead37-cdd9-4939-a674-47743fc2915c", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("last_name", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "453c3364-c103-49e1-bec4-7997b99f9eb3", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("company", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "637b7ae4-cacd-4f4a-a9c7-ca6f7e98c898", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("address_line1", db.ConstTypeVarchar, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0ebeb457-8af1-4f5a-80ae-b64c69c50a3e", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("address_line2", db.ConstTypeVarchar, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "137615db-7188-4e37-83f7-3976c852e409", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("country", db.TypeWPrecision(db.ConstTypeVarchar, 50), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d941e468-9413-427c-b7a6-3ef420359fd2", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("state", db.TypeWPrecision(db.ConstTypeVarchar, 2), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "79eaa03d-9167-4811-9586-66e88185881f", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("city", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8da041d8-aa41-46da-97c3-f16d10faa9f2", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("phone", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "26999082-fa7a-4d90-a1a9-fa2317fec353", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("zip_code", db.TypeWPrecision(db.ConstTypeVarchar, 10), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "67227598-08df-4b02-83dd-576e05a9dcf2", "unable to add column: "+err.Error())
	}

	return nil
}
