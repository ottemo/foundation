package mongo

import (
	"gopkg.in/mgo.v2"
	"time"

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
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "9cbde45b-17c0-4a45-b0cb-c261db261458", "Can't connect to DBEngine")
	}

	it.session = session
	it.database = session.DB(DBName)
	it.DBName = DBName
	it.collections = map[string]bool{}

	// if ConstMongoDebug {
	// 	mgo.SetDebug(true)
	// 	mgo.SetLogger(it)
	// }

	// timer routine to check connection state and reconnect by perforce
	ticker := time.NewTicker(ConstConnectionValidateInterval)
	go func() {
		for _ = range ticker.C {
			err := it.session.Ping()
			if err != nil {
				it.session.Refresh()
			}
		}
	}()

	if collectionsList, err := it.database.CollectionNames(); err == nil {
		for _, collection := range collectionsList {
			it.collections[collection] = true
		}
	}

	err = db.OnDatabaseStart()

	return err
}

// Output is a implementation of mgo.log_Logger interface
func (it *DBEngine) Output(calldepth int, s string) error {
	env.Log("mongo.log", "DEBUG", s)
	return nil
}
