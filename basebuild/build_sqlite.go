// +build !mongo

package basebuild

import (
	// SQLite based database service
	_ "github.com/ottemo/foundation/db/sqlite"
)

func init() {
	println("SQLite in use")
}
