package post

// DefaultBlogPost type implements:
//	- InterfaceBlogPost
//	- InterfaceModel
//	- InterfaceObject
// 	- InterfaceStorable

import (
	"strings"
	"time"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/blog/post"
)

// ------------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models")
// ------------------------------------------------------------------------------------

// GetModelName returns model name
func (it *DefaultBlogPost) GetModelName() string {
	return post.ConstModelNameBlogPost
}

// GetImplementationName returns model implementation name
func (it *DefaultBlogPost) GetImplementationName() string {
	return "Default" + it.GetModelName()
}

// New returns new instance of model implementation object
func (it *DefaultBlogPost) New() (models.InterfaceModel, error) {
	return new(DefaultBlogPost), nil
}

// ------------------------------------------------------------------------------------
// InterfaceStorable implementation (package "github.com/ottemo/foundation/app/models")
// ------------------------------------------------------------------------------------

// GetID returns current blog post id
func (it *DefaultBlogPost) GetID() string {
	return it.id
}

// SetID sets current blog post id
func (it *DefaultBlogPost) SetID(id string) error {
	it.id = id
	return nil
}

// Load loads blog post information from DB
func (it *DefaultBlogPost) Load(id string) error {

	collection, err := db.GetCollection(ConstBlogPostCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := collection.LoadByID(id)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "14c82544-03d3-438a-b3dd-cdd55a31b205", "Unable to find blog post by id; "+id)
	}

	if err = it.FromHashMap(dbRecord); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Delete removes current blog post from DB
func (it *DefaultBlogPost) Delete() error {
	collection, err := db.GetCollection(ConstBlogPostCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err = collection.DeleteByID(it.GetID()); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Save stores current blog post to DB
func (it *DefaultBlogPost) Save() error {
	collection, err := db.GetCollection(ConstBlogPostCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if it.GetIdentifier() == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b0cea03a-128f-4ad5-9673-578e3f4fd828", "identifier is not set")
	}

	// set dates
	currentTime := time.Now()
	if it.GetCreatedAt().IsZero() {
		if err := it.SetCreatedAt(currentTime); err != nil {
			return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fbd6b83d-1c91-4d84-bcfe-cf4053453cbd", "unable to set creation date")
		}
	}
	if err := it.SetUpdatedAt(currentTime); err != nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3146839c-9c42-4a21-be4a-2d2403714a4d", "unable to set update date")
	}

	// store model
	valuesToStore := it.ToHashMap()
	newID, err := collection.Save(valuesToStore)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err = it.SetID(newID); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// ----------------------------------------------------------------------------------
// InterfaceObject implementation (package "github.com/ottemo/foundation/app/models")
// ----------------------------------------------------------------------------------

// Get returns an object attribute value or nil
func (it *DefaultBlogPost) Get(attribute string) interface{} {

	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.GetID()
	case "identifier":
		return it.GetIdentifier()
	case "published":
		return it.IsPublished()
	case "title":
		return it.GetTitle()
	case "excerpt":
		return it.GetExcerpt()
	case "content":
		return it.GetContent()
	case "tags":
		return it.GetTags()
	case "featured_image":
		return it.GetFeaturedImage()
	case "created_at":
		return it.GetCreatedAt()
	case "updated_at":
		return it.GetUpdatedAt()
	}

	return nil
}

// Set will apply the given attribute value to the blog post or return an error
func (it *DefaultBlogPost) Set(attribute string, value interface{}) error {
	var err error
	lowerCaseAttribute := strings.ToLower(attribute)

	switch lowerCaseAttribute {
	case "_id", "id":
		err = it.SetID(utils.InterfaceToString(value))
	case "identifier":
		err = it.SetIdentifier(utils.InterfaceToString(value))
	case "published":
		err = it.SetPublished(utils.InterfaceToBool(value))
	case "title":
		err = it.SetTitle(utils.InterfaceToString(value))
	case "excerpt":
		err = it.SetExcerpt(utils.InterfaceToString(value))
	case "content":
		err = it.SetContent(utils.InterfaceToString(value))
	case "tags":
		err = it.SetTags(utils.InterfaceToArray(value))
	case "featured_image":
		err = it.SetFeaturedImage(utils.InterfaceToString(value))
	case "created_at":
		err = it.SetCreatedAt(utils.InterfaceToTime(value))
	case "updated_at":
		err = it.SetUpdatedAt(utils.InterfaceToTime(value))
	default:
		return env.ErrorNew(
			ConstErrorModule,
			ConstErrorLevel,
			"eb027b8f-118e-4036-a4e3-7b8abdb20ed8", "unknown attribute "+attribute+" for blog post")
	}

	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel,"9d3d9d30-4517-4be6-af09-af817dcd5215", "unable to set '"+attribute+"' value '"+utils.InterfaceToString(value)+"'")
	}

	return nil
}

// FromHashMap will populate object attributes from map[string]interface{}
func (it *DefaultBlogPost) FromHashMap(input map[string]interface{}) error {
	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.LogError(err)
		}
	}
	return nil
}

