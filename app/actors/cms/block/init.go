package block

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
	cmsBlockInstance := new(DefaultCMSBlock)
	var _ cms.InterfaceCMSBlock = cmsBlockInstance
	models.RegisterModel(cms.ConstModelNameCMSBlock, cmsBlockInstance)

	cmsBlockCollectionInstance := new(DefaultCMSBlockCollection)
	var _ cms.InterfaceCMSBlockCollection = cmsBlockCollectionInstance
	models.RegisterModel(cms.ConstModelNameCMSBlockCollection, cmsBlockCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)

	utils.RegisterTemplateFunction("block", blockTemplateDirective)

}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCmsBlockCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("identifier", db.ConstTypeVarchar, true)
	collection.AddColumn("content", db.ConstTypeText, false)
	collection.AddColumn("created_at", db.ConstTypeDatetime, false)
	collection.AddColumn("updated_at", db.ConstTypeDatetime, false)

	return nil
}

// blockTemplateDirective - text templates directive can be used to get block contents by identifier
//   use {{block "Identifier" .}} for recursive template processing; {{block "Identifier"}} - one level only
func blockTemplateDirective(identifier string, args ...interface{}) string {
	const contextStackKey = "blockDirectiveStack"

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

	// processing new identifier
	block, err := cms.LoadCMSBlockByIdentifier(identifier)
	if err != nil {
		return ""
	}
	blockContents := block.GetContent()

	if context == nil {
		return blockContents
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

	result, err := utils.TextTemplate(blockContents, context)
	if err != nil {
		return blockContents
	}

	return result
}
