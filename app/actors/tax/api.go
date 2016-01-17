package tax

import (
	"encoding/csv"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("taxes/csv", api.GET, APIDownloadTaxCSV)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("taxes/csv", api.POST, APIUploadTaxCSV)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIDownloadTaxCSV returns csv file with currently used tax rates
//   - returns not a JSON, but csv file
func APIDownloadTaxCSV(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	csvWriter := csv.NewWriter(context.GetResponseWriter())

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
			records, err := collection.Load()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			context.SetResponseContentType("text/csv")
			context.SetResponseSetting("Content-disposition", "attachment;filename=tax_rates.csv")

			csvWriter.Write([]string{"Code", "Country", "State", "Zip", "Rate"})
			csvWriter.Flush()

			for _, record := range records {
				csvWriter.Write([]string{
					utils.InterfaceToString(record["code"]),
					utils.InterfaceToString(record["country"]),
					utils.InterfaceToString(record["state"]),
					utils.InterfaceToString(record["zip"]),
					utils.InterfaceToString(record["rate"])})

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

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	csvFile := context.GetRequestFile("file")
	if csvFile == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "033b69a0-d33d-4bfe-b670-b469d3e86f90", "file unspecified")
	}

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ','

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
			collection.Delete()

			csvReader.Read() //skipping header
			for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
				if len(record) >= 5 {
					taxRecord := make(map[string]interface{})

					taxRecord["code"] = record[0]
					taxRecord["country"] = record[1]
					taxRecord["state"] = record[2]
					taxRecord["zip"] = record[3]
					taxRecord["rate"] = utils.InterfaceToFloat64(record[4])

					collection.Save(taxRecord)
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
