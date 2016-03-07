package mongo

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"time"

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
	var UseSSL = false
	var Cert = ""

	// root certificates
	roots := x509.NewCertPool()

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("mongodb.uri", DBUri); iniValue != "" {
			DBUri = iniValue
		}

		if iniValue := iniConfig.GetValue("mongodb.db", DBName); iniValue != "" {
			DBName = iniValue
		}
		if iniValue := iniConfig.GetValue("ssl.cert", Cert); iniValue != "" {
			UseSSL = true
			Cert = iniValue
			if ca, err := ioutil.ReadFile(Cert); err == nil {
				roots.AppendCertsFromPEM(ca)
			}
		}
	}

	// create session with db
	if !UseSSL {
		// use plain text
		if session, err := mgo.Dial(DBUri); err == nil {
			it.session = session
		} else {
			return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "9cbde45b-17c0-4a45-b0cb-c261db261458", "Can't connect to DBEngine")
		}
	} else {
		tlsConfig := &tls.Config{}
		tlsConfig.RootCAs = roots

		// make tls connection
		dialInfo, _ := mgo.ParseURL(DBUri)
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
		if session, err := mgo.DialWithInfo(dialInfo); err == nil {
			it.session = session
		} else {
			return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "948fda13-9579-4407-a22b-654ceda317c2", "Can't connect to DBEngine using SSL")
		}
	}

	it.database = it.session.DB(DBName)
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

	err := db.OnDatabaseStart()

	return err
}

// Output is a implementation of mgo.log_Logger interface
func (it *DBEngine) Output(calldepth int, s string) error {
	env.Log("mongo.log", "DEBUG", s)
	return nil
}
