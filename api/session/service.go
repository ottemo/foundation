package session

import (
	"crypto/rand"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"os"
	"time"
)

const (
	alphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890" // sessionID allowed symbols
)

// InitDefaultSessionService makes a new instance of DefaultSessionService
//   - makes internal fields initialization
func InitDefaultSessionService() *DefaultSessionService {
	sessionService := new(DefaultSessionService)
	sessionService.sessions = make(map[string]*DefaultSessionContainer)
	sessionService.storage = sessionService

	return sessionService
}

// GenerateSessionID returns new session id number
func GenerateSessionID() (string, error) {
	sessionID := make([]byte, 32)
	if _, err := rand.Read(sessionID); err != nil {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "aab09114-9726-4844-984c-772fb25dcb88", "can't generate sessionID")
	}

	for i := 0; i < 32; i++ {
		sessionID[i] = alphanumeric[sessionID[i]%62]
	}

	return string(sessionID), nil
}

// synchronized access to session container
// ----------------------------------------

func (it *DefaultSessionContainer) Set(key string, value interface{}) {
	it.mutex.Lock()
	it.Data[key] = value
	it.mutex.Unlock()
}

func (it *DefaultSessionContainer) Get(key string) interface{} {
	defer it.mutex.Unlock()
	it.mutex.Lock()
	if value, present := it.Data[key]; present {
		return value
	}
	return nil
}

func (it *DefaultSessionContainer) GetUpdatedAt() time.Time {
	defer it.mutex.Unlock()
	it.mutex.Lock()
	return it.UpdatedAt
}

func (it *DefaultSessionContainer) SetUpdatedAt(value time.Time) {
	it.mutex.Lock()
	it.UpdatedAt = value
	it.mutex.Unlock()
}

func (it *DefaultSessionContainer) GetID() string {
	defer it.mutex.Unlock()
	it.mutex.Lock()
	return it.id
}

func (it *DefaultSessionContainer) SetID(value string) {
	it.mutex.Lock()
	it.id = value
	it.mutex.Unlock()
}

// synchronized access to sessions map
// -----------------------------------

// syncCount returns amount of sessions map items
func (it *DefaultSessionService) syncCount() int {
	defer it.mutex.Unlock()
	it.mutex.Lock()
	return len(it.sessions)
}

// syncSet updates sessions map item for a specified session id
func (it *DefaultSessionService) syncSet(id string, session *DefaultSessionContainer) {
	it.mutex.Lock()
	it.sessions[id] = session
	it.mutex.Unlock()
}

// syncSet removes sessions map item by specified session id
func (it *DefaultSessionService) syncDel(id string) {
	it.mutex.Lock()
	if _, present := it.sessions[id]; present {
		delete(it.sessions, id)
	}
	it.mutex.Unlock()
}

// syncGet returns sessions map item by specified session id
func (it *DefaultSessionService) syncGet(id string) *DefaultSessionContainer {
	defer it.mutex.Unlock()
	it.mutex.Lock()
	if session, present := it.sessions[id]; present {
		return session
	}
	return nil
}

// syncLoop executes [action] for each sessions map item
func (it *DefaultSessionService) syncLoop(action func(*DefaultSessionContainer) bool) {
	it.mutex.Lock()
	var ids []string
	for id := range it.sessions {
		ids = append(ids, id)
	}
	it.mutex.Unlock()

	for _, id := range ids {
		if session := it.syncGet(id); session != nil {
			if action(session) {
				break
			}
		}
	}
}

// InterfaceServiceStorage implementation
// --------------------------------------

// GetStorageName returns storage implementation name for a session service
func (it *DefaultSessionService) GetStorageName() string {
	return "DefaultSessionService"
}

// LoadSession is a stub function for no action
func (it *DefaultSessionService) LoadSession(sessionID string) (*DefaultSessionContainer, error) {
	return nil, nil
}

// FlushSession is a stub function for no action
func (it *DefaultSessionService) FlushSession(sessionID string) error {
	return nil
}

// InterfaceSessionService implementation
// --------------------------------------

// GetName returns implementation name of session service
func (it *DefaultSessionService) GetName() string {
	return it.storage.GetStorageName()
}

// allocateSessionInstance is a routine to allocate a free memory for a session instance
//   - there is additional check to make sure of allowed sessions amount
func (it *DefaultSessionService) allocateSessionInstance(sessionInstance *DefaultSessionContainer) error {

	if ConstSessionUpdateTime == -1 && ConstSessionKeepInMemoryItems > 0 {
		_ = it.GC()

		numOfSessionsToClean := it.syncCount() - ConstSessionKeepInMemoryItems
		if numOfSessionsToClean >= 0 {
			it.syncLoop(
				func(item *DefaultSessionContainer) bool {
					if err := it.storage.FlushSession(item.id); err != nil {
						_ = env.ErrorDispatch(err)
					}
					numOfSessionsToClean--

					if numOfSessionsToClean <= 0 {
						return false
					}
					return true
				})
		}
	}

	it.syncSet(sessionInstance.id, sessionInstance)

	if ConstSessionUpdateTime <= 0 {
		if err := it.storage.FlushSession(sessionInstance.id); err != nil {
			_ = env.ErrorDispatch(err)
		}
	}

	return nil
}

