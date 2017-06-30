package vantagepoint

import (
	"strings"
	"regexp"
	"time"

	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/actors/other/vantagepoint/actors"
)

// -------------------------------------------------------------------------------------------------------------------

type envType struct {}

func (it *envType) ErrorDispatch(err error) error {
	return env.ErrorDispatch(err)
}

func (it *envType) ErrorNew(module string, level int, code string, message string) error {
	return env.ErrorNew(module, level, code, message)
}

func (it *envType) LogError(message string) {
	env.Log(ConstLogStorage, env.ConstLogPrefixError, message)
}

func (it *envType) LogWarn(message string) {
	env.Log(ConstLogStorage, env.ConstLogPrefixWarning, message)
}

func (it *envType) LogInfo(message string) {
	env.Log(ConstLogStorage, env.ConstLogPrefixInfo, message)
}

func (it *envType) LogDebug(message string) {
	env.Log(ConstLogStorage, env.ConstLogPrefixDebug, message)
}

// -------------------------------------------------------------------------------------------------------------------

type fileNameType struct {}

func (c *fileNameType) getPattern() string {
	return strings.ToLower("^Fera-(\\d+)-(\\d+)-(\\d+).csv$")
}

func (it *fileNameType) Valid(fileName string) (bool, error) {
	var matched, err = regexp.MatchString(it.getPattern(), strings.ToLower(fileName))
	if err != nil {
		return false, err
	} else if !matched {
		return false, nil
	}

	return true, nil
}

func (it *fileNameType) GetSortValue(fileName string) (string, error) {
	re := regexp.MustCompile(it.getPattern())
	values := re.FindAllStringSubmatch(strings.ToLower(fileName), -1)

	dateStr := values[0][1] + "-" + values[0][2] + "-" + values[0][3]
	fileTime, err := time.Parse("1-2-06", dateStr)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return utils.InterfaceToString(fileTime.Unix()), nil
}


// -------------------------------------------------------------------------------------------------------------------

func CheckNewUploads(params map[string]interface{}) error {
	// not required yet
	_ = params

	env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "check new uploads")
	defer env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "check new uploads done")

	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "383a1377-cf4b-40f9-af4a-dae7e4992fce", "can't obtain config")
	}

	localEnv := envType{}

	path := utils.InterfaceToString(config.GetValue(ConstConfigPathVantagePointUploadPath))
	storagePtr, err := actors.NewDiskStorage(path, &localEnv)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	inventoryProcessorPtr, err := actors.NewInventoryCSV(&localEnv)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	processor, err := actors.NewUploadsProcessor(&localEnv, storagePtr, &fileNameType{}, inventoryProcessorPtr)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err = processor.Process(); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
