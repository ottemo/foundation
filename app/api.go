package app

import (
	"bytes"
	"runtime"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("app/email", api.POST, restSendEmail)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app/login", api.GET, restLogin)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app/login", api.POST, restLogin)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app/logout", api.GET, restLogout)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app/rights", api.GET, restRightsInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app/status", api.GET, restStatusInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app/location", api.POST, setSessionTimeZone)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app/location", api.GET, getSessionTimeZone)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function login application with root rights
func restLogin(context api.InterfaceApplicationContext) (interface{}, error) {

	requestLogin := context.GetRequestArgument("login")
	requestPassword := context.GetRequestArgument("password")

	if requestLogin == "" || requestPassword == "" {

		requestData, err := api.GetRequestContentAsMap(context)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if !utils.KeysInMapAndNotBlank(requestData, "login", "password") {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fee28a56-adb1-44b9-a0e2-1c9be6bd6fdb", "login and password should be specified")
		}

		requestLogin = utils.InterfaceToString(requestData["login"])
		requestPassword = utils.InterfaceToString(requestData["password"])
	}

	rootLogin := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathStoreRootLogin))
	rootPassword := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathStoreRootPassword))

	if requestLogin == rootLogin && requestPassword == rootPassword {
		context.GetSession().Set(api.ConstSessionKeyAdminRights, true)

		return "ok", nil
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "68546aa8-a6be-4c31-ac44-ea4278dfbdb0", "wrong login or password")
}

// WEB REST API function logout application - session data clear
func restLogout(context api.InterfaceApplicationContext) (interface{}, error) {
	err := context.GetSession().Close()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	return "ok", nil
}

// restContactUs creates a new email message via a POST from the Contact Us form
//   - following attributes are required:
//   - "formLocation"
func restSendEmail(context api.InterfaceApplicationContext) (interface{}, error) {

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// formLocation is the only required parameter
	if !utils.KeysInMapAndNotBlank(requestData, "formLocation") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e47f0671-0e19-4bbd-a771-ae4fac56a714", "A required parameter is missing: 'formLocation'")
	}

	// remove form location from map
	frmLocation := utils.InterfaceToString(requestData["formLocation"])
	delete(requestData, "formLocation")

	recipient := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathContactUsRecipient))

	// create body of email
	var body bytes.Buffer

	body.WriteString("The form contained the following information: <br><br>")
	for key, val := range requestData {
		body.WriteString(key + ": " + utils.InterfaceToString(val) + "<br>")
	}

	headers := map[string]string{
		"To": recipient,
	}

	if replyToEmail, present := requestData["email"]; present {
		headers["Reply-To"] = utils.InterfaceToString(replyToEmail)
	}

	emailContext := map[string]interface{}{
		"Subject": "New Message from Form: " + frmLocation,
	}

	err = SendMailEx(headers, body.String(), emailContext)

	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// WEB REST API function to get info about current rights
func restRightsInfo(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})

	result["is_admin"] = utils.InterfaceToBool(context.GetSession().Get(api.ConstSessionKeyAdminRights))

	return result, nil
}

