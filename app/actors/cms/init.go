// Package cms is just a grouping container for sub-packages auto init
package cms

import (
	// self-initiabilizable sub-package
	_ "github.com/ottemo/foundation/app/actors/cms/block"

	// self-initiabilizable sub-package
	_ "github.com/ottemo/foundation/app/actors/cms/page"

	// self-initiabilizable sub-package
	_ "github.com/ottemo/foundation/app/actors/cms/media"
)
