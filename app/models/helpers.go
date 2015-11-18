package models

import (
	"strconv"
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetModelAndSetID retrieves current model implementation and sets its ID to some value
func GetModelAndSetID(modelName string, modelID string) (InterfaceStorable, error) {
	someModel, err := GetModel(modelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	storableModel, ok := someModel.(InterfaceStorable)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "652d1949-b661-438e-9097-231a52734feb", "model is not InterfaceStorable capable")
	}

	err = storableModel.SetID(modelID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return storableModel, nil
}

// LoadModelByID loads model data in current implementation
func LoadModelByID(modelName string, modelID string) (InterfaceStorable, error) {

	someModel, err := GetModel(modelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	storableModel, ok := someModel.(InterfaceStorable)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "56ec204a-dbb9-49fc-a5e9-9d43e1f19025", "model is not InterfaceStorable capable")
	}

	err = storableModel.Load(modelID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return storableModel, nil
}

// ApplyExtraAttributes modifies given model collection with adding extra attributes to list
//   - default attributes are specified in StructListItem as static fields
//   - StructListItem fields can be not a direct copy of model attribute,
//   - extra attributes are taken from model directly
func ApplyExtraAttributes(context api.InterfaceApplicationContext, collection InterfaceCollection) error {
	extra := context.GetRequestArgument("extra")
	if extra == "" {
		contentMap, err := api.GetRequestContentAsMap(context)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		if contentMapExtra, present := contentMap["extra"]; present && contentMapExtra != "" {
			extra = utils.InterfaceToString(contentMapExtra)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0bc07b3d-1443-4594-af82-9d15211ed179", "no extra attributes specified")
		}
	}

	extraAttributes := utils.Explode(utils.InterfaceToString(extra), ",")
	for _, attributeName := range extraAttributes {
		err := collection.ListAddExtraAttribute(attributeName)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// ApplyFilters modifies collection with applying filters from request URL
func ApplyFilters(context api.InterfaceApplicationContext, collection db.InterfaceDBCollection) error {

	// sets filter to particular attribute within collection
	addFilterToCollection := func(attributeName string, attributeValue string, groupName string) {
		if collection.HasColumn(attributeName) {

			filterOperator := "="
			for _, prefix := range []string{">=", "<=", "!=", ">", "<", "~"} {
				if strings.HasPrefix(attributeValue, prefix) {
					attributeValue = strings.TrimPrefix(attributeValue, prefix)
					filterOperator = prefix
				}
			}
			if filterOperator == "~" {
				filterOperator = "like"
			}

			switch {
			case strings.Contains(attributeValue, ".."):
				rangeValues := strings.Split(attributeValue, "..")
				if rangeValues[0] != "" {
					collection.AddGroupFilter(groupName, attributeName, ">=", rangeValues[0])
				}
				if rangeValues[1] != "" {
					collection.AddGroupFilter(groupName, attributeName, "<=", rangeValues[1])
				}

			case strings.Contains(attributeValue, ","):
				options := strings.Split(attributeValue, ",")
				if filterOperator == "=" {
					collection.AddGroupFilter(groupName, attributeName, "in", options)
				} else {
					filterGroupName := attributeName + "_inFilter"
					collection.SetupFilterGroup(filterGroupName, true, groupName)
					for _, optionValue := range options {
						collection.AddGroupFilter(filterGroupName, attributeName, filterOperator, optionValue)
					}
				}

			default:
				attributeType := collection.GetColumnType(attributeName)
				if attributeType != db.ConstTypeText && attributeType != db.ConstTypeID &&
					!strings.Contains(attributeType, db.ConstTypeVarchar) &&
					filterOperator == "like" {

					filterOperator = "="
				}

				if typedValue, err := utils.StringToType(attributeValue, attributeType); err == nil {
					// fix for NULL db boolean values filter (perhaps should be part of DB adapter)
					if attributeType == db.ConstTypeBoolean && typedValue == false {
						filterGroupName := attributeName + "_applyFilter"

						collection.SetupFilterGroup(filterGroupName, true, groupName)
						collection.AddGroupFilter(filterGroupName, attributeName, filterOperator, typedValue)
						collection.AddGroupFilter(filterGroupName, attributeName, "=", nil)
					} else {
						collection.AddGroupFilter(groupName, attributeName, filterOperator, typedValue)
					}
				} else {
					collection.AddGroupFilter(groupName, attributeName, filterOperator, attributeValue)
				}
			}
		}

	}

	collection.SetLimit(0, ConstCollectionListLimit)
	// checking arguments user set
	for attributeName, attributeValue := range context.GetRequestArguments() {
		switch attributeName {

		// collection limit required
		case "limit":
			collection.SetLimit(GetListLimit(context))

			// collection sort required
		case "sort":
			attributesList := strings.Split(attributeValue, ",")

			for _, attributeName := range attributesList {
				descOrder := false
				if attributeName[0] == '^' {
					descOrder = true
					attributeName = strings.Trim(attributeName, "^")
				}
				collection.AddSort(attributeName, descOrder)
			}

			// filter for any columns matches value required
		case "search":
			collection.SetupFilterGroup("search", true, "")

			// checking value type we are working with
			lookingFor := "text"
			if strings.HasPrefix(attributeValue, ">") || strings.HasPrefix(attributeValue, "<") || strings.Contains(attributeValue, "..") {
				lookingFor = "number"
			}
			if strings.HasPrefix(attributeValue, "~") {
				lookingFor = "text"
			}
			if lookingFor != "number" {
				searchValue := strings.TrimLeft(attributeValue, "><=~")
				if strings.Trim(searchValue, "1234567890.") == "" {
					lookingFor = "text,number"
				}
			}

			// looking for possible attributes to filter
			for attributeName, attributeType := range collection.ListColumns() {
				switch {
				case attributeType == db.ConstTypeText || strings.Contains(attributeType, db.ConstTypeVarchar):
					if strings.Contains(lookingFor, "text") {
						addFilterToCollection(attributeName, attributeValue, "search")
					}

				case attributeType == db.ConstTypeFloat ||
					attributeType == db.ConstTypeDecimal ||
					attributeType == db.ConstTypeMoney ||
					attributeType == db.ConstTypeInteger:

					if strings.Contains(lookingFor, "number") {
						addFilterToCollection(attributeName, attributeValue, "search")
					}
				}
			}

		default:
			addFilterToCollection(attributeName, attributeValue, "default")
		}
	}
	return nil
}

// GetListLimit returns (offset, limit, error) values based on request string value
//   "1,2" will return offset: 1, limit: 2, error: nil
//   "2" will return offset: 0, limit: 2, error: nil
//   "something wrong" will return offset: 0, limit: 0, error: [error msg]
func GetListLimit(context api.InterfaceApplicationContext) (int, int) {
	limitValue := ""

	if value := context.GetRequestArgument("limit"); value != "" {
		limitValue = utils.InterfaceToString(value)
	} else {
		contentMap, err := api.GetRequestContentAsMap(context)
		if err == nil {
			if value, isLimit := contentMap["limit"]; isLimit {
				if value, ok := value.(string); ok {
					limitValue = value
				}
			}
		}
	}
	// limitValue, _ = url.QueryUnescape(limitValue)

	splitResult := strings.Split(limitValue, ",")
	if len(splitResult) > 1 {
		offset, err := strconv.Atoi(strings.TrimSpace(splitResult[0]))
		if err != nil {
			return 0, 0
		}

		limit, err := strconv.Atoi(strings.TrimSpace(splitResult[1]))
		if err != nil {
			return 0, 0
		}

		return offset, limit
	} else if len(splitResult) > 0 {
		limit, err := strconv.Atoi(strings.TrimSpace(splitResult[0]))
		if err != nil {
			return 0, 0
		}

		return 0, limit
	}

	return 0, 0
}
