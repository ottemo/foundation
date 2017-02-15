package page

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine
func init() {
	cmsPageInstance := new(DefaultCMSPage)
	var _ cms.InterfaceCMSPage = cmsPageInstance
	if err := models.RegisterModel(cms.ConstModelNameCMSPage, cmsPageInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "aa7d8cab-b0b0-42e4-92b7-b45761c91b15", err.Error())
	}

	cmsPageCollectionInstance := new(DefaultCMSPageCollection)
	var _ cms.InterfaceCMSPageCollection = cmsPageCollectionInstance
	if err := models.RegisterModel(cms.ConstModelNameCMSPageCollection, cmsPageCollectionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6516dd1e-2773-4c77-bee7-13c974f635f2", err.Error())
	}

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)

	if err := utils.RegisterTemplateFunction("page", pageTemplateDirective); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "642eedcf-ea33-4bce-9149-7085cd9c4377", err.Error())
	}
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCmsPageCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddColumn("enabled", db.ConstTypeBoolean, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "831e8f1b-5f43-4d59-9a51-2ad8a90b4f81", err.Error())
	}
	if err := collection.AddColumn("identifier", db.ConstTypeVarchar, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "cfed807f-a6d7-4105-ac7c-c33c0e7a73eb", err.Error())
	}
	if err := collection.AddColumn("title", db.ConstTypeVarchar, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a17c2751-63ae-4e46-89de-745e39401a75", err.Error())
	}
	if err := collection.AddColumn("content", db.ConstTypeText, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "44e81d7d-5953-43db-b291-4417e2e41e4e", err.Error())
	}
	if err := collection.AddColumn("created_at", db.ConstTypeDatetime, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a6181911-3a49-48a7-ab6d-f67d9132adf4", err.Error())
	}
	if err := collection.AddColumn("updated_at", db.ConstTypeDatetime, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8b960b9c-cf86-4e0d-a170-63bb377c0e7f", err.Error())
	}

	return nil
}

// pageTemplateDirective - text templates directive can be used to get page contents by identifier
//   use {{page "Identifier" .}} for recursive template processing; {{page "Identifier"}} - one level only
func pageTemplateDirective(identifier string, args ...interface{}) string {
	const contextStackKey = "pageDirectiveStack"

	var context map[string]interface{}
	var stack []string

	if len(args) == 1 {
		if mapValue, ok := args[0].(map[string]interface{}); ok {
			context = mapValue

			if stackValue, present := context[contextStackKey]; present {
				if arrayValue, ok := stackValue.([]string); ok {
					stack = arrayValue
				}
			}
		}
	}

	page, err := cms.LoadCMSBlockByIdentifier(identifier)
	if err != nil {
		return ""
	}
	pageContents := page.GetContent()

	if context == nil {
		return pageContents
	}

	// prevents infinite loop
	for _, stackIdentifier := range stack {
		if stackIdentifier == identifier {
			context[contextStackKey] = []string{}
			return ""
		}
	}
	stack = append(stack, identifier)
	context[contextStackKey] = stack

	result, err := utils.TextTemplate(pageContents, context)
	if err != nil {
		return pageContents
	}

	return result
}
