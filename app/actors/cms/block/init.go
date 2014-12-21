package block

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
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
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCmsBlockCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("identifier", "varchar(255)", true)
	collection.AddColumn("content", "text", false)
	collection.AddColumn("created_at", "datetime", false)
	collection.AddColumn("updated_at", "datetime", false)

	return nil
}
