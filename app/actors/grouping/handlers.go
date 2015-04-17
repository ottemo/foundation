package grouping

import (
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
)

func updateCartHandler(event string, eventData map[string]interface{}) bool {

	if sessionInstance, ok := eventData["session"].(api.InterfaceSession); ok {

	}

	return true
}

