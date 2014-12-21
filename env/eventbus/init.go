package eventbus

import (
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultEventBus)
	instance.listeners = make(map[string][]env.FuncEventListener)

	var _ env.InterfaceEventBus = instance

	env.RegisterEventBus(instance)
}
