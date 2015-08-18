// +build !sqlite,!mysql

package basebuild

import (
	// MongoDB based database service
	_ "github.com/ottemo/foundation/db/mongo"
)
