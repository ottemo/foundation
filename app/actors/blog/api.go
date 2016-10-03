package blog

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/blog"
)

func SetupApi() error {
	service := api.GetRestService()

	// admin
	service.GET("/blog/posts", APIGetBlogList)
	service.GET("/blog/posts/:id", APIGetPublishedPost)
	service.POST("/blog/posts", api.IsAdmin(APICreateBlogPost))
	service.POST("/blog/posts/edit/:id", api.IsAdmin(APIPublishBlogPost))
	service.PUT("/blog/posts/edit/:id", api.IsAdmin(APIUpdateBlogPost))

	return nil
}

func APIGetBlogList(context api.InterfaceApplicationContext) (interface{}, error) {
	blogCollectionModel, err := blog.GetBlogCollectionModel()
	if err {
		return nil, env.ErrorDispatch(err)
	}

	// not allowing to see disabled articles if not admin
	if err := api.ValidateAdminRights(context); err != nil {
		blogCollectionModel.GetDBCollection().AddFilter("enabled", "=", true)
	}

	// limit parameter handle
	blogCollectionModel.ListLimit(models.GetListLimit(context))

	result := blogCollectionModel.List()

	return result, nil
}

func APIGetPublishedPost(context api.InterfaceApplicationContext) (interface{}, error) {
	blogID := context.GetRequestArgument("id")
	if blogID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "51a5455b-b4de-4a8f-b0c6-9d439b44bfcc", "blog id was not specified")
	}

	if blogModel, err := blog.LoadBlogByID(blogID); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := blogModel.ToHashMap()

	// check for admin rights to see particular article
	err := api.ValidateAdminRights(context)
	if err != nil && result["enabled"] == false {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "56959fe2-de88-4a0f-863c-834745a76ea2", "you are not allowed to see this page")
	} else {
		return result, nil
	}
}

func APICreateBlogPost(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	blogModel, err := blog.GetBlogModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		err := blogModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = blogModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return blogModel.ToHashMap(), nil
}

func APIPublishBlogPost(context api.InterfaceApplicationContext) (interface{}, error) {
	blogID := context.GetRequestArgument("id")
	if blogID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3ae535d3-4f2a-4cf1-afa4-effea6400580", "blog id was not specified")
	}

	blogModel, err := blog.LoadBlogByID(blogID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	blogModel.Set("enabled", true)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	blogModel.Save()

	result := blogModel.ToHashMap()

	return result, nil
}

func APIUpdateBlogPost(context api.InterfaceApplicationContext) (interface{}, error) {
	blogID := context.GetRequestArgument("id")
	if blogID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f77a5908-2322-4bdd-a826-470056f4cdea", "blog id was not specified")
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	blogModel, err := blog.LoadBlogByID(blogID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		err := blogModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = blogModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := blogModel.ToHashMap()

	return result, nil
}
