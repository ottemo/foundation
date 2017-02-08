package token

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	visitorCardInstance := new(DefaultVisitorCard)
	var _ visitor.InterfaceVisitorCard = visitorCardInstance
	if err := models.RegisterModel(visitor.ConstModelNameVisitorCard, visitorCardInstance); err != nil {
		_ = env.ErrorDispatch(err)
	}

	visitorCardCollectionInstance := new(DefaultVisitorCardCollection)
	var _ visitor.InterfaceVisitorCardCollection = visitorCardCollectionInstance
	if err := models.RegisterModel(visitor.ConstModelNameVisitorCardCollection, visitorCardCollectionInstance); err != nil {
		_ = env.ErrorDispatch(err)
	}

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameVisitorToken)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddColumn("visitor_id", db.ConstTypeID, true); err != nil { // ottemo vid
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0381c849-0269-470a-9578-370985f1d6bc", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("token_id", db.ConstTypeVarchar, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "064a7877-2b9d-45ad-8217-000cebc6a8e1", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("payment", db.TypeWPrecision(db.ConstTypeVarchar, 150), true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bbf94dfd-055b-436b-876d-47156cc13c73", "unable to add column: "+err.Error())
	}

	if err := collection.AddColumn("customer_id", db.ConstTypeVarchar, false); err != nil { // 3rd party vid
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4dda6806-60d0-447a-aee8-e7b549f91b40", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("type", db.TypeWPrecision(db.ConstTypeVarchar, 50), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "87cba2f7-20eb-4c92-ab65-12f0225860b4", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("number", db.TypeWPrecision(db.ConstTypeVarchar, 50), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "206346b2-e186-49e9-84c5-f71b7b36280c", "unable to add column: "+err.Error())
	}

	// collection.AddColumn("expiration_month", db.ConstTypeInteger, false)
	// collection.AddColumn("expiration_year", db.ConstTypeInteger, false)

	if err := collection.AddColumn("expiration_date", db.ConstTypeVarchar, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fdb1db0f-fb32-4f01-aa40-2ff1df851e52", "unable to add column: "+err.Error())
	}

	if err := collection.AddColumn("holder", db.ConstTypeVarchar, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "60b2c7cf-6c85-4098-a26f-c372b67adecf", "unable to add column: "+err.Error())
	}

	if err := collection.AddColumn("token_updated", db.ConstTypeDatetime, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "46fc52a5-4116-4294-8d17-c3d5cbc95849", "unable to add column: "+err.Error())
	}
	if err := collection.AddColumn("created_at", db.ConstTypeDatetime, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0381c849-0269-470a-9578-370985f1d6bc", "unable to add column: "+err.Error())
	}

	return nil
}
