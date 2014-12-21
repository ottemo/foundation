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

	err = api.GetRestService().RegisterAPI("tax", "GET", "download/csv", restTaxCSVDownload)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("tax", "POST", "upload/csv", restTaxCSVUpload)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function to download current tax rates in CSV format
func restTaxCSVDownload(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	csvWriter := csv.NewWriter(params.ResponseWriter)

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
			records, err := collection.Load()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			err = csvWriter.Write([]string{"Code", "Country", "State", "Zip", "Rate"})
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			params.ResponseWriter.Header().Set("Content-type", "text/csv")
			params.ResponseWriter.Header().Set("Content-disposition", "attachment;filename=tax_rates.csv")

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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4cc9eb51663e4ebc9b235e7dee56e078", "can't get DB engine")
	}

	return nil, nil
}

// WEB REST API function to upload tax rates into CSV
func restTaxCSVUpload(params *api.StructAPIHandlerParams) (interface{}, error) {

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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cc1f5af3fd1c4da8a4c2f40613eb682f", "can't get DB engine")
	}

	return "ok", nil
}
