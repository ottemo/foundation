package sqlite

import (
	"fmt"
	"io"
	"testing"

	"github.com/mxk/go-sqlite/sqlite3"

	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/models"
)

//--------------------------------------------------------------------------------------------------------------
// api.InterfaceSession test implementation
//--------------------------------------------------------------------------------------------------------------

type testSession struct {
	_test_data_ map[string]interface{}
}

func (it *testSession) Close() error {
	return nil
}
func (it *testSession) Get(key string) interface{} {
	return it._test_data_[key]
}
func (it *testSession) GetID() string {
	return "ApplicationSession GetID"
}
func (it *testSession) IsEmpty() bool {
	return true
}
func (it *testSession) Set(key string, value interface{}) {
	it._test_data_[key] = value
}
func (it *testSession) Touch() error {
	return nil
}

//--------------------------------------------------------------------------------------------------------------
// api.InterfaceApplicationContext test implementation
//--------------------------------------------------------------------------------------------------------------

type testContext struct {
	//ResponseWriter    http.ResponseWriter
	//Request           *http.Request
	Request string
	//RequestParameters map[string]string
	RequestArguments map[string]string
	RequestContent   interface{}
	//RequestFiles      map[string]io.Reader

	Session       api.InterfaceSession
	ContextValues map[string]interface{}
	//Result        interface{}
}

func (it *testContext) GetRequestArguments() map[string]string {
	return it.RequestArguments
}
func (it *testContext) GetContextValues() map[string]interface{} {
	return it.ContextValues
}
func (it *testContext) GetContextValue(key string) interface{} {
	return it.ContextValues[key]
}
func (it *testContext) GetRequest() interface{} {
	return it.Request
}
func (it *testContext) GetRequestArgument(name string) string {
	return it.RequestArguments[name]
}
func (it *testContext) GetRequestContent() interface{} {
	return it.RequestContent
}
func (it *testContext) GetRequestContentType() string {
	return "request content type"
}
func (it *testContext) GetRequestFile(name string) io.Reader {
	return nil
}
func (it *testContext) GetRequestFiles() map[string]io.Reader {
	return nil
}
func (it *testContext) GetRequestSettings() map[string]interface{} {
	return map[string]interface{}{}
}
func (it *testContext) GetRequestSetting(name string) interface{} {
	return "request setting"
}
func (it *testContext) GetResponse() interface{} {
	return "response"
}
func (it *testContext) GetResponseContentType() string {
	return "response content type"
}
func (it *testContext) GetResponseResult() interface{} {
	return "response result"
}
func (it *testContext) GetResponseSetting(name string) interface{} {
	return "response setting"
}
func (it *testContext) GetResponseWriter() io.Writer {
	return nil
}
func (it *testContext) GetSession() api.InterfaceSession {
	return it.Session
}
func (it *testContext) SetContextValue(key string, value interface{}) {
	//return it.Session
}
func (it *testContext) SetResponseContentType(mimeType string) error {
	return nil
}
func (it *testContext) SetResponseResult(value interface{}) error {
	return nil
}
func (it *testContext) SetResponseSetting(name string, value interface{}) error {
	return nil
}
func (it *testContext) SetResponseStatus(code int) {
	//return nil
}
func (it *testContext) SetResponseStatusBadRequest()          {}
func (it *testContext) SetResponseStatusForbidden()           {}
func (it *testContext) SetResponseStatusNotFound()            {}
func (it *testContext) SetResponseStatusInternalServerError() {}
func (it *testContext) SetSession(session api.InterfaceSession) error {
	it.Session = session
	return nil
}

//--------------------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------------------

func TestApplyFilters(t *testing.T) {
	var err error

	// init session
	session := new(testSession)
	session._test_data_ = map[string]interface{}{}

	// init context
	context := new(testContext)
	context.SetSession(session)

	// create test db in memory
	if newConnection, err := sqlite3.Open(":memory:"); err == nil {
		dbEngine.connection = newConnection
	} else {
		t.Error("sqlite3.Open", err)
	}

	// create helper table
	SQL := "CREATE TABLE IF NOT EXISTS " + ConstCollectionNameColumnInfo + ` (
		_id        INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		collection VARCHAR(255),
		column     VARCHAR(255),
		type       VARCHAR(255),
		indexed    NUMERIC)`

	err = dbEngine.connection.Exec(SQL)
	if err != nil {
		t.Error("dbEngine.connection.Exec", err)
	}

	// create test table
	if err := dbEngine.CreateCollection("testTable"); err != nil {
		t.Error("dbEngine.CreateCollection", err)
	}

	var dbCollection = &DBCollection{
		Name:         "testTable",
		FilterGroups: make(map[string]*StructDBFilterGroup),
	}

	err = dbCollection.AddColumn("type", "string", false)
	if err != nil {
		t.Error("dbCollection.AddColumn", err)
	}

	var wrongGroupName = "_initial_value_"
	var sqls = []string{}

	// Order of keys in this map is IMPORTANT!!!
	context.RequestArguments = map[string]string{
		"_id":    "!=58592a4d9ccee8613b5f16e8,58591b893792efc42e122da5",
		"action": "count",
	}

	// We should test functionality few times, because of map processing could take values
	// from map in different order. Initial loops count is "tryCount" - abstract value.
	var tryCount = 10
	for i := 0; i < tryCount; i++ {
		// Some of test runs will be executed with additional "request argument"
		// ApplyFilters will process "request arguments" in random order
		// That's why we use loop of "tryCount"
		// It will be visible by running test with "-v" option
		if i == 1 {
			context.RequestArguments["type"] = "!=configurable"
		}

		if err := dbCollection.ClearFilters(); err != nil {
			t.Error("dbCollection.ClearFilters", err)
			continue
		}

		if err := models.ApplyFilters(context, dbCollection); err != nil {
			t.Error("models.ApplyFilters", err)
			continue
		}

		for groupName := range dbCollection.FilterGroups {
			if groupName != ConstFilterGroupDefault {
				wrongGroupName = groupName
			}
		}

		sqls = append(sqls, dbCollection.getSQLFilters())
	}

	if wrongGroupName != "_initial_value_" {
		t.Error("Invalid group name '" + wrongGroupName + "'. Should be '" + ConstFilterGroupDefault + "'")
	}

	for _, SQL := range sqls {
		fmt.Println(SQL)
	}
}
