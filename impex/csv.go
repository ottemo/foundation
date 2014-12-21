package impex

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// MapToCSV converts map[string]interface{} to csv data
func MapToCSV(input []map[string]interface{}, csvWriter *csv.Writer) error {

	csvColumnHeaders := make(map[string]string)

	// recursuve functions for internal usage
	//----------------------------------------
	var collectColumns func(mapItem map[string]interface{}, path string)
	var getPathValue func(item map[string]interface{}, path []string) interface{}

	collectColumns = func(mapItem map[string]interface{}, path string) {
		for itemKey, itemKeyValue := range mapItem {
			currentPath := strings.Trim(path+"."+itemKey, ".")

			switch typedValue := itemKeyValue.(type) {
			case map[string]interface{}:
				collectColumns(typedValue, currentPath)
			case []map[string]interface{}:
				for _, typedValueListItem := range typedValue {
					collectColumns(typedValueListItem, currentPath)
				}
			case []interface{}:
				isMapItemsInside := false
				for _, typedValueListItem := range typedValue {
					if typedValueListItemAsMap, ok := typedValueListItem.(map[string]interface{}); ok {
						collectColumns(typedValueListItemAsMap, currentPath)
						isMapItemsInside = true
					} else {
						isMapItemsInside = false
						break
					}
				}
				if !isMapItemsInside {
					csvColumnHeaders[currentPath] = "^" + currentPath
				}
			default:
				csvColumnHeaders[currentPath] = currentPath
			}
		}
	}

	getPathValue = func(item map[string]interface{}, pathArray []string) interface{} {
		if len(pathArray) == 0 {
			return nil
		}

		if keyValue, present := item[pathArray[0]]; present {

			followingPath := pathArray[1:]
			if keyValueAsList, ok := keyValue.([]interface{}); ok {
				var result []interface{}
				for _, listValue := range keyValueAsList {
					if listValueAsMap, ok := listValue.(map[string]interface{}); ok {
						result = append(result, getPathValue(listValueAsMap, followingPath))
					} else {
						result = append(result, listValue)
					}
				}
				return result
			} else if keyValueAsMap, ok := keyValue.(map[string]interface{}); ok {
				return getPathValue(keyValueAsMap, followingPath)
			}

			if len(pathArray) == 1 {
				return keyValue
			}
		}

		return nil
	}

	// making header
	//---------------
	for _, mapItem := range input { // 1st loop - collecting information for header
		collectColumns(mapItem, "")
	}

	sortedPaths := make([]string, 0, len(csvColumnHeaders))
	for path := range csvColumnHeaders {
		sortedPaths = append(sortedPaths, path)
	}
	sort.Strings(sortedPaths)

	var csvHeader []string
	for _, currentPath := range sortedPaths {
		csvHeader = append(csvHeader, csvColumnHeaders[currentPath])
	}

	csvWriter.Write(csvHeader)
	csvWriter.Flush()

	// making contents
	//------------------
	numberOfColumns := len(csvColumnHeaders)

	for _, mapItem := range input { // 2nd loop - writing content rows
		// one record by default for item
		var itemCSVRecords [][]string
		itemCSVRecords = append(itemCSVRecords, make([]string, numberOfColumns))

		for columnIdx, columnPath := range sortedPaths {
			columnValue := getPathValue(mapItem, strings.Split(columnPath, "."))

			if arrayValue, ok := columnValue.([]interface{}); ok {
				for lineIdx, lineValue := range arrayValue {
					if len(itemCSVRecords) > lineIdx {
						itemCSVRecords[lineIdx][columnIdx] = utils.InterfaceToString(lineValue)
					} else {
						newCSVRecord := make([]string, numberOfColumns)

						if lineValue != nil {
							newCSVRecord[columnIdx] = utils.InterfaceToString(lineValue)
						}

						itemCSVRecords = append(itemCSVRecords, newCSVRecord)
					}
				}
			} else {
				itemCSVRecords[0][columnIdx] = utils.InterfaceToString(columnValue)
			}
		}

		for _, csvRecord := range itemCSVRecords {
			csvWriter.Write(csvRecord)
		}
		csvWriter.Flush()
	}

	return nil
}

