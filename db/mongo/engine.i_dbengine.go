package mongo

import (
	"strings"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GetName returns current DB engine name
func (it *DBEngine) GetName() string {
	return "MongoDB"
}

// HasCollection checks if collection already exists
func (it *DBEngine) HasCollection(CollectionName string) bool {
	CollectionName = strings.ToLower(CollectionName)

	return true
}

// CreateCollection creates cllection by it's name
func (it *DBEngine) CreateCollection(CollectionName string) error {
	CollectionName = strings.ToLower(CollectionName)

	err := it.database.C(CollectionName).Create(new(mgo.CollectionInfo))
	//it.database.C(CollectionName).EnsureIndex(mgo.Index{Key: []string{"_id"}, Unique: true})

	//CMD := "db.createCollection(\"" + CollectionName + "\")"
	//println(CMD)
	//err := it.database.Run(CMD, nil)

	return env.ErrorDispatch(err)
}

// GetCollection returns collection by name or creates new one
func (it *DBEngine) GetCollection(CollectionName string) (db.InterfaceDBCollection, error) {
	CollectionName = strings.ToLower(CollectionName)

	if _, present := it.collections[CollectionName]; !present {
		if err := it.CreateCollection(CollectionName); err != nil {
			return nil, err
		}
		it.collections[CollectionName] = true
	}

	mgoCollection := it.database.C(CollectionName)

	result := &DBCollection{
		Name:             CollectionName,
		FilterGroups:     make(map[string]*StructDBFilterGroup),
		ResultAttributes: make([]string, 0),
		database:         it.database,
		collection:       mgoCollection,
	}

	return result, nil
}

// RawQuery executes raw query for DB engine.
//   This function makes eval commang on mongo db (http://docs.mongodb.org/manual/reference/command/eval/#dbcmd.eval)
//   so if you are using "db.collection.find()" - it returns cursor object, do not forget to add ".toArray()", i.e.
//   "db.collection.find().toArray()"
func (it *DBEngine) RawQuery(query string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := it.database.Run(bson.D{{"eval", query}}, result)
	return result, err
}
