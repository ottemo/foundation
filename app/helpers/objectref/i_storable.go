package objectref

import "github.com/ottemo/foundation/env"

// GetID returns current object id
func (it *DBObjectRef) GetID() string {
	return it.id
}

// SetID sets new id to current object
func (it *DBObjectRef) SetID(id string) {
	it.id = id
}

// Save stores current product to DB
func (it *DBObjectRef) Save() error {
	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d41d76c0-fa63-4418-b90e-bde64864ecaa", "not implemented")
}

// Load loads information from DB
func (it *DBObjectRef) Load(id string) error {
	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d65b0e21-c51d-4460-a32a-e4f09a84f421", "not implemented")
}

// Delete removes current object instance from DB
func (it *DBObjectRef) Delete() error {
	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "61c0b53d-b2a4-499a-be8b-6f2101f1ab97", "not implemented")
}
