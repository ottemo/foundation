package vantagepoint

import (
	"github.com/ottemo/foundation/app/actors/other/vantagepoint/actors"
	"github.com/ottemo/foundation/env"
	"strings"
	"regexp"
	"time"
	"github.com/ottemo/foundation/utils"
	"io"
	"fmt"
)

// -------------------------------------------------------------------------------------------------------------------

type envType struct {}

func (it *envType) ErrorDispatch(err error) error {
	return env.ErrorDispatch(err)
}

func (it *envType) ErrorNew(module string, level int, code string, message string) error {
	return env.ErrorNew(module, level, code, message)
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

type tmpDataProcessor struct {}

func (it *tmpDataProcessor) Process(reader io.Reader) error {
	fmt.Println("tmpDataProcessor) Process")
	return nil
}

// -------------------------------------------------------------------------------------------------------------------

func CheckNewUploads(params map[string]interface{}) error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "383a1377-cf4b-40f9-af4a-dae7e4992fce", "can't obtain config")
	}

	if !utils.InterfaceToBool(config.GetValue(ConstConfigPathVantagePointEnabled)) && false { //TODO remove false
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "40f3e176-337d-4055-a4db-dfc200820a13", "VantagePoint CheckNewUploads called but not enabled")
		return nil
	}

	//func NewUploadsProcessor(env EnvInterface, storage StorageInterface, fileName FileNameInterface, dataProcessor DataProcessorInterface) (uploadsProcessor, error) {
	// TODO: use config value
	var path = "./vantagepoint/" //utils.InterfaceToBool(config.GetValue(ConstConfigPathVantagePointUploadPath))
	storagePtr, err := actors.NewDiskStorage(path)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	processor, err := actors.NewUploadsProcessor(&envType{}, storagePtr, &fileNameType{}, &tmpDataProcessor{})
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err = processor.Process(); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
