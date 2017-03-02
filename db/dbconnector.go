package db

import (
	"time"
	"fmt"

	"github.com/ottemo/foundation/utils"
)

// DBConnector takes logic of connection to DB and keeps it alive
type DBConnector struct {
	connector InterfaceDBConnector
}

// NewDBConnector returns new DBConnector instance
func NewDBConnector(connector InterfaceDBConnector) *DBConnector {
	return &DBConnector{connector: connector}
}

// ConnectAsync makes connection process asynchronously
func (it *DBConnector) ConnectAsync() error {
	go func() {
		if err := it.Connect(); err != nil {
			it.log("Internal error: unable to start "+it.connector.GetEngineName()+" DBEngine: "+err.Error())
		}
	}()
	return nil
}

// Connect implements connection process algorithm
func (it *DBConnector) Connect() error {
	var connectionParams = it.connector.GetConnectionParams()

	it.log("Connecting to "+it.connector.GetEngineName()+" DB. Timeout: 10 seconds.")

	var ticker = time.NewTicker(it.connector.GetValidationInterval())
	var isConnectMessageShown = false
	for !it.connector.IsConnected() {
		err := it.connector.Connect(connectionParams)

		if err != nil {
			if !isConnectMessageShown {
				it.log("Can't connect to DBEngine: " + err.Error())
				isConnectMessageShown = true;
			}

			it.log("Wait " + utils.InterfaceToString(10) + " seconds to reconnect.")
			_ = <-ticker.C
		} else {
			it.connector.SetConnected(true)
			it.log("DB connection established.")
		}
	}

	err := it.connector.AfterConnect(connectionParams)
	if err != nil {
		it.log("After DB connect error:")
	}

	go func() {
		for range ticker.C {
			if err := it.connector.Ping(); err != nil {
				it.connector.SetConnected(false)
				it.log("DB connection lost. Reconnect in 60 seconds.")
				if err := it.connector.Reconnect(connectionParams); err != nil {
					it.log("Unable to reconnect: "+err.Error())
				}
			} else if !it.connector.IsConnected() { // Show message only once
				it.connector.SetConnected(true)
				it.log("DB connection restored.")
			}
		}
	}()

	err = OnDatabaseStart()

	return err
}

// log outputs messages to stdout and connector endpoint
func (it *DBConnector) log(message string) {
	// output to stdout
	fmt.Println(time.Now().Format(time.RFC3339), message)
	it.connector.LogConnection(message)
}
