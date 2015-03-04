package session

// GetID returns current session id
func (it DefaultSession) GetID() string {
	return string(it)
}

// Get returns session value by a given key or nil - if not set
func (it DefaultSession) Get(key string) interface{} {
	return SessionService.GetKey(string(it), key)
}

// Set assigns value to session key
func (it DefaultSession) Set(key string, value interface{}) {
	SessionService.SetKey(string(it), key, value)
}

// Touch updates session last modification time to current moment
func (it DefaultSession) Touch() error {
	return SessionService.Touch(string(it))
}

// Close makes current session instance expired
func (it DefaultSession) Close() error {
	return SessionService.Close(string(it))
}
