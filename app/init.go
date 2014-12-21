package app

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)
}
