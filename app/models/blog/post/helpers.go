package post

import (
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models"
)

// GetBlogPostModel retrieves current InterfaceBlogPost model implementation
func GetBlogPostModel() (InterfaceBlogPost, error) {
	model, err := models.GetModel(ConstModelNameBlogPost)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	blogPostModel, ok := model.(InterfaceBlogPost)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d2fd8dc4-094e-4a32-8e89-d8be2f3bf1ea", "model "+model.GetImplementationName()+" is not 'InterfaceBlogPost' capable")
	}

	return blogPostModel, nil
}
