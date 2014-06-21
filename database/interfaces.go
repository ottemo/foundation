package database

// DBEngine is the interface that holds the necessary methods required for a database engine.
type DBEngine interface {
	GetName() string

	CreateCollection(Name string) error
	GetCollection(Name string) (DBCollection, error)
	HasCollection(Name string) bool
}

// DBCollection is the interface that contains all the necessary methods for a
// collection when working with extendable objects.
type DBCollection interface {
	Load() ([]map[string]interface{}, error)
	LoadByID(id string) (map[string]interface{}, error)

	Save(map[string]interface{}) (string, error)

	Delete() (int, error)
	DeleteByID(id string) error

	AddFilter(ColumnName string, Operator string, Value string) error
	ClearFilters() error

	ListColumns() map[string]string
	HasColumn(ColumnName string) bool

	AddColumn(ColumnName string, ColumnType string, indexed bool) error
	RemoveColumn(ColumnName string) error
}