// Get returns session object for given session id or nil of not currently exists
func (it *DefaultSessionService) Get(sessionID string, create bool) (api.InterfaceSession, error) {

	var resultError error
	var resultSession api.InterfaceSession

	if sessionID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "15fc38db-0848-4992-897e-82b93513f4c6", "blank session id")
	}
	replaceInstanceFlag := false

	// taking application instance of session
	sessionInstance := it.syncGet(sessionID)
	if sessionInstance != nil {
		// expiration check
		if time.Now().Sub(sessionInstance.GetUpdatedAt()).Seconds() >= ConstSessionLifeTime {
			sessionInstance = nil
		}
	}

	// session taking from storage for case of "immediate" mode and if no session in memory
	if sessionInstance == nil || ConstSessionUpdateTime == 0 {
		storedInstance, err := it.storage.LoadSession(sessionID)
		if storedInstance != nil && err == nil {
			// checking that loaded session is newer then we already have
			if sessionInstance == nil || storedInstance.UpdatedAt.After(sessionInstance.GetUpdatedAt()) {
				replaceInstanceFlag = true
				sessionInstance = storedInstance
			}
		}
	}

	// checking if session was found, if not - making new session for given id
	if create && sessionInstance == nil {
		sessionInstance = &DefaultSessionContainer{
			id:        sessionID,
			Data:      make(map[string]interface{}),
			UpdatedAt: time.Now()}

		err := it.allocateSessionInstance(sessionInstance)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		resultError = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "11670fe2-ee1c-45c9-a732-1349737b53f6", "new session created")
		replaceInstanceFlag = true
	}

	// updating application sessions
	if replaceInstanceFlag {
		it.syncSet(sessionID, sessionInstance)

		if ConstSessionUpdateTime <= 0 {
			if err := it.storage.FlushSession(sessionID); err != nil {
				_ = env.ErrorDispatch(err)
			}
		}
	}

	if sessionInstance != nil {
		resultSession = DefaultSession(sessionInstance.GetID())
	}

	return resultSession, resultError
}

// New initializes new session instance
func (it *DefaultSessionService) New() (api.InterfaceSession, error) {

	// receiving new session id
	sessionID, err := GenerateSessionID()
	if err != nil {
		return nil, err
	}

	// filling session structure
	sessionInstance := &DefaultSessionContainer{
		id:        sessionID,
		Data:      make(map[string]interface{}),
		UpdatedAt: time.Now()}

	err = it.allocateSessionInstance(sessionInstance)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return DefaultSession(sessionID), nil
}

// Touch updates session last modification time to current moment
func (it *DefaultSessionService) Touch(sessionID string) error {
	sessionInstance, err := it.Get(sessionID, false)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if sessionInstance != nil {
		if sessionInstance := it.syncGet(sessionID); sessionInstance != nil {
			sessionInstance.SetUpdatedAt(time.Now())

			if ConstSessionUpdateTime <= 0 {
				if err := it.storage.FlushSession(sessionID); err != nil {
					_ = env.ErrorDispatch(err)
				}
			}
		}
	}

	return nil
}

// Close makes current session instance expired
func (it *DefaultSessionService) Close(sessionID string) error {

	if it.syncGet(sessionID) != nil {

		// removing file
		filename := ConstStorageFolder + sessionID
		if _, err := os.Stat(filename); err == nil {
			err := os.Remove(filename)
			if err != nil {
				_ = env.ErrorDispatch(err)
			}
		}

		// releasing memory
		it.syncDel(sessionID)
	}

	return nil
}

// GetKey returns session value for a given key or nil - if not set
func (it *DefaultSessionService) GetKey(sessionID string, key string) interface{} {
	sessionInstance, err := it.Get(sessionID, false)
	if sessionInstance == nil || err != nil {
		return nil
	}

	if sessionInstance != nil {
		if sessionInstance := it.syncGet(sessionID); sessionInstance != nil {
			return sessionInstance.Get(key)
		}
	}
	return nil
}

// SetKey assigns value to session key
func (it *DefaultSessionService) SetKey(sessionID string, key string, value interface{}) {
	sessionInstance, _ := it.Get(sessionID, true)

	if sessionInstance != nil {
		if sessionInstance := it.syncGet(sessionID); sessionInstance != nil {
			sessionInstance.Set(key, value)
			sessionInstance.SetUpdatedAt(time.Now())

			if ConstSessionUpdateTime <= 0 {
				if err := it.storage.FlushSession(sessionID); err != nil {
					_ = env.ErrorDispatch(err)
				}
			}
		}
	}
}

// GC is a garbage collector for sessions, it removes expired sessions, flushes to storage, etc.
func (it *DefaultSessionService) GC() error {
	it.syncLoop(
		func(sessionInstance *DefaultSessionContainer) bool {
			secondsAfterLastUpdate := time.Now().Sub(sessionInstance.UpdatedAt).Seconds()

			// closing out of date sessions
			if secondsAfterLastUpdate >= ConstSessionLifeTime {
				if err := it.Close(sessionInstance.id); err != nil {
					_ = env.ErrorDispatch(err)
				}
				return false
			}

			// updating sessions information in a storage
			if secondsAfterLastUpdate > ConstSessionUpdateTime {
				err := it.storage.FlushSession(sessionInstance.id)
				if err != nil {
					_ = env.ErrorDispatch(err)
				}
			}
			return false
		})

	return nil
}

// IsEmpty checks if session contains data
func (it *DefaultSessionService) IsEmpty(sessionID string) bool {

	_, err := it.Get(sessionID, false)
	if err != nil {
		_ = env.ErrorDispatch(err)
	}

	if sessionInstance := it.syncGet(sessionID); sessionInstance != nil {
		return len(sessionInstance.Data) == 0
	}
	return true
}
