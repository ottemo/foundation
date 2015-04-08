// Copyright 2014 The Ottemo Authors. All rights reserved.

/*
Package rest is a default implementation of a RESTful server. That package provides "InterfaceRestService" functionality
declared in "github.com/ottemo/foundation/api" package.

Package stands on a top of "github.com/julienschmidt/httprouter" package.

RESTful server is a HTTP protocol accessible service you can use to interact with application. Refer to "http://[url-base]/"
URL fo a list of API function application build providing. ("http://localhost:300" by default). The default interaction
protocol for API functions is "application/json" (but not limited to, some API returns raw data, others "plain text", some
of them supports couple content types).

Ottemo API calls are supplied with ApplicationContext abstraction (refer "InterfaceApplicationContext"), to transfer API
call related "attributes" and "content". For REST server attributes are parameters specified in URL string whereas content
is a HTTP Request content. Notice that URL arguments are both - REST resource required parameters and optional URL parameters.

	Example:
	  endpoint: /category/:categoryID/products
	  api call: http://localhost/category/5488485b49c43d4283000067/products?action=count,sku=~10

	So, the arguments provided to API handler function are:
	  categoryID="5488485b49c43d4283000067", action="count", sku="~10" (with string values)

Session specification addressed to "OTTEMOSESSION=[sessionID]" COOKIE value. Each request with unspecified session will
be supplied with new one session. SessionID will be returned in mentioned COOKIE value.
*/
package rest
