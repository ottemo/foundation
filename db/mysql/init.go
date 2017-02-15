package mysql

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	dbEngine = new(DBEngine)
	dbEngine.attributeTypes = make(map[string]map[string]string)

	var _ db.InterfaceDBEngine = dbEngine

	var dbConnector = db.NewDBConnector(dbEngine)
	env.RegisterOnConfigIniStart(dbConnector.ConnectAsync)

	db.RegisterDBEngine(dbEngine)
}
