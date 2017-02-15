package page

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/utils"
)

// Get returns object attribute value or nil
func (it *DefaultCMSPage) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.GetID()
	case "enabled":
		return it.GetEnabled()
	case "identifier":
		return it.GetIdentifier()
	case "title":
		return it.GetTitle()
	case "content":
		return it.GetContent()
	case "created_at":
		return it.CreatedAt
	case "updated_at":
		return it.UpdatedAt
	}

	return nil
}

// Set sets attribute value to object or returns error
func (it *DefaultCMSPage) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		return it.SetID(utils.InterfaceToString(value))
	case "enabled":
		return it.SetEnabled(utils.InterfaceToBool(value))
	case "identifier":
		return it.SetIdentifier(utils.InterfaceToString(value))
	case "title":
		return it.SetTitle(utils.InterfaceToString(value))
	case "content":
		return it.SetContent(utils.InterfaceToString(value))
	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)
		return nil
	case "updated_at":
		it.UpdatedAt = utils.InterfaceToTime(value)
		return nil
	}

	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e6af9084-fddc-45bf-a90c-9c3d6ff88a57", "unknown attribute '"+attribute+"'")
}

// FromHashMap fills object attributes from map[string]interface{}
func (it *DefaultCMSPage) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			_ = env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap represents object as map[string]interface{}
func (it *DefaultCMSPage) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.GetID()
	result["enabled"] = it.Get("enabled")
	result["identifier"] = it.Get("identifier")
	result["title"] = it.Get("title")
	result["content"] = it.Get("content")
	result["created_at"] = it.Get("created_at")
	result["updated_at"] = it.Get("updated_at")

	return result
}

// GetAttributesInfo returns information about object attributes
func (it *DefaultCMSPage) GetAttributesInfo() []models.StructAttributeInfo {

	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
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
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "enabled",
			Type:       db.ConstTypeBoolean,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Enabled",
			Group:      "General",
			Editors:    "boolean",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "identifier",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Identifier",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
			Validators: "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "title",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Title",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "content",
			Type:       db.ConstTypeText,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Content",
			Group:      "General",
			Editors:    "html",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "created_at",
			Type:       db.ConstTypeDatetime,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Created At",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "updated_at",
			Type:       db.ConstTypeDatetime,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Updated At",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
	}

	return info
}
