package coupon

import (
	"encoding/csv"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("coupons", List)
	service.POST("coupons", Create)
	service.GET("csv/coupons", DownloadCSV)
	service.POST("csv/coupons", UploadCSV)
	service.POST("cart/coupons", Apply)
	service.DELETE("cart/coupons/:code", Revert)
	service.GET("coupons/:id", GetByID)
	service.PUT("coupons/:id", UpdateByID)
	service.DELETE("coupons/:id", DeleteByID)

	return nil
}

// List returns a list registered coupons
func List(context api.InterfaceApplicationContext) (interface{}, error) {

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

// Create will generate a new coupon code when supplied the following required keys,
// they are not required to match.
//   * "name" is the desired reference key for the coupon
//   * "code" is the text visitors must enter to apply a coupon in checkout
func Create(context api.InterfaceApplicationContext) (interface{}, error) {

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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "842d3ba9-3354-4470-a85f-cbaf909c3827", "Required fields, 'code' and 'name', cannot be blank.")
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

	valueTarget := checkout.ConstDiscountObjectCart
	if targetValue, present := postValues["target"]; present {
		target := strings.ToLower(utils.InterfaceToString(targetValue))
		if target != "" && !strings.Contains(target, checkout.ConstDiscountObjectCart) {
			valueTarget = target
		}
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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "34cb6cfe-fba3-4c1f-afc5-1ff7266a9a86", "A Discount with the provided code: '"+valueCode+"', already exists.")
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
		"target":  valueTarget,
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

// Apply will coupon code to the current checkout
//   - coupon code should be specified in "coupon" argument
func Apply(context api.InterfaceApplicationContext) (interface{}, error) {

	var couponCode string
	var present bool

	// check request context
	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// validate presence of code in post
	if _, present = postValues["code"]; present {
		couponCode = utils.InterfaceToString(postValues["code"])
	} else {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "085b8e25-7939-4b94-93f1-1007ada357d4", "Required key 'code' cannot have a  blank value.")
	}

	currentSession := context.GetSession()

	// get applied coupons array for current session
	appliedCoupons := utils.InterfaceToStringArray(currentSession.Get(ConstSessionKeyAppliedDiscountCodes))

	// check if coupon has already been applied
	if utils.IsInArray(couponCode, appliedCoupons) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "29c4c963-0940-4780-8ad2-9ed5ca7c97ff", "Coupon code, "+couponCode+" has already been applied.")
	}

	// load coupon for specified code
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

	// verify and apply obtained coupon
	if len(records) > 0 {
		discountCoupon := records[0]

		applyTimes := utils.InterfaceToInt(discountCoupon["times"])
		couponStart := utils.InterfaceToTime(discountCoupon["since"])
		couponEnd := utils.InterfaceToTime(discountCoupon["until"])

		currentTime := time.Now()

		// to be applicable, the coupon should satisfy following conditions:
		//   [applyTimes] should be -1 or >0 and [couponStart] >= currentTime <= [couponEnd] if set
		if (applyTimes == -1 || applyTimes > 0) &&
			(utils.IsZeroTime(couponStart) || couponStart.Unix() <= currentTime.Unix()) &&
			(utils.IsZeroTime(couponEnd) || couponEnd.Unix() >= currentTime.Unix()) {

			// TODO: applied coupons are lost with session clear, probably should be made on order creation,
			// or add an event handler to add to session # of times used
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
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "63442858-bd71-4f10-855a-b5975fc2dd16", "Coupon code, "+strings.ToUpper(couponCode)+", cannot be applied, exceeded usage limits.")
		}
	} else {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b2934505-06e9-4250-bb98-c22e4918799e", "Coupon code, "+strings.ToUpper(couponCode)+", is not a valid coupon code.")
	}

	return "ok", nil
}

// Revert will remove the coupon code and its value from the current checkout
//   * "coupon" key refers to the coupon code
//   * use a "*" as the coupon code to revert all discounts
func Revert(context api.InterfaceApplicationContext) (interface{}, error) {

	couponCode := context.GetRequestArgument("code")

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

// DownloadCSV returns a csv file with the current coupons and their configuration
//   * returns a csv file
func DownloadCSV(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// preparing csv writer
	csvWriter := csv.NewWriter(context.GetResponseWriter())

	context.SetResponseContentType("text/csv")
	context.SetResponseSetting("Content-disposition", "attachment;filename=discount_coupons.csv")

	csvWriter.Write([]string{"Code", "Name", "Amount", "Percent", "Times", "Since", "Until", "Limits", "Target"})
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
			utils.InterfaceToString(record["limits"]),
			utils.InterfaceToString(record["target"]),
		})

		csvWriter.Flush()
		return true
	})

	return nil, nil
}

// UploadCSV will overwrite and replace the current coupon configuration with the uploaded CSV
//   NOTE: the csv file should be provided in a "file" field when sent as a multipart form
func UploadCSV(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	csvFile := context.GetRequestFile("file")
	if csvFile == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3398f40a-726b-48ad-9f29-9dd390b7e952", "A file name must be specified.")
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
			record["limits"] = utils.InterfaceToMap(csvRecord[7])
			record["target"] = utils.InterfaceToString(csvRecord[8])

			collection.Save(record)
		}
	}

	return "ok", nil
}

// GetByID returns a coupon with the specified ID
// * coupon id should be specified in the "id" argument
func GetByID(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	id := context.GetRequestArgument("id")
	records, err := collection.LoadByID(id)

	return records, env.ErrorDispatch(err)
}

// UpdateByID updates existing coupon specified in the request argument
//   * coupon id should be specified in "couponID" argument
func UpdateByID(context api.InterfaceApplicationContext) (interface{}, error) {

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

	couponID := context.GetRequestArgument("id")
	record, err := collection.LoadByID(couponID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// if discount 'code' was changed - checking new value for duplicates
	if codeValue, present := postValues["code"]; present && codeValue != record["code"] {
		codeValue := utils.InterfaceToString(codeValue)

		collection.AddFilter("code", "=", codeValue)
		recordsNumber, err := collection.Count()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		if recordsNumber > 0 {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e49e5e01-4f6f-4ff0-bd28-dfb616308aa7", "A Discount with the provided code: '"+codeValue+"', already exists.")
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

	record["target"] = checkout.ConstDiscountObjectCart
	if targetValue, present := postValues["target"]; present {
		target := strings.ToLower(utils.InterfaceToString(targetValue))
		if target != "" && !strings.Contains(target, checkout.ConstDiscountObjectCart) {
			record["target"] = target
		}
	}

	if value, present := postValues["until"]; present {
		record["until"] = utils.InterfaceToTime(value)
	}

	if value, present := postValues["since"]; present {
		record["since"] = utils.InterfaceToTime(value)
	}

	// saving updates
	//---------------
	_, err = collection.Save(record)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return record, nil
}

// DeleteByID deletes specified SEO item
//   * discount id should be specified in the "couponID" argument
func DeleteByID(context api.InterfaceApplicationContext) (interface{}, error) {
	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(context.GetRequestArgument("id"))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}
