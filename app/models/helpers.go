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
			// Ignore error processing if attribute is wrong. Log it anyway as warning.
			// Attribute could be absent in collection declaration because of external attribute package is disabled.
			// For example, "sale_price" attribute could be absent, if "Sale Price" package disabled.
			env.Log(ConstErrorModule+".log", env.ConstLogPrefixWarning, "incorrect or disabled attribute '" + attributeName + "' added to collection list.")
		}
	}

	return nil
}

// ApplyFilters modifies collection with applying filters from request URL
func ApplyFilters(context api.InterfaceApplicationContext, collection db.InterfaceDBCollection) error {

	// sets filter to particular attribute within collection
	addFilterToCollection := func(attributeName string, attributeValue string, groupName string) error {
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

			attributeType := collection.GetColumnType(attributeName)

			switch {
			case strings.Contains(attributeValue, ".."):
				rangeValues := strings.Split(attributeValue, "..")
				if rangeValues[0] != "" {
					if err := collection.AddGroupFilter(groupName, attributeName, ">=", rangeValues[0]); err != nil {
						return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ceffd0a6-c61b-41db-8b04-2b8582a18030", "unable to add group filter: "+err.Error())
					}
				}
				if rangeValues[1] != "" {
					if err := collection.AddGroupFilter(groupName, attributeName, "<=", rangeValues[1]); err != nil {
						return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e1b24a05-c918-4673-8f46-c997518c0a29", "unable to add group filter: "+err.Error())
					}
				}

			// check also if value without commas but field type is array
			case strings.Contains(attributeValue, ",") || strings.HasPrefix(attributeType, "[]"):
				attributeValue = strings.Trim(attributeValue, ",")
				options := strings.Split(attributeValue, ",")
				if filterOperator == "=" {
					if err := collection.AddGroupFilter(groupName, attributeName, "in", options); err != nil {
						return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "94b0b03a-f608-4c1e-891a-4170338c6043", "unable to add group filter: "+err.Error())
					}
				} else {
					// OrSequence should be "false" because (a != 1 || a != 2) => have no sense
					// If conflict detected with other operators, like ">=", "<=", "!=", ">", "<", "~"
					// this functionality should be reviewed to add some kind of OrSequence flag to
					// function parameters or to check filterOperator

					for _, optionValue := range options {
						if err := collection.AddGroupFilter(groupName, attributeName, filterOperator, optionValue); err != nil {
							return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "afbc434e-3ca2-48a4-9674-f0c0eb67c8e3", "unable to add group filter: "+err.Error())
						}
					}
				}

			default:
				if attributeType != db.ConstTypeText && attributeType != db.ConstTypeID &&
					!strings.Contains(attributeType, db.ConstTypeVarchar) &&
					filterOperator == "like" {

					filterOperator = "="
				}

				if typedValue, err := utils.StringToType(attributeValue, attributeType); err == nil {
					// fix for NULL db boolean values filter (perhaps should be part of DB adapter)
					if attributeType == db.ConstTypeBoolean && typedValue == false {
						filterGroupName := attributeName + "_applyFilter"

						if err := collection.SetupFilterGroup(filterGroupName, true, groupName); err != nil {
							return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "84dfb177-c308-4bac-9d1d-accfe88f9e38", "unable to setup filter group: "+err.Error())
						}
						if err := collection.AddGroupFilter(filterGroupName, attributeName, filterOperator, typedValue); err != nil {
							return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6841cc89-dc41-42de-89f6-cb94f93af1d3", "unable to add group filter: "+err.Error())
						}
						if err := collection.AddGroupFilter(filterGroupName, attributeName, "=", nil); err != nil {
							return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2e93e11b-7f64-4346-bced-8195068f5951", "unable to add group filter: "+err.Error())
						}
					} else {
						if err := collection.AddGroupFilter(groupName, attributeName, filterOperator, typedValue); err != nil {
							return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "84f369ff-9385-4035-85d9-bd0f2c46df5b", "unable to add group filter: "+err.Error())
						}
					}
				} else {
					if err := collection.AddGroupFilter(groupName, attributeName, filterOperator, attributeValue); err != nil {
						return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7755261d-5574-4b92-9e2a-33c18427cb03", "unable to add group filter: "+err.Error())
					}
				}
			}
		}

		return nil
	}

	if err := collection.SetLimit(0, ConstCollectionListLimit); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "82b0665a-9d1a-42b8-a49a-e7fb41ea5f83", "unable to set default limit: "+err.Error())
	}
	// checking arguments user set
	for attributeName, attributeValue := range context.GetRequestArguments() {
		switch attributeName {

		// collection limit required
		case "limit":
			if err := collection.SetLimit(GetListLimit(context)); err != nil {
				return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d1339211-e72d-41c9-b079-c4128612704c", "unable to set limit: "+err.Error())
			}

			// collection sort required
		case "sort":
			attributesList := strings.Split(attributeValue, ",")

			for _, attributeName := range attributesList {
				descOrder := false
				if attributeName[0] == '^' {
					descOrder = true
					attributeName = strings.Trim(attributeName, "^")
				}
				if err := collection.AddSort(attributeName, descOrder); err != nil {
					return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e80eb8fc-cdf0-4fc1-928e-fec77a415835", "unable to add sort: "+err.Error())
				}
			}

			// filter for any columns matches value required
		case "search":
			if err := collection.SetupFilterGroup("search", true, ""); err != nil {
				return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5e8f6e31-f335-4660-97ef-8f886890fe9a", "unable to setup filter group: "+err.Error())
			}

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
						if err := addFilterToCollection(attributeName, attributeValue, "search"); err != nil {
							return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "45b820eb-07b6-4708-988e-a3c8236778cc", "Unable to add filter to collection:" + err.Error())						}
					}

				case attributeType == db.ConstTypeFloat ||
					attributeType == db.ConstTypeDecimal ||
					attributeType == db.ConstTypeMoney ||
					attributeType == db.ConstTypeInteger:

					if strings.Contains(lookingFor, "number") {
						if err := addFilterToCollection(attributeName, attributeValue, "search"); err != nil {
							return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "31e7e803-bd29-4f74-bfd3-f20de6979b1f", "Unable to add filter to collection:" + err.Error())
						}
					}
				}
			}

		default:
			if err := addFilterToCollection(attributeName, attributeValue, "default"); err != nil {
				return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ce1985e3-e96d-4724-ba4e-0381393e9e1a", "Unable to add filter to collection:" + err.Error())
			}
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
