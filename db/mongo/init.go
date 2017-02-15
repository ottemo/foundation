package mongo

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DBEngine)

	var dbConnector = db.NewDBConnector(instance)
	env.RegisterOnConfigIniStart(dbConnector.ConnectAsync)

	_ = db.RegisterDBEngine(instance)
}

// Output is a implementation of mgo.log_Logger interface
func (it *DBEngine) Output(calldepth int, s string) error {
	env.Log("mongo.log", "DEBUG", s)
	return nil
}
