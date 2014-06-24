package mongodb

import (
	"errors"
	"strings"

	config "github.com/ottemo/foundation/config"
	database "github.com/ottemo/foundation/database"
	"labix.org/v2/mgo"
)

// MongoDB struct holds the MongoDB database information.
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

// Startup initializes and creates a new MongoDb instance.
func (mdb *MongoDB) Startup() error {

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

	mdb.session = session
	mdb.database = session.DB(DBName)
	mdb.DBName = DBName
	mdb.collections = map[string]bool{}

	if collectionsList, err := mdb.database.CollectionNames(); err == nil {
		for _, collection := range collectionsList {
			mdb.collections[collection] = true
		}
	}

	database.OnDatabaseStart()

	return nil
}

// GetName returns the datastore type of MongoDB.
func (mdb *MongoDB) GetName() string { return "MongoDB" }

// HasCollection returns a boolean true if the requested collection exists.
func (mdb *MongoDB) HasCollection(CollectionName string) bool {
	CollectionName = strings.ToLower(CollectionName)

	return true
}

// CreateCollection will instantiate a new collection in MongoDB or return
// an error.
func (mdb *MongoDB) CreateCollection(CollectionName string) error {
	CollectionName = strings.ToLower(CollectionName)

	err := mdb.database.C(CollectionName).Create(new(mgo.CollectionInfo))
	//mdb.database.C(CollectionName).EnsureIndex(mgo.Index{Key: []string{"_id"}, Unique: true})

	//CMD := "db.createCollection(\"" + CollectionName + "\")"
	//println(CMD)
	//err := mdb.database.Run(CMD, nil)

	return err
}

// GetCollection returns the requested MongoDB collection.
func (mdb *MongoDB) GetCollection(CollectionName string) (database.DBCollection, error) {
	CollectionName = strings.ToLower(CollectionName)

	if _, present := mdb.collections[CollectionName]; !present {
		if err := mdb.CreateCollection(CollectionName); err != nil {
			return nil, err
		}
		mdb.collections[CollectionName] = true
	}

	mgoCollection := mdb.database.C(CollectionName)

	result := &MongoDBCollection{
		Selector:   map[string]interface{}{},
		Name:       CollectionName,
		database:   mdb.database,
		collection: mgoCollection}

	return result, nil
}
