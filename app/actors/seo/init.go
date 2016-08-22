package seo

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/seo"
)

// init makes package self-initialization routine
func init() {
	seoItemInstance := new(DefaultSEOItem)
	var _ seo.InterfaceSEOItem = seoItemInstance
	models.RegisterModel(seo.ConstModelNameSEOItem, seoItemInstance)

	api.RegisterOnRestServiceStart(setupAPI)
	db.RegisterOnDatabaseStart(setupDB)
}

// setupDB prepares system database for package usage
func setupDB() error {

	if collection, err := db.GetCollection(ConstCollectionNameURLRewrites); err == nil {
		collection.AddColumn("url", db.ConstTypeVarchar, true)
		collection.AddColumn("type", db.ConstTypeVarchar, true)
		collection.AddColumn("rewrite", db.ConstTypeVarchar, false)
		collection.AddColumn("title", db.ConstTypeVarchar, false)
		collection.AddColumn("meta_keywords", db.ConstTypeVarchar, false)
		collection.AddColumn("meta_description", db.ConstTypeVarchar, false)
	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}
