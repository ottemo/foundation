package ini

import (
	"io"
	"os"
	"sort"
	"strings"

	goini "github.com/vaughan0/go-ini"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultIniConfig)
	var _ env.InterfaceIniConfig = instance

	instance.iniFileValues = make(map[string]map[string]string)
	instance.keysToStore = make(map[string]bool)

	app.OnAppInit(instance.appInitEvent)
	app.OnAppEnd(instance.appEndEvent)

	if err := env.RegisterIniConfig(instance); err != nil {
		_ = env.ErrorDispatch(err)
	}
}

// routines before application end
func (it *DefaultIniConfig) appEndEvent() error {

	// checking that we have to store ini file
	if len(it.keysToStore) > 0 || it.storeAll {

		// Backup ini file.
		// Implemented by function to be sure all file operations are finished.
		// Do not copy empty file. Stop on error and ignore it.
		err := func() error {
			// check source content exists
			srcFileInfo, err := os.Stat(it.iniFilePath)
			if err != nil {
				return env.ErrorDispatch(err)
			}
			if srcFileInfo.Size() <= 0 {
				return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "35c6aaf5-119c-42e7-8c11-bf0650279398", "file size incorrect")
			}

			// open source file for reading
			srcFile, err := os.Open(it.iniFilePath)
			if err != nil {
				return env.ErrorDispatch(err)
			}
			defer func(c io.Closer){
				if err := c.Close(); err != nil {
					_ = env.ErrorDispatch(err)
				}
			}(srcFile)

			// create target file
			tgtFile, err := os.Create(it.iniFilePath + ConstBackupFileSuffix)
			if err != nil {
				return env.ErrorDispatch(err)
			}
			defer func(c io.Closer){
				if err := c.Close(); err != nil {
					_ = env.ErrorDispatch(err)
				}
			}(tgtFile)

			// copy content
			_, err = io.Copy(tgtFile, srcFile)
			if err != nil {
				return env.ErrorDispatch(err)
			}

			// make data persistent as soon as possible
			err = tgtFile.Sync()
			if err != nil {
				return env.ErrorDispatch(err)
			}

			return nil
		}()
		if err != nil {
			_ = env.ErrorDispatch(err)
		}

		// making alphabetically sorted section names
		sortedSections := make([]string, 0, len(it.iniFileValues))
		for sectionName := range it.iniFileValues {
			sortedSections = append(sortedSections, sectionName)
		}
		sort.Strings(sortedSections)

		var output = ""

		// loop over alphabetically sorted section names
		for _, sectionName := range sortedSections {
			sectionValues := it.iniFileValues[sectionName]

			// section creation, global section have no header, instead of others
			if sectionName != "" {
				output += "\n[" + sectionName + "]\n"
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
					output += key + "=" + value + "\n"
				}
			}
		}

		if len(output) > 0 {
			// opening ini file
			iniFile, err := os.OpenFile(it.iniFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			defer func(c io.Closer) {
				if err := c.Close(); err != nil {
					_ = env.ErrorDispatch(err)
				}
			}(iniFile)

			if err != nil {
				return env.ErrorDispatch(err)
			}

			// write whole file content by one operation
			_, err = iniFile.WriteString(output)
			if err != nil {
				return env.ErrorDispatch(err)
			}

			// make data persistent as soon as possible
			err = iniFile.Sync()
			if err != nil {
				return env.ErrorDispatch(err)
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
