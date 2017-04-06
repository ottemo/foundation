package api

import (
	"net/http"

	"strconv"
	"strings"
	"time"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// StartSession returns session object for request or creates new one.  To use
// a secure session cookie in HTTPS, please set the environment variable
// OTTEMOCOOKIE.  It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false,
// False. Any other value returns an error.
func StartSession(context InterfaceApplicationContext) (InterfaceSession, error) {

	request := context.GetRequest()
	// use secure cookies by default
	var flagSecure = true
	var tmpSecure = ""
	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("secure_cookie", tmpSecure); iniValue != "" {
			tmpSecure = iniValue
			flagSecure, _ = strconv.ParseBool(tmpSecure)
		}
	}

	// old method - HTTP specific
	if _, ok := request.(*http.Request); ok {
		responseWriter := context.GetResponseWriter()
		if responseWriter, ok := responseWriter.(http.ResponseWriter); ok {
			// check session-cookie or header
			if sessionID := context.GetRequestSetting(ConstSessionCookieName); sessionID != nil {
				sessionID := utils.InterfaceToString(sessionID)
				sessionInstance, err := currentSessionService.Get(sessionID, true)
				if err == nil {
					return sessionInstance, nil
				}
			}

			// session cookie is not set or expired, making new
			result, err := currentSessionService.New()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			// Session Cookie Declaration
			// - expires in 1 year
			// - Domain defaults to the full subdomain path
			cookieExpires := time.Now().Add(365 * 24 * time.Hour)
			var cookie = &http.Cookie{
				Name:     ConstSessionCookieName,
				Value:    result.GetID(),
				Path:     "/",
				Secure:   flagSecure,
				HttpOnly: true,
				Expires:  cookieExpires,
			}
			http.SetCookie(responseWriter, cookie)

			return result, nil
		}
	}

	// new approach - not HTTP related
	if sessionID := context.GetRequestSetting(ConstSessionCookieName); sessionID != nil {
		sessionID := utils.InterfaceToString(sessionID)
		sessionInstance, err := currentSessionService.Get(sessionID, true)
		if err == nil {
			// ignore non critical error
			_ = context.SetResponseSetting(ConstSessionCookieName, sessionInstance.GetID())
			return sessionInstance, nil
		}
	}

	// session id not found of was not specified - making new session
	sessionInstance, err := currentSessionService.New()
	if err == nil {
		// ignore non critical error
		_ = context.SetResponseSetting(ConstSessionCookieName, sessionInstance.GetID())
	}

	return sessionInstance, err
}

// NewSession returns new session instance
func NewSession() (InterfaceSession, error) {
	return currentSessionService.New()
}

// GetSessionByID returns session instance by id or nil
func GetSessionByID(sessionID string, create bool) (InterfaceSession, error) {
	sessionInstance, err := currentSessionService.Get(sessionID, create)

	// "(*session.DefaultSession)(nil)" is not "nil", and we want to have exact nil
	if sessionInstance == nil {
		return nil, err
	}

	return sessionInstance, err
}

// ValidateAdminRights returns nil if session contains admin rights
func ValidateAdminRights(context InterfaceApplicationContext) error {

	if IsAdminSession(context) {
		return nil
	}

	// it is un-secure as request can be intercepted by malefactor, so use it only if no other way to do auth
	// (we are using it for "gulp build" local tool, so all data within one host)
	if value := context.GetRequestArgument(ConstGETAuthParamName); value != "" {
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

	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2f3438ba-7fb7-4811-b8a5-7acf36910d3d", "no admin rights")
}

// IsAdminHandler returns middleware API Handler that checks admin rights
func IsAdminHandler(next FuncAPIHandler) FuncAPIHandler {
	return func(context InterfaceApplicationContext) (interface{}, error) {
		isAdminErr := ValidateAdminRights(context)

		if isAdminErr != nil {
			context.SetResponseStatusForbidden()
			return nil, isAdminErr
		}

		return next(context)
	}
}

// IsAdminSession returns true if session with admin rights
func IsAdminSession(context InterfaceApplicationContext) bool {
	return utils.InterfaceToBool(context.GetSession().Get(ConstSessionKeyAdminRights))
}

// GetRequestContentAsMap tries to represent HTTP request content in map[string]interface{} format
func GetRequestContentAsMap(context InterfaceApplicationContext) (map[string]interface{}, error) {

	result, ok := context.GetRequestContent().(map[string]interface{})
	if !ok {
		result = make(map[string]interface{})
	}

	return result, nil
}

// GetContentValue looks for a given name within request content (map only), returns value or nil if not found
func GetContentValue(context InterfaceApplicationContext, name string) interface{} {
	if contentMap, ok := context.GetRequestContent().(map[string]interface{}); ok {
		if value, present := contentMap[name]; present {
			return value
		}
	}

	return nil
}

// GetArgumentOrContentValue looks for a given name within request parameters and content (map only), returns first occurrence
// according to mentioned sequence or nil if not found.
func GetArgumentOrContentValue(context InterfaceApplicationContext, name string) interface{} {

	if value := context.GetRequestArgument(name); value != "" {
		return value
	}

	if contentMap, ok := context.GetRequestContent().(map[string]interface{}); ok {
		if value, present := contentMap[name]; present {
			return value
		}
	}

	return nil
}
