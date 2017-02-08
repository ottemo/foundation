package visitor

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	visitorInstance := new(DefaultVisitor)
	var _ visitor.InterfaceVisitor = visitorInstance
	if err := models.RegisterModel(visitor.ConstModelNameVisitor, visitorInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e3b3cefa-060e-458c-827f-ce26b7fe4c67", err.Error())
	}

	visitorCollectionInstance := new(DefaultVisitorCollection)
	var _ visitor.InterfaceVisitorCollection = visitorCollectionInstance
	if err := models.RegisterModel(visitor.ConstModelNameVisitorCollection, visitorCollectionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "49a81828-f8ca-4ff9-88a9-b56a47a1d6f6", err.Error())
	}

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
}

// setupDB prepares system database for package usage
func setupDB() error {

	collection, err := db.GetCollection(ConstCollectionNameVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddColumn("email", db.TypeWPrecision(db.ConstTypeVarchar, 150), true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e3aa6547-2d7e-4b84-a92b-52447d629b84", err.Error())
	}
	if err := collection.AddColumn("validate", db.TypeWPrecision(db.ConstTypeVarchar, 128), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ae1c7f07-a031-4146-a616-1759d53375a9", err.Error())
	}
	if err := collection.AddColumn("password", db.TypeWPrecision(db.ConstTypeVarchar, 128), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "76cd6cf3-1704-4f1d-a69f-41a451f1ba5e", err.Error())
	}
	if err := collection.AddColumn("first_name", db.TypeWPrecision(db.ConstTypeVarchar, 50), true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e73efac7-593d-423e-897e-5ca9e332d8f1", err.Error())
	}
	if err := collection.AddColumn("last_name", db.TypeWPrecision(db.ConstTypeVarchar, 50), true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5351c000-25e5-42e8-85ef-bef303b8a89e", err.Error())
	}

	if err := collection.AddColumn("facebook_id", db.TypeWPrecision(db.ConstTypeVarchar, 100), true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1ea9a741-095d-4502-9d76-f97e8eceac0a", err.Error())
	}
	if err := collection.AddColumn("google_id", db.TypeWPrecision(db.ConstTypeVarchar, 100), true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a21d1ac1-049e-4b69-b443-6af6520575f7", err.Error())
	}

	if err := collection.AddColumn("billing_address_id", db.ConstTypeID, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1fd9272b-e4e4-42fe-be0b-bcd47aa1b612", err.Error())
	}
	if err := collection.AddColumn("shipping_address_id", db.ConstTypeID, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4e1ccd68-cad0-4f3d-b8a7-309ecb6366f7", err.Error())
	}

	if err := collection.AddColumn("token_id", db.ConstTypeID, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "251c1e38-d74d-427f-ae5e-68faf727abe5", err.Error())
	}

	if err := collection.AddColumn("is_admin", db.ConstTypeBoolean, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ef298b26-1972-478a-839b-70f4679d34ea", err.Error())
	}
	if err := collection.AddColumn("created_at", db.ConstTypeDatetime, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a0689e88-acab-4134-b8eb-713345d07ff5", err.Error())
	}

	return nil
}
