package saleprice

import (
	"github.com/ottemo/foundation/env"
)

// newErrorHelper produce new module level error is declared to minimize repeatable code
func newErrorHelper(msg, code string) error {
	return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, code, msg)
}

// logWarnHelper writes WARN level message to global errors.log file
func logWarnHelper(msg string) {
	env.GetLogger().Log("errors.log", env.ConstLogPrefixWarning, msg)
}
