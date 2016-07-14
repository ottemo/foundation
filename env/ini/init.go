package ini

import (
	"os"
	"sort"
	"strings"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"

	goini "github.com/vaughan0/go-ini"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultIniConfig)
	var _ env.InterfaceIniConfig = instance

	instance.iniFileValues = make(map[string]map[string]string)
	instance.keysToStore = make(map[string]bool)

	app.OnAppInit(instance.appInitEvent)
	app.OnAppEnd(instance.appEndEvent)

	env.RegisterIniConfig(instance)
}

// routines before application end
func (it *DefaultIniConfig) appEndEvent() error {

	// checking that we have to store ini file
	if len(it.keysToStore) > 0 || it.storeAll {
		// opening ini file
		iniFile, err := os.OpenFile(it.iniFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		defer iniFile.Close()

		if err != nil {
			return env.ErrorDispatch(err)
		}

		// making alphabetically sorted section names
		sortedSections := make([]string, 0, len(it.iniFileValues))
		for sectionName := range it.iniFileValues {
			sortedSections = append(sortedSections, sectionName)
		}
		sort.Strings(sortedSections)

		// loop over alphabetically sorted section names
		for _, sectionName := range sortedSections {
			sectionValues := it.iniFileValues[sectionName]

			// section creation, global section have no header, instead of others
			if sectionName != "" {
				_, err := iniFile.WriteString("\n[" + sectionName + "]\n")
				if err != nil {
					return env.ErrorDispatch(err)
				}
			}

			var storingValueNames []string

			if it.storeAll {
				storingValueNames = make([]string, 0, len(sectionValues))
				for iniItem := range sectionValues {
					storingValueNames = append(storingValueNames, iniItem)
				}
			} else {
				storingValueNames = make([]string, 0, len(it.keysToStore))
				for valueName, store := range it.keysToStore {
					if store {
						storingValueNames = append(storingValueNames, valueName)
					}
				}
			}

			sort.Strings(storingValueNames)

			// loop over alphabetically sorted section values
			for _, key := range storingValueNames {
				if value, present := sectionValues[key]; present {
					_, err := iniFile.WriteString(key + "=" + value + "\n")
					if err != nil {
						env.ErrorDispatch(err)
					}
				}
			}
		}
	}

	return nil
}

// routines before application start (on init phase)
func (it *DefaultIniConfig) appInitEvent() error {

	// checking for environment variable for ini location
	iniFilePath := os.Getenv(ConstEnvironmentIniFile)
	if iniFilePath == "" {
		iniFilePath = ConstDefaultIniFile
	}
	it.iniFilePath = iniFilePath

	it.currentSection = ConstIniGlobalSection
	it.iniFileValues[ConstIniGlobalSection] = make(map[string]string)

	if envSectionName := os.Getenv(ConstEnvironmentIniSection); envSectionName != "" {
		it.currentSection = envSectionName
	}

	// checking command line args for additional parameters
	for _, arg := range os.Args {
		if arg == ConstCmdArgStoreAllFlag {
			it.storeAll = true
		}

		if arg == ConstCmdArgTestFlag {
			it.currentSection = ConstTestSectionName
		}

		if strings.HasPrefix(arg, ConstCmdArgSectionName) {
			argValue := strings.TrimPrefix(arg, ConstCmdArgSectionName)
			it.currentSection = argValue
		}
	}

	// loading values from ini file
	iniFile, _ := goini.LoadFile(iniFilePath)
	for sectionName, sectionValue := range iniFile {
		it.iniFileValues[sectionName] = sectionValue

		// so all the keys we read from file should be stored back
		for valueName := range sectionValue {
			it.keysToStore[valueName] = true
		}
	}

	// firing event for other packages waiting for ini
	err := env.OnConfigIniStart()

	return err
}
