package db

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Package global constants
const (
	ConstTypeID       = utils.ConstDataTypeID
	ConstTypeBoolean  = utils.ConstDataTypeBoolean
	ConstTypeVarchar  = utils.ConstDataTypeVarchar
	ConstTypeText     = utils.ConstDataTypeText
	ConstTypeInteger  = utils.ConstDataTypeInteger
	ConstTypeDecimal  = utils.ConstDataTypeDecimal
	ConstTypeMoney    = utils.ConstDataTypeMoney
	ConstTypeFloat    = utils.ConstDataTypeFloat
	ConstTypeDatetime = utils.ConstDataTypeDatetime
	ConstTypeJSON     = utils.ConstDataTypeJSON

	ConstErrorModule = "db"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// InterfaceDBEngine represents interface to access database engine
type InterfaceDBEngine interface {
	GetName() string

	CreateCollection(Name string) error
	GetCollection(Name string) (InterfaceDBCollection, error)
	HasCollection(Name string) bool

	RawQuery(query string) (map[string]interface{}, error)
}

// InterfaceDBCollection interface to access particular table/collection of database
type InterfaceDBCollection interface {
	Load() ([]map[string]interface{}, error)
	LoadByID(id string) (map[string]interface{}, error)

	Save(map[string]interface{}) (string, error)

	Delete() (int, error)
	DeleteByID(id string) error

	Iterate(iteratorFunc func(record map[string]interface{}) bool) error

	Count() (int, error)
	Distinct(columnName string) ([]interface{}, error)

	SetupFilterGroup(groupName string, orSequence bool, parentGroup string) error
	RemoveFilterGroup(groupName string) error
	AddGroupFilter(groupName string, columnName string, operator string, value interface{}) error

	AddStaticFilter(columnName string, operator string, value interface{}) error
	AddFilter(columnName string, operator string, value interface{}) error

	ClearFilters() error

	AddSort(columnName string, Desc bool) error
	ClearSort() error

	SetResultColumns(columns ...string) error

	SetLimit(offset int, limit int) error

	ListColumns() map[string]string
	GetColumnType(columnName string) string
	HasColumn(columnName string) bool

	AddColumn(columnName string, columnType string, indexed bool) error
	RemoveColumn(columnName string) error
}
