// Package session is a default implementation of InterfaceSession
// declared in "github.com/ottemo/foundation/api" package
package session

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ALPHANUMERIC = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890" // sessionID allowed symbols

	ConstSessionCookieName = "OTTEMOSESSION" // cookie name which should contain sessionID

	ConstErrorModule = "api/session"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// Package global variables
var (
	Sessions      = make(map[string]*Session) // active session set
	sessionsMutex sync.RWMutex                // syncronization on Sessions variable modification

	gcRate int64 = 10 // garbage collection rate
)

// StartSession returns session object for request or creates new one
func StartSession(request *http.Request, responseWriter http.ResponseWriter) (*Session, error) {

	// check session-cookie
	cookie, err := request.Cookie(ConstSessionCookieName)
	if err == nil {
		// looking for cookie-based session
		sessionID := cookie.Value
		if session, ok := Sessions[sessionID]; ok == true {
			return session, nil
		}
	} else {
		if err != http.ErrNoCookie {
			return nil, err
		}
	}

	// cookie session is not set or expired, making new
	result, err := NewSession()
	if err != nil {
		return nil, err
	}

	// storing session id to cookie
	cookie = &http.Cookie{Name: ConstSessionCookieName, Value: result.GetID(), Path: "/"}
	http.SetCookie(responseWriter, cookie)

	return result, nil
}

// GetSessionByID returns session object for given id or nil
func GetSessionByID(sessionID string) (*Session, error) {
	if sessionInstance, present := Sessions[sessionID]; present {
		return sessionInstance, nil
	}
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e370399684f34c258996625154161d4b", "session not found")
}

// NewSession initializes new session
func NewSession() (*Session, error) {

	// receiving new session id
	sessionID, err := newSessionID()
	if err != nil {
		return nil, err
	}

	// initializing session structure
	sessionID = url.QueryEscape(sessionID)
	Sessions[sessionID] = &Session{
		id:     sessionID,
		values: make(map[string]interface{}),
		time:   time.Now()}

	// garbage collecting
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(gcRate))
	if err == nil && randomNumber.Cmp(big.NewInt(1)) == 0 {
		Gc()
	}

	return Sessions[sessionID], nil
}

// returns new session number
func newSessionID() (string, error) {
	sessionID := make([]byte, 32)
	if _, err := rand.Read(sessionID); err != nil {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "aab0911497264844984c772fb25dcb88", "can't generate sessionID")
	}

	for i := 0; i < 32; i++ {
		sessionID[i] = ALPHANUMERIC[sessionID[i]%62]
	}

	return string(sessionID), nil
}

// Gc removes expired sessions
func Gc() {
	for _, session := range Sessions {
		if time.Now().Sub(session.time).Seconds() > 3600 {
			session.Close()
		}
	}
}
