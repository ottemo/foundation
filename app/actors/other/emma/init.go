package emma

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/api"
)

func init() {
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)
}

