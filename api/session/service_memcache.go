// +build memcache

// "service_memcache.go" is a memcache based session storage implementation - "memcache" build tag should be specified in order to use it

package session

import (
	"bytes"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"io"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// MemcacheSessionService is a memcache based storage implementer based on "DefaultSessionService"
type MemcacheSessionService struct {
	memcacheClient *memcache.Client
	*DefaultSessionService
}

// init makes package self-initialization routine
func init() {

	memcacheService := new(MemcacheSessionService)
	memcacheService.DefaultSessionService = InitDefaultSessionService()
	memcacheService.DefaultSessionService.storage = memcacheService

	SessionService = memcacheService

	// starting timer if session update time specified and service supports garbage collection
	if ConstSessionUpdateTime > 0 {
		timerInterval := time.Second * ConstSessionUpdateTime
		ticker := time.NewTicker(timerInterval)
		go func() {
			for _ = range ticker.C {
				SessionService.GC()
			}
		}()
	}

	// service registration within system
	api.RegisterSessionService(SessionService)

	app.OnAppStart(startup)
	app.OnAppEnd(shutdown)
}

// Startup is a MemcacheSessionService initialization routines
func startup() error {

	memcacheService, ok := SessionService.(*MemcacheSessionService)
	if !ok {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "846e84a1-06cf-4dda-b441-5ebb6c7b3413", "unexpected session service instance")
	}

	serversList := "127.0.0.1:11211"

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("memcache.servers", serversList); iniValue != "" {
			serversList = iniValue
		}
	}

	memcacheService.memcacheClient = memcache.New(serversList)

	return nil
}

// shutdown is a MemcacheSessionService shutdown routines
func shutdown() error {
	return nil
}

// InterfaceServiceStorage implementation
// --------------------------------------

// GetStorageName returns storage implementation name for a session service
func (it *MemcacheSessionService) GetStorageName() string {
	return "MemcacheSessionService"
}

// LoadSession de-serializes file from memcache server, returns nil on error
func (it *MemcacheSessionService) LoadSession(sessionID string) (*DefaultSessionContainer, error) {

	// making new session holder instance
	sessionInstance := &DefaultSessionContainer{id: sessionID}

	// checking existence
	item, err := it.memcacheClient.Get(sessionID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// reading and un-serializing data
	var reader io.Reader = bytes.NewReader(item.Value)
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

	// checking that session instance object is ok (fixing if not)
	if sessionInstance.Data == nil {
		sessionInstance.Data = make(map[string]interface{})
	}

	if sessionInstance.UpdatedAt.IsZero() {
		sessionInstance.UpdatedAt = time.Now().Add(time.Duration(-ConstSessionLifeTime + item.Expiration))
	}

	return sessionInstance, nil
}

// FlushSession stores serialized session on memcache server
func (it *MemcacheSessionService) FlushSession(sessionID string) error {
	// checking session existence
	sessionInstance := it.syncGet(sessionID)
	if sessionInstance == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5eabd4f1-2ebb-4601-92de-188a6ef7da4e", "session "+sessionID+" not found")
	}

	// skipping flush for empty sessions
	if SessionService.IsEmpty(sessionID) {
		return nil
	}

	// preparing session data writer
	var buffer bytes.Buffer

	var err error
	var writer io.Writer = &buffer
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
	value := buffer.Bytes()
	sessionInstance.mutex.Unlock()

	// storing session item
	item := &memcache.Item{Key: sessionID, Value: value, Expiration: ConstSessionLifeTime}
	it.memcacheClient.Set(item)

	// releasing application memory
	it.syncDel(sessionID)

	return nil
}
