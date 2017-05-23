package product

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/stock"
)

// Package global variables
var (
	registeredStock stock.InterfaceStock
)

// UnRegisterStock removes stock management from system
func UnRegisterStock() error {
	registeredStock = nil
	return nil
}

// RegisterStock registers given stock manager in system
func RegisterStock(stock stock.InterfaceStock) error {
	if registeredStock != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "80dc18f2-da63-418e-8430-5a832a4c3bd2", "Already registered")
	}
	registeredStock = stock

	return nil
}

// GetRegisteredStock returns currently used stack manager or nil
func GetRegisteredStock() stock.InterfaceStock {
	return registeredStock
}
