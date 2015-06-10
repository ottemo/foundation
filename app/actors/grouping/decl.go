// Package grouping implements products set grouping into another set
package grouping

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstErrorModule = "grouping"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstGroupingConfigPath = "general.stock.grouprules"
)

// Package global variables
var (
	currentRules = make([]interface{}, 0)
)
