package shipstation

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
)

func init() {
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)
}
