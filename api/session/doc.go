// Copyright 2014 The Ottemo Authors. All rights reserved.

/*
Package session is a default implementation of a application Session manager. That package provides "InterfaceSessionService"
functionality declared in "github.com/ottemo/foundation/api" package.

Sessions are API call related storage. So, each call to API function supplied with own separated storage which can holds
a values related to that particular action. By default sessions have a lifetime, within that period application routines
can hold information fo future usage for either themselves or other API calls. In order to use previously created session
API call should specify sessionID within application context.
*/
package session