// WEB REST API function to get info about current application status
func restStatusInfo(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})

	result["Ottemo"] = GetVerboseVersion()
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		dbEngineName := dbEngine.GetName()
		result["Ottemo.DBEngine"] = dbEngineName

		if iniConfig := env.GetIniConfig(); iniConfig != nil {

			iniConfigDBKey := "mongodb.db"
			switch dbEngineName {
			case "Sqlite3":
				iniConfigDBKey = "db.sqlite3.uri"
				break
			case "MySQL":
				iniConfigDBKey = "db.mysql.db"
				break
			}

			if iniValue := iniConfig.GetValue(iniConfigDBKey, "ottemo"); iniValue != "" {
				result["Ottemo.DBName"] = iniValue
			}
		}
	}

	result["Ottemo.VersionMajor"] = ConstVersionMajor
	result["Ottemo.VersionMinor"] = ConstVersionMinor
	result["Ottemo.BuildTags"] = buildTags

	result["StartTime"] = startTime
	result["Uptime"] = (time.Now().Truncate(time.Second).Sub(startTime)).Seconds()

	if buildNumber != "" {
		result["Ottemo.BuildNumber"] = utils.InterfaceToInt(buildNumber)
	}
	if buildDate != "" {
		result["Ottemo.BuildDate"] = utils.InterfaceToTime(buildDate).UTC()
	}
	if buildBranch != "" {
		result["Ottemo.BuildBranch"] = buildBranch
	}
	if buildHash != "" {
		result["Ottemo.BuildHash"] = buildHash
	}

	result["GO"] = runtime.Version()
	result["NumGoroutine"] = runtime.NumGoroutine()
	result["NumCPU"] = runtime.NumCPU()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// General statistics.
	result["memStats.Alloc"] = memStats.Alloc           // bytes allocated and still in use
	result["memStats.TotalAlloc"] = memStats.TotalAlloc // bytes allocated (even if freed)
	result["memStats.Sys"] = memStats.Sys               // bytes obtained from system (sum of XxxSys below)
	result["memStats.Lookups"] = memStats.Lookups       // number of pointer lookups
	result["memStats.Mallocs"] = memStats.Mallocs       // number of mallocs
	result["memStats.Frees"] = memStats.Frees           // number of frees

	// Main allocation heap statistics.
	result["memStats.HeapAlloc"] = memStats.HeapAlloc       // bytes allocated and still in use
	result["memStats.HeapSys"] = memStats.HeapSys           // bytes obtained from system
	result["memStats.HeapIdle"] = memStats.HeapIdle         // bytes in idle spans
	result["memStats.HeapInuse"] = memStats.HeapInuse       // bytes in non-idle span
	result["memStats.HeapReleased"] = memStats.HeapReleased // bytes released to the OS
	result["memStats.HeapObjects"] = memStats.HeapObjects   // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	// (Inuse is bytes used now.; Sys is bytes obtained from system.)
	result["memStats.StackInuse"] = memStats.StackInuse // bytes used by stack allocator
	result["memStats.StackSys"] = memStats.StackSys
	result["memStats.MSpanInuse"] = memStats.MSpanInuse // mspan structures
	result["memStats.MSpanSys"] = memStats.MSpanSys
	result["memStats.MCacheInuse"] = memStats.MCacheInuse // mcache structures
	result["memStats.MCacheSys"] = memStats.MCacheSys
	result["memStats.BuckHashSys"] = memStats.BuckHashSys // profiling bucket hash table
	result["memStats.GCSys"] = memStats.GCSys             // GC metadata
	result["memStats.OtherSys"] = memStats.OtherSys       // other system allocations

	// Garbage collector statistics.
	result["memStats.NextGC"] = memStats.NextGC // next collection will happen when HeapAlloc ≥ this amount
	result["memStats.LastGC"] = memStats.LastGC // end time of last collection (nanoseconds since 1970)
	result["memStats.PauseTotalNs"] = memStats.PauseTotalNs
	result["memStats.NumGC"] = memStats.NumGC
	result["memStats.EnableGC"] = memStats.EnableGC
	result["memStats.DebugGC"] = memStats.DebugGC

	return result, nil
}

// getSessionTimeZone return a time zone of session
func getSessionTimeZone(context api.InterfaceApplicationContext) (interface{}, error) {

	return context.GetSession().Get(api.ConstSessionKeyTimeZone), nil
}

// setSessionTimeZone validate time zone and set it to session
func setSessionTimeZone(context api.InterfaceApplicationContext) (interface{}, error) {

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if requestTimeZone := utils.GetFirstMapValue(requestData, "timeZone", "time_zone", "time"); requestTimeZone != nil {

		result, err := SetSessionTimeZone(context.GetSession(), utils.InterfaceToString(requestTimeZone))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		return result, nil
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a33f0f5d-2110-4208-8fb6-023da3ffd241", "time zone should be specified")
}
