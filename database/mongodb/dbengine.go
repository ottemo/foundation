package mongodb

import (
	"errors"
	"strings"

	config "github.com/ottemo/foundation/config"
	database "github.com/ottemo/foundation/database"
	"labix.org/v2/mgo"
)

type MongoDB struct {
	database *mgo.Database
	session  *mgo.Session

	DBName      string
	collections map[string]bool
}

func init() {
	instance := new(MongoDB)

	config.RegisterOnConfigIniStart(instance.Startup)
	database.RegisterDBEngine(instance)
}

func (it *MongoDB) Startup() error {

	var DBUri = "mongodb://localhost:27017/ottemo"
	var DBName = "ottemo"

	if iniConfig := config.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("mongodb.uri"); iniValue != "" {
			DBUri = iniValue
		}

		if iniValue := iniConfig.GetValue("mongodb.db"); iniValue != "" {
			DBName = iniValue
		}
	}

	session, err := mgo.Dial(DBUri)
	if err != nil {
		return errors.New("Can't connect to MongoDB")
	}

	it.session = session
	it.database = session.DB(DBName)
	it.DBName = DBName
	it.collections = map[string]bool{}

	if collectionsList, err := it.database.CollectionNames(); err == nil {
		for _, collection := range collectionsList {
			it.collections[collection] = true
		}
	}

	database.OnDatabaseStart()

	return nil
}

func (it *MongoDB) GetName() string { return "MongoDB" }

func (it *MongoDB) HasCollection(CollectionName string) bool {
	CollectionName = strings.ToLower(CollectionName)

	return true
}

func (it *MongoDB) CreateCollection(CollectionName string) error {
	CollectionName = strings.ToLower(CollectionName)

	err := it.database.C(CollectionName).Create(new(mgo.CollectionInfo))
	//it.database.C(CollectionName).EnsureIndex(mgo.Index{Key: []string{"_id"}, Unique: true})

	//CMD := "db.createCollection(\"" + CollectionName + "\")"
	//println(CMD)
	//err := it.database.Run(CMD, nil)

	return err
}

func (it *MongoDB) GetCollection(CollectionName string) (database.DBCollection, error) {
	CollectionName = strings.ToLower(CollectionName)

	if _, present := it.collections[CollectionName]; !present {
		if err := it.CreateCollection(CollectionName); err != nil {
			return nil, err
		}
		it.collections[CollectionName] = true
	}

	mgoCollection := it.database.C(CollectionName)

	result := &MongoDBCollection{
		Selector:   map[string]interface{}{},
		Name:       CollectionName,
		database:   it.database,
		collection: mgoCollection}

	return result, nil
}
