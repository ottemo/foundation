package session

import (
	"github.com/ottemo/foundation/env"
	"time"
)

// Session is a default implementer of InterfaceSession declared in
// "github.com/ottemo/foundation/api" package
type Session struct {
	id string

	values map[string]interface{}

	time time.Time
}

// GetID returns current session id
func (it *Session) GetID() string {
	return it.id
}

// Get returns session value by a given key or nil - if not set
func (it *Session) Get(key string) interface{} {
	if value, ok := it.values[key]; ok == true {
		return value
	}
	return nil
}

// Set assigns value to session key
func (it *Session) Set(key string, value interface{}) {
	it.values[key] = value
}

// Close clears session data
func (it *Session) Close() error {

	eventData := map[string]interface{}{"session": it, "sessionID": it.GetID()}
	env.Event("session.close", eventData)

	sessionsMutex.Lock()
	delete(Sessions, it.GetID())
	sessionsMutex.Unlock()

	return nil
}
