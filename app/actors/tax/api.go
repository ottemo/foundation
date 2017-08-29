package tax

import (
	"encoding/csv"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/impex"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("taxes/csv", api.IsAdminHandler(APIDownloadTaxCSV))
	service.POST("taxes/csv", api.IsAdminHandler(
		impex.ImportStartHandler(
			api.AsyncHandler(APIUploadTaxCSV, impex.ImportResultHandler))))

	return nil
}

// APIDownloadTaxCSV returns csv file with currently used tax rates
//   - returns not a JSON, but csv file
func APIDownloadTaxCSV(context api.InterfaceApplicationContext) (interface{}, error) {

	csvWriter := csv.NewWriter(context.GetResponseWriter())

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
			records, err := collection.Load()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			if err := context.SetResponseContentType("text/csv"); err != nil {
				_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "68abd3a8-349d-4a4a-881d-6b09892d80c7", err.Error())
			}
			if err := context.SetResponseSetting("Content-disposition", "attachment;filename=tax_rates.csv"); err != nil {
				_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5569f985-877b-4a7f-929d-ee2ec3c00e62", err.Error())
			}

			if err := csvWriter.Write([]string{"Code", "Country", "State", "Zip", "Rate"}); err != nil {
				_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "54cfe05e-7dc6-433e-b566-6b8ed3712a68", err.Error())
			}
			csvWriter.Flush()

			for _, record := range records {
				if err := csvWriter.Write([]string{
					utils.InterfaceToString(record["code"]),
					utils.InterfaceToString(record["country"]),
					utils.InterfaceToString(record["state"]),
					utils.InterfaceToString(record["zip"]),
					utils.InterfaceToString(record["rate"])}); err != nil {
					_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2aa78843-c751-47d6-8f13-669455a2ece1", err.Error())
				}

				csvWriter.Flush()
			}

			return nil, nil
		}
	} else {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4cc9eb51-663e-4ebc-9b23-5e7dee56e078", "can't get DB engine")
	}

	return nil, nil
}

// APIUploadTaxCSV replaces currently used discount coupons with data from provided in csv file
//   - csv file should be provided in "file" field
func APIUploadTaxCSV(context api.InterfaceApplicationContext) (interface{}, error) {

	csvFileName := context.GetRequestArgument("file")
	if csvFileName == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6995fa26-f738-4408-99a3-515c519f7a0f", "A file name must be specified.")
	}

	csvFile := context.GetRequestFile(csvFileName)
	if csvFile == nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c0f00a48-17fc-4ed9-b7e8-f38cc097316c", "A file must be specified.")
	}

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ','

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
			if _, err := collection.Delete(); err != nil {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "683e97e0-a35a-4b99-834b-95bc999dcf2a", err.Error())
			}

			if _, err := csvReader.Read(); err != nil { //skipping header
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "7a3b3349-504f-4635-8e7d-0ed41dc4c05e", err.Error())
			}
			for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
				if len(record) >= 5 {
					taxRecord := make(map[string]interface{})

					taxRecord["code"] = record[0]
					taxRecord["country"] = record[1]
					taxRecord["state"] = record[2]
					taxRecord["zip"] = record[3]
					taxRecord["rate"] = utils.InterfaceToFloat64(record[4])

					if _, err := collection.Save(taxRecord); err != nil {
						return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "90c64ccd-aae5-465d-a721-409cd4197137", err.Error())
					}
				}
			}
		} else {
			return nil, env.ErrorDispatch(err)
		}
	} else {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cc1f5af3-fd1c-4da8-a4c2-f40613eb682f", "can't get DB engine")
	}

	return "ok", nil
}
