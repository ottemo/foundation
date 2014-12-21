package page

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
)

// GetCollection returns collection of current instance type
func (it *DefaultCMSPage) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(cms.ConstModelNameCMSPageCollection)
	if result, ok := model.(cms.InterfaceCMSPageCollection); ok {
		return result
	}

	return nil
}
