// Copyright 2014 The Ottemo Authors. All rights reserved.

/*

Package app represents Ottemo application object.

That package contains routines which allows other components to register callbacks on application start/end, API
functions for administrator login, system configuration values, etc. So, this package contains the code addressed to
application instance. Ottemo packages should address this package to interact with running application instance but not
to "github.com/ottemo/foundaton" package".

*/
package app
