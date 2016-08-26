package seo

// DefaultSEOItem type implements:
// 	- InterfaceSEOItem
// 	- InterfaceModel
// 	- InterfaceObject
// 	- InterfaceStorable

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/seo"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// ---------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models")
// ---------------------------------------------------------------------------------

// GetModelName returns model name
func (it *DefaultSEOItem) GetModelName() string {
	return seo.ConstModelNameSEOItem
}

// GetImplementationName returns model implementation name
func (it *DefaultSEOItem) GetImplementationName() string {
	return "Default" + seo.ConstModelNameSEOItem
}

// New returns new instance of model implementation object
func (it *DefaultSEOItem) New() (models.InterfaceModel, error) {
	newInstance := new(DefaultSEOItem)

	return newInstance, nil
}

// -------------------------------------------------------------------------------------------
// InterfaceSEOItem implementation (package "github.com/ottemo/foundation/app/models/seo")
// -------------------------------------------------------------------------------------------

// GetURL returns url for the given seo
func (it *DefaultSEOItem) GetURL() string {
	return it.Url
}

// SetURL sets url for the given seo
func (it *DefaultSEOItem) SetURL(value string) error {
	it.Url = value

	return nil
}

// GetRewrite returns object id for the given seo
func (it *DefaultSEOItem) GetRewrite() string {
	return it.Rewrite
}

// GetType returns object type for the given seo
func (it *DefaultSEOItem) GetType() string {
	return it.Type
}

// GetTitle returns title for the given seo
func (it *DefaultSEOItem) GetTitle() string {
	return it.Title
}

// GetMetaDescription returns description for the given seo
func (it *DefaultSEOItem) GetMetaDescription() string {
	return it.MetaDescription
}

// GetMetaKeywords returns keywords for the given seo
func (it *DefaultSEOItem) GetMetaKeywords() string {
	return it.MetaKeywords
}

// ------------------------------------------------------------------------------------
// InterfaceStorable implementation (package "github.com/ottemo/foundation/app/models")
// ------------------------------------------------------------------------------------

// GetID returns current seo id
func (it *DefaultSEOItem) GetID() string {
	return it.id
}

// SetID sets current seo id
func (it *DefaultSEOItem) SetID(id string) error {
	it.id = id

	return nil
}

// Load loads seo information from DB
func (it *DefaultSEOItem) Load(id string) error {

	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := collection.LoadByID(id)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bd4b64e9-ed73-4203-b987-ab7d72f991a0", "Unable to find seo item by id; "+id)
	}

	err = it.FromHashMap(dbRecord)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Delete removes current seo item from DB
func (it *DefaultSEOItem) Delete() error {
	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Save stores current seo item to DB
func (it *DefaultSEOItem) Save() error {
	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if it.GetURL() == "" || it.GetRewrite() == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fcfdf3cb-3ffe-4399-8678-499abaa880bf", "url and rewrite should be specified")
	}

	valuesToStore := it.ToHashMap()

	newID, err := collection.Save(valuesToStore)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.SetID(newID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// ----------------------------------------------------------------------------------
// InterfaceObject implementation (package "github.com/ottemo/foundation/app/models")
// ----------------------------------------------------------------------------------

// Get returns an object attribute value or nil
func (it *DefaultSEOItem) Get(attribute string) interface{} {

	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "url":
		return it.Url
	case "rewrite":
		return it.Rewrite
	case "type":
		return it.Type
	case "title":
		return it.Title
	case "meta_keywords":
		return it.MetaKeywords
	case "meta_description":
		return it.MetaDescription
	}

	return nil
}

// Set will apply the given attribute value to the seo or return an error
func (it *DefaultSEOItem) Set(attribute string, value interface{}) error {
	lowerCaseAttribute := strings.ToLower(attribute)

	switch lowerCaseAttribute {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)
	case "url":
		it.Url = utils.InterfaceToString(value)
	case "rewrite":
		it.Rewrite = utils.InterfaceToString(value)
	case "type":
		it.Type = utils.InterfaceToString(value)
	case "title":
		it.Title = utils.InterfaceToString(value)
	case "meta_keywords":
		it.MetaKeywords = utils.InterfaceToString(value)
	case "meta_description":
		it.MetaDescription = utils.InterfaceToString(value)
	default:
		return env.ErrorNew(
			ConstErrorModule,
			ConstErrorLevel,
			"03648446-637b-4249-8cdc-4560f2ed3c58", "unknown attribute "+attribute+" for SEOItem")
	}

	return nil
}

// FromHashMap will populate object attributes from map[string]interface{}
func (it *DefaultSEOItem) FromHashMap(input map[string]interface{}) error {
	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.LogError(err)
		}
	}
	return nil
}

// ToHashMap returns a map[string]interface{}
func (it *DefaultSEOItem) ToHashMap() map[string]interface{} {
	result := make(map[string]interface{})

	result["_id"] = it.id

	result["url"] = it.Url
	result["rewrite"] = it.Rewrite

	result["type"] = it.Type
	result["title"] = it.Title
	result["meta_keywords"] = it.MetaKeywords
	result["meta_description"] = it.MetaDescription

	return result
}

// GetAttributesInfo returns the requested object attributes
func (it *DefaultSEOItem) GetAttributesInfo() []models.StructAttributeInfo {
	result := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: ConstCollectionNameURLRewrites,
			Attribute:  "_id",
			Type:       db.ConstTypeID,
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: ConstCollectionNameURLRewrites,
			Attribute:  "url",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Url",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: ConstCollectionNameURLRewrites,
			Attribute:  "rewrite",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Rewrite",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: ConstCollectionNameURLRewrites,
			Attribute:  "type",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Type",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: ConstCollectionNameURLRewrites,
			Attribute:  "title",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Title",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: ConstCollectionNameURLRewrites,
			Attribute:  "meta_keywords",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Meta Keywrods",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: ConstCollectionNameURLRewrites,
			Attribute:  "meta_description",
			Type:       db.ConstTypeText,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Meta Description",
			Group:      "General",
			Editors:    "multiline_text",
			Options:    "",
			Default:    "",
		},
	}

	return result
}
