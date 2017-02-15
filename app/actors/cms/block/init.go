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
	if err := models.RegisterModel(cms.ConstModelNameCMSBlock, cmsBlockInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b3834a39-b6c1-4b2d-9051-3c81f76216c9", err.Error())
	}

	cmsBlockCollectionInstance := new(DefaultCMSBlockCollection)
	var _ cms.InterfaceCMSBlockCollection = cmsBlockCollectionInstance
	if err := models.RegisterModel(cms.ConstModelNameCMSBlockCollection, cmsBlockCollectionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "49336ec2-3204-4fe4-b530-975b347dbd0e", err.Error())
	}

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)

	if err := utils.RegisterTemplateFunction("block", blockTemplateDirective); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "43876927-a8f4-4c5e-8ced-e8daa7faed7c", err.Error())
	}

}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCmsBlockCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddColumn("identifier", db.ConstTypeVarchar, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "28029fa1-02d7-4ec6-a482-924e6121ab5f", err.Error())
	}
	if err := collection.AddColumn("content", db.ConstTypeText, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2013333b-0f9d-4b4b-8cc7-ed9ae6d214e7", err.Error())
	}
	if err := collection.AddColumn("created_at", db.ConstTypeDatetime, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0d8b786b-4382-45a1-80cd-f8585309f43d", err.Error())
	}
	if err := collection.AddColumn("updated_at", db.ConstTypeDatetime, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5170d746-127f-40ff-ab30-05804931b84d", err.Error())
	}

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