// CSVToMap converts csv data to map[string]interface{} and sends to processorFunc
func CSVToMap(csvReader *csv.Reader, processorFunc func(item map[string]interface{}) bool) error {

	// reading header/columns information
	//------------------------------------
	csvColumns, err := csvReader.Read()
	if err != nil {
		return err
	}
	csvColumnsNumber := len(csvColumns)

	csvColumnFlags := make([]string, csvColumnsNumber)
	csvColumnPath := make([]string, csvColumnsNumber)
	csvColumnMemorizeAs := make([]string, csvColumnsNumber)
	csvColumnMemorizeTo := make([]string, csvColumnsNumber)
	csvColumnType := make([]string, csvColumnsNumber)
	csvColumnConvertors := make([]string, csvColumnsNumber)

	allColumnsBlankFlag := true
	// extracting column header parts
	for idx, column := range csvColumns {
		regexpGroups := ConstCSVColumnRegexp.FindStringSubmatch(column)

		if len(regexpGroups) == 0 { // un-recognized column header
			if strings.TrimSpace(column) != "" { // unless it is blank, considering as path
				csvColumnPath[idx] = column
				allColumnsBlankFlag = false
			}
			continue
		}
		allColumnsBlankFlag = false

		csvColumnFlags[idx] = regexpGroups[1]
		csvColumnPath[idx] = regexpGroups[2]

		tmp := strings.Split(strings.TrimSpace(regexpGroups[3]), " ")
		if len(tmp) == 2 {
			csvColumnMemorizeAs[idx] = strings.TrimSpace(tmp[0])
			csvColumnMemorizeTo[idx] = strings.TrimSpace(tmp[1])
		}

		csvColumnType[idx] = strings.TrimSpace(regexpGroups[4])
		csvColumnConvertors[idx] = strings.TrimSpace(regexpGroups[5])
	}

	if allColumnsBlankFlag {
		processorFunc(make(map[string]interface{}))
		return nil
	}

	// reading CSV contents
	//----------------------
	var result []map[string]interface{}
	csvRecordMap := make(map[string]interface{})
	csvMemorize := make(map[string]interface{})

	// column path can be non static, there we storing currently calculated value
	columnPath := make([]string, csvColumnsNumber)
	columnPathArray := make([][]string, csvColumnsNumber)

	csvRecordNum := 1
	for csvRecord, err := csvReader.Read(); err == nil; csvRecord, err = csvReader.Read() { // csv records loop

		// If path collapse flag is set - we have a new object on that csv row for given path
		// each sub-path have own collapse indicator, if lower path collapses, it not affects top path
		// but not vice versa, if top one collapses, lower should also
		pathCollapseFlag := make(map[string]bool)

		// 1st loop - looking for collapsing paths
		//-----------------------------------------
		blankRecordFlag := true
		for columnIdx, value := range csvRecord { // 1st loop
			// skipping blank values and columns without header information
			if value == "" || csvColumnPath[columnIdx] == "" {
				continue
			}

			blankRecordFlag = false

			isArrayColumn := (csvColumnFlags[columnIdx] == "^")
			isMaybeArrayColumn := (csvColumnFlags[columnIdx] == "?")
			isIgnoreColumn := (csvColumnFlags[columnIdx] == "~")

			keyPath := csvColumnPath[columnIdx]
			keyPathArray := strings.Split(keyPath, ".")

			// updating variable paths (with @) if needed,
			// and storing calculations for following data processing loop
			keyPathWasUpdated := false
			for idx, pathValue := range keyPathArray {
				if strings.HasPrefix(pathValue, "@") {
					if value, present := csvMemorize[pathValue[1:]]; present {
						if stringValue, ok := value.(string); ok {
							keyPathArray[idx] = stringValue
							keyPathWasUpdated = true
						}
					}
				}
			}
			if keyPathWasUpdated {
				keyPath = strings.Join(keyPathArray, ".")
				keyPathArray = strings.Split(keyPath, ".")
			}
			columnPath[columnIdx] = keyPath
			columnPathArray[columnIdx] = keyPathArray

			// storing column value or path+value as variable for variable paths
			switch csvColumnMemorizeAs[columnIdx] {
			case "=":
				csvMemorize[csvColumnMemorizeTo[columnIdx]] = keyPath + "." + value
			case ">":
				csvMemorize[csvColumnMemorizeTo[columnIdx]] = value
			}

			// setting collapse condition for path path if it not array, and value was changed
			if !isArrayColumn && !isMaybeArrayColumn && !isIgnoreColumn && csvColumnMemorizeAs[columnIdx] == "" {
				path := ""
				if arrayLen := len(keyPathArray); arrayLen > 1 {
					path = strings.Join(keyPathArray[:arrayLen-1], ".")
				}
				pathCollapseFlag[path] = true
			}
		}

		// check for blank record
		//------------------------
		if blankRecordFlag {
			break // if we found blank record - this means end data chunk
		}

		// checking for top level collapse
		//---------------------------------
		if value, present := pathCollapseFlag[""]; present && value {
			pathCollapseFlag = make(map[string]bool)

			if csvRecordNum != 1 { // on first row there are blank map
				if !processorFunc(csvRecordMap) {
					return nil // so, processor says to stop
				}
			}
			csvRecordMap = make(map[string]interface{})
		}

		// processing record values (2nd loop)
		//-------------------------------------
		for columnIdx, value := range csvRecord { // 2nd loop
			// skipping blank values, columns without header, path folding columns
			if value == "" || csvColumnPath[columnIdx] == "" || csvColumnMemorizeAs[columnIdx] == "=" {
				continue
			}

			isArrayColumn := (csvColumnFlags[columnIdx] == "^")
			isMaybeArrayColumn := (csvColumnFlags[columnIdx] == "?")

			keyPathArray := columnPathArray[columnIdx]
			lastPathIdx := len(keyPathArray) - 1

			var prevPathMap map[string]interface{}
			var prevPathKey string
			var prevPathValue interface{} = result

			currentPathMap := csvRecordMap
			currentPath := ""

			// moving down through column path
			for idx, key := range keyPathArray { // path loop

				// check we need new object for path
				if value, present := pathCollapseFlag[currentPath]; present && value {
					pathCollapseFlag[currentPath] = false

					newMapValue := make(map[string]interface{})
					if prevPathValueAsList, isList := prevPathValue.([]interface{}); isList {
						prevPathMap[prevPathKey] = append(prevPathValueAsList, newMapValue)
					} else {
						if prevPathValueAsMap, isMap := prevPathValue.(map[string]interface{}); isMap && len(prevPathValueAsMap) > 0 {
							prevPathMap[prevPathKey] = []interface{}{prevPathValue, newMapValue}
						} else {
							prevPathMap[prevPathKey] = newMapValue
						}
					}
					currentPathMap = newMapValue
				}

				if idx == lastPathIdx { // we are at end of key path (i.e. on x for key like a.b.c.d.x)

					// trying to convert string value to supposed type
					var typedValue interface{}

					if columnType := csvColumnType[columnIdx]; columnType != "" {
						if result, err := utils.StringToType(value, columnType); err == nil {
							typedValue = result
						}
					} else {
						typedValue = utils.StringToInterface(value)
					}

					currentKeyValue, present := currentPathMap[key]
					if present && (isArrayColumn || isMaybeArrayColumn) {
						if currentValueAsArray, ok := currentKeyValue.([]interface{}); ok {
							currentPathMap[key] = append(currentValueAsArray, typedValue)
						} else {
							currentPathMap[key] = []interface{}{currentKeyValue, typedValue}
						}
					} else {

						if isArrayColumn {
							currentPathMap[key] = []interface{}{typedValue}
						} else {
							currentPathMap[key] = typedValue
						}
					}

				} else { // still looping through path (i.e. on i for path like i.i.i.i.x)

					prevPathMap = currentPathMap
					prevPathKey = key

					currentKeyValue, present := currentPathMap[key]
					prevPathValue = currentKeyValue

					if currentKeyValueAsMap, isMap := currentKeyValue.(map[string]interface{}); present && isMap {
						currentPathMap = currentKeyValueAsMap
					} else if currentKeyValueAsList, isList := currentKeyValue.([]interface{}); present && isList {
						makeNewMapFlag := true
						if len(currentKeyValueAsList) > 0 {
							prevPathValue = currentKeyValueAsList
							lastItemValue := currentKeyValueAsList[len(currentKeyValueAsList)-1]
							if lasItemValueAsMap, isMap := lastItemValue.(map[string]interface{}); isMap {
								currentPathMap = lasItemValueAsMap
								makeNewMapFlag = false
							}
						}
						if makeNewMapFlag {
							newMapValue := make(map[string]interface{})
							prevPathValue = append(currentKeyValueAsList, newMapValue)
							currentPathMap[key] = prevPathValue
							currentPathMap = newMapValue
						}
					} else {
						newMapValue := make(map[string]interface{})
						prevPathValue = newMapValue
						currentPathMap[key] = newMapValue
						currentPathMap = newMapValue
					}
				}

				// updating path
				if currentPath != "" {
					currentPath += "."
				}
				currentPath += key
			}
		}

		csvRecordNum++
	}
	processorFunc(csvRecordMap)

	return nil
}

