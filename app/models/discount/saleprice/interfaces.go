package saleprice

import (
	"time"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstErrorModule = "saleprice"
	ConstErrorLevel  = env.ConstErrorLevelModel
)

// InterfaceCategory represents interface to access business layer implementation of category object
type InterfaceSalePrice interface {
	GetAmount() float64
	SetAmount(float64) error

	GetEndDatetime() time.Time
	SetEndDatetime(time.Time) error

	GetProductID() string
	SetProductID(string) error

	GetStartDatetime() time.Time
	SetStartDatetime(time.Time) error

	models.InterfaceObject
	models.InterfaceStorable
}
