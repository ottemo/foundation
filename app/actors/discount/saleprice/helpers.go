package saleprice

import (
	"github.com/ottemo/foundation/env"
)

// newErrorHelper produce new module level error is declared to minimize repeatable code
func newErrorHelper(msg, code string) error {
	return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, code, msg)
}
