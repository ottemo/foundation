package discount

import (
	"encoding/csv"

	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("discount", "GET", "apply/:code", restDiscountApply)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("discount", "GET", "neglect/:code", restDiscountNeglect)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("discount", "GET", "download/csv", restDiscountCSVDownload)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("discount", "POST", "upload/csv", restDiscountCSVUpload)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function to apply discount code to current checkout
func restDiscountApply(params *api.StructAPIHandlerParams) (interface{}, error) {

	couponCode := params.RequestURLParams["code"]

	// getting applied coupons array for current session
	var appliedCoupons []string
	if sessionValue, ok := params.Session.Get(ConstSessionKeyAppliedDiscountCodes).([]string); ok {
		appliedCoupons = sessionValue
	}

	// checking if coupon was already applied
	if utils.IsInArray(couponCode, appliedCoupons) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "29c4c963-0940-4780-8ad2-9ed5ca7c97ff", "coupon code already applied")
	}

	// loading coupon for specified code
	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	err = collection.AddFilter("code", "=", couponCode)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	records, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// checking and applying obtained coupon
	if len(records) > 0 {
		discountCoupon := records[0]

		applyTimes := utils.InterfaceToInt(discountCoupon["times"])
		workSince := utils.InterfaceToTime(discountCoupon["since"])
		workUntil := utils.InterfaceToTime(discountCoupon["until"])

		currentTime := time.Now()

		// to be applicable coupon should satisfy following conditions:
		//   [applyTimes] should be -1 or >0 and [workSince] >= currentTime <= [workUntil] if set
		if (applyTimes == -1 || applyTimes > 0) &&
			(utils.IsZeroTime(workSince) || workSince.Unix() <= currentTime.Unix()) &&
			(utils.IsZeroTime(workUntil) || workUntil.Unix() >= currentTime.Unix()) {

			// TODO: coupon loosing with session clear, probably should be made on order creation, or have event on session
			// times used decrease
			if applyTimes > 0 {
				discountCoupon["times"] = applyTimes - 1
				_, err := collection.Save(discountCoupon)
				if err != nil {
					return nil, env.ErrorDispatch(err)
				}
			}

			// coupon is working - applying it
			appliedCoupons = append(appliedCoupons, couponCode)
			params.Session.Set(ConstSessionKeyAppliedDiscountCodes, appliedCoupons)

		} else {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "63442858-bd71-4f10-855a-b5975fc2dd16", "coupon is not applicable")
		}
	} else {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b2934505-06e9-4250-bb98-c22e4918799e", "coupon code not found")
	}

	return "ok", nil
}

// WEB REST API function to neglect(un-apply) discount code to current checkout
//   - use "*" as code to neglect all discounts
func restDiscountNeglect(params *api.StructAPIHandlerParams) (interface{}, error) {

	couponCode := params.RequestURLParams["code"]

	if couponCode == "*" {
		params.Session.Set(ConstSessionKeyAppliedDiscountCodes, make([]string, 0))
		return "ok", nil
	}

	if appliedCoupons, ok := params.Session.Get(ConstSessionKeyAppliedDiscountCodes).([]string); ok {
		var newAppliedCoupons []string
		for _, value := range appliedCoupons {
			if value != couponCode {
				newAppliedCoupons = append(newAppliedCoupons, value)
			}
		}
		params.Session.Set(ConstSessionKeyAppliedDiscountCodes, newAppliedCoupons)

		// times used increase
		collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		err = collection.AddFilter("code", "=", couponCode)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		records, err := collection.Load()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		if len(records) > 0 {
			applyTimes := utils.InterfaceToInt(records[0]["times"])
			if applyTimes >= 0 {
				records[0]["times"] = applyTimes + 1

				_, err := collection.Save(records[0])
				if err != nil {
					return nil, env.ErrorDispatch(err)
				}
			}
		}
	}

	return "ok", nil
}

// WEB REST API function to download current tax rates in CSV format
func restDiscountCSVDownload(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// preparing csv writer
	csvWriter := csv.NewWriter(params.ResponseWriter)

	params.ResponseWriter.Header().Set("Content-type", "text/csv")
	params.ResponseWriter.Header().Set("Content-disposition", "attachment;filename=discount_coupons.csv")

	csvWriter.Write([]string{"Code", "Name", "Amount", "Percent", "Times", "Since", "Until"})
	csvWriter.Flush()

	// loading records from DB and writing them in csv format
	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = collection.Iterate(func(record map[string]interface{}) bool {
		csvWriter.Write([]string{
			utils.InterfaceToString(record["code"]),
			utils.InterfaceToString(record["name"]),
			utils.InterfaceToString(record["amount"]),
			utils.InterfaceToString(record["percent"]),
			utils.InterfaceToString(record["times"]),
			utils.InterfaceToString(record["since"]),
			utils.InterfaceToString(record["until"]),
		})

		csvWriter.Flush()
		return true
	})

	return nil, nil
}

// WEB REST API function to upload tax rates into CSV
func restDiscountCSVUpload(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	csvFile, _, err := params.Request.FormFile("file")
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ','

	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	collection.Delete()

	csvReader.Read() //skipping header
	for csvRecord, err := csvReader.Read(); err == nil; csvRecord, err = csvReader.Read() {
		if len(csvRecord) >= 7 {
			record := make(map[string]interface{})

			code := utils.InterfaceToString(csvRecord[0])
			name := utils.InterfaceToString(csvRecord[1])
			if code == "" || name == "" {
				continue
			}

			times := utils.InterfaceToInt(csvRecord[4])
			if csvRecord[4] == "" {
				times = -1
			}

			record["code"] = code
			record["name"] = name
			record["amount"] = utils.InterfaceToFloat64(csvRecord[2])
			record["percent"] = utils.InterfaceToFloat64(csvRecord[3])
			record["times"] = times
			record["since"] = utils.InterfaceToTime(csvRecord[5])
			record["until"] = utils.InterfaceToTime(csvRecord[6])

			collection.Save(record)
		}
	}

	return "ok", nil
}
