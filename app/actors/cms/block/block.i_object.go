package block

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Get returns object attribute value or nil
func (it *DefaultCMSBlock) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.GetID()
	case "identifier":
		return it.GetIdentifier()
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
func (it *DefaultCMSBlock) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		return it.SetID(utils.InterfaceToString(value))
	case "identifier":
		return it.SetIdentifier(utils.InterfaceToString(value))
	case "content":
		return it.SetContent(utils.InterfaceToString(value))
	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)
		return nil
	case "updated_at":
		it.UpdatedAt = utils.InterfaceToTime(value)
		return nil
	}

	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ab4ceb20-cb87-429f-a8fb-296136515077", "unknown attribute '"+attribute+"'")
}

// FromHashMap represents object as map[string]interface{}
func (it *DefaultCMSBlock) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.LogError(err)
		}
	}

	return nil
}

// ToHashMap fills object attributes from map[string]interface{}
func (it *DefaultCMSBlock) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["identifier"] = it.Get("identifier")
	result["content"] = it.Get("content")
	result["created_at"] = it.Get("created_at")
	result["updated_at"] = it.Get("updated_at")

	return result
}

// GetAttributesInfo returns information about object attributes
func (it *DefaultCMSBlock) GetAttributesInfo() []models.StructAttributeInfo {

	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSBlock,
			Collection: ConstCmsBlockCollectionName,
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
			Model:      cms.ConstModelNameCMSBlock,
			Collection: ConstCmsBlockCollectionName,
			Attribute:  "identifier",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Identifier",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
			Validators: "sku",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSBlock,
			Collection: ConstCmsBlockCollectionName,
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
			Model:      cms.ConstModelNameCMSBlock,
			Collection: ConstCmsBlockCollectionName,
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
			Model:      cms.ConstModelNameCMSBlock,
			Collection: ConstCmsBlockCollectionName,
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
