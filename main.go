package main

import (
	"fmt"

	app "github.com/ottemo/foundation/app"

	_ "github.com/ottemo/foundation/config"

	//_ "github.com/ottemo/foundation/database/sqlite"
	_ "github.com/ottemo/foundation/database/mongodb"

	_ "github.com/ottemo/foundation/rest_service"
	_ "github.com/ottemo/foundation/rest_service/negroni"

	_ "github.com/ottemo/foundation/models"
	_ "github.com/ottemo/foundation/models/product/defaultproduct"
	_ "github.com/ottemo/foundation/models/visitor/default_address"
	_ "github.com/ottemo/foundation/models/visitor/default_visitor"
)

func main() {
	if err := app.Start(); err != nil {
		fmt.Println(err.Error())
	}

	app.Serve()

}
