package session

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"os"
	"time"
)

// GetName returns implementation name of session service
func (it *DefaultSessionService) GetName() string {
	return it.Storage.GetStorageName()
}

// allocateSessionInstance is a routine to allocate a free memory for a session instance
//   - there is additional check to make sure of allowed sessions amount
func (it *DefaultSessionService) allocateSessionInstance(sessionInstance *DefaultSessionContainer) error {

	if ConstSessionUpdateTime == -1 && ConstSessionKeepInMemoryItems > 0 {
		it.GC()

		numOfSessionsToClean := len(it.Sessions) - ConstSessionKeepInMemoryItems
		for sessionID := range it.Sessions {
			it.Storage.FlushSession(sessionID)

			numOfSessionsToClean--
			if numOfSessionsToClean == 0 {
				break
			}
		}
	}

	sessionID := string(sessionInstance.id)

	it.sessionsMutex.Lock()
	it.Sessions[sessionID] = sessionInstance
	it.sessionsMutex.Unlock()

	if ConstSessionUpdateTime <= 0 {
		it.Storage.FlushSession(sessionID)
	}

	return nil
}

// Get returns session object for given session id or nil of not currently exists
func (it *DefaultSessionService) Get(sessionID string) (api.InterfaceSession, error) {

	replaceInstanceFlag := false

	// taking application instance of session
	sessionInstance, present := it.Sessions[sessionID]
	if present {
		// expiration check
		if time.Now().Sub(sessionInstance.UpdatedAt).Seconds() >= ConstSessionLifeTime {
			sessionInstance = nil
		}
	}

	// session taking from storage for case of "immediate" mode and if no session in memory
	if sessionInstance == nil || ConstSessionUpdateTime == 0 {
		storedInstance, err := it.Storage.LoadSession(sessionID)
		if storedInstance != nil && err == nil {
			// checking that loaded session is newer then we already have
			if sessionInstance == nil || storedInstance.UpdatedAt.After(sessionInstance.UpdatedAt) {
				replaceInstanceFlag = true
				sessionInstance = storedInstance
			}
		}
	}

	// checking if session was found, if not - making new session for given id
	if sessionInstance == nil {
		sessionInstance = new(DefaultSessionContainer)
		sessionInstance.id = DefaultSession(sessionID)
		sessionInstance.Data = make(map[string]interface{})
		sessionInstance.UpdatedAt = time.Now()

		err := it.allocateSessionInstance(sessionInstance)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		replaceInstanceFlag = true
	}

	// updating application sessions
	if replaceInstanceFlag {
		it.sessionsMutex.Lock()
		it.Sessions[sessionID] = sessionInstance
		it.sessionsMutex.Unlock()

		if ConstSessionUpdateTime <= 0 {
			it.Storage.FlushSession(sessionID)
		}
	}

	return sessionInstance.id, nil
}

// New initializes new session instance
func (it *DefaultSessionService) New() (api.InterfaceSession, error) {

	// receiving new session id
	sessionID, err := GenerateSessionID()
	if err != nil {
		return nil, err
	}

	// filling session structure
	sessionInstance := new(DefaultSessionContainer)
	sessionInstance.id = DefaultSession(sessionID)
	sessionInstance.Data = make(map[string]interface{})
	sessionInstance.UpdatedAt = time.Now()

	err = it.allocateSessionInstance(sessionInstance)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return sessionInstance.id, nil
}

// Touch updates session last modification time to current moment
func (it *DefaultSessionService) Touch(sessionID string) error {
	sessionInstance, err := it.Get(sessionID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if sessionInstance != nil {
		if sessionInstance, present := it.Sessions[sessionID]; present {
			sessionInstance.UpdatedAt = time.Now()

			if ConstSessionUpdateTime <= 0 {
				it.Storage.FlushSession(sessionID)
			}
		}
	}

	return nil
}

// Close makes current session instance expired
func (it *DefaultSessionService) Close(sessionID string) error {

	if _, present := it.Sessions[sessionID]; present {

		// removing file
		filename := ConstStorageFolder + sessionID
		if _, err := os.Stat(filename); err == nil {
			err := os.Remove(filename)
			if err != nil {
				env.ErrorDispatch(err)
			}
		}

		// releasing memory
		it.sessionsMutex.Lock()
		delete(it.Sessions, sessionID)
		it.sessionsMutex.Unlock()
	}

	return nil
}

// GetKey returns session value for a given key or nil - if not set
func (it *DefaultSessionService) GetKey(sessionID string, key string) interface{} {
	sessionInstance, err := it.Get(sessionID)
	if sessionInstance == nil || err != nil {
		return nil
	}

	if sessionInstance != nil {
		if sessionInstance, present := it.Sessions[sessionID]; present {

			// looking for a key in session
			if value, present := sessionInstance.Data[key]; present {
				return value
			}
		}
	}
	return nil
}

// SetKey assigns value to session key
func (it *DefaultSessionService) SetKey(sessionID string, key string, value interface{}) {
	sessionInstance, _ := it.Get(sessionID)

	if sessionInstance != nil {
		if sessionInstance, present := it.Sessions[sessionID]; present {
			sessionInstance.Data[key] = value
			sessionInstance.UpdatedAt = time.Now()

			if ConstSessionUpdateTime <= 0 {
				it.Storage.FlushSession(sessionID)
			}
		}
	}
}

// GC is a garbage collector for sessions, it removes expired sessions, flushes to storage, etc.
func (it *DefaultSessionService) GC() error {
	for sessionID, sessionInstance := range it.Sessions {
		secondsAfterLastUpdate := time.Now().Sub(sessionInstance.UpdatedAt).Seconds()

		// closing out of date sessions
		if secondsAfterLastUpdate >= ConstSessionLifeTime {
			it.Close(sessionID)
			continue
		}

		// updating sessions information in a storage
		if secondsAfterLastUpdate > ConstSessionUpdateTime {
			err := it.Storage.FlushSession(sessionID)
			if err != nil {
				env.ErrorDispatch(err)
			}
		}
	}

	return nil
}
