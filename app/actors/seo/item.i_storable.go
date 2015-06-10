package seo

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetID returns id for SEO item
func (it *DefaultSEOItem) GetID() string {
	return it.id
}

// SetID sets id for SEO item
func (it *DefaultSEOItem) SetID(newID string) error {
	it.id = newID
	return nil
}

// Load loads SEO item information from DB
func (it *DefaultSEOItem) Load(id string) error {
	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbValues, err := collection.LoadByID(id)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.SetID(utils.InterfaceToString(dbValues["_id"]))

	it.URL = utils.InterfaceToString(dbValues["url"])
	it.Type = utils.InterfaceToString(dbValues["type"])

	it.Rewrite = utils.InterfaceToString(dbValues["rewrite"])
	it.Title = utils.InterfaceToString(dbValues["title"])

	it.MetaKeywords = utils.InterfaceToString(dbValues["meta_keywords"])
	it.MetaDescription = utils.InterfaceToString(dbValues["meta_description"])

	return nil
}

// Delete removes current SEO item from DB
func (it *DefaultSEOItem) Delete() error {
	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return env.ErrorDispatch(err)
}

// Save stores current SEO item to DB
func (it *DefaultSEOItem) Save() error {
	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	storingValues := make(map[string]interface{})

	storingValues["_id"] = it.GetID()
	storingValues["url"] = it.URL
	storingValues["type"] = it.Type
	storingValues["rewrite"] = it.Rewrite
	storingValues["title"] = it.Title
	storingValues["meta_keywords"] = it.MetaKeywords
	storingValues["meta_description"] = it.MetaDescription

	newID, err := collection.Save(storingValues)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	it.SetID(newID)

	return nil
}
