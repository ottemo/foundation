package session

import (
	"crypto/rand"
	"github.com/ottemo/foundation/env"
)

const (
	alphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890" // sessionID allowed symbols
)

// InitDefaultSessionService makes a new instance of DefaultSessionService
//   - makes internal fields initialization
func InitDefaultSessionService() *DefaultSessionService {
	sessionService := new(DefaultSessionService)
	sessionService.Sessions = make(map[string]*DefaultSessionContainer)
	sessionService.Storage = sessionService

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
