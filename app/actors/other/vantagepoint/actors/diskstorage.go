package actors

import (
	"github.com/ottemo/foundation/env"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"compress/gzip"
)

type Disk struct {
	path string
}

func NewDiskStorage(path string) (*Disk, error) {
	disk := &Disk{
		path: path,
	}

	return disk, nil
}

func (s *Disk) ListFiles() ([]string, error) {
	var result = []string{}

	var fileInfos, err = ioutil.ReadDir(s.path)
	if err != nil {
		env.ErrorDispatch(err)
	}

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			result = append(result, fileInfo.Name())
		}
	}

	return result, nil
}

func (s *Disk) Archive(fileName string) error {
	archive := func () error {
		srcFile, err := os.Open(s.getFullFileName(fileName))
		if err != nil {
			env.ErrorDispatch(err)
		}
		defer srcFile.Close()

		tgtFile, err := os.OpenFile(s.getFullFileName(fileName + ".gz"), os.O_CREATE | os.O_WRONLY, 0660)
		if err != nil {
			env.ErrorDispatch(err)
		}
		defer tgtFile.Close()

		archiver, err := gzip.NewWriterLevel(tgtFile, gzip.BestCompression)
		if err != nil {
			env.ErrorDispatch(err)
		}

		_, err = io.Copy(archiver, srcFile)
		if err != nil {
			env.ErrorDispatch(err)
		}

		return archiver.Close()
	}

	if err := archive(); err != nil {
		env.ErrorDispatch(err)
	}

	// TODO: uncomment
	//if err := os.Remove(s.getFullFileName(fileName)); err != nil {
	//	env.ErrorDispatch(err)
	//}

	return nil
}

func (s *Disk) GetReadCloser(fileName string) (io.ReadCloser, error) {
	var file, err = os.Open(s.getFullFileName(fileName))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return file, nil
}

func (s *Disk) getFullFileName(fileName string) string {
	return filepath.Join(s.path, fileName)
}
