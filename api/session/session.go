package session

// InterfaceSession implementation
// -------------------------------

// GetID returns current session id
func (it DefaultSession) GetID() string {
	return it.id
}

// Get returns session value by a given key or nil - if not set
func (it DefaultSession) Get(key string) interface{} {
	return SessionService.GetKey(it.id, key)
}

// Set assigns value to session key
func (it DefaultSession) Set(key string, value interface{}) {
	SessionService.SetKey(it.id, key, value)
}

// IsEmpty checks if session contains data
func (it DefaultSession) IsEmpty() bool {
	return SessionService.IsEmpty(it.GetID())
}

// Touch updates session last modification time to current moment
func (it DefaultSession) Touch() error {
	return SessionService.Touch(it.id)
}

// Close makes current session instance expired
func (it DefaultSession) Close() error {
	return SessionService.Close(it.id)
}
