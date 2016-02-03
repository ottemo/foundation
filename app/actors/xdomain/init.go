package xdomain

import "github.com/ottemo/foundation/api"

// init performs self-initialization routine before app start
func init() {
	api.RegisterOnRestServiceStart(setupAPI)
}
