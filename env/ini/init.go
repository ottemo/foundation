package ini

import (
	"os"
	"sort"
	"strings"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"

	goini "github.com/vaughan0/go-ini"
)

const (
	CMD_ARG_STORE_ALL_FLAG = "--iniStoreAll"
	CMD_ARG_SECTION_NAME = "--iniSection="
	CMD_ARG_TEST_FLAG = "--test"

	ENVIRONMENT_INI_FILE = "OTTEMO_INI"
	ENVIRONMENT_INI_SECTION = "OTTEMO_MODE"

	TEST_SECTION_NAME = "test"
	DEFAULT_INI_FILE = "ottemo.ini"
)

// module entry point before app start
func init() {
	instance := new(DefaultIniConfig)
	var _ env.I_IniConfig = instance

	instance.iniFileValues = make(map[string]map[string]string)
	instance.keysToStore = make([]string, 0)

	app.OnAppInit(instance.appInitEvent)
	app.OnAppEnd(instance.appEndEvent)

	env.RegisterIniConfig(instance)
}

// routines before application end
func (it *DefaultIniConfig) appEndEvent() error {

	// checking that we have to store ini file
	if len(it.keysToStore) > 0 || it.storeAll {
		// opening ini file
		iniFile, err := os.OpenFile(it.iniFilePath, os.O_CREATE | os.O_WRONLY, 0644)
		defer iniFile.Close()

		if err != nil {
			return env.ErrorDispatch(err)
		}

		// making alphabetically sorted section names
		sortedSections := make([]string, 0, len(it.iniFileValues))
		for sectionName, _ := range it.iniFileValues {
			sortedSections = append(sortedSections, sectionName)
		}
		sort.Strings(sortedSections)

		// loop over alphabetically sorted section names
		for _, sectionName := range sortedSections {
			sectionValues := it.iniFileValues[sectionName]

			storingKeys := it.keysToStore
			if it.storeAll {
				allKeys := make([]string, 0, len(sectionValues))
				for iniItem, _ := range sectionValues {
					allKeys = append(allKeys, iniItem)
				}
				storingKeys = allKeys
			}
			sort.Strings(storingKeys)

			// loop over alphabetically sorted section values
			for _, key := range it.keysToStore {
				if value, present := sectionValues[key]; present {
					_, err := iniFile.WriteString(key + "=" + value + "\n")
					if err != nil {
						env.ErrorDispatch(err)
					}
				}
			}

			// global section have no header, instead of others
			if sectionName != "" {
				_, err := iniFile.WriteString("\n[" + sectionName + "]\n")
				if err != nil {
					return env.ErrorDispatch(err)
				}
			}
		}
	}

	return nil
}

// routines before application start (on init phase)
func (it *DefaultIniConfig) appInitEvent() error {

	// checking for environment variable for ini location
	iniFilePath := os.Getenv(ENVIRONMENT_INI_FILE)
	if iniFilePath == "" {
		iniFilePath = DEFAULT_INI_FILE
	}
	it.iniFilePath = iniFilePath

	if envSectionName := os.Getenv(ENVIRONMENT_INI_SECTION); envSectionName != "" {
		it.currentSection = envSectionName
	}

	// checking command line args for additional parameters
	for _, arg := range os.Args {
		if arg == CMD_ARG_STORE_ALL_FLAG {
			it.storeAll = true
		}

		if arg == CMD_ARG_TEST_FLAG {
			it.currentSection = TEST_SECTION_NAME
		}

		if strings.HasPrefix(arg, CMD_ARG_SECTION_NAME) {
			argValue := strings.TrimPrefix(arg, CMD_ARG_SECTION_NAME)
			it.currentSection = argValue
		}
	}

	// loading values from ini file
	iniFile, _ := goini.LoadFile(iniFilePath)
	for sectionName, sectionValue := range iniFile {
		it.iniFileValues[sectionName] = sectionValue
	}

	// firing event for other packages waiting for ini
	err := env.OnConfigIniStart()

	return err
}
