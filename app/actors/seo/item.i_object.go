package seo

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/seo"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/utils"
)

// Get returns object attribute value or nil
func (it *DefaultSEOItem) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.GetID()
	case "type":
		return it.GetType()
	case "url":
		return it.Get
	case "rewrite", "object_id":
		return it.GetObjectID()
	case "title":
		return it.Title
	case "meta_keywords":
		return it.MetaKeywords
	case "meta_description":
		return it.MetaDescription
	}

	return nil
}

// Set sets attribute value to object or returns error
func (it *DefaultSEOItem) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		return it.SetID(utils.InterfaceToString(value))
	case "type":
		return it.SetType(utils.InterfaceToString(value))
	case "url":
		return it.SetURL(utils.InterfaceToString(value))
	case "rewrite", "object_id":
		return it.SetObjectID(utils.InterfaceToString(value))
	case "title":
		it.Title = utils.InterfaceToString(value)
		return nil
	case "meta_keywords":
		it.MetaKeywords = utils.InterfaceToString(value)
		return nil
	case "meta_description":
		it.MetaDescription = utils.InterfaceToString(value)
		return nil
	}

	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e6af9084-fddc-45bf-a90c-9c3d6ff88a57", "unknown attribute '"+attribute+"'")
}

// FromHashMap fills object attributes from map[string]interface{}
func (it *DefaultSEOItem) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap represents object as map[string]interface{}
func (it *DefaultSEOItem) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.GetID()
	result["type"] = it.Get("type")
	result["url"] = it.Get("url")
	result["rewrite"] = it.Get("rewrite")
	result["title"] = it.Get("title")
	result["meta_keywords"] = it.Get("meta_keywords")
	result["meta_description"] = it.Get("meta_description")

	return result
}

// GetAttributesInfo returns information about object attributes
func (it *DefaultSEOItem) GetAttributesInfo() []models.StructAttributeInfo {

	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: "",
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
			Collection: "",
			Attribute:  "type",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Type",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: "",
			Attribute:  "url",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "URL",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
			Validators: "",
		},
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: "",
			Attribute:  "rewrite",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Rewrite",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
			Validators: "",
		},
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: "",
			Attribute:  "title",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Title",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
			Validators: "",
		},
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: "",
			Attribute:  "meta_keywords",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Meta Keywords",
			Group:      "General",
			Editors:    "multiline_text",
			Options:    "",
			Default:    "",
			Validators: "",
		},
		models.StructAttributeInfo{
			Model:      seo.ConstModelNameSEOItem,
			Collection: "",
			Attribute:  "meta_keywords",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Meta Description",
			Group:      "General",
			Editors:    "multiline_text",
			Options:    "",
			Default:    "",
			Validators: "",
		},
	}

	return info
}