// ImportCSV imports csv data using command->data csv format
func ImportCSV(csvReader *csv.Reader) error {

	// impex csv file should contain command preceding data
	commandLine := ""
	appendFlag := false
	for csvRecord, err := csvReader.Read(); err == nil; csvRecord, err = csvReader.Read() { // csv records loop

		// reading csv command line
		//--------------------------
		csvLine := ""
		for columnIdx, csvColumn := range csvRecord {
			csvColumn = strings.TrimSpace(csvColumn)
			if columnIdx != 0 && strings.Contains(csvColumn, " ") {
				csvColumn = "\"" + strings.Replace(csvColumn, "\"", "\\\"", -1) + "\""
			}

			csvLine += " " + csvColumn
		}
		csvLine = strings.TrimSpace(csvLine)

		// skipping blank lines
		if csvLine == "" {
			continue
		}

		// checking that command line not suppose to read following csv record
		//---------------------------------------------------------------------
		if strings.HasPrefix(csvLine, "|") {
			if appendFlag {
				commandLine += " " + csvLine
				appendFlag = false
			} else {
				commandLine = csvLine + " " + commandLine
				continue
			}
		}

		if strings.HasSuffix(csvLine, "...") {
			csvLine = strings.TrimSuffix(csvLine, "...")
			commandLine += " " + csvLine
			appendFlag = true
			continue
		}
		commandLine = csvLine + " " + commandLine
		commandLine = strings.TrimSpace(commandLine)

		if ConstImpexLog || ConstDebugLog {
			env.Log("impex.log", env.ConstLogPrefixDebug, fmt.Sprintf("Command line: %s", commandLine))
		}

		// looking for required commands and preparing them to process
		//-------------------------------------------------------------
		exchangeDict := make(map[string]interface{})
		var commandsChain []InterfaceImpexImportCmd

		for _, command := range utils.SplitQuotedStringBy(commandLine, '|') {
			command = strings.TrimSpace(command)
			args := utils.SplitQuotedStringBy(command, ' ', '\n', '\t')

			if len(args) > 0 {
				if cmd, present := importCmd[args[0]]; present {
					if err := cmd.Init(args, exchangeDict); err == nil {
						commandsChain = append(commandsChain, cmd)
					} else {
						return env.ErrorDispatch(err)
					}
				} else {
					return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "004e9f7bbb974356bbc25e084736983b", "unknows cmd '"+args[0]+"'")
				}
			}
		}

		if len(commandsChain) == 0 {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "84ac22e2c2154375af15a9eed636b16c", "There are no commands for csv data processing")
		}

		// making csv data processor based on received commands
		//------------------------------------------------------
		dataProcessor := func(itemData map[string]interface{}) bool {
			if ConstImpexLog || ConstDebugLog {
				env.Log("impex.log", env.ConstLogPrefixDebug, fmt.Sprintf("Processing: %s", utils.EncodeToJSONString(itemData)))
			}

			var input interface{}
			for _, command := range commandsChain {
				if ConstDebugLog {
					env.Log("impex.log", env.ConstLogPrefixDebug, fmt.Sprintf("Command: %T", command))
					env.Log("impex.log", env.ConstLogPrefixDebug, fmt.Sprintf("Input: %#v", input))
					env.Log("impex.log", env.ConstLogPrefixDebug, fmt.Sprintf("itemData: %s", utils.EncodeToJSONString(itemData)))
					env.Log("impex.log", env.ConstLogPrefixDebug, fmt.Sprintf("Exchange: %s", utils.EncodeToJSONString(exchangeDict)))
				}

				input, err = command.Process(itemData, input, exchangeDict)
				if err != nil {
					if ConstImpexLog || ConstDebugLog {
						env.Log("impex.log", env.ConstLogPrefixDebug, fmt.Sprintf("Error: %s", err.Error()))
					}
					env.ErrorDispatch(err)
					return true
				}
				if ConstImpexLog || ConstDebugLog {
					env.Log("impex.log", env.ConstLogPrefixDebug, "Finished ok")
				}
			}
			return true
		}

		// passing control to data block reader
		//--------------------------------------
		err = CSVToMap(csvReader, dataProcessor)
		if err != nil && err != io.EOF {
			return env.ErrorDispatch(err)
		}

		commandLine = ""
	}

	return nil
}
