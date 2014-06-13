package main

import (
	"fmt"

	app "github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/models"

	_ "github.com/ottemo/foundation/config/default_config"
	_ "github.com/ottemo/foundation/config/default_ini_config"

	//_ "github.com/ottemo/foundation/database/sqlite"
	_ "github.com/ottemo/foundation/database/mongodb"

	_ "github.com/ottemo/foundation/rest_service"
	_ "github.com/ottemo/foundation/rest_service/negroni"

	_ "github.com/ottemo/foundation/models/custom_attributes"
	_ "github.com/ottemo/foundation/models/product/default_product"
	_ "github.com/ottemo/foundation/models/visitor/default_address"
	_ "github.com/ottemo/foundation/models/visitor/default_visitor"
)

func main() {
	if err := app.Start(); err != nil {
		fmt.Println(err.Error())
	}

	app.Serve()

	// CreateNewProductAttribute("x")
	// CreateNewProductAttribute("y")
}

// CreateNewProductAttribute allows for the creation of additional product attributes at runtime.
func CreateNewProductAttribute(attrName string) {
	model, err := models.GetModel("Product")
	if err != nil {
		fmt.Println("Product model not found: " + err.Error())
	}

	attribute := models.T_AttributeInfo{
		Model:      "product",
		Collection: "product",
		Attribute:  attrName,
		Type:       "text",
		Label:      "Test Attribute",
		Group:      "General",
		Editors:    "text",
		Options:    "",
		Default:    "",
	}

	if prod, ok := model.(models.I_CustomAttributes); ok {
		if err := prod.AddNewAttribute(attribute); err != nil {
			fmt.Println("Product new attribute error: " + err.Error())
		}
	} else {
		fmt.Println("product model is not I_CustomAttributes")
	}
}
