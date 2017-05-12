package model

import "github.com/ottemo/foundation/app/models"

type InterfaceProductTic interface {
	GetTicID() int
	SetTicID(int) error

	GetProductID() string
	SetProductID(string) error

	models.InterfaceObject
	models.InterfaceStorable
}

