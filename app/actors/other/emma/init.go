package emma

import (
	"github.com/ottemo/foundation/env"
)

func init() {
	env.RegisterOnConfigStart(setupConfig)
}

