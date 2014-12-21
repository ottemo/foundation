package impex

import (
	"github.com/ottemo/foundation/api"
)

// init makes package self-initialization routine
func init() {
	api.RegisterOnRestServiceStart(setupAPI)

	RegisterImportCommand("INSERT", new(ImportCmdInsert))
	RegisterImportCommand("UPDATE", new(ImportCmdUpdate))
	RegisterImportCommand("DELETE", new(ImportCmdDelete))

	RegisterImportCommand("STORE", new(ImportCmdStore))
	RegisterImportCommand("MEDIA", new(ImportCmdMedia))

	RegisterImportCommand("ATTRIBUTE_ADD", new(ImportCmdAttributeAdd))
}
