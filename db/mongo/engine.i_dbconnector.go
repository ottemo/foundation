package mongo

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net"
	"time"

	"github.com/ottemo/foundation/env"
)

// ------------------------------------------------------------------------------------
// InterfaceDBConnector implementation (package "github.com/ottemo/foundation/db/interfaces")
// ------------------------------------------------------------------------------------

// GetConnectionParams returns configured DB connection params
func (it *DBEngine) GetConnectionParams() interface{} {
	var connectionParams = connectionParamsType{
		DBUri:       "mongodb://localhost:27017/ottemo",
		DBName:      "ottemo",
		UseSSL:      false,
		CertPoolPtr: x509.NewCertPool(),
	}
	var Cert = ""

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("mongodb.uri", connectionParams.DBUri); iniValue != "" {
			connectionParams.DBUri = iniValue
		}

		if iniValue := iniConfig.GetValue("mongodb.db", connectionParams.DBName); iniValue != "" {
			connectionParams.DBName = iniValue
		}
		if iniValue := iniConfig.GetValue("ssl.cert", Cert); iniValue != "" {
			connectionParams.UseSSL = true
			Cert = iniValue
			if ca, err := ioutil.ReadFile(Cert); err == nil {
				connectionParams.CertPoolPtr.AppendCertsFromPEM(ca)
			}
		}
	}

	return connectionParams
}

// Connect establishes DB connection
func (it *DBEngine) Connect(srcConnectionParams interface{}) error {
	connectionParams, ok := srcConnectionParams.(connectionParamsType)
	if !ok {
		return errors.New("Wrong connection parameters type.")
	}

	if !connectionParams.UseSSL {
		// use plain text
		if session, err := mgo.Dial(connectionParams.DBUri); err == nil {
			it.session = session
		} else {
			return err
		}
	} else {
		tlsConfig := &tls.Config{}
		tlsConfig.RootCAs = connectionParams.CertPoolPtr

		// make tls connection
		dialInfo, _ := mgo.ParseURL(connectionParams.DBUri)
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
		if session, err := mgo.DialWithInfo(dialInfo); err == nil {
			it.session = session
		} else {
			return err
		}
	}

	return nil
}

// AfterConnect makes initialization of DB engine
func (it *DBEngine) AfterConnect(srcConnectionParams interface{}) error {
	connectionParams, ok := srcConnectionParams.(connectionParamsType)
	if !ok {
		return errors.New("Wrong connection parameters type.")
	}

	it.database = it.session.DB(connectionParams.DBName)
	it.DBName = connectionParams.DBName
	it.collections = map[string]bool{}

	if collectionsList, err := it.database.CollectionNames(); err == nil {
		for _, collection := range collectionsList {
			it.collections[collection] = true
		}
	}

	return nil
}

// Reconnect tries to reconnect to DB
func (it *DBEngine) Reconnect(connectionParams interface{}) error {
	it.session.Refresh()
	return nil
}

// IsConnected returns connection status
func (it *DBEngine) IsConnected() bool {
	return it.isConnected
}

// SetConnected sets connection status
func (it *DBEngine) SetConnected(connected bool) {
	it.isConnected = connected
}

// Ping checks connection alive
func (it *DBEngine) Ping() error {
	return it.session.Ping()
}

// GetValidationInterval returns delay between Ping
func (it *DBEngine) GetValidationInterval() time.Duration {
	return ConstConnectionValidateInterval
}

// GetEngineName returns DBEngine name (InterfaceDBConnector)
func (it *DBEngine) GetEngineName() string {
	return it.GetName()
}

// LogConnection outputs message to log
func (it *DBEngine) LogConnection(message string) {
	// ignore error
	_ = it.Output(0, message)
}
