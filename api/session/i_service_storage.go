package session

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
