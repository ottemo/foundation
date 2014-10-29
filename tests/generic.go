package tests

import (
	"os"
	"strings"

	"github.com/ottemo/foundation/env/ini"
	"github.com/ottemo/foundation/app"
)

// switches ini config to use value from test section instead of general
func SwitchToTestIniSection() {
	os.Setenv(ini.ENVIRONMENT_INI_SECTION, ini.TEST_SECTION_NAME)
}

// modifies current working directory to be same for all packages
func UpdateWorkingDirectory() error {
	// was specified environment variable
	if value := os.Getenv("OTTEMO_PATH"); value != "" {
		return os.Chdir(value)
	}

	// for Ottemo internal packages
	workingDirectory, _ := os.Getwd()
	if value := strings.Index(workingDirectory, "/src/github.com/ottemo/foundation"); value >0 {
		return os.Chdir(value[0:value])
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

// validates currently using database
func CheckDB() error {

}

// prepares database to be used for tests
func SetupDB() {

}


// you should use that function in your package GO tests to run application and init modules
func StartAppTestingMode() {
	UpdateWorkingDirectory()
	SwitchToTestIniSection()

	app.Start()

	CheckDB()
}
