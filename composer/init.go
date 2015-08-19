package impex

import (
	"github.com/ottemo/foundation/api"
)

// init makes package self-initialization routine
func init() {
	api.RegisterOnRestServiceStart(setupAPI)
}
