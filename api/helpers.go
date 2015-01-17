package api

import (
	"strconv"
	"strings"

	"net/http"
	"net/url"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// StartSession returns session object for request or creates new one
func StartSession(params *StructAPIHandlerParams) (InterfaceSession, error) {

	request := params.Request
	responseWriter := params.ResponseWriter

	// check session-cookie
	cookie, err := request.Cookie(ConstSessionCookieName)
	if err == nil {
		// looking for cookie-based session
		sessionID := cookie.Value

		sessionInstance, err := currentSessionService.Get(sessionID)
		if err == nil {
			return sessionInstance, nil
		}
	} else {
		if err != http.ErrNoCookie {
			return nil, err
		}
	}

	// session cookie is not set or expired, making new
	result, err := currentSessionService.New()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// storing session id to cookie
	cookie = &http.Cookie{Name: ConstSessionCookieName, Value: result.GetID(), Path: "/"}
	http.SetCookie(responseWriter, cookie)

	return result, nil
}

// NewSession returns new session instance
func NewSession() (InterfaceSession, error) {
	return currentSessionService.New()
}

// GetSessionByID returns session instance by id or nil
func GetSessionByID(sessionID string) (InterfaceSession, error) {
	sessionInstance, err := currentSessionService.Get(sessionID)

	// "(*session.DefaultSession)(nil)" is not "nil", and we want to have exact nil
	if sessionInstance == nil {
		return nil, err
	}

	return sessionInstance, err
}

// ValidateAdminRights returns nil if admin rights allowed for current session
func ValidateAdminRights(params *StructAPIHandlerParams) error {

	if value := params.Session.Get(ConstSessionKeyAdminRights); value != nil {
		if value, ok := value.(bool); ok && value {
			return nil
		}
	}

	// it is un-secure as request can be intercepted by malefactor, so use it only if no other way to do auth
	// (we are using it for "gulp build" local tool, so all data within one host)
	if value, present := params.RequestGETParams[ConstGETAuthParamName]; present {
		if splited := strings.Split(value, ":"); len(splited) > 1 {
			login := splited[0]
			password := splited[1]

			rootLogin := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathStoreRootLogin))
			rootPassword := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathStoreRootPassword))

			if login == rootLogin && password == rootPassword {
				return nil
			}
		}
	}

	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0bc07b3d-1443-4594-af82-9d15211ed179", "no admin rights")
}

// GetRequestContentAsMap tries to represent HTTP request content in map[string]interface{} format
func GetRequestContentAsMap(params *StructAPIHandlerParams) (map[string]interface{}, error) {

	result, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		if params.Request.Method == "POST" {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ab71d4a9-0304-4d48-ae78-556aa56ac53c", "unexpected request content")
		}
		result = make(map[string]interface{})
	}

	return result, nil
}

// ApplyFilters modifies collection with applying filters from request URL
func ApplyFilters(params *StructAPIHandlerParams, collection db.InterfaceDBCollection) error {

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
				if attributeType != db.ConstTypeText &&
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

	// checking arguments user set
	for attributeName, attributeValue := range params.RequestGETParams {
		switch attributeName {

		// collection limit required
		case "limit":
			collection.SetLimit(GetListLimit(params))

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
func GetListLimit(params *StructAPIHandlerParams) (int, int) {
	limitValue := ""

	if value, isLimit := params.RequestURLParams["limit"]; isLimit {
		limitValue = value
	} else if value, isLimit := params.RequestGETParams["limit"]; isLimit {
		limitValue = value
	} else {
		contentMap, err := GetRequestContentAsMap(params)
		if err == nil {
			if value, isLimit := contentMap["limit"]; isLimit {
				if value, ok := value.(string); ok {
					limitValue = value
				}
			}
		}
	}
	limitValue, _ = url.QueryUnescape(limitValue)

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
