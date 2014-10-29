package tests

import (
	"os"
	"strings"

	"github.com/ottemo/foundation/env"
)

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
	env.IniValue()
}

// prepares database to be used for tests
func SetupDB() {

}
