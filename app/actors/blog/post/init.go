package post

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/blog/post"
)

// init makes package self-initialization routine
func init() {
	blogPostInstance := new(DefaultBlogPost)
	var _ post.InterfaceBlogPost = blogPostInstance
	models.RegisterModel(post.ConstModelNameBlogPost, blogPostInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstBlogPostCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("identifier", db.ConstTypeVarchar, true)
	collection.AddColumn("published", db.ConstTypeBoolean, true)
	collection.AddColumn("title", db.ConstTypeVarchar, false)
	collection.AddColumn("excerpt", db.ConstTypeText, false)
	collection.AddColumn("content", db.ConstTypeText, false)
	collection.AddColumn("created_at", db.ConstTypeDatetime, false)
	collection.AddColumn("updated_at", db.ConstTypeDatetime, false)
	collection.AddColumn("tags", "[]"+db.ConstTypeText, false)
	collection.AddColumn("featured_image", db.ConstTypeVarchar, false)

	return nil
}
