package tests

import (
	"os"
	"strings"
	"sync"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/env/ini"

	// using standard set of packages
	_ "github.com/ottemo/foundation/basebuild"
)

// Package global variables
var (
	startAppFlag  bool
	startAppMutex sync.RWMutex
)

// SwitchToTestIniSection switches ini config to use value from test section instead of general
func SwitchToTestIniSection() error {
	os.Setenv(ini.ConstEnvironmentIniSection, ini.ConstTestSectionName)

	return nil
}

// UpdateWorkingDirectory modifies current working directory to be same for all packages
func UpdateWorkingDirectory() error {

	// was specified environment variable
	if value := os.Getenv("OTTEMO_PATH"); value != "" {
		return os.Chdir(value)
	}

	// for Ottemo internal packages
	workingDirectory, _ := os.Getwd()
	if value := strings.Index(workingDirectory, "/src/github.com/ottemo/foundation"); value > 0 {
		return os.Chdir(workingDirectory[0:value])
	}

	// for other packages
	goPathList := strings.Split(os.Getenv("GOPATH"), ";")
	for _, currentPath := range goPathList {
		if currentPath == "" {
			currentPath = "."
		}

		_, err := os.Stat(currentPath + "/src/github.com/ottemo/foundation")
		if os.IsExist(err) {
			return os.Chdir(currentPath)
		}
	}
	return nil
}

// CheckTestIniDefaults prepares database to be used for tests
func CheckTestIniDefaults() error {

	// we need to init iniConfig before check
	err := app.Init()
	if err != nil {
		return err
	}

	// checking default test mode values
	iniConfig := env.GetIniConfig()
	iniConfig.SetWorkingSection(ini.ConstTestSectionName)

	changesMade := false

	// checking test ini section for sqlite
	if iniConfig.GetSectionValue(ini.ConstTestSectionName, "db.sqlite3.uri", "") == "" {
		iniConfig.SetValue("db.sqlite3.uri", "ottemo_test.db")

		changesMade = true
	}

	// checking test ini section for mongodb
	if iniConfig.GetSectionValue(ini.ConstTestSectionName, "mongodb.uri", "") == "" {
		uriValue := strings.Trim(iniConfig.GetValue("mongodb.uri", "mongodb://localhost:27017/ottemo"), "/") + "_test"
		iniConfig.SetValue("mongodb.uri", uriValue)

		changesMade = true
	}

	if iniConfig.GetSectionValue(ini.ConstTestSectionName, "mongodb.db", "") == "" {
		dbValue := iniConfig.GetValue("mongodb.db", "ottemo") + "_test"
		iniConfig.SetValue("mongodb.db", dbValue)

		changesMade = true
	}

	// checking test ini section for mysql
	if iniConfig.GetSectionValue(ini.ConstTestSectionName, "db.mysql.uri", "") == "" {
		uriValue := iniConfig.GetValue("db.mysql.uri", "/")
		uriValue = uriValue[0 : strings.LastIndex(uriValue, "/")+1]
		iniConfig.SetValue("db.mysql.uri", uriValue)
		changesMade = true
	}

	if iniConfig.GetSectionValue(ini.ConstTestSectionName, "db.mysql.db", "") == "" {
		dbValue := iniConfig.GetValue("db.mysql.db", "ottemo") + "_test"
		iniConfig.SetValue("db.mysql.db", dbValue)

		changesMade = true
	}

	// if ini default values were updated
	if changesMade {
		err = app.End()
		if err != nil {
			return err
		}

		err = app.Init()
		if err != nil {
			return err
		}
	}

	envConfig := env.GetConfig()
	envConfig.SetValue(app.ConstConfigPathMailPort, nil)

	return nil
}

// StartAppInTestingMode starts application in "test mode" (you should use that function for your package tests)
func StartAppInTestingMode() error {
	startAppMutex.Lock()
	defer startAppMutex.Unlock()

	if !startAppFlag {
		err := UpdateWorkingDirectory()
		if err != nil {
			return err
		}

		err = SwitchToTestIniSection()
		if err != nil {
			return err
		}

		err = CheckTestIniDefaults()
		if err != nil {
			return err
		}

		err = app.Start()
		if err != nil {
			return err
		}

		startAppFlag = true
	}

	return nil
}
