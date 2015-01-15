// +build mongo

package basebuild

import (
	// MongoDB based database service
	_ "github.com/ottemo/foundation/db/mongo"
)

func init() {
	println("MongoDB in use")
}
