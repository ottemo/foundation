// Copyright 2019 Ottemo. All rights reserved.

/*
Package api represents API services abstraction. It provides a set of interfaces and helpers for a packages to extend
application with new API endpoints as well as to interact with other components of application through them.

As application components (actor packages) should not interact between them directly them should use this package to make
indirect calls between the. Interaction with a session manager also should happen through this package.
*/
package api
