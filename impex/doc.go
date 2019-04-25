// Copyright 2019 Ottemo. All rights reserved.

/*
Package impex represents import/export service. It provides a set of interfaces and helpers for extending it's functionality
as well as implements basic ideas of import/export functionality.

Importing and exporting happens through CSV files you should specify within API request in order to ue this service. CSV file
converts into map[string]interface{} objects which applies to model instances.

Refer to "Ottemo Import/Export Manual" for detail on service usage.
*/
package impex
