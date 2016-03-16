package shipstation

import (
	"github.com/ottemo/foundation/env"
)

const (
	ConstErrorModule = "shipstation"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// struct goes here
type Order struct {
	Name string `xml:"name"`
}

type Orders struct {
	Orders []Order `xml:"Order"`
}
