package saleprice

import (
	"github.com/ottemo/foundation/env"
)

// Helper to produce new module level error
func newErrorHelper(msg, code string) error {
	return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, code, msg)
}

func logDebugHelper(msg string) {
	env.GetLogger().Log("errors.log", env.ConstLogPrefixDebug, msg)
}
