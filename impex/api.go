package impex

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("impex/models", restImpexListModels)
	service.GET("impex/import/status", restImpexImportStatus)
	service.GET("impex/export/:model", restImpexExportModel)
	service.POST("impex/import/:model", restImpexImportModel)
	service.POST("impex/import", restImpexImport)

	service.POST("impex/test/import", restImpexTestImport)
	service.POST("impex/test/mapping", restImpexTestCsvToMap)

	return nil
}

// WEB REST API used to list available models for Impex system
func restImpexListModels(context api.InterfaceApplicationContext) (interface{}, error) {
	var result []string

	for modelName, modelInstance := range models.GetDeclaredModels() {

		if _, present := impexModels[modelName]; present {
			continue
		}

		_, isObject := modelInstance.(models.InterfaceObject)
		_, isStorable := modelInstance.(models.InterfaceStorable)

		if isObject && isStorable {
			result = append(result, modelName)
		}
	}

	for modelName := range impexModels {
		result = append(result, modelName)
	}

	return result, nil
}

// WEB REST API used export specific model data from system
func restImpexExportModel(context api.InterfaceApplicationContext) (interface{}, error) {

	modelName := context.GetRequestArgument("model")

	var records []map[string]interface{}

	if model, present := impexModels[modelName]; present {
		exportIterator := func(item map[string]interface{}) bool {
			records = append(records, item)
			return true
		}
		err := model.Export(exportIterator)
		if err != nil {
			return nil, err
		}
	} else {

		model, err := models.GetModel(modelName)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		listable, isListable := model.(models.InterfaceListable)
		object, isObject := model.(models.InterfaceObject)

		if isListable && isObject {
			collection := listable.GetCollection()
			if collection == nil {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "12f978f5-3a90-438b-a3f4-475b34a97884", "can't obtain model collection")
			}

			var attributes []string
			for _, attribute := range object.GetAttributesInfo() {
				attributes = append(attributes, attribute.Attribute)
				if err := collection.ListAddExtraAttribute(attribute.Attribute); err != nil {
					return nil, env.ErrorDispatch(err)
				}
			}

			list, err := collection.List()
			if err != nil {
				return nil, err
			}

			for _, item := range list {
				records = append(records, item.Extra)
			}
		}
	}

	// preparing csv writer
	csvWriter := csv.NewWriter(context.GetResponseWriter())
	csvWriter.Comma = ','

	exportFilename := strings.ToLower(modelName) + "_export_" + time.Now().Format(time.RFC3339) + ".csv"

	if err := context.SetResponseContentType("text/csv"); err != nil {
		_ = env.ErrorDispatch(err)
	}
	if err := context.SetResponseSetting("Content-disposition", "attachment;filename="+exportFilename); err != nil {
		_ = env.ErrorDispatch(err)
	}

	err := MapToCSV(records, csvWriter)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return nil, nil
}

// WEB REST API for getting status of current import
func restImpexImportStatus(context api.InterfaceApplicationContext) (interface{}, error) {

	result := make(map[string]interface{})

	result["status"] = "idle"

	if importingFile != nil {
		result["status"] = "processing"
		result["name"] = importingFile.name
		result["size"] = importingFile.size

		if seeker, ok := importingFile.reader.(io.Seeker); ok {
			if position, err := seeker.Seek(0, 1); err == nil {
				result["position"] = position
			}
		}
	}

	return result, nil
}

