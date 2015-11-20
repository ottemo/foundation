package mailchimp

import (
	"github.com/ottemo/foundation/env"
)

func init() {
	env.RegisterOnConfigStart(setupConfig)
}
