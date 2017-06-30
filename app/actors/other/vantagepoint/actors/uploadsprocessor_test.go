package actors

import (
	"testing"
	"github.com/ottemo/foundation/env"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"io"
)

type testConfig struct {}

func (c *testConfig) RegisterItem(Item env.StructConfigItem, Validator env.FuncConfigValueValidator) error {
	return nil
}
func (c *testConfig) UnregisterItem(Path string) error {
	return nil
}

func (c *testConfig) ListPathes() []string {
	return []string{}
}
func (c *testConfig) GetValue(Path string) interface{} {
	switch Path {
	case "general.vantagepoint.enabled":
		return true
	}

	return nil
}
func (c *testConfig) SetValue(Path string, Value interface{}) error {
	return nil
}

func (c *testConfig) GetGroupItems() []env.StructConfigItem {
	return []env.StructConfigItem{}
}
func (c *testConfig) GetItemsInfo(Path string) []env.StructConfigItem {
	return []env.StructConfigItem{}
}

func (c *testConfig) Load() error {
	return nil
}
func (c *testConfig) Reload() error {
	return nil
}

var config = &testConfig{}

// --------------------------------------------------------------------------------------------------------------------

type testEnv struct {}

func (e *testEnv) ErrorDispatch(err error) error {
	fmt.Println("ERROR", err)
	return err
}
func (e *testEnv) ErrorNew(module string, level int, code string, message string) error {
	fmt.Println(code, message)
	return errors.New(code+" "+message)
}
func (e *testEnv) GetConfig() env.InterfaceConfig {
	return config
}

func (it *testEnv) LogError(message string) {
	fmt.Println("LOG ERROR", message)
}

func (it *testEnv) LogWarn(message string) {
	fmt.Println("LOG WARn", message)
}

func (it *testEnv) LogInfo(message string) {
	fmt.Println("LOG INfo", message)
}

func (it *testEnv) LogDebug(message string) {
	fmt.Println("LOG Debug", message)
}

// --------------------------------------------------------------------------------------------------------------------

var testReadClosed int

type testReadCloser struct {
	io.Reader
}

func (c *testReadCloser) Read(b []byte) (n int, err error){
	n, err = strings.NewReader("String value").Read(b)
	return n, err
}

func (c *testReadCloser) Close() error {
	testReadClosed++
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

type testStorage struct {}

func (s *testStorage) ListFiles() ([]string, error) {
	return []string{"Prefix-1-25-18.csv", "Prefix-1-24-17.csv", "Prefix-1-26-17.csv", "Prefix-2-26-16.csv", "Prefix-1-26-17.xlsx"}, nil
}
func (s *testStorage) Archive(fileName string) error {
	return nil
}
func (s *testStorage) GetReadCloser(fileName string) (io.ReadCloser, error) {
	return &testReadCloser{}, nil
}

// --------------------------------------------------------------------------------------------------------------------

type testFileName struct {}

func (c *testFileName) getPattern() string {
	return strings.ToLower("^Prefix-(\\d+)-(\\d+)-(\\d+).csv$")
}

func (c *testFileName) Valid(fileName string) (bool, error) {
	var matched, err = regexp.MatchString(c.getPattern(), strings.ToLower(fileName))
	if err != nil {
		return false, err
	} else if !matched {
		return false, nil
	}

	return true, nil
}

func (c *testFileName) GetSortValue(fileName string) (string, error) {
	return fileName, nil
}

// --------------------------------------------------------------------------------------------------------------------

type testDataProcessor struct {
	processed int
}

func (d *testDataProcessor) Process(reader io.Reader) error {
	d.processed++
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func TestPrepareFileNames(t *testing.T) {
	processor, err := NewUploadsProcessor(&testEnv{}, &testStorage{}, &testFileName{}, &testDataProcessor{})
	if err != nil {
		t.Fatal(err)
	}

	fileNames, err := processor.prepareFileNames()
	if err != nil {
		t.Fatal(err)
	}

	expectFileNames := []string{"Prefix-1-24-17.csv", "Prefix-1-25-18.csv", "Prefix-1-26-17.csv", "Prefix-2-26-16.csv"}
	if !isStringSlicesEqual(fileNames, expectFileNames) {
		t.Log(fileNames)
		t.Fatal("Prepared fileNames wrong")
	}
}

func TestProcessFiles(t *testing.T) {
	dataProcessor := testDataProcessor{}

	processor, err := NewUploadsProcessor(&testEnv{}, &testStorage{}, &testFileName{}, &dataProcessor)
	if err != nil {
		t.Fatal(err)
	}

	fileNames := []string{"Prefix-2-26-16.csv", "Prefix-1-24-17.csv", "Prefix-1-26-17.csv", "Prefix-1-25-18.csv"}
	for _, fileName := range(fileNames) {
		if err := processor.processFile(fileName); err != nil {
			t.Fatal(err)
		}
	}

	if dataProcessor.processed != len(fileNames) {
		t.Fatal("Not all files processed")
	}

	if testReadClosed != len(fileNames) {
		t.Fatal("Not all files closed")
	}
}

func isStringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for idx, valueA := range(a) {
		if (idx+1) > len(b) {
			return false
		} else if valueA != b[idx] {
			return false
		}
	}

	return true
}
