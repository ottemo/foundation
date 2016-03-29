package reporting

import (
	"github.com/ottemo/foundation/api"
	// "github.com/ottemo/foundation/app"
	// "github.com/ottemo/foundation/app/models"
	// "github.com/ottemo/foundation/app/models/subscription"
	// "github.com/ottemo/foundation/db"
	// "github.com/ottemo/foundation/env"
	// "github.com/ottemo/foundation/utils"
)

func init() {
	api.RegisterOnRestServiceStart(setupAPI)

}
