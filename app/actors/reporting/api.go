package reporting

import (
	"github.com/ottemo/foundation/api"
	"time"
	// "github.com/ottemo/foundation/app"
	// "github.com/ottemo/foundation/app/models"
	// "github.com/ottemo/foundation/app/models/order"
	// "github.com/ottemo/foundation/app/models/subscription"
	// "github.com/ottemo/foundation/app/models/visitor"
	// "github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("reporting/product-performance", listProductPerformance)
	return nil
}

func listProductPerformance(context api.InterfaceApplicationContext) (interface{}, error) {

	startDateS := context.GetRequestArgument("start_date")
	endDateS := context.GetRequestArgument("end_date")
	if startDateS == "" || endDateS == "" {
		//todo: err
	}
	endDate := utils.InterfaceToTime(endDateS)
	startDate := utils.InterfaceToTime(startDateS)

	// Debugging
	endDate = time.Now()
	startDate = time.Now()

	some := []time.Time{
		startDate,
		endDate,
	}
	// oiModel, _ := order.GetOrderItemCollectionModel()
	// oiModel.ListFilterAdd(attribute, operator, value)
	// some, _ := oiModel.List()

	return some, nil
}
