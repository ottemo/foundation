package main

import (
	"fmt"

	config "github.com/ottemo/foundation/config"

	app "github.com/ottemo/foundation/app"

	//_ "github.com/ottemo/foundation/database/sqlite"
	_ "github.com/ottemo/foundation/database/mongodb"

	_ "github.com/ottemo/foundation/models"
	_ "github.com/ottemo/foundation/models/product/defaultproduct"
	_ "github.com/ottemo/foundation/models/visitor/default_address"
	_ "github.com/ottemo/foundation/models/visitor/default_visitor"
)

func main() {
	iniConfig := config.NewDefaultIniConfig()
	app.OnAppStart(iniConfig.Startup)
	config.RegisterIniConfig(iniConfig)

	if err := app.Start(); err != nil {
		fmt.Println(err.Error())
	}

	app.Serve()

}
