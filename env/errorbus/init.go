package errorbus

import (
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := &DefaultErrorBus{listeners: make([]env.FuncErrorListener, 0)}
	var _ env.InterfaceErrorBus = instance

	var _ env.InterfaceOttemoError = new(OttemoError)

	env.RegisterErrorBus(instance)
	env.RegisterOnConfigStart(setupConfig)
}
