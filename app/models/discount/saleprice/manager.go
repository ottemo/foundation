package saleprice

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	registeredSalePrice InterfaceSalePrice
)

// UnRegisterSalePrice removes sale price manager from system
func UnRegisterSalePrice() error {
	registeredSalePrice = nil
	return nil
}

// RegisterSalePrice registers given sale price manager in system
func RegisterSalePrice(salePrice InterfaceSalePrice) error {
	if registeredSalePrice != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f6485d14-7af9-4451-86e9-2543d2620a52", "already registered")
	}
	registeredSalePrice = salePrice

	return nil
}

// GetRegisteredSalePrice returns currently used sale price manager or nil
func GetRegisteredSalePrice() InterfaceSalePrice {
	return registeredSalePrice
}
