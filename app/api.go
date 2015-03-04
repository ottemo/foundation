package app

import (
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

	err = api.GetRestService().RegisterAPI("app/login", api.ConstRESTOperationGet, restLogin)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app/login", api.ConstRESTOperationCreate, restLogin)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app/logout", api.ConstRESTOperationGet, restLogout)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app/rights", api.ConstRESTOperationGet, restRightsInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app/status", api.ConstRESTOperationGet, restStatusInfo)
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

// WEB REST API function to get info about current rights
func restRightsInfo(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})

	result["is_admin"] = utils.InterfaceToBool(context.GetSession().Get(api.ConstSessionKeyAdminRights))

	return result, nil
}

// WEB REST API function to get info about current rights
func restStatusInfo(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})

	result["Ottemo"] = GetVerboseVersion()
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		result["Ottemo.DBEngine"] = dbEngine.GetName()
	}
	result["Ottemo.VersionMajor"] = ConstVersionMajor
	result["Ottemo.VersionMainor"] = ConstVersionMinor
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
