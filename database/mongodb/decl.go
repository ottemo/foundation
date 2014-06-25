package mongodb

import (
	"labix.org/v2/mgo"
)

type MongoDBCollection struct {
	  database *mgo.Database
	collection *mgo.Collection
	      Name string

	Selector map[string]interface{}
	Sort []string

	Limit int
	Offset int

}

type MongoDB struct {
	database *mgo.Database
	session  *mgo.Session

	DBName string
	collections map[string]bool
}
