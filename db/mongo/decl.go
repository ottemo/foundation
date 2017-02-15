package mongo

import (
	"crypto/x509"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"time"

	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	attributeTypes      = make(map[string]map[string]string) // cached values of collection attribute types
	attributeTypesMutex sync.RWMutex                         // syncronization for attributeTypes modification
)

// Package global constants
const (
	ConstConnectionValidateInterval = time.Second * 10 // timer interval to ping connection and refresh it by perforce

	ConstMongoDebug = false // flag which indicates to perform log on each operation

	ConstFilterGroupStatic  = "static"  // name for static filter, ref. to AddStaticFilter(...)
	ConstFilterGroupDefault = "default" // name for default filter, ref. to by AddFilter(...)

	ConstCollectionNameColumnInfo = "collection_column_info" // collection name to hold Ottemo types of attributes

	ConstErrorModule = "db/mongo"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// StructDBFilterGroup is a structure to hold information of named collection filter
type StructDBFilterGroup struct {
	Name         string
	FilterValues []bson.D
	ParentGroup  string
	OrSequence   bool
}

// DBCollection is a implementer of InterfaceDBCollection
type DBCollection struct {
	database   *mgo.Database
	collection *mgo.Collection

	subcollections []*DBCollection
	subresults     []*bson.Raw

	Name string

	FilterGroups map[string]*StructDBFilterGroup

	Sort []string

	ResultAttributes []string

	Limit  int
	Offset int
}

// DBEngine is a implementer of InterfaceDBEngine
type DBEngine struct {
	database *mgo.Database
	session  *mgo.Session

	DBName      string
	collections map[string]bool

	isConnected bool
}

// connectionParamsType describes params required to connect to DB
type connectionParamsType struct {
	UseSSL      bool
	DBUri       string
	DBName      string
	CertPoolPtr *x509.CertPool
}
