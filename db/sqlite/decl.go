// Package sqlite is a SQLite implementation of interfaces declared in
// "github.com/ottemo/foundation/db" package
package sqlite

import (
	"regexp"
	"sync"
	"time"

	"github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstConnectionValidateInterval = time.Second * 10 // timer interval to ping connection and refresh it by perforce

	ConstUseUUIDids = true  // flag which indicates to use UUID "_id" column type instead of default integer
	ConstDebugSQL   = false // flag which indicates to perform log on each SQL operation

	ConstFilterGroupStatic  = "static"  // name for static filter, ref. to AddStaticFilter(...)
	ConstFilterGroupDefault = "default" // name for default filter, ref. to by AddFilter(...)

	ConstCollectionNameColumnInfo = "collection_column_info" // table name to hold Ottemo types of columns

	ConstErrorModule = "db/sqlite"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// Package global variables
var (
	// dbEngine is an instance of database engine (one per application)
	dbEngine *DBEngine

	// ConstSQLNameValidator is a regex expression used to check names used within SQL queries
	ConstSQLNameValidator = regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")
)

// StructDBFilterGroup is a structure to hold information of named collection filter
type StructDBFilterGroup struct {
	Name         string
	FilterValues []string
	ParentGroup  string
	OrSequence   bool
}

// DBCollection is a InterfaceDBCollection implementer
type DBCollection struct {
	Name string

	ResultColumns []string
	FilterGroups  map[string]*StructDBFilterGroup
	Order         []string

	Limit string
}

// DBEngine is a InterfaceDBEngine implementer
type DBEngine struct {
	connection      *sqlite3.Conn
	connectionMutex sync.RWMutex

	attributeTypes      map[string]map[string]string
	attributeTypesMutex sync.RWMutex
}
