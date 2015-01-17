package session

import (
	"crypto/rand"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"os"
	"time"
)

// updateSession reloads session data from storage if it newer
//   - you can send sessionID or session instance to function, if have no session instance use nil
func (it *DefaultSessionService) updateSession(sessionID string, sessionInstance *DefaultSession) error {
	if sessionInstance != nil {
		sessionID = sessionInstance.GetID()
	} else {
		if value, present := it.Sessions[sessionID]; present {
			sessionInstance = value
		} else {
			sessionInstance = new(DefaultSession)
		}
	}

	fileName := ConstStorageFolder + sessionID
	stats, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	if stats.ModTime().After(sessionInstance.UpdatedAt) {
		sessionInstance.Load(sessionID)
	}

	return nil
}

// flushSession writes session data to shared storage and forgets about
//   - you can send sessionID or session instance to function, if have no session instance use nil
func (it *DefaultSessionService) flushSession(sessionID string, sessionInstance *DefaultSession) error {
	if sessionInstance != nil {
		sessionID = sessionInstance.GetID()
	}

	if value, present := it.Sessions[sessionID]; present {
		sessionInstance = value

		it.sessionsMutex.Lock()
		delete(it.Sessions, sessionID)
		it.sessionsMutex.Unlock()
	}

	if sessionInstance != nil {
		return sessionInstance.Save()
	}

	return nil
}

// GetName returns implementation name of session service
func (it *DefaultSessionService) GetName() string {
	return "DefaultSessionService"
}

// Get returns session object for given session id or nil of not currently exists
func (it *DefaultSessionService) Get(sessionID string) (api.InterfaceSession, error) {
	if sessionInstance, present := it.Sessions[sessionID]; present {
		if time.Now().Sub(sessionInstance.UpdatedAt).Seconds() < ConstSessionLifeTime {
			return sessionInstance, nil
		}
	}

	sessionInstance := new(DefaultSession)
	err := sessionInstance.Load(sessionID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return sessionInstance, nil
}

// New initializes new session instance
func (it *DefaultSessionService) New() (api.InterfaceSession, error) {

	// "keep in memory" mode - checking for allowed amount of sessions
	if ConstSessionUpdateTime == -1 && ConstSessionKeepInMemoryItems > 0 {
		it.gc()

		numOfSessionsToClean := len(it.Sessions) - ConstSessionKeepInMemoryItems
		for sessionID, sessionInstance := range it.Sessions {
			it.flushSession(sessionID, sessionInstance)

			numOfSessionsToClean--
			if numOfSessionsToClean == 0 {
				break
			}
		}
	}

	// receiving new session id
	sessionID, err := it.generateSessionID()
	if err != nil {
		return nil, err
	}

	sessionInstance := new(DefaultSession)
	sessionInstance.id = sessionID
	sessionInstance.Data = make(map[string]interface{})
	sessionInstance.UpdatedAt = time.Now()

	it.Sessions[sessionID] = sessionInstance

	return it.Sessions[sessionID], nil
}

// generateSessionID returns new session id number
func (it *DefaultSessionService) generateSessionID() (string, error) {
	sessionID := make([]byte, 32)
	if _, err := rand.Read(sessionID); err != nil {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "aab09114-9726-4844-984c-772fb25dcb88", "can't generate sessionID")
	}

	for i := 0; i < 32; i++ {
		sessionID[i] = ALPHANUMERIC[sessionID[i]%62]
	}

	return string(sessionID), nil
}

// gc is a garbage collector for sessions, it removes expired sessions, flushes to storage, etc.
func (it *DefaultSessionService) gc() {
	for sessionID, sessionInstance := range it.Sessions {
		secondsAfterLastUpdate := time.Now().Sub(sessionInstance.UpdatedAt).Seconds()

		// closing out of date sessions
		if secondsAfterLastUpdate > ConstSessionLifeTime {
			sessionInstance.Close()
			continue
		}

		// updating sessions information in a storage
		if secondsAfterLastUpdate > ConstSessionUpdateTime {
			err := it.flushSession(sessionID, nil)
			if err != nil {
				env.ErrorDispatch(err)
			}
		}
	}
}
