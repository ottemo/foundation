package quickbooks

import (
	"github.com/ottemo/foundation/api"
)

// init makes package self-initialization routine before app start
func init() {
	api.RegisterOnRestServiceStart(setupAPI)
}
