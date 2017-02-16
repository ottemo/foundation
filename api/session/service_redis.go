// +build redis

// "service_redis.go" is a redis based session storage - "redis" build tag should be specified in order to use it

package session

import (
	"bytes"
	"encoding/json"
	"github.com/fiorix/go-redis/redis"
	"io"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// RedisSessionService is a memcache based storage implementer based on "DefaultSessionService"
type RedisSessionService struct {
	redisClient *redis.Client
	*DefaultSessionService
}

// init makes package self-initialization routine
func init() {

	redisService := new(RedisSessionService)
	redisService.DefaultSessionService = InitDefaultSessionService()
	redisService.DefaultSessionService.storage = redisService

	SessionService = redisService

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

// Startup is a RedisSessionService initialization routines
func startup() error {

	redisService, ok := SessionService.(*RedisSessionService)
	if !ok {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "281f883d-1ec8-47d4-a916-0107e9de997c", "unexpected session service instance")
	}

	serversList := "127.0.0.1:6379"

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("redis.servers", serversList); iniValue != "" {
			serversList = iniValue
		}
	}

	redisService.redisClient = redis.New(serversList)
	error := redisService.redisClient.Ping()

	return error
}

// shutdown is a RedisSessionService shutdown routines
func shutdown() error {
	return nil
}

// InterfaceServiceStorage implementation
// --------------------------------------

// GetStorageName returns storage implementation name for a session service
func (it *RedisSessionService) GetStorageName() string {
	return "RedisSessionService"
}

// LoadSession de-serializes file from memcache server, returns nil on error
func (it *RedisSessionService) LoadSession(sessionID string) (*DefaultSessionContainer, error) {

	// making new session holder instance
	sessionInstance := &DefaultSessionContainer{id: sessionID}

	// checking existence
	item, err := it.redisClient.Get(sessionID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// reading and un-serializing data
	var reader io.Reader = bytes.NewBufferString(item)
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
		it.redisClient.Expire(sessionID, ConstSessionLifeTime)
		sessionInstance.UpdatedAt = time.Now()
	}

	return sessionInstance, nil
}

// FlushSession stores serialized session on memcache server
func (it *RedisSessionService) FlushSession(sessionID string) error {
	// checking session existence
	sessionInstance := it.syncGet(sessionID)
	if sessionInstance == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5ec1d9b8-48e4-43e7-a1b6-a658a576c286", "session "+sessionID+" not found")
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

	value := buffer.String()
	sessionInstance.mutex.Unlock()

	// storing session item
	it.redisClient.SetEx(sessionID, ConstSessionLifeTime, value)

	// releasing application memory
	it.syncDel(sessionID)

	return nil
}
