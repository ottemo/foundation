package impex

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
)

// setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("impex", "GET", "export/:model", restImpexExportModel)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("impex", "POST", "import/:model", restImpexImportModel)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("impex", "POST", "import", restImpexImport)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("impex", "GET", "tstImport", restImpexTstImport)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("impex", "GET", "tstExport", restImpexTstExport)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API used export specific model data from system
func restImpexExportModel(params *api.StructAPIHandlerParams) (interface{}, error) {

	model, err := models.GetModel(params.RequestURLParams["model"])
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	listable, isListable := model.(models.InterfaceListable)
	object, isObject := model.(models.InterfaceObject)

	if isListable && isObject {
		collection := listable.GetCollection()

		var attributes []string
		for _, attribute := range object.GetAttributesInfo() {
			attributes = append(attributes, attribute.Attribute)
			collection.ListAddExtraAttribute(attribute.Attribute)
		}

		// preparing csv writer
		csvWriter := csv.NewWriter(params.ResponseWriter)
		csvWriter.Comma = ','

		params.ResponseWriter.Header().Set("Content-type", "text/csv")
		params.ResponseWriter.Header().Set("Content-disposition", "attachment;filename=export_"+time.Now().Format(time.RFC3339)+".csv")

		// csvWriter.Write(attributes)
		// csvWriter.Flush()

		var records []map[string]interface{}

		list, _ := collection.List()
		for _, item := range list {
			records = append(records, item.Extra)

			/*record := make([]string, len(attributes))
			for idx, attribute := range attributes {
				record[idx] = utils.InterfaceToString(item.Extra[attribute])
			}
			csvWriter.Write(record)
			csvWriter.Flush()*/
		}

		MapToCSV(records, csvWriter)
	}

	return nil, nil
}

// WEB REST API used process scv file in impex format
func restImpexImport(params *api.StructAPIHandlerParams) (interface{}, error) {

	filesProcessed := 0
	for _, fileInfoArray := range params.Request.MultipartForm.File {
		for _, fileInfo := range fileInfoArray {
			// if utils.IsAmongStr(fileInfo.Header.Get("Content-Type"), "application/csv", "text/csv") {
			attachedFile, err := fileInfo.Open()
			defer attachedFile.Close()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			// preparing csv reader
			csvReader := csv.NewReader(attachedFile)
			csvReader.Comma = ','

			err = ImportCSV(csvReader)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			filesProcessed++
		}
	}

	return fmt.Sprintf("%d file(s) processed", filesProcessed), nil
}

// WEB REST API used import data to system
func restImpexImportModel(params *api.StructAPIHandlerParams) (interface{}, error) {

	modelName := params.RequestURLParams["model"]
	model, err := models.GetModel(modelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	object, isObject := model.(models.InterfaceObject)
	_, isStorable := model.(models.InterfaceStorable)
	if !isObject || !isStorable {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e4a8bb15704143caa29d1c053409efef", modelName+" not implements InterfaceObject or InterfaceStorable interface")
	}

	attributes := make(map[string]models.StructAttributeInfo)
	for _, attribute := range object.GetAttributesInfo() {
		attributes[attribute.Attribute] = attribute
	}

	// start reading csv
	csvFile, _, err := params.Request.FormFile("file")
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ','

	// reading header
	csvColumns, err := csvReader.Read()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for _, csvColumn := range csvColumns {
		if _, ok := attributes[csvColumn]; !ok {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "de6474369f7945d880f9139b52d1ba23", "there is no attribute "+csvColumn)
		}
	}

	for csvRecord, err := csvReader.Read(); err == nil; csvRecord, err = csvReader.Read() {
		model, _ = model.New()
		object, _ = model.(models.InterfaceObject)
		for idx, value := range csvRecord {
			object.Set(csvColumns[idx], value)
		}
		object.(models.InterfaceStorable).Save()
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "0a00d776415843a2845102f4344380d6", "not implemented")
}

// WEB REST API to test import csv file
func restImpexTstImport(params *api.StructAPIHandlerParams) (interface{}, error) {
	csvFile, err := os.OpenFile("test.csv", os.O_RDONLY, 0666)
	defer csvFile.Close()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	var result []map[string]interface{}
	processor := func(item map[string]interface{}) bool {
		result = append(result, item)
		return true
	}

	reader := csv.NewReader(csvFile)
	reader.Comma = ','
	err = CSVToMap(reader, processor)

	return result, err
}

// WEB REST API to test export
func restImpexTstExport(params *api.StructAPIHandlerParams) (interface{}, error) {

	params.ResponseWriter.Header().Set("Content-Type", "application/csv")

	var data []map[string]interface{}

	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for _, productItem := range productCollection.ListProducts() {
		data = append(data, productItem.ToHashMap())
	}

	csvWriter := csv.NewWriter(params.ResponseWriter)
	csvWriter.Comma = ','

	err = MapToCSV(data, csvWriter)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return []byte{}, nil
}
