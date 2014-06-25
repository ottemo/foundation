package sqlite

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/ottemo/foundation/database"
)

type SQLiteCollection struct {
	Connection *sqlite3.Conn
	TableName string
	Columns map[string]string

	Filters []string
	Order []string

	Limit string
}

type SQLite struct {
	Connection *sqlite3.Conn
}

var collections = map[string]database.IDBCollection{}
