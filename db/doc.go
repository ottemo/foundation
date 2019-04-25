// Copyright 2019 Ottemo. All rights reserved.

/*
Package db represents database services abstraction layer. It provides a set of interfaces and helpers to interact with
database storage services.

So, when package decides to store something in database it should work with this package instead of trying to interact
with concrete database engine (as Ottemo supposing ability to work with different databases).

Providing Ottemo with a new database engine supposes implementation of "InterfaceDBEngine" with following registration
for db package.

For database type specification you should refer this package for possible types.

	Example:
	--------
	collection, err := db.GetCollection( myCollectionName )
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("customer_id", db.ConstTypeID, false)
	collection.AddColumn("customer_email", db.TypeWPrecision(db.ConstTypeVarchar, 100), true)
	collection.AddColumn("bonus_code", db.ConstTypeInteger, false)
	collection.AddColumn("bonus_amount", db.ConstTypeInteger, false)

*/
package db
