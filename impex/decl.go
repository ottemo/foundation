// Package impex is a implementation of import/export service
package impex

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"regexp"
)

// Package global constants
const (
	ConstErrorModule = "impex"
	ConstErrorLevel  = env.ConstErrorLevelService

	ConstLogFileName = "impex.log"
)

// Package global variables
var (
	ConstImpexLog = true  // flag indicates to make log of values going to be processed
	ConstDebugLog = false // flag indicates to have extra log information

	// ConstCSVColumnRegexp is a regular expression used to grab csv column information
	//
	//	column format: [flags]path [memorize] [type] [convertors]
	//                 [~|^|?][@a.b.c.]path [={name}|>{name}] [<{type}>]
	//
	//	flags - optional column modificator
	//		format: [~|^|?]
	//		"~" - ignore column on collapse lookup
	//		"^" - array column
	//		"?" - maybe array column
	//
	//	path - attribute name in result map sub-levels separated by "."
	//		format: [@a.b.c.]d
	//		"@a" - memorized value
	//
	//	memorize - marks column to hold value in memorize map, these values can be used in path like "item.@value.label"
	//		format: ={name} | >{name}
	//		{name}  - alphanumeric value
	//		={name} - saves {column path} + {column value} to memorize map
	//		>{name}	- saves {column value} to memorize map
	//
	//	type - optional type for column
	//		format: <{type}>
	//		{type} - int | float | bool
	//
	//	convertors - text template modifications you can apply to value before use it
	//		format: see (http://golang.org/pkg/text/template/)
	ConstCSVColumnRegexp = regexp.MustCompile(`^\s*([~^?])?((?:@?\w+\.)*@?\w+)(\s+(?:=|>)\s*\w+)?(?:\s+<([^>]+)>)?\s*(.*)$`)

	ConversionFuncs = map[string]interface{}{}

	// set of service import commands
	importCmd = make(map[string]InterfaceImpexImportCmd)

	impexModels = make(map[string]InterfaceImpexModel)
)

// ImportCmdAttributeAdd is a implementer of InterfaceImpexImportCmd
//  - command allows to create custom attributes on model
type ImportCmdAttributeAdd struct {
	model     models.InterfaceModel
	attribute models.StructAttributeInfo
}

// ImportCmdImport is a implementer of InterfaceImpexImportCmd
//  - command allows to work with InterfaceImpexModel instances
type ImportCmdImport struct {
	model      InterfaceImpexModel
	attributes map[string]bool
}

// ImportCmdInsert is a implementer of InterfaceImpexImportCmd
//  - command allows to upload data in system through model item abstraction
type ImportCmdInsert struct {
	model      models.InterfaceModel
	attributes map[string]bool
	skipErrors bool
}

// ImportCmdUpdate is a implementer of InterfaceImpexImportCmd
//  - command allows to update model item data in system
type ImportCmdUpdate struct {
	model      models.InterfaceModel
	attributes map[string]bool
	idKey      string
}

// ImportCmdDelete is a implementer of InterfaceImpexImportCmd
//  - command allows to delete model items from the system
type ImportCmdDelete struct {
	model models.InterfaceModel
	idKey string
}

// ImportCmdMedia is a implementer of InterfaceImpexImportCmd
//  - command allows to assign media content to models item
type ImportCmdMedia struct {
	mediaField string
	mediaType  string
	mediaName  string
}

// ImportCmdStore is a implementer of InterfaceImpexImportCmd
//  - command allows temporary load/store previous command results during import process
type ImportCmdStore struct {
	storeObjectAs string
	storeValueAs  map[string]string

	prefix    string
	prefixKey string
}

// ImportCmdAlias is a implementer of InterfaceImpexImportCmd
//  - command allows to make record field alias to object attribute value
type ImportCmdAlias struct {
	aliases map[string]string
}
