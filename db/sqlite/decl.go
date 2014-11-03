package sqlite

import (
	"regexp"
	"sync"

	"github.com/mxk/go-sqlite/sqlite3"
)

const (
	UUID_ID = true

	DEBUG_SQL = false

	FILTER_GROUP_STATIC  = "static"
	FILTER_GROUP_DEFAULT = "default"

	COLLECTION_NAME_COLUMN_INFO = "collection_column_info"
)

var SQL_NAME_VALIDATOR = regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")

type T_DBFilterGroup struct {
	Name         string
	FilterValues []string
	ParentGroup  string
	OrSequence   bool
}

type SQLiteCollection struct {
	Name string

	ResultColumns []string
	FilterGroups  map[string]*T_DBFilterGroup
	Order         []string

	Limit string
}

type SQLite struct {
	connection      *sqlite3.Conn
	connectionMutex sync.RWMutex

	attributeTypes      map[string]map[string]string
	attributeTypesMutex sync.RWMutex
}
