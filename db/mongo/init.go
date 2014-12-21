package mongo

import (
	"gopkg.in/mgo.v2"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DBEngine)

	env.RegisterOnConfigIniStart(instance.Startup)
	db.RegisterDBEngine(instance)
}

// Startup is a database engine startup routines
func (it *DBEngine) Startup() error {

	var DBUri = "mongodb://localhost:27017/ottemo"
	var DBName = "ottemo"

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("mongodb.uri", DBUri); iniValue != "" {
			DBUri = iniValue
		}

		if iniValue := iniConfig.GetValue("mongodb.db", DBName); iniValue != "" {
			DBName = iniValue
		}
	}

	session, err := mgo.Dial(DBUri)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "9cbde45b17c04a45b0cbc261db261458", "Can't connect to DBEngine")
	}

	it.session = session
	it.database = session.DB(DBName)
	it.DBName = DBName
	it.collections = map[string]bool{}

	if ConstMongoDebug {
		mgo.SetDebug(true)
		mgo.SetLogger(it)
	}

	if collectionsList, err := it.database.CollectionNames(); err == nil {
		for _, collection := range collectionsList {
			it.collections[collection] = true
		}
	}

	db.OnDatabaseStart()

	return nil
}

// Output is a implementation of mgo.log_Logger interface
func (it *DBEngine) Output(calldepth int, s string) error {
	env.Log("mongo.log", "DEBUG", s)
	return nil
}
