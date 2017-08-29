package impex

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
)

// ImportStartHandler is a middleware to init importState for async "next" procedure.
func ImportStartHandler(next api.FuncAPIHandler) api.FuncAPIHandler {
	return func(context api.InterfaceApplicationContext) (interface{}, error) {
		if importStatus.state != constImportStateIdle {
			additionalMessage := ""
			if importStatus.file != nil {
				additionalMessage = " Currently processing " + importStatus.file.name
			}
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4bec46b6-b6b0-4821-8978-44d0f051750d", "Another import is in progres." + additionalMessage)
		}

		importStatus.state = constImportStateProcessing
		delete(importStatus.sessions, context.GetSession().GetID())
		return next(context)
	}
}

// ImportResultHandler will process import call's result
// It return no values, because of async handler result processing.
var ImportResultHandler = func(context api.InterfaceApplicationContext, result interface{}, err error) {
	importStatus.state = constImportStateIdle

	importStatus.sessions[context.GetSession().GetID()] = map[string]interface{}{
		"result": result,
		"err":    err,
	}
}

