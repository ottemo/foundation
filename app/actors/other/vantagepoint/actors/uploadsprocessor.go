package actors

import (
	"sort"
	"io"
)

// --------------------------------------------------------------------------------------------------------------------

type FileNameInterface interface {
	Valid(fileName string) (bool, error)
	GetSortValue(fileName string) (string, error)
}

type EnvInterface interface {
	ErrorDispatch(err error) error
	ErrorNew(module string, level int, code string, message string) error
}

type StorageInterface interface {
	ListFiles() ([]string, error)
	Archive(fileName string) error
	GetReadCloser(fileName string) (io.ReadCloser, error)
}

type DataProcessorInterface interface {
	Process(reader io.Reader) error
}

// --------------------------------------------------------------------------------------------------------------------

type FileNamesSorter struct {
	Items []string
	FileName FileNameInterface

	Err error
}

func (fn *FileNamesSorter) Len() int {
	return len(fn.Items)
}

func (fn *FileNamesSorter) Swap(i, j int) {
	fn.Items[i], fn.Items[j] = fn.Items[j], fn.Items[i]
}

func (fn *FileNamesSorter) Less(i, j int) bool {
	iValue, err := fn.FileName.GetSortValue(fn.Items[i])
	if err != nil {
		fn.Err = err
		return false
	}
	jValue, err := fn.FileName.GetSortValue(fn.Items[j])
	if err != nil {
		fn.Err = err
		return false
	}

	return iValue < jValue
}

// --------------------------------------------------------------------------------------------------------------------

type uploadsProcessor struct {
	env      EnvInterface
	storage  StorageInterface
	fileName FileNameInterface
	dataProcessor DataProcessorInterface
}

func NewUploadsProcessor(env EnvInterface, storage StorageInterface, fileName FileNameInterface, dataProcessor DataProcessorInterface) (uploadsProcessor, error) {
	var newActor = uploadsProcessor{
		env : env,
		storage: storage,
		fileName: fileName,
		dataProcessor: dataProcessor,
	}

	return newActor, nil
}

func (a *uploadsProcessor) Process() error {
	fileNames, err := a.prepareFileNames()
	if err != nil {
		return a.env.ErrorNew(ConstErrorModule, ConstErrorLevel, "40f3e176-337d-4055-a4db-dfc200820a13", err.Error())
	}

	for _, fileName := range(fileNames) {
		if err = a.processFile(fileName); err != nil {
			return a.env.ErrorNew(ConstErrorModule, ConstErrorLevel, "606b6a0b-ba97-44ce-b901-2a3deba64482", err.Error())
		}

		// TODO: check config period setting
		if err := a.storage.Archive(fileName); err != nil {
			return a.env.ErrorDispatch(err)
		}
	}

	return nil
}

func (a *uploadsProcessor) prepareFileNames() ([]string, error) {
	// get file names
	fileNames, err := a.storage.ListFiles()
	if err != nil {
		return []string{}, err
	}

	// filter
	filteredFileNames := []string{}
	for _, fileName := range(fileNames) {
		valid, err := a.fileName.Valid(fileName)
		if err != nil {
			return []string{}, a.env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8237fc61-180d-401f-9ef5-2c69b6e32de1", err.Error())
		} else if valid {
			filteredFileNames = append(filteredFileNames, fileName)
		}
	}

	// sort
	fileNamesSorter := &FileNamesSorter{
		Items: filteredFileNames,
		FileName: a.fileName,
	}
	sort.Sort(fileNamesSorter)
	if fileNamesSorter.Err != nil {
		return []string{}, a.env.ErrorDispatch(fileNamesSorter.Err)
	}

	return fileNamesSorter.Items, nil
}

func (a *uploadsProcessor) processFile(fileName string) error {
	readCloser, err := a.storage.GetReadCloser(fileName)
	if err != nil {
		return a.env.ErrorDispatch(err)
	}
	defer readCloser.Close()

	return a.dataProcessor.Process(readCloser)
}
