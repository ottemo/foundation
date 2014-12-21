package block

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
)

// GetCollection returns collection of current instance type
func (it *DefaultCMSBlock) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(cms.ConstModelNameCMSBlockCollection)
	if result, ok := model.(cms.InterfaceCMSBlockCollection); ok {
		return result
	}

	return nil
}