// WEB REST API used import data to system
func restImpexImportModel(context api.InterfaceApplicationContext) (interface{}, error) {

	if importingFile != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "12f978f5-3a90-438b-a3f4-475b34a97457", "Another import is in progres. Currently processing "+importingFile.name)
	}

	var commandLine string
	exchangeDict := make(map[string]interface{})

	modelName := context.GetRequestArgument("model")

	if _, present := impexModels[modelName]; present {
		commandLine = "IMPORT " + modelName
	} else {
		commandLine = "UPDATE " + modelName
	}

	filesProcessed := 0
	additionalMessage := ""

	for fileName, attachedFile := range context.GetRequestFiles() {

		importingFile = &StructImportingFile{reader: attachedFile, name: fileName}

		if seeker, ok := attachedFile.(io.Seeker); ok {
			if fileSize, err := seeker.Seek(0, 2); err == nil {
				importingFile.size = fileSize
				_, _ = seeker.Seek(0, 0)
			}
		}

		csvReader := csv.NewReader(attachedFile)
		csvReader.Comma = ','

		err := ImportCSVData(commandLine, exchangeDict, csvReader, nil, false)
		if err != nil && additionalMessage == "" {
			_ = env.ErrorDispatch(err)
			additionalMessage += "with errors"
		}

		filesProcessed++
	}

	importingFile = nil

	return fmt.Sprintf("%d file(s) processed %s", filesProcessed, additionalMessage), nil
}

// WEB REST API used to test csv file before import
func restImpexTestImport(context api.InterfaceApplicationContext) (interface{}, error) {

	if err := context.SetResponseContentType("text/plain"); err != nil {
		_ = env.ErrorDispatch(err)
	}

	filesProcessed := 0
	additionalMessage := ""
	for fileName, attachedFile := range context.GetRequestFiles() {

		importingFile = &StructImportingFile{reader: attachedFile, name: fileName}

		if seeker, ok := attachedFile.(io.Seeker); ok {
			if fileSize, err := seeker.Seek(0, 2); err == nil {
				importingFile.size = fileSize
				_, _ = seeker.Seek(0, 0)
			}
		}

		csvReader := csv.NewReader(attachedFile)
		csvReader.Comma = ','

		err := ImportCSVScript(csvReader, context.GetResponseWriter(), true)
		if err != nil && additionalMessage == "" {
			_ = env.ErrorDispatch(err)
			additionalMessage += "with errors"
		}
		filesProcessed++
	}

	importingFile = nil

	return []byte(fmt.Sprintf("%d file(s) processed %s", filesProcessed, additionalMessage)), nil
}

// WEB REST API used to process csv file script in impex format
func restImpexImport(context api.InterfaceApplicationContext) (interface{}, error) {

	if importingFile != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "12f978f5-3a90-438b-a3f4-475b34a97888", "Another import is in progres. Currently processing "+importingFile.name)
	}

	filesProcessed := 0
	additionalMessage := ""
	for fileName, attachedFile := range context.GetRequestFiles() {

		importingFile = &StructImportingFile{reader: attachedFile, name: fileName}

		if seeker, ok := attachedFile.(io.Seeker); ok {
			if fileSize, err := seeker.Seek(0, 2); err == nil {
				importingFile.size = fileSize
				_, _ = seeker.Seek(0, 0)
			}
		}

		csvReader := csv.NewReader(attachedFile)
		csvReader.Comma = ','

		err := ImportCSVScript(csvReader, nil, false)
		if err != nil && additionalMessage == "" {
			_ = env.ErrorDispatch(err)
			additionalMessage += "with errors"
		}

		filesProcessed++
	}

	importingFile = nil

	return fmt.Sprintf("%d file(s) processed %s", filesProcessed, additionalMessage), nil
}

// WEB REST API to test conversion from csv to json / map[string]interface{}
func restImpexTestCsvToMap(context api.InterfaceApplicationContext) (interface{}, error) {

	var attachedFile io.Reader

	for fileName, requestFile := range context.GetRequestFiles() {
		if strings.HasSuffix(fileName, ".csv") {
			attachedFile = requestFile
		}
	}

	// file was not found
	if attachedFile == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6f5a582f-46f1-4b46-834d-e6b39da1ca68", "no csv file was attached")
	}

	// preparing csv reader
	var result []map[string]interface{}
	processor := func(item map[string]interface{}) bool {
		result = append(result, item)
		return true
	}

	reader := csv.NewReader(attachedFile)
	reader.Comma = ','

	// making csv processing
	exchangeDict := make(map[string]interface{})
	err := CSVToMap(reader, processor, exchangeDict)

	return result, err
}
