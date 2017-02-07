// +build !redis,!memcache,!memsession

// "service_filesystem.go" is a filesystem based session storage implementation - default option if no tags specified

package session

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// FilesystemSessionService is a filesystem based storage implementer based on "DefaultSessionService"
type FilesystemSessionService struct {
	*DefaultSessionService
}

// init makes package self-initialization routine
func init() {

	filesystemService := new(FilesystemSessionService)
	filesystemService.DefaultSessionService = InitDefaultSessionService()
	filesystemService.DefaultSessionService.storage = filesystemService

	SessionService = filesystemService

	// starting timer if session update time specified and service supports garbage collection
	if ConstSessionUpdateTime > 0 {
		timerInterval := time.Second * ConstSessionUpdateTime
		ticker := time.NewTicker(timerInterval)
		go func() {
			for _ = range ticker.C {
				if err := SessionService.GC(); err != nil {
					_ = env.ErrorDispatch(err)
				}
			}
		}()
	}

	// service registration within system
	if err := api.RegisterSessionService(SessionService); err != nil {
		_ = env.ErrorDispatch(err)
	}

	app.OnAppStart(startup)
	app.OnAppEnd(shutdown)
}

// Startup is a FilesystemSessionService initialization routines
func startup() error {

	// checking session storage directory existence, creating if not exists
	if _, err := os.Stat(ConstStorageFolder); !os.IsExist(err) {
		err := os.MkdirAll(ConstStorageFolder, os.ModePerm)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	// collecting files within session storage folder
	files, err := ioutil.ReadDir(ConstStorageFolder)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	currentTime := time.Now()

	// expired sessions cleanup, sessions loading
	for _, fileInfo := range files {

		// removing expired sessions
		if currentTime.Sub(fileInfo.ModTime()).Seconds() >= ConstSessionLifeTime {
			err := os.Remove(ConstStorageFolder + fileInfo.Name())
			if err != nil {
				_ = env.ErrorDispatch(err)
			}
			continue
		}
	}

	return nil
}

// Shutdown is a FilesystemSessionService shutdown routines
func shutdown() error {

	filesystemService, ok := SessionService.(*FilesystemSessionService)
	if !ok {
		err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a95b6e53-1fb5-480a-a8c6-a6169a7da9fa", "unexpected session service instance")
		return env.ErrorDispatch(err)
	}

	currentTime := time.Now()

	// saving all session to storage
	filesystemService.syncLoop(
		func(sessionInstance *DefaultSessionContainer) bool {
			// session expiration check
			if currentTime.Sub(sessionInstance.UpdatedAt).Seconds() >= ConstSessionLifeTime {
				return false
			}

			// flushing session
			if err := filesystemService.FlushSession(sessionInstance.id); err != nil {
				_ = env.ErrorDispatch(err)
			}
			return false
		})

	return nil
}

// InterfaceServiceStorage implementation
// --------------------------------------

// GetStorageName returns storage implementation name for a session service
func (it *FilesystemSessionService) GetStorageName() string {
	return "MemcacheSessionService"
}

// LoadSession de-serializes file from filesystem storage, returns nil on error
func (it *FilesystemSessionService) LoadSession(sessionID string) (*DefaultSessionContainer, error) {

	// making new session holder instance
	sessionInstance := &DefaultSessionContainer{id: sessionID}

	// checking file exists in file system
	filename := ConstStorageFolder + sessionID
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "363cd5a8-1a3d-4163-a7d3-cb96dbaff01c", "session "+sessionID+" not found")
	}

	// checking file modification time - expired session case
	if time.Now().Sub(fileInfo.ModTime()).Seconds() >= ConstSessionLifeTime {
		err := os.Remove(filename)
		if err != nil {
			return nil, err
		}
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "7aee9352-a08b-420a-a725-7f32a17495a8", "session "+sessionID+" expired")
	}

	// file not expired - loading data from it
	sessionFile, err := os.OpenFile(filename, os.O_RDONLY, 0660)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// defer file close with error handling
	defer func(c io.Closer){
		if err := sessionFile.Close(); err != nil {
			_ = env.ErrorDispatch(err)
		}
	}(sessionFile)

	var reader io.Reader = sessionFile
	if ConstCryptSession {
		reader, err = utils.EncryptReader(reader)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	jsonDecoder := json.NewDecoder(reader)
	err = jsonDecoder.Decode(sessionInstance)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if sessionInstance.Data == nil {
		sessionInstance.Data = make(map[string]interface{})
	}

	if sessionInstance.UpdatedAt.IsZero() {
		sessionInstance.UpdatedAt = fileInfo.ModTime()
	}

	return sessionInstance, nil
}

// FlushSession serializes session into filesystem storage
//   - routine not checks session expiration or modification time - it just flushes data to storage
func (it *FilesystemSessionService) FlushSession(sessionID string) error {
	sessionInstance := it.syncGet(sessionID)
	if sessionInstance == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "363cd5a8-1a3d-4163-a7d3-cb96dbaff01c", "session "+sessionID+" not found")
	}

	// skipping flush for empty sessions
	if SessionService.IsEmpty(sessionID) {
		return nil
	}

	// serializing session data to file
	sessionFile, err := os.OpenFile(ConstStorageFolder+sessionID, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	defer func() {
		if err := sessionFile.Close(); err != nil {
			_ = env.ErrorDispatch(err)
		}
		updatedAt := sessionInstance.GetUpdatedAt()
		if err := os.Chtimes(sessionFile.Name(), updatedAt, updatedAt); err != nil {
			_ = env.ErrorDispatch(err)
		}
	}()

	var writer io.Writer = sessionFile
	if ConstCryptSession {
		writer, err = utils.EncryptWriter(writer)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	sessionInstance.mutex.Lock()
	jsonEncoder := json.NewEncoder(writer)
	err = jsonEncoder.Encode(sessionInstance)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	sessionInstance.mutex.Unlock()

	// releasing application memory
	it.syncDel(sessionID)

	return nil
}
