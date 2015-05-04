package page

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/app/models/seo"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine
func init() {
	cmsPageInstance := new(DefaultCMSPage)
	var _ cms.InterfaceCMSPage = cmsPageInstance
	models.RegisterModel(cms.ConstModelNameCMSPage, cmsPageInstance)

	cmsPageCollectionInstance := new(DefaultCMSPageCollection)
	var _ cms.InterfaceCMSPageCollection = cmsPageCollectionInstance
	models.RegisterModel(cms.ConstModelNameCMSPageCollection, cmsPageCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)

	utils.RegisterTemplateFunction("page", pageTemplateDirective)

	seo.RegisterSEOType(ConstSEOTypeCMSPage, "cms/page/", cms.ConstModelNameCMSPage)
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCmsPageCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("enabled", db.ConstTypeBoolean, true)
	collection.AddColumn("identifier", db.ConstTypeVarchar, true)
	collection.AddColumn("title", db.ConstTypeVarchar, false)
	collection.AddColumn("content", db.ConstTypeText, false)
	collection.AddColumn("created_at", db.ConstTypeDatetime, false)
	collection.AddColumn("updated_at", db.ConstTypeDatetime, false)

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