// ToHashMap returns a map[string]interface{}
func (it *DefaultBlogPost) ToHashMap() map[string]interface{} {
	result := make(map[string]interface{})

	result["_id"] = it.GetID()

	result["identifier"] = it.GetIdentifier()
	result["published"] = it.IsPublished()
	result["title"] = it.GetTitle()
	result["excerpt"] = it.GetExcerpt()
	result["content"] = it.GetContent()
	result["tags"] = it.GetTags()
	result["featured_image"] = it.GetFeaturedImage()
	result["created_at"] = it.GetCreatedAt()
	result["updated_at"] = it.GetUpdatedAt()

	return result
}

// GetAttributesInfo returns the requested object attributes
func (it *DefaultBlogPost) GetAttributesInfo() []models.StructAttributeInfo {
	result := []models.StructAttributeInfo{
		{
			Model:      post.ConstModelNameBlogPost,
			Collection: ConstBlogPostCollectionName,
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
		{
			Model:      post.ConstModelNameBlogPost,
			Collection: ConstBlogPostCollectionName,
			Attribute:  "identifier",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Identifier",
			Group:      "General",
			Editors:    "text",
			Options:    "",
			Default:    "",
		},
		{
			Model:      post.ConstModelNameBlogPost,
			Collection: ConstBlogPostCollectionName,
			Attribute:  "published",
			Type:       db.ConstTypeBoolean,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Published",
			Group:      "General",
			Editors:    "boolean",
			Options:    "",
			Default:    "",
		},
		{
			Model:      post.ConstModelNameBlogPost,
			Collection: ConstBlogPostCollectionName,
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
		{
			Model:      post.ConstModelNameBlogPost,
			Collection: ConstBlogPostCollectionName,
			Attribute:  "tags",
			Type:       "[]" + db.ConstTypeText,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Tags",
			Group:      "General",
			Editors:    "string",
			Options:    "",
			Default:    "",
		},
		{
			Model:      post.ConstModelNameBlogPost,
			Collection: ConstBlogPostCollectionName,
			Attribute:  "featured_image",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Featured Image",
			Group:      "General",
			Editors:    "text",
			Options:    "",
			Default:    "",
		},
		{
			Model:      post.ConstModelNameBlogPost,
			Collection: ConstBlogPostCollectionName,
			Attribute:  "content",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Content",
			Group:      "General",
			Editors:    "html",
			Options:    "",
			Default:    "",
		},
		{
			Model:      post.ConstModelNameBlogPost,
			Collection: ConstBlogPostCollectionName,
			Attribute:  "excerpt",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Excerpt",
			Group:      "General",
			Editors:    "html",
			Options:    "",
			Default:    "",
		},
		{
			Model:      post.ConstModelNameBlogPost,
			Collection: ConstBlogPostCollectionName,
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
		{
			Model:      post.ConstModelNameBlogPost,
			Collection: ConstBlogPostCollectionName,
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

	return result
}

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceBlogPost implementation (package "github.com/ottemo/foundation/app/models/blog/post/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// SetIdentifier : identifier setter
func (it *DefaultBlogPost) SetIdentifier(value string) error {
	it.identifier = value
	return nil
}

// GetIdentifier : identifier getter
func (it *DefaultBlogPost) GetIdentifier() string {
	return it.identifier
}

// SetPublished : published setter
func (it *DefaultBlogPost) SetPublished(value bool) error {
	it.published = value
	return nil
}

// IsPublished : published getter
func (it *DefaultBlogPost) IsPublished() bool {
	return it.published
}

// SetTitle : title setter
func (it *DefaultBlogPost) SetTitle(value string) error {
	it.title = value
	return nil
}

// GetTitle : title getter
func (it *DefaultBlogPost) GetTitle() string {
	return it.title
}

// SetExcerpt : excerpt setter
func (it *DefaultBlogPost) SetExcerpt(value string) error {
	it.excerpt = value
	return nil
}

// GetExcerpt : excerpt getter
func (it *DefaultBlogPost) GetExcerpt() string {
	return it.excerpt
}

// SetContent : content setter
func (it *DefaultBlogPost) SetContent(value string) error {
	it.content = value
	return nil
}

// GetContent : content getter
func (it *DefaultBlogPost) GetContent() string {
	return it.content
}

// SetTags : tags setter
func (it *DefaultBlogPost) SetTags(value []interface{}) error {
	it.tags = value
	return nil
}

// GetTags : tags getter
func (it *DefaultBlogPost) GetTags() []interface{} {
	return it.tags
}

// SetFeaturedImage : featuredImage setter
func (it *DefaultBlogPost) SetFeaturedImage(value string) error {
	it.featuredImage = value
	return nil
}

// GetFeaturedImage : featuredImage getter
func (it *DefaultBlogPost) GetFeaturedImage() string {
	return it.featuredImage
}

// SetCreatedAt : createdAt setter
func (it *DefaultBlogPost) SetCreatedAt(value time.Time) error {
	it.createdAt = value
	return nil
}

// GetCreatedAt : createdAt getter
func (it *DefaultBlogPost) GetCreatedAt() time.Time {
	return it.createdAt
}

// SetUpdatedAt : updatedAt setter
func (it *DefaultBlogPost) SetUpdatedAt(value time.Time) error {
	it.updatedAt = value
	return nil
}

// GetUpdatedAt : updatedAt getter
func (it *DefaultBlogPost) GetUpdatedAt() time.Time {
	return it.updatedAt
}
