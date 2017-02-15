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
	if err := models.RegisterModel(seo.ConstModelNameSEOItem, seoItemInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6a3daec9-ac8c-47be-a837-046c8ece7b06", err.Error())
	}

	seoItemCollectionInstance := new(DefaultSEOCollection)
	var _ seo.InterfaceSEOCollection = seoItemCollectionInstance
	if err := models.RegisterModel(ConstCollectionNameURLRewrites, seoItemCollectionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "42552c35-a4ef-40e6-aa38-5d69d4e92578", err.Error())
	}

	api.RegisterOnRestServiceStart(setupAPI)
	db.RegisterOnDatabaseStart(setupDB)
}

// setupDB prepares system database for package usage
func setupDB() error {

	if collection, err := db.GetCollection(ConstCollectionNameURLRewrites); err == nil {
		if err := collection.AddColumn("url", db.ConstTypeVarchar, true); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d521d854-39af-4ea8-9e35-89c1c5897d04", err.Error())
		}
		if err := collection.AddColumn("type", db.ConstTypeVarchar, true); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "43dc89e3-e2bf-4323-86e2-b58c3a81c110", err.Error())
		}
		if err := collection.AddColumn("rewrite", db.ConstTypeVarchar, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "75ec0324-68a8-4175-867b-69c4f97dcfe9", err.Error())
		}
		if err := collection.AddColumn("title", db.ConstTypeVarchar, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "48ae744d-b992-4eca-9721-e01f8e1a5945", err.Error())
		}
		if err := collection.AddColumn("meta_keywords", db.ConstTypeVarchar, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "839ec5c1-7c2b-4948-9076-20b6631195f5", err.Error())
		}
		if err := collection.AddColumn("meta_description", db.ConstTypeVarchar, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e779151d-c27d-4ccd-8ad1-5d0777b0d9df", err.Error())
		}
	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}
