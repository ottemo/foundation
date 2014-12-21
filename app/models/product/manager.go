package product

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	registeredStock InterfaceStock
)

// UnRegisterStock removes stock management from system
func UnRegisterStock() error {
	registeredStock = nil
	return nil
}

// RegisterStock registers given stock manager in system
func RegisterStock(stock InterfaceStock) error {
	if registeredStock != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "80dc18f2da63418e84305a832a4c3bd2", "Already registered")
	}
	registeredStock = stock

	return nil
}

// GetRegisteredStock returns currently used stack manager or nil
func GetRegisteredStock() InterfaceStock {
	return registeredStock
}
