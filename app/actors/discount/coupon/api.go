package coupon

import (
	"encoding/csv"

	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("discount/:coupon/apply", api.ConstRESTOperationGet, APIApplyDiscount)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("discount/:coupon/neglect", api.ConstRESTOperationGet, APINeglectDiscount)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("discounts/csv", api.ConstRESTOperationGet, APIDownloadDiscountCSV)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("discounts/csv", api.ConstRESTOperationCreate, APIUploadDiscountCSV)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("coupons", api.ConstRESTOperationCreate, APICreateDiscount)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("coupons/:couponID", api.ConstRESTOperationGet, APIGetDiscount)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("coupons/:couponID", api.ConstRESTOperationUpdate, APIUpdateDiscount)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("coupons/:couponID", api.ConstRESTOperationDelete, APIDeleteDiscount)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("coupons", api.ConstRESTOperationGet, APIListDiscounts)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIApplyDiscount applies discount code promotion to current checkout
//   - coupon code should be specified in "coupon" argument
func APIApplyDiscount(context api.InterfaceApplicationContext) (interface{}, error) {

	couponCode := context.GetRequestArgument("coupon")

	currentSession := context.GetSession()

	// getting applied coupons array for current session
	appliedCoupons := utils.InterfaceToStringArray(currentSession.Get(ConstSessionKeyAppliedDiscountCodes))
	usedCodes := utils.InterfaceToStringArray(currentSession.Get(ConstSessionKeyUsedDiscountCodes))

	// checking if coupon was already applied
	if utils.IsInArray(couponCode, appliedCoupons) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "29c4c963-0940-4780-8ad2-9ed5ca7c97ff", "coupon code already applied")
	}

	if utils.IsInArray(couponCode, usedCodes) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "32315c6c-c932-4ad4-a1a1-5eaf86f1dcdc", "coupon code already used")
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
			currentSession.Set(ConstSessionKeyAppliedDiscountCodes, appliedCoupons)

		} else {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "63442858-bd71-4f10-855a-b5975fc2dd16", "coupon is not applicable")
		}
	} else {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b2934505-06e9-4250-bb98-c22e4918799e", "coupon code not found")
	}

	return "ok", nil
}

// APINeglectDiscount neglects (un-apply) discount code promotion to current checkout
//   - coupon code should be specified in "coupon" argument
//   - use "*" as coupon code to neglect all discounts
func APINeglectDiscount(context api.InterfaceApplicationContext) (interface{}, error) {

	couponCode := context.GetRequestArgument("coupon")

	if couponCode == "*" {
		context.GetSession().Set(ConstSessionKeyAppliedDiscountCodes, make([]string, 0))
		return "ok", nil
	}

	appliedCoupons := utils.InterfaceToStringArray(context.GetSession().Get(ConstSessionKeyAppliedDiscountCodes))
	if len(appliedCoupons) > 0 {
		var newAppliedCoupons []string
		for _, value := range appliedCoupons {
			if value != couponCode {
				newAppliedCoupons = append(newAppliedCoupons, value)
			}
		}
		context.GetSession().Set(ConstSessionKeyAppliedDiscountCodes, newAppliedCoupons)

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

// APIDownloadDiscountCSV returns csv file with currently used discount coupons
//   - returns not a JSON, but csv file
func APIDownloadDiscountCSV(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// preparing csv writer
	csvWriter := csv.NewWriter(context.GetResponseWriter())

	context.SetResponseContentType("text/csv")
	context.SetResponseSetting("Content-disposition", "attachment;filename=discount_coupons.csv")

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

// APIUploadDiscountCSV replaces currently used discount coupons with data from provided in csv file
//   - csv file should be provided in "file" field
func APIUploadDiscountCSV(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	csvFile := context.GetRequestFile("file")
	if csvFile == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3398f40a-726b-48ad-9f29-9dd390b7e952", "file unspecified")
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

// APIListDiscounts returns a list registered discounts
func APIListDiscounts(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	records, err := collection.Load()

	return records, env.ErrorDispatch(err)
}

// APIGetDiscount - returns discount item for a specified id
func APIGetDiscount(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	id := context.GetRequestArgument("couponID")
	records, err := collection.LoadByID(id)

	return records, env.ErrorDispatch(err)
}

// APICreateDiscount - creates new discount item
//   - "code" and "name" attributes are required
func APICreateDiscount(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// checking request context
	//------------------------
	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(postValues, "code", "name") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "842d3ba9-3354-4470-a85f-cbaf909c3827", "'code' or 'name' value is not specified")
	}

	valueCode := utils.InterfaceToString(postValues["code"])
	valueName := utils.InterfaceToString(postValues["name"])

	timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))

	valueUntil := time.Now()
	if value, present := postValues["until"]; present {
		valueUntil, _ = utils.MakeUTCTime(utils.InterfaceToTime(value), timeZone)
	}

	valueSince := time.Now()
	if value, present := postValues["since"]; present {
		valueSince, _ = utils.MakeUTCTime(utils.InterfaceToTime(value), timeZone)
	}

	valueLimits := make(map[string]interface{})
	if value, present := postValues["limits"]; present {
		valueLimits = utils.InterfaceToMap(value)
	}

	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection.AddFilter("code", "=", valueCode)
	recordsNumber, err := collection.Count()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	if recordsNumber > 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "34cb6cfe-fba3-4c1f-afc5-1ff7266a9a86", "discount with such code: '"+valueCode+"', already exists")
	}

	// making new record and storing it
	//---------------------------------
	newRecord := map[string]interface{}{
		"code":    valueCode,
		"name":    valueName,
		"amount":  0,
		"percent": 0,
		"times":   -1,
		"since":   valueSince,
		"until":   valueUntil,
		"limits":  valueLimits,
	}

	attributes := []string{"amount", "percent", "times"}
	for _, attribute := range attributes {
		if value, present := postValues[attribute]; present {
			newRecord[attribute] = value
		}
	}

	newID, err := collection.Save(newRecord)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	newRecord["_id"] = newID

	return newRecord, nil
}

// APIUpdateDiscount updates existing discount
//   - discount id should be specified in "couponID" argument
func APIUpdateDiscount(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	couponID := context.GetRequestArgument("couponID")
	record, err := collection.LoadByID(couponID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// if discount 'code' was changed - checking new value for duplicates
	//-----------------------------------------------------------------
	if codeValue, present := postValues["code"]; present && codeValue != record["code"] {
		codeValue := utils.InterfaceToString(codeValue)

		collection.AddFilter("code", "=", codeValue)
		recordsNumber, err := collection.Count()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		if recordsNumber > 0 {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e49e5e01-4f6f-4ff0-bd28-dfb616308aa7", "discount with such code: '"+codeValue+"', already exists")
		}

		record["code"] = codeValue
	}

	// updating other attributes
	//--------------------------
	attributes := []string{"amount", "percent", "times", "limits"}
	for _, attribute := range attributes {
		if value, present := postValues[attribute]; present {
			record[attribute] = value
		}
	}

	timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))
	if value, present := postValues["until"]; present {
		record["until"], _ = utils.MakeUTCTime(utils.InterfaceToTime(value), timeZone)
	}

	if value, present := postValues["since"]; present {
		record["since"], _ = utils.MakeUTCTime(utils.InterfaceToTime(value), timeZone)
	}

	// saving updates
	//---------------
	_, err = collection.Save(record)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return record, nil
}

// APIDeleteDiscount deletes specified SEO item
//   - discount id should be specified in "couponID" argument
func APIDeleteDiscount(context api.InterfaceApplicationContext) (interface{}, error) {
	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(context.GetRequestArgument("couponID"))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}
